package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
