package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			panic(fmt.Sprintf("open file %s error, %s", path, err))
		}
		defer file.Close()
	}
}

func CreateDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}
}

func SaveList2Txt(data []string, path string) {
	file, err := os.OpenFile(path, os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("open file %s error, %s", path, err))
	}

	for _, line := range data {
		line += "\n\n"
		_, err := file.WriteString(line)
		if err != nil {
			panic(fmt.Sprintf("write file %s error, %s", path, err))
		}
	}
	defer file.Close()
}

func ReadTxt2Set(path string) (res map[string]bool) {
	res = make(map[string]bool)
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		panic(fmt.Sprintf("open file %s error, %s", path, err))
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("read file error: %s", err))
		}
		line = strings.TrimRight(line, "\n")
		if line != "" {
			res[line] = true
		}
	}
	return
}

func RemoveContents(dir string) (n int, err error) {
	d, err := os.Open(dir)
	if err != nil {
		return 0, err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return 0, err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return n, err
		}
		n++
	}
	return
}
