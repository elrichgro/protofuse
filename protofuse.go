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

//  Mount marshalled protocol buffers as a FUSE filesystem.
//  command line arguments:
//		mount location
//		marshalled protocol buffer
//		descriptor .proto file
// 		message name

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/elrichgro/protofuse/mount"
	"github.com/gogo/protobuf/parser"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

var fileDesc *google_protobuf.FileDescriptorProto

func main() {

	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s MOUNT_LOCATION, MARSHALLED_PROTOCOL_BUFFER, PROTO_FILE_LOCATION, PACKAGE_NAME, MESSAGE_NAME\n", os.Args[0])
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

	fileDescSet, err := parser.ParseFile(filename, filename[:strings.LastIndex(filename, "/")])
	CheckError(err)
	var packageName string = os.Args[4]
	var messageName string = os.Args[5]

	err = mount.Mount(buf, fileDescSet, packageName, messageName, mountpoint)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
