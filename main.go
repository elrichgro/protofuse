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

	"github.com/gogo/protobuf/proto"
	"elrich/protobuf/tutorial"
)

func main() {
	
	if len(os.Args) != 5 {
       fmt.Printf("Usage: %s MOUNT_LOCATION, MARSHALLED_PROTOCOL_BUFFER, PROTO_FILE, MESSAGE_NAME\n", os.Args[0])
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
   	file, err = os.Open(os.Args[3])
   	CheckError(err)

   	fi, err = file.Stat()
   	CheckError(err)

   	buffer := make([]byte, fi.Size())
   	_, err = io.ReadFull(file, buffer)
   	file.Close()

   	file_descriptor_set := &google_protobuf.FileDescriptorSet{}
  	err = proto.Unmarshal(buffer, file_descriptor_set)
  	CheckError(err)

  	// TODO: read message name and use it to get DescriptorProto



	// PT, err := buildTree(address_book)
	// // DoNothing(PT)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = fs.Serve(c, &PT)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // check if the mount process has an error to report
	// <-c.Ready
	// if err := c.MountError; err != nil {
	// 	log.Fatal(err)
	// }
}