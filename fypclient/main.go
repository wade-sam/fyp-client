package main

import (
	"fmt"

	"github.com/wade-sam/fypclient/filescan"
)

func main() {
	fmt.Println("Hello you yute!")
	subDirToSkip := "golib"
	head := "/backup/Documents"
	fileScanResult := filescan.DirectoryScan(head, subDirToSkip)
	//fileScanResult := filescanDirectoryScan(head, subDirToSkip)
	//fmt.Println(fileScanResult)
	for key, value := range fileScanResult.Filepath {
		fmt.Println(key, value.Filename, value.Checksum)
	}
}
