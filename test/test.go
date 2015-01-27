package test;

import (
	"os"
	"io"

	// "github.com/elrichgro/protofuse/unmarshal/test"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
    "github.com/gogo/protobuf/proto"
)

type f struct {
	F1 string 
	F2 []int32
	F3 int64 
	F4 uint32 
	F5 uint64 
	F6 int32
	F7 bool
	F8 uint64 
	F9 int64 
	F10 float64 
	F11 []byte 
	F12 Bar
	F13 uint32 
	F14 int32 
	F15 float32
	F16 []*FooBaz
}

func GenerateFull() ([]byte, *google_protobuf.FileDescriptorSet, string, string, error) {
	name := "BAR"
	names := []string{"name", "name2"}
	id := int32(123)
	f := f{"one", []int32{1, 2, 3, 4}, 3, 4, 5, 6, true, 8, 9, 10.0, []byte{11, 11}, Bar{Id: &id}, 13, 14, 15.0,
	[]*FooBaz{&FooBaz{F1: &names[0], F2: Foo_e1.Enum(), F3: &FooBazFoobaz{Name: &names[0]}}, 
		&FooBaz{F1: &names[1], F2: Foo_e2.Enum(), F3: &FooBazFoobaz{Name: &names[1]}}}}
	foo := &Foo{F1:&f.F1, F2:f.F2, F3:&f.F3, F4:&f.F4, F5:&f.F5, F6:&f.F6, F7:&f.F7, F8:&f.F8,
		F9:&f.F9, F10:&f.F10, F11:f.F11, F12:&f.F12, F13:&f.F13, F14:&f.F14, F15:&f.F15, F16:f.F16}

	err := proto.SetExtension(foo.GetF12(), E_Name, &name)
	if err != nil {
		return nil, nil, "", "", err
	}
	err = proto.SetExtension(foo, E_F121, &id)
	if err != nil {
		return nil, nil, "", "", err
	}

	buf, err := proto.Marshal(foo)
	if err != nil {
		return nil, nil, "", "", err
	}

	fileDesc, err := getFileDescriptorSet("../test/test.desc")
	if err != nil {
		return nil, nil, "", "", err
	}

	packageName := "test"
	messageName := "foo"

	return buf, fileDesc, packageName, messageName, nil
}

func getFileDescriptorSet(filename string) (*google_protobuf.FileDescriptorSet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, fi.Size())
	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return nil, err
	}
	file.Close()

	fDesc := &google_protobuf.FileDescriptorSet{}
	err = proto.Unmarshal(buffer, fDesc)
	if err != nil {
		return nil, err
	}

	return fDesc, nil
}