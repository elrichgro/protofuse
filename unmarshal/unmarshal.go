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

//	Package unmarshal provides functions to unmarshal marshaled protocol buffers
//	into a pfuse.ProtoTree structure for mounting.
package unmarshal

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/elrichgro/protofuse/fuse"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

var fileDesc *google_protobuf.FileDescriptorSet

func Unmarshal(fDesc *google_protobuf.FileDescriptorSet, packageName string, messageName string, buf [][]byte) (*pfuse.ProtoTree, error) {
	fileDesc = fDesc
	PT := &pfuse.ProtoTree{}
	PT.Dir.Nodes = []pfuse.TreeNode{}
	msg := fileDesc.GetMessage(packageName, messageName)
	if msg == nil {
		return nil, fmt.Errorf("Could not find message %s in package %s\n", messageName, packageName)
	}

	// unmarshal messages
	for i, buffer := range buf {
		PT.Dir.Nodes = append(PT.Dir.Nodes, pfuse.TreeNode{Name: fmt.Sprintf("Message_%d", i+1), FieldNumber: 0, Type: google_protobuf.FieldDescriptorProto_TYPE_MESSAGE})
		err := unmarshalMessage(msg, bytes.NewBuffer(buffer), &PT.Dir.Nodes[i], packageName)
		if err != nil {
			return nil, err
		}
	}
	return PT, nil
}

func unmarshalMessage(msg *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, packageName string) error {
	var repNum int32 = 0
	var m map[int32]int32 = make(map[int32]int32)
	dir := &pfuse.Dir{}

	for buf.Len() != 0 {
		tN := &pfuse.TreeNode{}
		wireType, fieldNumber, err := decodeKey(buf)
		if err != nil {
			return err
		}
		tN.FieldNumber = fieldNumber

		var field *google_protobuf.FieldDescriptorProto

		// check if field is an extension
		if isExtension(msg, fieldNumber) {
			_, field = fileDesc.FindExtensionByFieldNumber(packageName, msg.GetName(), fieldNumber)
			if field == nil {
				return fmt.Errorf("Could not find extension: %d, of message %s\n", fieldNumber, msg.GetName())
			}
		} else {
			field, err = getField(msg, fieldNumber)
		}

		// handle repeated fields
		if field.GetLabel() == google_protobuf.FieldDescriptorProto_LABEL_REPEATED {
			m[fieldNumber] += 1
			repNum = m[fieldNumber]
		} else {
			repNum = 0
		}

		// handle packed repeated types
		var packed bool = false
		if field.GetOptions() != nil {
			if field.GetOptions().GetPacked() {
				packed = true
			}
		}

		if packed {
			len, n := binary.Uvarint(buf.Bytes())
			if n <= 0 {
				return fmt.Errorf("decodeVarint n = %d", n)
			}
			buf.Next(n)
			p := bytes.NewBuffer(buf.Next(int(len)))
			for p.Len() != 0 {
				tN = &pfuse.TreeNode{}
				err = unmarshalPacked(field, p, tN, repNum)
				if err != nil {
					return err
				}
				m[fieldNumber] += 1
				repNum = m[fieldNumber]
				dir.Nodes = append(dir.Nodes, *tN)
			}
		} else {
			err = unmarshalField(wireType, field, buf, tN, repNum)
			if err != nil {
				return err
			}
			dir.Nodes = append(dir.Nodes, *tN)
		}
	}
	t.Node = dir

	return nil
}

func unmarshalField(wireType int8, field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, tN *pfuse.TreeNode, repNum int32) error {
	switch wireType {
	case 0:
		err := unmarshal0(field, buf, tN, repNum)
		if err != nil {
			return err
		}
	case 1:
		err := unmarshal1(field, buf, tN, repNum)
		if err != nil {
			return err
		}
	case 2:
		err := unmarshal2(field, buf, tN, repNum)
		if err != nil {
			return err
		}
	case 3:
		err := unmarshal3(field, buf, tN, repNum)
		if err != nil {
			return err
		}
	case 4:
		err := unmarshal4(field, buf, tN, repNum)
		if err != nil {
			return err
		}
	case 5:
		err := unmarshal5(field, buf, tN, repNum)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Invalid wire type: %d\n", wireType))
	}
	return nil
}

// Decodes a key and returns the wiretype and field number.
func decodeKey(buf *bytes.Buffer) (int8, int32, error) {
	x, n := binary.Uvarint(buf.Bytes())
	if n <= 0 {
		return 0, 0, fmt.Errorf("decodeVarint n = %d", n)
	}
	buf.Next(n)
	return int8(x & 7), int32(x >> 3), nil
}

