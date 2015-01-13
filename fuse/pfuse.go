package pfuse

import (
	"os"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type ProtoTree struct {
	Dir
}

func (t *ProtoTree) Root() (fs.Node, fuse.Error) {
	return &t.Dir, nil
}

// treeNode represents each node in the tree
type TreeNode struct {
	Name string
	FieldNumber uint64
	// TODO: add type as a field
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