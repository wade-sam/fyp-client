package main

import (
	"fmt"

	"github.com/wade-sam/fypclient/filescan"
	"github.com/wade-sam/fypclient/filescan/writetree"
)

func main() {
	subDirToSkip := "golib"
	head := "/backup/Documents"
	fileScanResult := filescan.InitialDirectoryScan(head, subDirToSkip)

	//writetree.WriteToFile(fileScanResult)
	//time.Sleep(20 * time.Second)
	//fmt.Println("Checking for differences")
	differences := writetree.CompareJsonFile(fileScanResult)
	fmt.Println(differences)
	//readfile := writetree.ReadInJsonFile()

	//fileScanResult := filescanDirectoryScan(head, subDirToSkip)
	//fmt.Println(fileScanResult)
	//for key, value := range fileScanResult.Filepath {
	//	fmt.Println(key, value.Filename, value.Checksum)
	//}

}
