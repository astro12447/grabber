package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type Files struct {
	name      string
	extension string
	size      int64
}

var s []Files

func (ob *Files) print() {
	fmt.Println("Name:", ob.name, "Type:", ob.extension, "FileSize/byte", ob.size)
}
func getFilePathFromCommand(root string, sort string) (string, string, error) {
	var sourcepath *string
	var sortflag *string
	sourcepath = flag.String(root, "None", "")
	sortflag = flag.String(sort, "None", "")
	flag.Parse()
	return *sourcepath, *sortflag, nil
}
func rootExist(root string) (bool, error) {
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		fmt.Println("Root не существует...!")
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
		return "", errors.New("Root  пуст!")
	}
	return root + "/" + filename, nil
}
func getAllFromDir(path string) ([]Files, error) {
	err := filepath.Walk(path, func(p string, inf os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		size, err := getsize(p)
		if err != nil {
			fmt.Println(err)
		}
		Ext, err := getFileExtension2(p)
		if err != nil {
			fmt.Println(err)
		}
		element := Files{name: p, extension: Ext, size: size}
		s = append(s, element)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return s, nil
}
func getsize(filename string) (int64, error) {
	f, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}
	return f.Size(), nil
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
		return "Каталог", nil
	}
	return "файл", nil
}
func getFileExtension2(filename string) (string, error) {
	f, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}
	if f.IsDir() {
		return "Каталог", nil
	}
	return "файл", nil
}
func getFilesFromDirectory(pathName string) ([]Files, error) {
	fi, err := os.Open(pathName)
	if err != nil {
		log.Fatal(err, fi.Name())
	}
	defer fi.Close()
	files, err := os.ReadDir(pathName)
	if err != nil {
		fmt.Print("Невозможно прочитать из каталога!", err)
	}
	for _, item := range files {
		p, err := getFileLocation(pathName, item.Name())
		f, err := os.Stat(p)
		if err != nil {
			panic(err)
		}
		Ext, err := getFileExtension(pathName, item.Name())
		name := pathName + "/" + f.Name()
		element := Files{name: name, extension: Ext, size: f.Size()}
		s = append(s, element)
	}
	return s, nil
}
func sortAsc(arr []Files) {
	if len(arr) < 0 {
		fmt.Println("Массив пуст!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].size < arr[j].size
	})
}
func sortDesc(arr []Files) {
	if len(arr) < 0 {
		fmt.Println("Массив пуст!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].size > arr[j].size
	})
}

func main() {
	//root := "/Users/ismaelnvo/Desktop/sort/pathfiles"
	//getAllFromDir("/Users/ismaelnvo/Desktop/sort/temp")
	rootflag := "root"
	sortflag := "sort"
	root, sort, err := getFilePathFromCommand(rootflag, sortflag)
	if err != nil {
		fmt.Println(err)
	}
	if root != "None" && sort == "None" {
		list, err := getFilesFromDirectory(root)
		if err != nil {
			panic(err)
		}
		sortAsc(list)
		for i := 0; i < len(list); i++ {
			list[i].print()
		}
	} else if sort == "Desc" && root != "None" {
		list, err := getFilesFromDirectory(root)
		if err != nil {
			panic(err)
		}
		sortDesc(list)
		for i := 0; i < len(list); i++ {
			list[i].print()
		}
	} else {
		fmt.Println("->Введите правильную командную строку:(--root=/pathfile  --sort=Desc) or --root=/pathfile")
	}

}
