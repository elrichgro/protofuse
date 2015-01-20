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

// import (
// 	"testing"
// 	// "fmt"

// 	// "elrich/protofuse/fuse"
// )

// func compareProtoTree(pt1 *pfuse.ProtoTree, pt2 *pfuse.ProtoTree, t *testing.T) {
// 	if len(*pt1.Dir) != len(*pt2.Dir) {
// 		t.Error("ProtoTree.Dir lengths don't match")
// 	}
// 	for i := range(pt1.Dir) {
// 		compareTreeNode(pt1.Dir[i], pt2.Dir[i], t)
// 	}
// }

// func compareTreeNode(tn1 *pfuse.TreeNode, tn2 *pfuse.TreeNode, t *testing.T) {
// 	if tn1.Name != tn2.Name {
// 		t.Error(fmt.Sprintf("TreeNode names don't match: %s != %s\n", tn1.Name, tn2.Name))
// 	}
// 	if tn1.FieldNumber != tn2.FieldNumber {
// 		t.Error(fmt.Sprintf("TreeNode field numbers don't match: %d != %d\n", tn1.FieldNumber, tn2.FieldNumber))
// 	}
// 	if tn1.Type != tn2.Type {
// 		t.Error(fmt.Sprintf("TreeNode types don't match: %s != %s\n", google_protobuf.FieldDescriptorProto_Type_name[int32(tn1.Type)], google_protobuf.FieldDescriptorProto_Type_name[int32(tn2.Type)]))
// 	}
// }

// func compareDir(dir1 *pfuse.Dir, dir2 *pfuse.Dir, t *testing.T) {

// }

// func compareFile(f1 *pfuse.File, f2 *pfuse.File, t *testing.T) {

// }
