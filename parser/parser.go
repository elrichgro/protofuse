// Parser for marshalled protocol buffers. Builds a FUSE filsesystem structure.
package protoparser

import (
	"bytes"

	// "github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"elrich/protofuse/fuse"
)

func Parse(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, PT *pfuse.ProtoTree) {
	PT.Dir.Nodes = []pfuse.TreeNode{pfuse.TreeNode{}}
	parseMessage(desc, buf, &PT.Dir.Nodes[0])
}

func parseMessage(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode) {
	dir := &pfuse.Dir{}
	for buf.Len() != 0 {
		tN := &pfuse.TreeNode{}
		parseField(desc, buf, tN)
		dir.Nodes = append(dir.Nodes, *tN)
	} 
	t.Name = desc.GetName()
	t.Node = dir
}

func parseField(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, t *pfuse.TreeNode) {
	t.Name = "hello_world.txt"
	t.Node = &pfuse.File{"Hello, World!"}
	buf.Reset()
}
