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

package mount

import (
	"fmt"
	"log"
	"strings"

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
	desc, err = GetDescriptorProto(messageName, fileDesc)
	if err != nil {
		return err
	}

	PT, err := unmarshal.Unmarshal(fileDesc, desc, [][]byte{marshaled})
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
	desc, err = GetDescriptorProto(messageName, fileDesc)
	if err != nil {
		return err
	}

	PT, err := unmarshal.Unmarshal(fileDesc, desc, marshaled)
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

func GetDescriptorProto(name string, fileDesc *google_protobuf.FileDescriptorSet) (*google_protobuf.DescriptorProto, error) {
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