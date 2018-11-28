package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Document struct {
	Id   string
	Name string
	Size int64
}

const PathDir = "./Files/"

func readFilesFromPath() []os.FileInfo {
	files, eventError := ioutil.ReadDir(PathDir)
	if eventError != nil {
		log.Fatal(eventError)
	}
	return files
}
func getIdFromRequest(req *http.Request) string {
	vars := mux.Vars(req)
	id := vars["Id"]
	return id
}
func getDocument(IdFile string) (Document, error) {
	files := readFilesFromPath()
	var foundDoc Document

	for _, file := range files {
		if strings.Compare(getMD5ChecksumOfAFile(PathDir+file.Name()), IdFile) == 0 {
			foundDoc = Document{Id: IdFile, Name: file.Name(), Size: file.Size()}
			return foundDoc, nil
		}
	}
	return foundDoc, errors.New("No found document with Id:" + IdFile)
}
func getMD5ChecksumOfAFile(PathDirectory string) string {
	fileOpen, eventError := os.Open(PathDirectory)
	if eventError != nil {
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

	files := readFilesFromPath()
	for _, file := range files {
		docs = append(docs,
			Document{Id: getMD5ChecksumOfAFile(PathDir + file.Name()), Name: file.Name(), Size: file.Size()})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}
func getDocumentById(w http.ResponseWriter, r *http.Request) {
	id := getIdFromRequest(r)

	foundDoc, eventError := getDocument(id)
	if eventError == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(foundDoc)
	} else {
		fmt.Fprintln(w, eventError)
	}
}
func createDocument(w http.ResponseWriter, r *http.Request) {
	file, handler, eventError := r.FormFile("file")
	fileName := handler.Filename

	if eventError != nil {
		log.Fatal(eventError)
	}
	defer file.Close()
	fileOpen, eventError := os.OpenFile(PathDir+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if eventError != nil {
		log.Fatal(eventError)
	}
	defer fileOpen.Close()
	_, err := io.Copy(fileOpen, file)
	if err != nil {
		log.Fatal("Error to create the document. Check your write access privilege. " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Document uploaded successfully!. File Name:"+fileName)
}

func deleteDocumentById(w http.ResponseWriter, r *http.Request) {

	id := getIdFromRequest(r)
	foundDoc, eventError := getDocument(id)
	if eventError == nil {
		eventError := os.Remove(PathDir + foundDoc.Name)
		if eventError != nil {
			fmt.Println(eventError)
		} else {
			fmt.Fprintln(w, "Document has been removed")
		}
	} else {
		fmt.Fprintln(w, eventError)
	}
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	router.HandleFunc("/documents/{Id}", getDocumentById).Methods("GET")
	router.HandleFunc("/documents", createDocument).Methods("POST")
	router.HandleFunc("/documents/{Id}", deleteDocumentById).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9000", router))
}
