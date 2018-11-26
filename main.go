package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Document struct {
	Id   string
	Name string
	Size int64
}
func getMD5ChecksumOfAFile(PathDirectory string)string{
	fileOpen, eventError := os.Open(PathDirectory)
	if eventError != nil{
		log.Fatal(eventError)
	}
	defer fileOpen.Close()
	hash := md5.New()
	if _, eventError := io.Copy(hash, fileOpen); eventError != nil {
		log.Fatal(eventError)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
func getDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []Document
	/*docs = append(docs,
		Document{Id: "doc-1", Name: "Report.docx", Size: 1500})
	docs = append(docs,
		Document{Id: "doc-2", Name: "Sheet.xlsx", Size: 5000})
	docs = append(docs,
		Document{Id: "doc-3", Name: "Container.tar", Size: 50000})*/
    var PathDir = "./Files/"
	files, eventError := ioutil.ReadDir(PathDir)
	if eventError != nil {
		log.Fatal(eventError)
	}
	for _, file := range files {
		docs = append(docs,
			Document{Id: getMD5ChecksumOfAFile(PathDir + file.Name()), Name: file.Name(), Size: file.Size()})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}
