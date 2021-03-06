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

// 	The mount package provides methods to mount marshaled protocol buffers as a filesystem.
package mount

import (
	"fmt"
	"log"
	"strings"
	"os/signal"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/elrichgro/protofuse/unmarshal"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

//	Mounts a marshaled protocol buffer as a filesytem. 
func Mount(marshaled []byte, fileDesc *google_protobuf.FileDescriptorSet, packageName string, messageName string, mountPoint string) error {
	// mount
	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("protofuse"),
		fuse.Subtype("protofs"),
		fuse.LocalVolume(),
		fuse.VolumeName("ProtoFS"),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	// unmount filesystem in the event of an interrupt
	c1 := make(chan os.Signal, 1)
	signal.Notify(c1, os.Interrupt)
	go func(){
    for sig := range c1 {
    	log.Printf("captured %v, unmounting filesystem", sig)
        err = fuse.Unmount(mountPoint)
		if err != nil {
			log.Println(err)
		}
    }
	}()

	// create the filesystem structure
	PT, err := unmarshal.Unmarshal(fileDesc, packageName, messageName, [][]byte{marshaled})
	if err != nil {
		return err
	}

	// serve
	err = fs.Serve(c, PT)
	if err != nil {
		return err
	}
	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil
}

// Mounts a list of marshaled protocol buffers as a filesystem.
func MountList(marshaled [][]byte, fileDesc *google_protobuf.FileDescriptorSet, packageName string, messageName string, mountPoint string) error {
	// mount
	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("protofuse"),
		fuse.Subtype("protofs"),
		fuse.LocalVolume(),
		fuse.VolumeName("ProtoFS"),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	// unmount the filesystem in the event of an interrupt
	c1 := make(chan os.Signal, 1)
	signal.Notify(c1, os.Interrupt)
	go func(){
    for sig := range c1 {
    	log.Printf("captured %v, unmounting filesystem", sig)
        err = fuse.Unmount(mountPoint)
		if err != nil {
			log.Println(err)
		}
    }
	}()

	// create the filesystem structure
	PT, err := unmarshal.Unmarshal(fileDesc, packageName, messageName, marshaled)
	if err != nil {
		return err
	}

	// serve
	err = fs.Serve(c, PT)
	if err != nil {
		return err
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil
}

// Tries to unmount the filesystem at dir.
func Unmount(dir string) error {
	err := fuse.Unmount(dir)
	if err != nil {
		return err
	}
	return nil
}

// Finds the google_protobuf.DescriptorProto for the message name.
func getDescriptorProto(name string, fileDesc *google_protobuf.FileDescriptorSet) (*google_protobuf.DescriptorProto, error) {
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
