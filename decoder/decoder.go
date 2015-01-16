// decoder for marshalled protocol buffers. Builds a FUSE filsesystem structure.
package decoder

import (
	"bytes"
	"strings"
	"fmt"
	"os"
	"errors"
	"encoding/binary"
	"math"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"elrich/protofuse/fuse"
)

var fileDesc *google_protobuf.FileDescriptorProto

func Decode(fDesc *google_protobuf.FileDescriptorProto, desc *google_protobuf.DescriptorProto, buf *bytes.Buffer) (*pfuse.ProtoTree, error) {
	fileDesc = fDesc
	PT := &pfuse.ProtoTree{}
	PT.Dir.Nodes = []pfuse.TreeNode{pfuse.TreeNode{}}
	err := decodeMessage(desc, buf, &PT.Dir.Nodes[0], 0)
	return PT, err
}

func CheckError(err error) {
   if err != nil {
       fmt.Println(err.Error())
       os.Exit(-1)
   }
}

func decodeMessage(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	var repNum int32 = 0
	var m map[uint64]int32 = make(map[uint64]int32)
	dir := &pfuse.Dir{}
	for buf.Len() != 0 {
		tN := &pfuse.TreeNode{}
	    wireType, fieldNumber := decodeKey(buf)
	    tN.FieldNumber = fieldNumber

	    if *desc.Field[fieldNumber-1].Label == google_protobuf.FieldDescriptorProto_LABEL_REPEATED {
	    	m[fieldNumber] += 1
	    	repNum = m[fieldNumber]
	    } else {
	    	repNum = 0
	    }

	    switch wireType {
	    case 0:
	    	decode0(desc, buf, tN, repNum)
	    case 1:
	    	decode1(desc, buf, tN, repNum)
	    case 2:
	    	decode2(desc, buf, tN, repNum)
	    case 3:
	    	decode3(desc, buf, tN, repNum)
		case 4:
			decode4(desc, buf, tN, repNum)
		case 5:
			decode5(desc, buf, tN, repNum)
		default:
			return errors.New(fmt.Sprintf("Invalid wire type: %d\n", wireType))
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

	return nil
}

func decodeKey(buf *bytes.Buffer) (int8, uint64) {
  	x := decodeVarint(buf)
  	return int8(x & 7), x >> 3
}

func decodeVarint(buf *bytes.Buffer) uint64{
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
func decode0(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
  	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
	var contents string
	x := decodeVarint(buf)
	if rN != 0 {
		t.Name = fmt.Sprintf(*field.Name + "%d", rN)
	} else {
		t.Name = *field.Name
	}
	t.Type = *field.Type

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_INT32:
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_INT64:
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_UINT32:
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_UINT64:
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_SINT32:
		// TODO: decodeZigZag
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_SINT64:
		// TODO: decodeZigzag
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_BOOL:
		if x == 0 {
			contents = "False"
		} else {
			contents = "True"
		}
	case google_protobuf.FieldDescriptorProto_TYPE_ENUM:
		// TODO: find enum and get string value
		contents = fmt.Sprintf("%d\n", x)
	default:
		contents = fmt.Sprintf("%d\n", x)
	}
	t.Node = &pfuse.File{contents}
	return nil
}

func decode1(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
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

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_DOUBLE:
		x, err := decodeDouble(p)
		CheckError(err)
		t.Node = &pfuse.File{fmt.Sprintf("%.6f", x)}
	case google_protobuf.FieldDescriptorProto_TYPE_FIXED64:
		x, err := decodeFixed64(p)
		CheckError(err)
		t.Node = &pfuse.File{string(x)}
	case google_protobuf.FieldDescriptorProto_TYPE_SFIXED64:
		x, err := decodeSFixed64(p)
		CheckError(err)
		t.Node = &pfuse.File{string(x)}
	default:
		t.Node = &pfuse.File{string(p)}
	}
	return nil
}

// TODO: add switch block for different field types
func decode2(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
  	len := decodeVarint(buf)
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

  	switch *field.Type {
  	case google_protobuf.FieldDescriptorProto_TYPE_STRING:
  		t.Node = &pfuse.File{string(p)}
  	case google_protobuf.FieldDescriptorProto_TYPE_BYTES:
  		var buffer bytes.Buffer
  		for _, b := range p {
  			buffer.WriteString(fmt.Sprintf(" %x", b))
  		}
  		t.Node = &pfuse.File{buffer.String()}
  	case google_protobuf.FieldDescriptorProto_TYPE_MESSAGE:
  		var messageName string = (*field.TypeName)[strings.LastIndex(*field.TypeName, ".")+1:]
  		messageDesc, err := GetDescriptorProto(messageName, desc)
  		CheckError(err)
    	decodeMessage(messageDesc, bytes.NewBuffer(p), t, rN)
    default:
    	t.Node = &pfuse.File{string(p)}
  	// Packed repeat types?
  	}

  	return nil
}

func decode3(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	return errors.New("Groups are not supported")
}

func decode4(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	return errors.New("Groups are not supported")
}

func decode5(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
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

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_FLOAT:
		x, err := decodeFloat(p)
		CheckError(err)
		t.Node = &pfuse.File{fmt.Sprintf("%.6f", x)}
	case google_protobuf.FieldDescriptorProto_TYPE_FIXED32:
		x, err := decodeFixed32(p)
		CheckError(err)
		t.Node = &pfuse.File{string(x)}
	case google_protobuf.FieldDescriptorProto_TYPE_SFIXED32:
		x, err := decodeSFixed32(p)
		CheckError(err)
		t.Node = &pfuse.File{string(x)}
	default:
		t.Node = &pfuse.File{string(p)}
	}

	return nil
}

func decodeFloat(buf []byte) (float32, error) {
	b := binary.LittleEndian.Uint64(buf)
    f := float32(math.Float64frombits(b))
    return f, nil
}

func decodeFixed32(buf []byte) (uint32, error) {
	var x uint32
	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
	return x, err
}

func decodeSFixed32(buf []byte) (int32, error) {
	var x int32
	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
	return x, err
}

func decodeDouble(buf []byte) (float64, error) {
	b := binary.LittleEndian.Uint64(buf)
    f := math.Float64frombits(b)
    return f, nil
}

func decodeFixed64(buf []byte) (uint64, error) {
	var x uint64
	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
	return x, err
}

func decodeSFixed64(buf []byte) (int64, error) {
	var x int64
	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
	return x, err
}

// TODO: update to work with fully qualified message names
func GetDescriptorProto(name string, messageDesc *google_protobuf.DescriptorProto) (*google_protobuf.DescriptorProto, error) {
	if messageDesc != nil {
		for _, message := range messageDesc.NestedType {
			if *message.Name == name {
				return message, nil
			}
		}
	}

	for _, message := range fileDesc.MessageType {
		if *message.Name == name {
			return message, nil
		}
	}

	//TODO: throw error
	return nil, errors.New("Cannot find message: " + name)
}
