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

package unmarshal

import (
	"testing"
	"reflect"
	"fmt"

	"github.com/elrichgro/protofuse/fuse"
	"github.com/elrichgro/protofuse/test"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

func TestUnmarshal(t *testing.T) {

	buf, fDesc, packageName, messageName, err := test.GenerateFull()
	if err != nil {
		t.Fatal(err)
	}

	PT1, err := Unmarshal(fDesc, packageName, messageName, [][]byte{buf})
	if err != nil {
		t.Fatal(err)
	}

	PT2 := &pfuse.ProtoTree{pfuse.Dir{[]pfuse.TreeNode{pfuse.TreeNode{Name:"Message_1", FieldNumber:0, Type: google_protobuf.FieldDescriptorProto_TYPE_MESSAGE, 
	Node: &pfuse.Dir{[]pfuse.TreeNode{pfuse.TreeNode{Name:"f121", FieldNumber:121, Type: google_protobuf.FieldDescriptorProto_TYPE_INT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"123"}}, pfuse.TreeNode{Name:"f1", FieldNumber:1, Type: google_protobuf.FieldDescriptorProto_TYPE_STRING, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REQUIRED, Node:&pfuse.File{Contents:"one"}}, pfuse.TreeNode{Name:"f2_1", FieldNumber:2, Type: google_protobuf.FieldDescriptorProto_TYPE_INT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"1"}}, pfuse.TreeNode{Name:"f2_2", FieldNumber:2, Type: google_protobuf.FieldDescriptorProto_TYPE_INT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"2"}}, pfuse.TreeNode{Name:"f2_3", FieldNumber:2, Type: google_protobuf.FieldDescriptorProto_TYPE_INT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"3"}}, pfuse.TreeNode{Name:"f2_4", FieldNumber:2, Type: google_protobuf.FieldDescriptorProto_TYPE_INT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"4"}}, pfuse.TreeNode{Name:"f3", FieldNumber:3, Type: google_protobuf.FieldDescriptorProto_TYPE_INT64, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"3"}}, pfuse.TreeNode{Name:"f4", FieldNumber:4, Type: google_protobuf.FieldDescriptorProto_TYPE_UINT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"4"}}, pfuse.TreeNode{Name:"f5", FieldNumber:5, Type: google_protobuf.FieldDescriptorProto_TYPE_UINT64, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"5"}}, pfuse.TreeNode{Name:"f6", FieldNumber:6, Type: google_protobuf.FieldDescriptorProto_TYPE_SINT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"6"}}, pfuse.TreeNode{Name:"f7", FieldNumber:7, Type: google_protobuf.FieldDescriptorProto_TYPE_BOOL, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"True"}}, pfuse.TreeNode{Name:"f8", FieldNumber:8, Type: google_protobuf.FieldDescriptorProto_TYPE_FIXED64, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"8"}}, pfuse.TreeNode{Name:"f9", FieldNumber:9, Type: google_protobuf.FieldDescriptorProto_TYPE_SFIXED64, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"9"}}, pfuse.TreeNode{Name:"f10", FieldNumber:10, Type: google_protobuf.FieldDescriptorProto_TYPE_DOUBLE, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"10.000000"}}, pfuse.TreeNode{Name:"f11", FieldNumber:11, Type: google_protobuf.FieldDescriptorProto_TYPE_BYTES, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"0b0b"}}, pfuse.TreeNode{Name:"f12", FieldNumber:12, Type: google_protobuf.FieldDescriptorProto_TYPE_MESSAGE, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.Dir{[]pfuse.TreeNode{pfuse.TreeNode{Name:"name", FieldNumber:100, Type:google_protobuf.FieldDescriptorProto_TYPE_STRING, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REQUIRED, Node:&pfuse.File{Contents:"BAR"}}, pfuse.TreeNode{Name:"id", FieldNumber:1, Type:google_protobuf.FieldDescriptorProto_TYPE_INT32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REQUIRED, Node:&pfuse.File{Contents:"123"}}}}}, pfuse.TreeNode{Name:"f13", FieldNumber:13, Type: google_protobuf.FieldDescriptorProto_TYPE_FIXED32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"13"}}, pfuse.TreeNode{Name:"f14", FieldNumber:14, Type: google_protobuf.FieldDescriptorProto_TYPE_SFIXED32, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"14"}}, pfuse.TreeNode{Name:"f15", FieldNumber:15, Type: google_protobuf.FieldDescriptorProto_TYPE_FLOAT, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"15.000000"}}, pfuse.TreeNode{Name:"f16_1", FieldNumber:16, Type: google_protobuf.FieldDescriptorProto_TYPE_MESSAGE, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REPEATED, Node:&pfuse.Dir{[]pfuse.TreeNode{pfuse.TreeNode{Name:"f1", FieldNumber:1, Type: google_protobuf.FieldDescriptorProto_TYPE_STRING, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REQUIRED, Node:&pfuse.File{Contents:"name"}}, pfuse.TreeNode{Name:"f2", FieldNumber:2, Type: google_protobuf.FieldDescriptorProto_TYPE_ENUM, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"e1"}}, pfuse.TreeNode{Name:"f3", FieldNumber:3, Type: google_protobuf.FieldDescriptorProto_TYPE_MESSAGE, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.Dir{[]pfuse.TreeNode{pfuse.TreeNode{Name:"name", FieldNumber:1, Type: google_protobuf.FieldDescriptorProto_TYPE_STRING, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REQUIRED, Node:&pfuse.File{Contents:"name"}}}}}}}}, pfuse.TreeNode{Name:"f16_2", FieldNumber:16, Type: google_protobuf.FieldDescriptorProto_TYPE_MESSAGE, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REPEATED, Node:&pfuse.Dir{[]pfuse.TreeNode{pfuse.TreeNode{Name:"f1", FieldNumber:1, Type: google_protobuf.FieldDescriptorProto_TYPE_STRING, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REQUIRED, Node:&pfuse.File{Contents:"name2"}}, pfuse.TreeNode{Name:"f2", FieldNumber:2, Type: google_protobuf.FieldDescriptorProto_TYPE_ENUM, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.File{Contents:"e2"}}, pfuse.TreeNode{Name:"f3", FieldNumber:3, Type: google_protobuf.FieldDescriptorProto_TYPE_MESSAGE, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL, Node:&pfuse.Dir{[]pfuse.TreeNode{pfuse.TreeNode{Name:"name", FieldNumber:1, Type: google_protobuf.FieldDescriptorProto_TYPE_STRING, 
	Label: google_protobuf.FieldDescriptorProto_LABEL_REQUIRED, Node:&pfuse.File{Contents:"name2"}}}}}}}}}}}}}}

	compareProtoTree(PT1, PT2, t)
}

