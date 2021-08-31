package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var index_response_string []byte

func main() {
	var err error
	index_response_string, err = ioutil.ReadFile("./static/index.html")
	if err != nil {
		log.Fatal("Missing index.html: ", err)
	}
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/upload", upload)
	err = http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%s", index_response_string)
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Hello %s", r.Method) // send data to client side
	fmt.Println("method:", r.Method)
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		form := r.MultipartForm
		files := form.File["uploadfile[]"]
		if files == nil {
			fmt.Println("files is empty")
			return
		}

		for _, file := range files {
			log.Println(file.Filename)
			f, err := file.Open()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			local_file, err := os.OpenFile(file.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer local_file.Close()
			io.Copy(local_file, f)
			fmt.Fprintf(w, "<div>File saved successfully %s</div>", file.Filename)
		}
	}
}
