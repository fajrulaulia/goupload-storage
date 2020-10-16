package gouploadstorage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type response struct {
	FilenameSource string `json:"filename_source"`
	Filename       string `json:"filename"`
	Size           int    `json:"size"`
}

func upload(w http.ResponseWriter, r *http.Request) {
	var res response
	w.Header().Set("Content-Type", "application/json")
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file_source")
	if err != nil {
		w.WriteHeader(500)
		log.Println("r.FormFile", err)
		return
	}
	defer file.Close()
	res.FilenameSource = handler.Filename
	res.Size = int(handler.Size)
	tempFile, err := ioutil.TempFile("storage", "upload-*.jpeg")
	if err != nil {
		w.WriteHeader(500)
		log.Println("ioutil.TempFile", err)
		return
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		log.Println("ioutil.ReadAl", err)
		return
	}
	tempFile.Write(fileBytes)
	path := strings.Split(tempFile.Name(), "/")
	filename := path[len(path)-1]
	res.Filename = filename
	output, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		log.Println("json.Marshal", err)
		return
	}
	w.Write(output)
}

func preview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	img, err := os.Open("storage/" + vars["filename"])
	if err != nil {
		w.WriteHeader(404)
		log.Println("json.Marshal", err)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}

//SetupHandler Should be exported
func SetupHandler() {

	r := mux.NewRouter()
	r.HandleFunc("/upload", upload)
	r.HandleFunc("/image/{filename}", preview)
	log.Fatal(http.ListenAndServe(":8080", r))
}
