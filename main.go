// Mount marshalled protocol buffers as a FUSE filesystem.
// command line arguments:
//		mount location
//		marshalled protocol buffer
//		descriptor .proto file
// 		message name

package main

import (
	"fmt"
	"os"
	"log"
	"io"
	"bytes"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/gogo/protobuf/parser"
	"elrich/protofuse/fuse"
	"elrich/protofuse/parser"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

var fileDesc google_protobuf.FileDescriptorProto

func main() {
	
	if len(os.Args) != 5 {
       fmt.Printf("Usage: %s MOUNT_LOCATION, MARSHALLED_PROTOCOL_BUFFER, PROTO_FILE_LOCATION, MESSAGE_NAME\n", os.Args[0])
       os.Exit(-1)
   	}

	mountpoint := os.Args[1]

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("ProtoFuse"),
		fuse.Subtype("protofuse"),
		fuse.LocalVolume(),
		fuse.VolumeName("ProtoFuse"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Read protocol buffer and create buffer [copied from tutorial]
   	file, err := os.Open(os.Args[2])
   	CheckError(err)

   	fi, err := file.Stat()
   	CheckError(err)

   	buf := make([]byte, fi.Size())
   	_, err = io.ReadFull(file, buf)
   	file.Close()
   	// [end copy]

   	// Read .proto file and generate FileDescriptorSet
   	filename := string(os.Args[3])

   	file_descriptor_set, err := parser.ParseFile(filename, ".")
  	CheckError(err)

  	fileDesc = *file_descriptor_set.File[0]

  	// Get file descriptor proto
  	var messageName string = os.Args[4]
  	desc := &google_protobuf.DescriptorProto{}
  	GetDescriptorProto(desc, messageName, nil) // TODO: err = GetDescriptorProto(desc, messageName)
  	CheckError(err)

  	PT := &pfuse.ProtoTree{}

  	// Parse the FileDescriptorProto
  	// TODO: make parser return err 
  	protoparser.Parse(desc, bytes.NewBuffer(buf), PT)
  	CheckError(err)

  	// Start FUSE serve loop
  	err = fs.Serve(c, PT)
  	CheckError(err)

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

// copied from tutorial
func CheckError(err error) {
   if err != nil {
       fmt.Println(err.Error())
       os.Exit(-1)
   }
}

func GetDescriptorProto(desc *google_protobuf.DescriptorProto, name string, messageDesc *google_protobuf.DescriptorProto) {
	if messageDesc != nil {
		for _, message := range messageDesc.NestedType {
			if *message.Name == name {
				*desc = *message
				return
			}
		}
	}

	for _, message := range fileDesc.MessageType {
		if *message.Name == name {
			*desc = *message
			return
		}
	}

	//TODO: else throw error message could not be found
	// *desc = nil
	fmt.Println("Can't find message\n")
}