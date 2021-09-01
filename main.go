package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/speedata/optionparser"
)

var index_response_string []byte

var image_directory string

var video_directory string

func main() {
	image_directory = "static/images"
	video_directory = "static/videos"
	op := optionparser.NewOptionParser()
	op.On("-i", "--images VAL", "set images directory", &image_directory)
	op.On("-v", "--videos VAL", "set videos directory", &video_directory)

	err := op.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("images directory is %s\n", image_directory)
	fmt.Printf("videos directory is %s\n", video_directory)
	err = os.MkdirAll(image_directory, os.ModePerm)
	if err != nil {
		log.Fatal("Can't create directory ", image_directory, err)
	}
	err = os.MkdirAll(video_directory, os.ModePerm)
	if err != nil {
		log.Fatal("Can't create directory ", video_directory, err)
	}

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
			fmt.Print(file.Filename)

			f, err := file.Open()
			if err != nil {
				fmt.Println(err)
				return
			}
			mime_type, err := mimetype.DetectReader(f)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("	mime type", mime_type)
			defer f.Close()

			if strings.Contains(mime_type.String(), "image") {
				local_file, err := os.OpenFile(path.Join(image_directory, file.Filename), os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
					return
				}
				io.Copy(local_file, f)
				defer local_file.Close()
			}
			if strings.Contains(mime_type.String(), "video") {
				local_file, err := os.OpenFile(path.Join(video_directory, file.Filename), os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
					return
				}
				io.Copy(local_file, f)
				defer local_file.Close()
			}

			fmt.Fprintf(w, "<div>File saved successfully %s</div>", file.Filename)
		}
	}
}
