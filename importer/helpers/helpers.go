package helpers

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
	pb "github.com/ipfs/go-unixfs/pb"
)

// BlockSizeLimit specifies the maximum size an imported block can have.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/importer/helpers.BlockSizeLimit
var BlockSizeLimit = 1048576 // 1 MB

// rough estimates on expected sizes
var roughLinkBlockSize = 1 << 13 // 8KB
var roughLinkSize = 34 + 8 + 5   // sha256 multihash + size + no name + protobuf framing

// DefaultLinksPerBlock governs how the importer decides how many links there
// will be per block. This calculation is based on expected distributions of:
//   - the expected distribution of block sizes
//   - the expected distribution of link sizes
//   - desired access speed
//
// For now, we use:
//
//	var roughLinkBlockSize = 1 << 13 // 8KB
//	var roughLinkSize = 34 + 8 + 5   // sha256 multihash + size + no name
//	                                 // + protobuf framing
//	var DefaultLinksPerBlock = (roughLinkBlockSize / roughLinkSize)
//	                         = ( 8192 / 47 )
//	                         = (approximately) 174
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/importer/helpers.DefaultLinksPerBlock
var DefaultLinksPerBlock = roughLinkBlockSize / roughLinkSize

// ErrSizeLimitExceeded signals that a block is larger than BlockSizeLimit.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/importer/helpers.ErrSizeLimitExceeded
var ErrSizeLimitExceeded = fmt.Errorf("object size limit exceeded")

type Helper interface {
	Done() bool
	Next() ([]byte, error)
	GetDagServ() ipld.DAGService
	GetCidBuilder() cid.Builder
	NewLeafNode(data []byte, fsNodeType pb.Data_DataType) (ipld.Node, error)
	FillNodeLayer(node *FSNodeOverDag) error
	NewLeafDataNode(fsNodeType pb.Data_DataType) (node ipld.Node, dataSize uint64, err error)
	ProcessFileStore(node ipld.Node, dataSize uint64) ipld.Node
	Add(node ipld.Node) error
	Maxlinks() int
	NewFSNodeOverDag(fsNodeType pb.Data_DataType) *FSNodeOverDag
	NewFSNFromDag(nd *dag.ProtoNode) (*FSNodeOverDag, error)
}
