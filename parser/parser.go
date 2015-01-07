// Parser for marshalled protocol buffers. Builds a FUSE filsesystem structure.
package parser

import (
	"github.com/gogo/protobuf/proto"
	"elrich/protofuse/fuse"
)

func Parse(desc *google_protobuf.DescriptorProto, buf *bytes.Buffer, PT *pfuse.ProtoTree) {
	parseMessage(desc, buf, PT.Dir)
}

