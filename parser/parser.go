// Parser for marshalled protocol buffers. Builds a FUSE filsesystem structure.
package protoparser

import (
	"bytes"
	"strings"
	"fmt"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"elrich/protofuse/fuse"
)

var fileDesc *google_protobuf.FileDescriptorProto

func Parse(fDesc *google_protobuf.FileDescriptorProto, desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, PT *pfuse.ProtoTree) {
	fileDesc = fDesc
	PT.Dir.Nodes = []pfuse.TreeNode{pfuse.TreeNode{}}
	parseMessage(desc, buf, &PT.Dir.Nodes[0], 0)
}

func parseMessage(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	var repNum int32 = 0
	var m map[uint64]int32 = make(map[uint64]int32)
	dir := &pfuse.Dir{}
	for buf.Len() != 0 {
		tN := &pfuse.TreeNode{}
	    wireType, fieldNumber := parseKey(buf)
	    tN.FieldNumber = fieldNumber

	    if *desc.Field[fieldNumber-1].Label == google_protobuf.FieldDescriptorProto_LABEL_REPEATED {
	    	m[fieldNumber] += 1
	    	repNum = m[fieldNumber]
	    } else {
	    	repNum = 0
	    }

	    switch wireType {
	    case 0:
	    	parse0(desc, buf, tN, repNum)
	    case 1:
	    	parse1(desc, buf, tN, repNum)
	    case 2:
	    	parse2(desc, buf, tN, repNum)
	    case 3:
	    	parse3(desc, buf, tN, repNum)
		case 4:
			parse4(desc, buf, tN, repNum)
		case 5:
			parse5(desc, buf, tN, repNum)
	    }
		dir.Nodes = append(dir.Nodes, *tN)
	} 

  	// Set directory name
	if rN != 0 {
		t.Name = fmt.Sprintf(desc.GetName() + "%d", rN)
	} else {
		t.Name = desc.GetName()
	}
	t.Node = dir
}

func parseKey(buf *bytes.Buffer) (int8, uint64) {
  x := parseVarint(buf)
  return int8(x & 7), x >> 3
}

func parseVarint(buf *bytes.Buffer) uint64{
  b := byte(0xff)
  var x uint64 = 0
  i := uint(0)
  for b >> 7 > 0 {
    b,_ = buf.ReadByte()
    y := uint64(b)
    y = y & 0x7f
    y = y << (7*i)
    x = x | y
    i += 1
  }
  return x
}

// TODO: add switch block for different varint types
func parse0(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
  	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
	x := parseVarint(buf)
	if rN != 0 {
		t.Name = fmt.Sprintf(*field.Name + "%d", rN)
	} else {
		t.Name = *field.Name
	}
	t.Node = &pfuse.File{fmt.Sprintf("%d\n", x)}

	t.Type = *field.Type

	// switch *field.Type {
	// case google_protobuf.FieldDescriptorProto_TYPE_INT32:
	// case google_protobuf.FieldDescriptorProto_TYPE_INT64:
	// case google_protobuf.FieldDescriptorProto_TYPE_UINT32:
	// case google_protobuf.FieldDescriptorProto_TYPE_UINT64:
	// case google_protobuf.FieldDescriptorProto_TYPE_SINT32:
	// case google_protobuf.FieldDescriptorProto_TYPE_SINT64:
	// case google_protobuf.FieldDescriptorProto_TYPE_BOOL:
	// case google_protobuf.FieldDescriptorProto_TYPE_ENUM:
	// }
}

// TODO: finish implementation
func parse1(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	// fmt.Println("Fixed types not yet supported")
  	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
  	p := make([]byte, 8)
  	buf.Read(p)
  	// Set file name
  	if rN != 0 {
  		t.Name = fmt.Sprintf(*field.Name + "%d", rN)
  	} else {
  		t.Name = *field.Name
  	}

	t.Type = *field.Type


	// switch *field.Type {
	// case google_protobuf.FieldDescriptorProto_TYPE_DOUBLE:

	// case google_protobuf.FieldDescriptorProto_TYPE_FIXED64:

	// case google_protobuf.FieldDescriptorProto_TYPE_SFIXED64:

	// }
}

// TODO: add switch block for different field types
func parse2(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
  	len := parseVarint(buf)
  	p := make([]byte, len)
  	buf.Read(p)
  	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
  	// Set file name
  	if rN != 0 {
  		t.Name = fmt.Sprintf(*field.Name + "%d", rN)
  	} else {
  		t.Name = *field.Name
  	}
	t.Type = *field.Type


  	// switch *field.Type {
  	// case google_protobuf.FieldDescriptorProto_TYPE_STRING:

  	// case google_protobuf.FieldDescriptorProto_TYPE_BYTES:

  	// case google_protobuf.FieldDescriptorProto_TYPE_MESSAGE:

  	// // Packed repeat types?
  	// }
  	if *field.Type == google_protobuf.FieldDescriptorProto_TYPE_MESSAGE {
  		var messageName string = (*field.TypeName)[strings.LastIndex(*field.TypeName, ".")+1:]
    	parseMessage(GetDescriptorProto(messageName, desc), bytes.NewBuffer(p), t, rN)
  	} else {
    	t.Node = &pfuse.File{string(p)}
  	}
}

func parse3(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	fmt.Println("Groups not supported")
}

func parse4(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	fmt.Println("Groups not supported")
}

// TODO: finish switch block implementation
func parse5(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	// fmt.Println("Fixed types not yet supported")
	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
  	p := make([]byte, 4)
  	buf.Read(p)
  	// Set file name
  	if rN != 0 {
  		t.Name = fmt.Sprintf(*field.Name + "%d", rN)
  	} else {
  		t.Name = *field.Name
  	}
	t.Type = *field.Type

	// switch *field.Type {
	// case google_protobuf.FieldDescriptorProto_TYPE_FLOAT:

	// case google_protobuf.FieldDescriptorProto_TYPE_FIXED32:

	// case google_protobuf.FieldDescriptorProto_TYPE_SFIXED32:

	// }
}

// TODO: update to work with fully qualified message names
func GetDescriptorProto(name string, messageDesc *google_protobuf.DescriptorProto) *google_protobuf.DescriptorProto {
	if messageDesc != nil {
		for _, message := range messageDesc.NestedType {
			if *message.Name == name {
				return message
			}
		}
	}

	for _, message := range fileDesc.MessageType {
		if *message.Name == name {
			return message
		}
	}

	//TODO: throw error
	fmt.Println("Can't find message\n")
	return nil
}
