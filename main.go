package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

type Files struct {
	name      string
	extension string
	size      int64
}

func (ob *Files) print() {
	fmt.Println("FileName:", ob.name, "FileExt:", ob.extension, "FileSize/byte", ob.size)
}
func getFilePathFromCommand(temp string, sort string) (string, string, error) {
	var sourcepath *string
	var sortflag *string
	sourcepath = flag.String(temp, "None", "")
	sortflag = flag.String(sort, "None", "")
	flag.Parse()
	return *sourcepath, *sortflag, nil
}
func rootExist(root string) (bool, error) {
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		fmt.Println("Root does not exist...!")
	}
	return true, nil
}
func getRoot(root string) (string, error) {
	var rootflag *string
	rootflag = flag.String(root, "None", "")
	flag.Parse()
	_, err := rootExist(*rootflag)
	if err != nil {
		panic(err)
	}
	return *rootflag, nil
}

func (ob *Files) getSize() int64 {
	return ob.size
}
func (ob *Files) getName() string {
	return ob.name
}
func (ob *Files) getExtension() string {
	return ob.extension
}
func getFilesRecurvise(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return true, nil
}
func getFileLocation(root string, filename string) (string, error) {
	if root == "" {
		return "", errors.New("Root is Empty")
	}
	return root + "/" + filename, nil
}
func getFileExtension(root string, filename string) (string, error) {
	f, err := getFileLocation(root, filename)
	if err != nil {
		fmt.Println(err)
	}
	st, err := os.Stat(f)
	if err != nil {
		fmt.Println(err)
	}
	if st.IsDir() {
		return "Directory", nil
	}
	return "file", nil
}
func getFilesFromDirectory(pathName string) ([]Files, error) {
	fi, err := os.Open(pathName)
	if err != nil {
		log.Fatal(err, fi.Name())
	}
	defer fi.Close()
	files, err := os.ReadDir(pathName)
	if err != nil {
		fmt.Print("Can't Read from the Directory!", err)
	}
	var s []Files
	for _, item := range files {
		p, err := getFileLocation(pathName, item.Name())
		f, err := os.Stat(p)
		if err != nil {
			panic(err)
		}
		Ext, err := getFileExtension(pathName, item.Name())
		element := Files{name: f.Name(), extension: Ext, size: f.Size()}
		s = append(s, element)
	}
	return s, nil
}
func sortAsc(arr []Files) error {
	if len(arr) < 0 {
		fmt.Println("Array is empty!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].size > arr[j].size
	})
	return nil
}
func sortDesc(arr []Files) error {
	if len(arr) < 0 {
		fmt.Println("Array is empty!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].size < arr[j].size
	})
	return nil
}

func main() {
	//root := "/Users/ismaelnvo/Desktop/sort/pathfiles"
	rootflag := "root"
	sortflag := "sort"
	root, sort, err := getFilePathFromCommand(rootflag, sortflag)
	if err != nil {
		fmt.Println(err)
	}
	if sort == "None" {
		list, err := getFilesFromDirectory(root)
		if err != nil {
			panic(err)
		}
		sortAsc(list)
		for i := 0; i < len(list); i++ {
			list[i].print()
		}
	} else {
		list, err := getFilesFromDirectory(root)
		if err != nil {
			panic(err)
		}
		sortDesc(list)
		for i := 0; i < len(list); i++ {
			list[i].print()
		}
	}

}
