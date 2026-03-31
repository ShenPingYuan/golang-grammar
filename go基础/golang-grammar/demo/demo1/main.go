package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	content, err := readFile("test.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Print(content)
}

// 编写一个通用函数，用于打开文件并读取其内容。该函数应该接受文件路径作为参数，并返回文件的内容和可能的错误。
func readFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var out []byte
	buf := make([]byte, 32*1024)
	for {
		n, rerr := f.Read(buf)
		if n > 0 {
			out = append(out, buf[:n]...)
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return "", rerr
		}
	}
	return string(out), nil
}
