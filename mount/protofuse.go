package protofuse

import (
	"log"
	"errors"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"elrich/protofuse/unmarshal"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

func Mount(marshaled []byte, fileDesc *google_protobuf.FileDescriptorSet, messageName string, mountPoint string) error {
	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("ProtoFuse"),
		fuse.Subtype("protofuse"),
		fuse.LocalVolume(),
		fuse.VolumeName("ProtoFuse"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

  	desc := &google_protobuf.DescriptorProto{}
  	desc, err = GetDescriptorProto(messageName, nil, fileDesc.File[0])
  	if err != nil {
  		return err
  	}

  	PT, err := unmarshal.Unmarshal(fileDesc.File[0], desc, [][]byte{marshaled})
  	if err != nil {
  		return err
  	}

  	err = fs.Serve(c, PT)
  	if err != nil {
  		return err
  	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}

	return nil
}

func MountList(marshaled [][]byte, fileDesc *google_protobuf.FileDescriptorSet, messageName string, mountPoint string) error {
	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("ProtoFuse"),
		fuse.Subtype("protofuse"),
		fuse.LocalVolume(),
		fuse.VolumeName("ProtoFuse"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

  	desc := &google_protobuf.DescriptorProto{}
  	desc, err = GetDescriptorProto(messageName, nil, fileDesc.File[0])
  	if err != nil {
  		return err
  	}

  	PT, err := unmarshal.Unmarshal(fileDesc.File[0], desc, marshaled)
  	if err != nil {
  		return err
  	}

  	err = fs.Serve(c, PT)
  	if err != nil {
  		return err
  	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}

	return nil
}


func GetDescriptorProto(name string, messageDesc *google_protobuf.DescriptorProto, fileDesc *google_protobuf.FileDescriptorProto) (*google_protobuf.DescriptorProto, error) {
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

	return nil, errors.New("Cannot find message: " + name)
}