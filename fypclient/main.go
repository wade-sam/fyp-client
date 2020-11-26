package main

import (
	"github.com/wade-sam/fypclient/filescan"
	"github.com/wade-sam/fypclient/filescan/writetree"
)

func main() {
	subDirToSkip := "golib"
	head := "/backup/Documents/"
	fileScanResult := filescan.DirectoryScan(head, subDirToSkip)
	writetree.WriteJSON(fileScanResult)
	//fileScanResult := filescanDirectoryScan(head, subDirToSkip)
	//fmt.Println(fileScanResult)
	//for key, value := range fileScanResult.Filepath {
	//	fmt.Println(key, value.Filename, value.Checksum)
	//}

}