func unmarshal0(field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	var contents string
	if rN != 0 {
		t.Name = fmt.Sprintf(field.GetName()+"_%d", rN)
	} else {
		t.Name = field.GetName()
	}
	t.Type = field.GetType()
	t.Label = field.GetLabel()

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_INT32:
		x, n, err := decodeInt32(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d", x)
	case google_protobuf.FieldDescriptorProto_TYPE_INT64:
		x, n, err := decodeInt64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d", x)
	case google_protobuf.FieldDescriptorProto_TYPE_UINT32:
		x, n, err := decodeUint32(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d", x)
	case google_protobuf.FieldDescriptorProto_TYPE_UINT64:
		x, n, err := decodeUint64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d", x)
	case google_protobuf.FieldDescriptorProto_TYPE_SINT32:
		x, n, err := decodeSint32(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d", x)
	case google_protobuf.FieldDescriptorProto_TYPE_SINT64:
		x, n, err := decodeSint64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		contents = fmt.Sprintf("%d", x)
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
		e, err := getEnumDescriptorProto(field.GetTypeName())
		if err != nil {
			return err
		}
		x, n, err := decodeUint64(buf.Bytes())
		buf.Next(n)
		if err != nil {
			return err
		}
		var found bool = false
		for _, value := range e.GetValue() {
			if uint64(value.GetNumber()) == x {
				contents = value.GetName()
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Invalid enum value: %d, for enum: %s", x, e.GetName())
		}
	default:
		return fmt.Errorf("Invalid wire type")
	}
	t.Node = &pfuse.File{contents}
	return nil
}

func unmarshal1(field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	p := make([]byte, 8)
	buf.Read(p)
	// Set file name
	if rN != 0 {
		t.Name = fmt.Sprintf(field.GetName()+"_%d", rN)
	} else {
		t.Name = field.GetName()
	}
	t.Type = field.GetType()
	t.Label = field.GetLabel()

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
		t.Node = &pfuse.File{fmt.Sprintf("%d", x)}
	case google_protobuf.FieldDescriptorProto_TYPE_SFIXED64:
		x, err := decodeSfixed64(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{fmt.Sprintf("%d", x)}
	default:
		t.Node = &pfuse.File{fmt.Sprintf("%x", p)}
	}
	return nil
}

func unmarshal2(field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	len, n := binary.Uvarint(buf.Bytes())
	if n <= 0 {
		return fmt.Errorf("decodeVarint n = %d", n)
	}
	buf.Next(n)
	p := make([]byte, len)
	buf.Read(p)
	// Set file name
	if rN != 0 {
		t.Name = fmt.Sprintf(field.GetName()+"_%d", rN)
	} else {
		t.Name = field.GetName()
	}
	t.Type = field.GetType()
	t.Label = field.GetLabel()

	switch *field.Type {
	case google_protobuf.FieldDescriptorProto_TYPE_STRING:
		t.Node = &pfuse.File{string(p)}
	case google_protobuf.FieldDescriptorProto_TYPE_BYTES:
		t.Node = &pfuse.File{hex.EncodeToString(p)}
	case google_protobuf.FieldDescriptorProto_TYPE_MESSAGE:
		var messageName string = field.GetTypeName()
		packageName := strings.Split(messageName, ".")[1]
		messageDesc, err := getDescriptorProto(messageName)
		if err != nil {
			return err
		}
		unmarshalMessage(messageDesc, bytes.NewBuffer(p), t, packageName)
	default:
		t.Node = &pfuse.File{string(p)}
	}

	return nil
}

func unmarshal3(field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	return errors.New("Groups are not supported")
}

func unmarshal4(field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	return errors.New("Groups are not supported")
}

func unmarshal5(field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	p := make([]byte, 4)
	buf.Read(p)
	// Set file name
	if rN != 0 {
		t.Name = fmt.Sprintf(field.GetName()+"_%d", rN)
	} else {
		t.Name = field.GetName()
	}
	t.Type = field.GetType()
	t.Label = field.GetLabel()

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
		t.Node = &pfuse.File{fmt.Sprintf("%d", x)}
	case google_protobuf.FieldDescriptorProto_TYPE_SFIXED32:
		x, err := decodeSfixed32(p)
		if err != nil {
			return err
		}
		t.Node = &pfuse.File{fmt.Sprintf("%d", x)}
	default:
		t.Node = &pfuse.File{fmt.Sprintf("%x", p)}
	}

	return nil
}

func unmarshalPacked(field *google_protobuf.FieldDescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode, rN int32) error {
	ft := field.GetType()
	if ft == google_protobuf.FieldDescriptorProto_TYPE_INT32 || ft == google_protobuf.FieldDescriptorProto_TYPE_INT64 ||
		ft == google_protobuf.FieldDescriptorProto_TYPE_UINT32 || ft == google_protobuf.FieldDescriptorProto_TYPE_UINT64 ||
		ft == google_protobuf.FieldDescriptorProto_TYPE_SINT32 || ft == google_protobuf.FieldDescriptorProto_TYPE_SINT64 ||
		ft == google_protobuf.FieldDescriptorProto_TYPE_BOOL || ft == google_protobuf.FieldDescriptorProto_TYPE_ENUM {
			t.FieldNumber = field.GetNumber()
			unmarshal0(field, buf, t, rN)
	} else if ft == google_protobuf.FieldDescriptorProto_TYPE_FIXED64 || ft == google_protobuf.FieldDescriptorProto_TYPE_SFIXED64 ||
		ft == google_protobuf.FieldDescriptorProto_TYPE_DOUBLE {
			t.FieldNumber = field.GetNumber()
			unmarshal1(field, buf, t, rN)
	} else if ft == google_protobuf.FieldDescriptorProto_TYPE_FIXED32 || ft == google_protobuf.FieldDescriptorProto_TYPE_SFIXED32 ||
		ft == google_protobuf.FieldDescriptorProto_TYPE_FLOAT {
			t.FieldNumber = field.GetNumber()
			unmarshal5(field, buf, t, rN)
	} else {
		return fmt.Errorf("Invalid packed type\n")
	}
	return nil
}

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

// Finds the google_protobuf.DescriptorProto for the message name.
func getDescriptorProto(name string) (*google_protobuf.DescriptorProto, error) {
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
						return nil, fmt.Errorf("Cannot find message: %s", name)
					}
				}
			}
		}
		return nil, fmt.Errorf("Cannot find message: %s", name)
	}
	return nil, fmt.Errorf("Message name not fully qualified: %s", name)
}

