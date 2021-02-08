package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/wade-sam/fypclient/backup"

	"github.com/wade-sam/fypclient/backup"
	"github.com/wade-sam/fypclient/filescan"
	"github.com/wade-sam/fypclient/filescan/writetree"
)

func Filescan(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Filescan() called")
	w.Header().Set("Content-Type", "application/json")
	subDirToSkip := "golib"
	head := "/home/dev"
	filescan := filescan.InitialDirectoryScan(head, subDirToSkip)
	//fileScanResult := filescan.InitialDirectoryScan(head, subDirToSkip)
	response := writetree.ObjectToJson(filescan)

	json.NewEncoder(w).Encode(response)
	fmt.Println("Finished")
	//fmt.Fprintf(w, "Endpoint called:  test page")
}

func FBackup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FullBackup called")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Full backup started")
	subDirToSkip := ""
	head := "/home/dev"
	filescan := filescan.InitialDirectoryScan(head, subDirToSkip)
	backup.FullBackup(filescan)
	fmt.Println("FullBackup called")
}

func IBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Incremental backup started")
	subDirToSkip := "golib"
	head := "/home/dev"
	filescan := filescan.InitialDirectoryScan(head, subDirToSkip)
	backup.IncrementalBackup(filescan)
	fmt.Println("Incremental() started")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/filescan", Filescan).Methods("GET")
	router.HandleFunc("/full", FBackup).Methods("GET")
	router.HandleFunc("/incremental", IBackup).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {

	handleRequests()
	//subDirToSkip := "golib"
	//head := "/backup/Documents"
	//fileScanResult := filescan.InitialDirectoryScan(head, subDirToSkip)
	//for key, value := range fileScanResult.Filepath {
	//	fmt.Println(key, value.Filename, value.Checksum, value.Permissions.Ownership)
	//	}
	//backup.FullBackup(fileScanResult)
	//backup.IncrementalBackup(fileScanResult)
	//writetree.WriteToFile(fileScanResult)
	//time.Sleep(20 * time.Second)
	//	fmt.Println("Checking for differences")
	//differences := writetree.CompareJsonFile(fileScanResult)
	//fmt.Println(differences)
	//readfile := writetree.ReadInJsonFile()

	//fileScanResult := filescanDirectoryScan(head, subDirToSkip)
	//fmt.Println(fileScanResult)
	//for key, value := range fileScanResult.Filepath {
	//	fmt.Println(key, value.Filename, value.Checksum)
	//}

}
