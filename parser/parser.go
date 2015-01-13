// Parser for marshalled protocol buffers. Builds a FUSE filsesystem structure.
package protoparser

import (
	"bytes"
	"strings"
	"fmt"

	// "github.com/gogo/protobuf/proto"
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
	dir := &pfuse.Dir{}
	var repNum int32 = 0
	var m map[uint64]int32 = make(map[uint64]int32)
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

		// parseField(desc, buf, tN)
		dir.Nodes = append(dir.Nodes, *tN)
	} 

	if rN != 0 {
		t.Name = fmt.Sprintf(desc.GetName() + "%d", rN)
	} else {
		t.Name = desc.GetName()
	}
	t.Node = dir
}

// func parseField(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode) {
//     wireType, fieldNumber := parseKey(buf)
//     t.FieldNumber = fieldNumber
//     switch wireType {
//     case 0:
//     	parse0(desc, buf, t)
//     case 1:
//     	parse1(desc, buf, t)
//     case 2:
//     	parse2(desc, buf, t)
//     case 3:
//     	parse3(desc, buf, t)
// 	case 4:
// 		parse4(desc, buf, t)
// 	case 5:
// 		parse5(desc, buf, t)
//     }
// }

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
  // fmt.Println(x)
  return x
}

func parse0(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	x := parseVarint(buf)
	if rN != 0 {
		t.Name = fmt.Sprintf(*desc.Field[t.FieldNumber-1].Name + "%d", rN)
	} else {
		t.Name = *desc.Field[t.FieldNumber-1].Name
	}
	t.Node = &pfuse.File{fmt.Sprintf("%d\n", x)}
}

func parse1(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	fmt.Println("Fixed types not yet supported")
}

func parse2(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
  	len := parseVarint(buf)
  	p := make([]byte, len)
  	buf.Read(p)
  	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
  	if rN != 0 {
  		t.Name = fmt.Sprintf(*field.Name + "%d", rN)
  	} else {
  		t.Name = *field.Name
  	}
  	if desc != nil && *field.Type == google_protobuf.FieldDescriptorProto_TYPE_MESSAGE {
  		var messageName string = (*field.TypeName)[strings.LastIndex(*field.TypeName, ".")+1:]
    	parseMessage(GetDescriptorProto(messageName, desc), bytes.NewBuffer(p), t, rN)
  	} else {
    	// fmt.Println("Parse string:", string(p))
    	t.Node = &pfuse.File{string(p)}
  	}
}

func parse3(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	fmt.Println("Groups not supported")
}

func parse4(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	fmt.Println("Groups not supported")
}

func parse5(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) {
	fmt.Println("Fixed types not yet supported")
}

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

	//TODO: else throw error message could not be found
	fmt.Println("Can't find message\n")
	return nil
}