func compareProtoTree(pt1 *pfuse.ProtoTree, pt2 *pfuse.ProtoTree, t *testing.T) {
	if len(pt1.Dir.Nodes) != len(pt2.Dir.Nodes) {
		t.Fatal("ProtoTree.Dir lengths don't match")
	}
	for i := range(pt1.Dir.Nodes) {
		compareTreeNode(pt1.Dir.Nodes[i], pt2.Dir.Nodes[i], t)
	}
}

func compareTreeNode(tn1 pfuse.TreeNode, tn2 pfuse.TreeNode, t *testing.T) {
	fmt.Printf("Comparing TreeNode %s with %s\n", tn1.Name, tn2.Name)
	if tn1.Name != tn2.Name {
		fmt.Println("Error")
		t.Error(fmt.Sprintf("TreeNode names don't match: %s != %s", tn1.Name, tn2.Name))
	}
	if tn1.FieldNumber != tn2.FieldNumber {
		fmt.Println("Error")
		t.Error(fmt.Sprintf("TreeNode field numbers don't match: %d != %d, for %s", tn1.FieldNumber, tn2.FieldNumber, tn1.Name))
	}
	if tn1.Type != tn2.Type {
		fmt.Println("Error")
		t.Error(fmt.Sprintf("TreeNode types don't match: %s != %s, for %s", google_protobuf.FieldDescriptorProto_Type_name[int32(tn1.Type)], google_protobuf.FieldDescriptorProto_Type_name[int32(tn2.Type)], tn1.Name))
	}
	if reflect.TypeOf(tn1.Node) != reflect.TypeOf(tn2.Node) {
		t.Fatal(fmt.Sprintf("TreeNode Node types don't match: %s != %s, for %s", reflect.TypeOf(tn1.Node).String(), reflect.TypeOf(tn2.Node).String(), tn1.Name))
	}
	if reflect.TypeOf(tn1.Node) == reflect.TypeOf(&pfuse.Dir{}) {
		compareDir(tn1.Node.(*pfuse.Dir), tn2.Node.(*pfuse.Dir), t)
	} else {
		compareFile(tn1.Node.(*pfuse.File), tn2.Node.(*pfuse.File), t)
	}
}

func compareDir(dir1 *pfuse.Dir, dir2 *pfuse.Dir, t *testing.T) {
	if len(dir1.Nodes) != len(dir2.Nodes) {
		t.Fatal(fmt.Sprintf("Dir.Nodes lengths don't match: %d != %d", len(dir1.Nodes), len(dir2.Nodes)))
	}
	for i := range(dir1.Nodes) {
		compareTreeNode(dir1.Nodes[i], dir2.Nodes[i], t)
	}
}

func compareFile(f1 *pfuse.File, f2 *pfuse.File, t *testing.T) {
	if f1.Contents != f2.Contents {
		fmt.Println("Error")
		t.Error(fmt.Sprintf("File contents don't match: %s != %s", f1.Contents, f2.Contents))
	}
}
