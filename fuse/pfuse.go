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

package pfuse

import (
	"os"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

type ProtoTree struct {
	Dir
}

func (t *ProtoTree) Root() (fs.Node, fuse.Error) {
	return &t.Dir, nil
}

// treeNode represents each node in the filesystem tree
type TreeNode struct {
	Name string
	FieldNumber uint64
	Type google_protobuf.FieldDescriptorProto_Type
	Node fs.Node
}

// Dir implements both Node and Handle for the directories.
type Dir struct {
	Nodes []TreeNode
}

func (dir *Dir) Attr() fuse.Attr {
	return fuse.Attr{Mode: os.ModeDir | 0555}
}

func (dir *Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	for _, treenode := range dir.Nodes {
		if name == treenode.Name {
			return treenode.Node, nil
		}
	}
	return nil, fuse.ENOENT
}

func (dir *Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	var dirs []fuse.Dirent
	for _, treenode := range dir.Nodes {
		dirs = append(dirs, fuse.Dirent{Name: treenode.Name})
	}
	return dirs, nil
}

// File implements both Node and Handle for the files.
type File struct{
	Contents string
}

func (file *File) Attr() fuse.Attr {
	return fuse.Attr{Mode: 0444, Size: uint64(len(file.Contents))}
}

func (file *File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	return []byte(file.Contents), nil
}