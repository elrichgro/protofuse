//  Copyright 2015 Elrich Groenewald
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package unmarshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"elrich/protofuse/fuse"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

var fileDesc *google_protobuf.FileDescriptorSet

func Unmarshal(fDesc *google_protobuf.FileDescriptorSet, desc *google_protobuf.DescriptorProto, buf [][]byte) (*pfuse.ProtoTree, error) {
	fileDesc = fDesc
	PT := &pfuse.ProtoTree{}
	PT.Dir.Nodes = []pfuse.TreeNode{}

	if len(buf) > 1 {
		for i, buffer := range buf {
			PT.Dir.Nodes = append(PT.Dir.Nodes, pfuse.TreeNode{})
			err := unmarshalMessage(desc, bytes.NewBuffer(buffer), &PT.Dir.Nodes[i], int32(i+1))
			if err != nil {
				return nil, err
			}
		}
	} else {
		PT.Dir.Nodes = append(PT.Dir.Nodes, pfuse.TreeNode{})
		err := unmarshalMessage(desc, bytes.NewBuffer(buf[0]), &PT.Dir.Nodes[0], 0)
		if err != nil {
			return nil, err
		}
	}
	return PT, nil
}

func unmarshalMessage(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	var repNum int32 = 0
	var m map[uint64]int32 = make(map[uint64]int32)
	dir := &pfuse.Dir{}
	for buf.Len() != 0 {
		tN := &pfuse.TreeNode{}
		wireType, fieldNumber, err := decodeKey(buf)
		if err != nil {
			return err
		}
		tN.FieldNumber = fieldNumber

		if *desc.Field[fieldNumber-1].Label == google_protobuf.FieldDescriptorProto_LABEL_REPEATED {
			m[fieldNumber] += 1
			repNum = m[fieldNumber]
		} else {
			repNum = 0
		}

		switch wireType {
		case 0:
			err := unmarshal0(desc, buf, tN, repNum)
			if err != nil {
				return err
			}
		case 1:
			err := unmarshal1(desc, buf, tN, repNum)
			if err != nil {
				return err
			}
		case 2:
			err := unmarshal2(desc, buf, tN, repNum)
			if err != nil {
				return err
			}
		case 3:
			err := unmarshal3(desc, buf, tN, repNum)
			if err != nil {
				return err
			}
		case 4:
			err := unmarshal4(desc, buf, tN, repNum)
			if err != nil {
				return err
			}
		case 5:
			err := unmarshal5(desc, buf, tN, repNum)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("Invalid wire type: %d\n", wireType))
		}
		dir.Nodes = append(dir.Nodes, *tN)
	}

	// Set directory name
	if rN != 0 {
		t.Name = fmt.Sprintf(desc.GetName()+"%d", rN)
	} else {
		t.Name = desc.GetName()
	}
	t.Node = dir

	return nil
}

func decodeKey(buf *bytes.Buffer) (int8, uint64, error) {
	x, n := binary.Uvarint(buf.Bytes())
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	buf.Next(n)
	return int8(x & 7), x >> 3, nil
}

// func decodeVarint(buf *bytes.Buffer) uint64{
//   	b := byte(0xff)
//   	var x uint64 = 0
//   	i := uint(0)
//   	for b >> 7 > 0 {
//   	  	b,_ = buf.ReadByte()
//   	  	y := uint64(b)
//   	  	y = y & 0x7f
//   	  	y = y << (7*i)
//   	  	x = x | y
//   	  	i += 1
//  	}
//   	return x
// }

