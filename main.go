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

package main

import (
	"fmt"
	"os"
	"io"

	"github.com/gogo/protobuf/parser"
	"elrich/protofuse/mount"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

var fileDesc *google_protobuf.FileDescriptorProto

func main() {
	
	if len(os.Args) != 5 {
       fmt.Printf("Usage: %s MOUNT_LOCATION, MARSHALLED_PROTOCOL_BUFFER, PROTO_FILE_LOCATION, MESSAGE_NAME\n", os.Args[0])
       os.Exit(-1)
   	}

	mountpoint := os.Args[1]

   	file, err := os.Open(os.Args[2])
   	CheckError(err)

   	fi, err := file.Stat()
   	CheckError(err)

   	buf := make([]byte, fi.Size())
   	_, err = io.ReadFull(file, buf)
   	file.Close()

   	filename := string(os.Args[3])

   	fileDescSet, err := parser.ParseFile(filename, ".")
  	CheckError(err)
  	var messageName string = os.Args[4]

	err = protofuse.Mount(buf, fileDescSet, messageName, mountpoint)
	CheckError(err)
}

func CheckError(err error) {
   if err != nil {
       fmt.Println(err.Error())
       os.Exit(-1)
   }
}
