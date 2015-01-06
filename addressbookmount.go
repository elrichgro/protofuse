// Mount marshalled protocol buffers as a filesystem
// command line arguments:
//		marshaled protobuf
// 		mount location
package main

import (
	"fmt"
	"os"
	// "flag"
	"log"
	"io"
	// "reflect"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/gogo/protobuf/proto"
	"elrich/protobuf/tutorial"
)

func main() {
	
	if len(os.Args) != 3 {
       fmt.Printf("Usage: %s ADDRESS_BOOK_FILE, MOUNT_LOCATION\n", os.Args[0])
       os.Exit(-1)
   	}

	mountpoint := os.Args[2]

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("ProtoFuse"),
		fuse.Subtype("protofuse"),
		fuse.LocalVolume(),
		fuse.VolumeName("ProtoFuse"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Read file [copied from tutorial]
   	file, err := os.Open(os.Args[1])
   	CheckError(err)

   	fi, err := file.Stat()
   	CheckError(err)

   	buffer := make([]byte, fi.Size())
   	_, err = io.ReadFull(file, buffer)
   	file.Close()

   	address_book := &tutorial.AddressBook{}
   	err = proto.Unmarshal(buffer, address_book)
   	CheckError(err)
   	// [end copy]

	PT, err := buildTree(address_book)
	// DoNothing(PT)
	if err != nil {
		log.Fatal(err)
	}

	err = fs.Serve(c, &PT)
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

// func DoNothing(pt ProtoTree) {

// }

// copied from tutorial
func CheckError(err error) {
   if err != nil {
       fmt.Println(err.Error())
       os.Exit(-1)
   }
}

// TODO: build tree representing protobuf
func buildTree(address_book *tutorial.AddressBook) (ProtoTree, fuse.Error) {
	// return ProtoTree{Dir{nodes: []treeNode{{name: "dir1", node: &Dir{}}, {name:"file1", node:&File{}}}}}, nil
	prototree := ProtoTree{Dir{}}
	dir, err := parseAddressBook(address_book)
	CheckError(err)
	prototree.Dir.nodes = append(prototree.Dir.nodes, treeNode{name: "AddressBook", node: dir})
	
	return prototree, nil
} 

func parseAddressBook(address_book *tutorial.AddressBook) (*Dir, fuse.Error) {
	dir := Dir{}
	i := 0
	for _, person := range address_book.Person {
		i += 1
		personDir, err := parsePerson(person)
		CheckError(err)
		dir.nodes = append(dir.nodes, treeNode{name: fmt.Sprintf("Person_%d", i), node: personDir})
	}
	return &dir, nil
}

func parsePerson(person *tutorial.Person) (*Dir, fuse.Error) {
	personDir := Dir{}
	personDir.nodes = append(personDir.nodes, treeNode{name: "Name", node: &File{person.GetName()}})
	personDir.nodes = append(personDir.nodes, treeNode{name: "Id", node: &File{fmt.Sprintf("%d", person.GetId())}})
	personDir.nodes = append(personDir.nodes, treeNode{name: "Email", node: &File{person.GetEmail()}})
	i := 0
	for _, phone := range person.Phone {
		i += 1
		phoneDir, err := parsePhone(phone)
		CheckError(err)
		personDir.nodes = append(personDir.nodes, treeNode{name: fmt.Sprintf("Phone_%d", i), node: phoneDir})
	}

	return &personDir, nil
}

func parsePhone(phone *tutorial.Person_PhoneNumber) (*Dir, fuse.Error) {
	phoneDir := Dir{}
	phoneDir.nodes = append(phoneDir.nodes, treeNode{name: "Number", node: &File{phone.GetNumber()}})
	phoneDir.nodes = append(phoneDir.nodes, treeNode{name: "Type", node: &File{tutorial.Person_PhoneType_name[int32(phone.GetType())]}})
	return &phoneDir, nil
}

// func buildTree(address_book *tutorial.AddressBook) (ProtoTree, fuse.Error) {
// 	prototree := ProtoTree{}
// 	// prototree.Dir.nodes = append(prototree.Dir.nodes, parseMessage(address_book))
// 	parseMessage(address_book)
// 	return prototree, nil
// }

// func parseMessage(message interface{}) *Dir {
// 	// fmt.Println("Type:", reflect.TypeOf(message))
// 	dir := Dir{}
// 	v := reflect.ValueOf(message).Elem()
// 	for i := 0; i < v.NumField()-1; i++ {
// 		//if message type, parse(message)
// 		//if normal type, append to dir
// 		//if array type, parse one by one and add _number
//     	f := v.Field(i)
//     	fmt.Println(f.Type())
// 	}
// 	return &dir
// }

type ProtoTree struct {
	Dir
}

func (t *ProtoTree) Root() (fs.Node, fuse.Error) {
	return &t.Dir, nil
}

// treeNode represents each node in the tree
type treeNode struct {
	name string
	node fs.Node
}

// Dir implements both Node and Handle for the directories.
type Dir struct {
	nodes []treeNode
}

func (dir *Dir) Attr() fuse.Attr {
	return fuse.Attr{Mode: os.ModeDir | 0555}
}

func (dir *Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	for _, treenode := range dir.nodes {
		if name == treenode.name {
			return treenode.node, nil
		}
	}
	return nil, fuse.ENOENT
}

func (dir *Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	var dirs []fuse.Dirent
	for _, treenode := range dir.nodes {
			dirs = append(dirs, fuse.Dirent{Name: treenode.name})
		}
	return dirs, nil
}

// File implements both Node and Handle for the files.
type File struct{
	contents string
}

func (file *File) Attr() fuse.Attr {
	return fuse.Attr{Mode: 0444, Size: uint64(len(file.contents))}
}

func (file *File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	return []byte(file.contents), nil
}