func unmarshal0(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
	var contents string
	if rN != 0 {
		t.Name = fmt.Sprintf(*field.Name+"%d", rN)
	} else {
		t.Name = *field.Name
	}
	t.Type = *field.Type

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_INT32:
		x, n, err := decodeInt32(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_INT64:
		x, n, err := decodeInt64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_UINT32:
		x, n, err := decodeUint32(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_UINT64:
		x, n, err := decodeUint64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_SINT32:
		x, n, err := decodeSint32(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_SINT64:
		x, n, err := decodeSint64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d\n", x)
	case google_protobuf.FieldDescriptorProto_TYPE_BOOL:
		x, n, err := decodeBool(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		if x {
			contents = "True"
		} else {
			contents = "False"
		}
	case google_protobuf.FieldDescriptorProto_TYPE_ENUM:
		x, n, err := decodeUint64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d\n", x)
	default:
		return fmt.Errorf("Invalid wire type")
	}
	t.Node = &pfuse.File{contents}
	return nil
}

func unmarshal1(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
	p := make([]byte, 8)
	buf.Read(p)
	// Set file name
	if rN != 0 {
		t.Name = fmt.Sprintf(*field.Name+"%d", rN)
	} else {
		t.Name = *field.Name
	}
	t.Type = *field.Type

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_DOUBLE:
		x, err := decodeFloat64(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{fmt.Sprintf("%.6f", x)}
	case google_protobuf.FieldDescriptorProto_TYPE_FIXED64:
		x, err := decodeFixed64(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{string(x)}
	case google_protobuf.FieldDescriptorProto_TYPE_SFIXED64:
		x, err := decodeSfixed64(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{string(x)}
	default:
		t.Node = &pfuse.File{string(p)}
	}
	return nil
}

func unmarshal2(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	len, n := binary.Uvarint(buf.Bytes())
	if n <= 0 {
		return fmt.Errorf("decodeVarint n = %d", n)
	}
	buf.Next(n)
	p := make([]byte, len)
	buf.Read(p)
	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
	// Set file name
	if rN != 0 {
		t.Name = fmt.Sprintf(*field.Name+"%d", rN)
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
		var messageName string = (*field.TypeName)
		messageDesc, err := GetDescriptorProto(messageName)
		if err != nil {
			return err
		}
		unmarshalMessage(messageDesc, bytes.NewBuffer(p), t, rN)
	default:
		t.Node = &pfuse.File{string(p)}
		// Packed repeat types?
	}

	return nil
}

func unmarshal3(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	return errors.New("Groups are not supported")
}

func unmarshal4(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	return errors.New("Groups are not supported")
}

func unmarshal5(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	var field *google_protobuf.FieldDescriptorProto = desc.Field[t.FieldNumber-1]
	p := make([]byte, 4)
	buf.Read(p)
	// Set file name
	if rN != 0 {
		t.Name = fmt.Sprintf(*field.Name+"%d", rN)
	} else {
		t.Name = *field.Name
	}
	t.Type = *field.Type

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_FLOAT:
		x, err := decodeFloat32(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{fmt.Sprintf("%.6f", x)}
	case google_protobuf.FieldDescriptorProto_TYPE_FIXED32:
		x, err := decodeFixed32(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{string(x)}
	case google_protobuf.FieldDescriptorProto_TYPE_SFIXED32:
		x, err := decodeSfixed32(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{string(x)}
	default:
		t.Node = &pfuse.File{string(p)}
	}

	return nil
}

// func decodeFloat(buf []byte) (float32, error) {
// 	b := binary.LittleEndian.Uint64(buf)
//     f := float32(math.Float64frombits(b))
//     return f, nil
// }

// func decodeFixed32(buf []byte) (uint32, error) {
// 	var x uint32
// 	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
// 	return x, err
// }

// func decodeSFixed32(buf []byte) (int32, error) {
// 	var x int32
// 	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
// 	return x, err
// }

// func decodeDouble(buf []byte) (float64, error) {
// 	b := binary.LittleEndian.Uint64(buf)
//     f := math.Float64frombits(b)
//     return f, nil
// }

// func decodeFixed64(buf []byte) (uint64, error) {
// 	var x uint64
// 	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
// 	return x, err
// }

// func decodeSFixed64(buf []byte) (int64, error) {
// 	var x int64
// 	err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &x)
// 	return x, err
// }

func decodeBool(buf []byte) (bool, int, error) {
	v, n := binary.Uvarint(buf)
	if n <= 0 {
		return false, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	return v != 0, n, nil
}

func decodeFloat64(buf []byte) (float64, error) {
	if len(buf) < 8 {
		return 0, fmt.Errorf("Double: buffer too short")
	}
	return *(*float64)(unsafe.Pointer(&buf[0])), nil
}

func decodeFloat32(buf []byte) (float32, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("Float: buffer too short")
	}
	return *(*float32)(unsafe.Pointer(&buf[0])), nil
}

func decodeInt64(buf []byte) (int64, int, error) {
	v, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	return int64(v), n, nil
}

func decodeUint64(buf []byte) (uint64, int, error) {
	v, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	return v, n, nil
}

func decodeInt32(buf []byte) (int32, int, error) {
	v, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	return int32(v), n, nil
}

func decodeFixed64(buf []byte) (uint64, error) {
	if len(buf) < 8 {
		return 0, fmt.Errorf("decodeFixed64: buffer too short")
	}
	return *(*uint64)(unsafe.Pointer(&buf[0])), nil
}

func decodeFixed32(buf []byte) (uint32, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("decodeFixed32: buffer too short")
	}
	return *(*uint32)(unsafe.Pointer(&buf[0])), nil
}

func decodeUint32(buf []byte) (uint32, int, error) {
	v, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	return uint32(v), n, nil
}

func decodeSfixed32(buf []byte) (int32, error) {
	if len(buf) < 4 {
		return 0, fmt.Errorf("decodeDouble: buffer too short")
	}
	return *(*int32)(unsafe.Pointer(&buf[0])), nil
}

func decodeSfixed64(buf []byte) (int64, error) {
	if len(buf) < 8 {
		return 0, fmt.Errorf("decodeSfixed64: buffer too short")
	}
	return *(*int64)(unsafe.Pointer(&buf[0])), nil
}

func decodeSint32(buf []byte) (int32, int, error) {
	v, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	return int32((uint32(v) >> 1) ^ uint32(((v&1)<<31)>>31)), n, nil
}

func decodeSint64(buf []byte) (int64, int, error) {
	v, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	return int64((v >> 1) ^ uint64((int64(v&1)<<63)>>63)), n, nil
}

func GetDescriptorProto(name string) (*google_protobuf.DescriptorProto, error) {
	if string(name[0]) == "." {
		s := strings.Split(name, ".")
		slen := len(s)
		for _, file := range fileDesc.File {
			if file.GetPackage() == s[1] {
				for _, message := range file.MessageType {
					if message.GetName() == s[2] {
						if slen <= 3 {
							return message, nil
						}
						var m *google_protobuf.DescriptorProto = message

						for i := 3; i < slen; i++ {
							for _, d := range m.GetNestedType() {
								if d.GetName() == s[i] {
									if i >= slen-1 {
										return d, nil
									}
									m = d
								}
							}
						}
						return nil, fmt.Errorf("Cannot find message1: %s", name)
					}
				}
			}
		}
		return nil, fmt.Errorf("Cannot find message: %s", name)
	}
	return nil, fmt.Errorf("Message name not fully qualified: %s", name)
}
