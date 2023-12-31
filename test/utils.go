package testu

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	ft "github.com/ipfs/go-unixfs"
	h "github.com/ipfs/go-unixfs/importer/helpers"
	trickle "github.com/ipfs/go-unixfs/importer/trickle"

	cid "github.com/ipfs/go-cid"
	chunker "github.com/ipfs/go-ipfs-chunker"
	u "github.com/ipfs/go-ipfs-util"
	ipld "github.com/ipfs/go-ipld-format"
	mdag "github.com/ipfs/go-merkledag"
	mdagmock "github.com/ipfs/go-merkledag/test"
	mh "github.com/multiformats/go-multihash"
)

// SizeSplitterGen creates a generator.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.SizeSplitterGen
func SizeSplitterGen(size int64) chunker.SplitterGen {
	return func(r io.Reader) chunker.Splitter {
		return chunker.NewSizeSplitter(r, size)
	}
}

// GetDAGServ returns a mock DAGService.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.GetDAGServ
func GetDAGServ() ipld.DAGService {
	return mdagmock.Mock()
}

// NodeOpts is used by GetNode, GetEmptyNode and GetRandomNode
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.NodeOpts
type NodeOpts struct {
	Prefix cid.Prefix
	// ForceRawLeaves if true will force the use of raw leaves
	ForceRawLeaves bool
	// RawLeavesUsed is true if raw leaves or either implicitly or explicitly enabled
	RawLeavesUsed bool
}

// Some shorthands for NodeOpts.
var (
	// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.UseProtoBufLeaves
	UseProtoBufLeaves = NodeOpts{Prefix: mdag.V0CidPrefix()}
	// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.UseRawLeaves
	UseRawLeaves = NodeOpts{Prefix: mdag.V0CidPrefix(), ForceRawLeaves: true, RawLeavesUsed: true}
	// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.UseCidV1
	UseCidV1 = NodeOpts{Prefix: mdag.V1CidPrefix(), RawLeavesUsed: true}
	// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.UseBlake2b256
	UseBlake2b256 NodeOpts
)

func init() {
	UseBlake2b256 = UseCidV1
	UseBlake2b256.Prefix.MhType = mh.Names["blake2b-256"]
	UseBlake2b256.Prefix.MhLength = -1
}

// GetNode returns a unixfs file node with the specified data.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.GetNode
func GetNode(t testing.TB, dserv ipld.DAGService, data []byte, opts NodeOpts) ipld.Node {
	in := bytes.NewReader(data)

	dbp := h.DagBuilderParams{
		Dagserv:    dserv,
		Maxlinks:   h.DefaultLinksPerBlock,
		CidBuilder: opts.Prefix,
		RawLeaves:  opts.RawLeavesUsed,
	}

	db, err := dbp.New(SizeSplitterGen(500)(in))
	if err != nil {
		t.Fatal(err)
	}
	node, err := trickle.Layout(db)
	if err != nil {
		t.Fatal(err)
	}

	return node
}

// GetEmptyNode returns an empty unixfs file node.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.GetEmptyNode
func GetEmptyNode(t testing.TB, dserv ipld.DAGService, opts NodeOpts) ipld.Node {
	return GetNode(t, dserv, []byte{}, opts)
}

// GetRandomNode returns a random unixfs file node.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.GetRandomNode
func GetRandomNode(t testing.TB, dserv ipld.DAGService, size int64, opts NodeOpts) ([]byte, ipld.Node) {
	in := io.LimitReader(u.NewTimeSeededRand(), size)
	buf, err := io.ReadAll(in)
	if err != nil {
		t.Fatal(err)
	}

	node := GetNode(t, dserv, buf, opts)
	return buf, node
}

// ArrComp checks if two byte slices are the same.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.ArrComp
func ArrComp(a, b []byte) error {
	if len(a) != len(b) {
		return fmt.Errorf("arrays differ in length. %d != %d", len(a), len(b))
	}
	for i, v := range a {
		if v != b[i] {
			return fmt.Errorf("arrays differ at index: %d", i)
		}
	}
	return nil
}

// PrintDag pretty-prints the given dag to stdout.
//
// Deprecated: use github.com/ipfs/boxo/ipld/unixfs/test.PrintDag
func PrintDag(nd *mdag.ProtoNode, ds ipld.DAGService, indent int) {
	fsn, err := ft.FSNodeFromBytes(nd.Data())
	if err != nil {
		panic(err)
	}

	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("{size = %d, type = %s, children = %d", fsn.FileSize(), fsn.Type().String(), fsn.NumChildren())
	if len(nd.Links()) > 0 {
		fmt.Println()
	}
	for _, lnk := range nd.Links() {
		child, err := lnk.GetNode(context.Background(), ds)
		if err != nil {
			panic(err)
		}
		PrintDag(child.(*mdag.ProtoNode), ds, indent+1)
	}
	if len(nd.Links()) > 0 {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
	}
	fmt.Println("}")
}