// Gets the google_protobuf.EnumDescriptorProto for name
func getEnumDescriptorProto(name string) (*google_protobuf.EnumDescriptorProto, error) {
	if string(name[0]) == "." {
		s := strings.Split(name, ".")
		slen := len(s)
		for _, file := range fileDesc.File {
			if file.GetPackage() == s[1] {
				if slen <= 3 {
					return getEnum(s[2], file.GetEnumType())
				}
				for _, message := range file.MessageType {
					if message.GetName() == s[2] {
						if slen <= 4 {
							return getEnum(s[3], message.GetEnumType())
						}
						var m *google_protobuf.DescriptorProto = message

						for i := 3; i < slen-1; i++ {
							for _, d := range m.GetNestedType() {
								if d.GetName() == s[i] {
									if i >= slen-2 {
										return getEnum(s[slen-1], d.GetEnumType())
									}
									m = d
								}
							}
						}
						return nil, fmt.Errorf("Cannot find enum: %s", name)
					}
				}
			}
		}
		return nil, fmt.Errorf("Cannot find enum: %s", name)
	}
	return nil, fmt.Errorf("Enum name not fully qualified: %s", name)
}

func getEnum(name string, enums []*google_protobuf.EnumDescriptorProto) (*google_protobuf.EnumDescriptorProto, error) {
	for _, enum := range enums {
		if enum.GetName() == name {
			return enum, nil
		}
	}
	return nil, fmt.Errorf("Cannot find enum: %s", name)
}

func isExtension(msg *google_protobuf.DescriptorProto, fieldNumber int32) bool {
	if len(msg.GetExtensionRange()) > 0 {
		for _, r := range msg.GetExtensionRange() {
			if fieldNumber >= r.GetStart() && fieldNumber <= r.GetEnd() {
				return true
			}
		}
	}
	return false
}

func getField(msg *google_protobuf.DescriptorProto, fieldNumber int32) (*google_protobuf.FieldDescriptorProto, error) {
	for _, field := range msg.GetField() {
		if field.GetNumber() == fieldNumber {
			return field, nil
		}
	}
	return nil, fmt.Errorf("Could not find field %d in message %s\n", fieldNumber, msg.GetName())
}
