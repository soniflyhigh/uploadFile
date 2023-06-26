package main

import (
	// "fmt"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	// "log"
	"net/http"
	// "os"
)

var temp *template.Template
const chunkSize = 1024 * 1024 // 1MB

func init(){
	temp = template.Must(template.ParseFiles("template/index.html"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request){
	filePath := `C:\Users\Admin\Downloads\saveData`
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()
    
	out, err := os.Create(filepath.Join(filePath, header.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()
	
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if n == 0 {
			break
		}

		// Write the chunk to the output file
		if _, err := out.Write(buf[:n]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fmt.Fprintf(w, "File uploaded successfully!")
}

func HandleFunc(w http.ResponseWriter, r *http.Request){
	http.HandleFunc("/upload", uploadHandler)
	temp.ExecuteTemplate(w, "index.html",nil)
}

func main(){
	http.HandleFunc("/",HandleFunc)
	http.ListenAndServe(":9999",nil)
}


