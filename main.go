package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Tree struct calls the os.FileInfo library method to list file details
type Tree struct {
	File     os.FileInfo
	childNodes []Tree
}

// Name used the Name() method from the ioutil package that finds the pathname of the file. It also calls the below Size() function to include the size of the file
func (t Tree) Name() string {
	if t.File.IsDir() {
		return t.File.Name()
	}
	return fmt.Sprintf("%s (%s)", t.File.Name(), t.Size())

}

// Size uses the Size() method from the ioutil package to get file size
func (t Tree) Size() string {
	if t.File.Size() > 0 {
		return fmt.Sprintf("%db", t.File.Size())
	}
	return "empty"
}

// readFiles compiles file details into a list based on the Tree struct
func readFiles(path string, hasFiles bool) ([]Tree, error) {
	// ReadDir returns a list of directory entries by filename
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var nodes []Tree
	for _, file := range files {
		if !hasFiles && !file.IsDir() {
			continue
		}

		node := Tree {
			File: file,
		}

		if file.IsDir() {
			children, err := readFiles(path+string(os.PathSeparator)+file.Name(), hasFiles)
			if err != nil {
				return nil, err
			}
			node.childNodes = children
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// printTree prints out the file details, adjusting the ASCI prefixes to denote 
// root folder and subtrees
func printTree(out io.Writer, nodes []Tree, prefix string) {

	lastFile := len(nodes) -1
	var parentPrefix = "├───"
	var childPrefix = "│\t"

	// Alter the prefix vars within the loop if necessary
	for i, node := range nodes {
		if i == lastFile {
			parentPrefix = "└───"
			childPrefix = "\t"
		}
		
		fmt.Fprint(out, prefix, parentPrefix, node.Name(), "\n")

		if node.File.IsDir() {
			printTree(out, node.childNodes, prefix+childPrefix)
			}
		}
	}


func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	nodes, err := readFiles(path, printFiles)
		if err != nil {
			return 
		}

		printTree(out, nodes, "")
		return 
	}


func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
