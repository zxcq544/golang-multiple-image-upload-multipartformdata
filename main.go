package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/speedata/optionparser"
)

var index_response_string []byte

var image_directory string

var video_directory string

func main() {
	image_directory = "image"
	video_directory = "video"
	op := optionparser.NewOptionParser()
	op.On("-i", "--image VAL", "set image directory", &image_directory)
	op.On("-v", "--video VAL", "set video directory", &video_directory)

	err := op.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("image directory is %s\n", image_directory)
	fmt.Printf("video directory is %s\n", video_directory)

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
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Fatal("Error parsing form", err)
		}
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
