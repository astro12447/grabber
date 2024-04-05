package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// определение структуры файла
type Files struct {
	name      string
	extension string
	size      int64
}

var s []Files

// определение функции для ввода информации классы Files в консоль
func (ob *Files) print() {
	fmt.Println("Name:", ob.name, "Type:", ob.extension, "FileSize/byte", ob.size)
}

// определение функции для получения строк через консоль
func getFilePathFromCommand(root string, sort string) (string, string, error) {
	if root == "None" || sort == "None" {
		fmt.Println("->Введите правильную командную строку:(--root=/pathfile  --sort=Desc) or --root=/pathfile")
	}
	var sourcepath *string
	var sortflag *string
	sourcepath = flag.String(root, "None", "")
	sortflag = flag.String(sort, "None", "")
	flag.Parse()
	return *sourcepath, *sortflag, nil
}

// функция для проверкаи попки
func rootExist(root string) (bool, error) {
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		fmt.Println("Root не существует...!")
	}
	return true, nil
}

// метод для получения значения size класса
func (ob *Files) getSize() int64 {
	return ob.size
}

// метод для получения значения name класса
func (ob *Files) getName() string {
	return ob.name
}

// метод для получения значения Extension класса
func (ob *Files) getExtension() string {
	return ob.extension
}

// метод для получение информации о файлах
func getFilesRecurvise(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return true, nil
}

// метод для получение информации католога файлы
func getFileLocation(root string, filename string) (string, error) {
	if root == "" {
		return "", errors.New("Root  пуст!")
	}
	return root + "/" + filename, nil
}

// Получение все файл из котолога
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

// функция для получения значения  size
func getsize(filename string) (int64, error) {
	f, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}
	return f.Size(), nil
}

// функция для получения значения  Extension
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

// функция для получения значения  Extension2
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

// функция для получения значения  файлы из католога
func getFilesFromDirectory(pathName string) ([]Files, error) {
	fi, err := os.Open(pathName)
	if err != nil {
		log.Fatal(err, fi.Name())
	}
	defer fi.Close()
	files, err := os.ReadDir(pathName)
	if err != nil {
		fmt.Print("Невозможно прочитать каталога!", err)
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

// функция для Обработки сортировки по Убывающий
func sortAsc(arr []Files) {
	if len(arr) < 0 {
		fmt.Println("Массив пуст!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].size < arr[j].size
	})
}

// функция для Обработки сортировки по возврастающий
func sortDesc(arr []Files) {
	if len(arr) < 0 {
		fmt.Println("Массив пуст!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].size > arr[j].size
	})
}

// Чтение файлов из католога(Root)
func readDir(root string) ([]fs.FileInfo, error) {
	arrayfiles, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}
	return arrayfiles, nil
}
func main() {
	rootflag := "root"
	sortflag := "sort"
	root, sort, err := getFilePathFromCommand(rootflag, sortflag)
	if err != nil {
		fmt.Println(err)
	}
	arrayfiles, err := readDir(root)
	if err != nil {
		fmt.Println(err)
		return
	}
	filesArr := []Files{}
	var wg sync.WaitGroup
	for _, item := range arrayfiles {
		wg.Add(1)
		go func(item os.FileInfo) {
			defer wg.Done()
			if err != nil {
				fmt.Println(err)
			}

			if item.IsDir() {
				var size int64 = 0
				err := filepath.Walk(item.Name(), func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					size += info.Size()
					return nil
				})
				filesArr = append(filesArr, Files{
					name:      item.Name(),
					extension: "Каталог",
					size:      size,
				})
				if err != nil {
					log.Println(err)
				}
			} else {
				filesArr = append(filesArr, Files{
					name:      item.Name(),
					extension: "Файл",
					size:      item.Size(),
				})
			}
		}(item)
	}
	wg.Wait()

	if root != "None" && sort == "None" {
		sortAsc(filesArr)
		for i := 0; i < len(filesArr); i++ {
			filesArr[i].print()
		}
	} else if sort == "Desc" && root != "None" {
		if err != nil {
			panic(err)
		}
		sortDesc(filesArr)
		for i := 0; i < len(filesArr); i++ {
			filesArr[i].print()
		}
	}

}
