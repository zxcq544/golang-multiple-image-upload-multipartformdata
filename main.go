package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/speedata/optionparser"
)

//go:embed "static/index/*"
var embeddedFS embed.FS

var image_directory string

var video_directory string

func main() {
	serverRoot, err := fs.Sub(embeddedFS, "static/index")
	if err != nil {
		log.Println("Root directory not found: ", err)
	}
	image_directory = "static/images"
	video_directory = "static/videos"
	op := optionparser.NewOptionParser()
	op.On("-i", "--images VAL", "set images directory", &image_directory)
	op.On("-v", "--videos VAL", "set videos directory", &video_directory)

	err = op.Parse()
	if err != nil {
		log.Println("Error parsing options: ", err)
	}
	fmt.Printf("images directory is %s\n", image_directory)
	fmt.Printf("videos directory is %s\n", video_directory)
	err = os.MkdirAll(image_directory, os.ModePerm)
	if err != nil {
		log.Println("Can't create directory ", image_directory, err)
	}
	err = os.MkdirAll(video_directory, os.ModePerm)
	if err != nil {
		log.Println("Can't create directory ", video_directory, err)
	}

	http.Handle("/", http.FileServer(http.FS(serverRoot)))
	http.HandleFunc("/upload", upload)
	err = http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s", r.Method) // send data to client side
	log.Println("method:", r.Method)
	if r.Method == "POST" {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Println("Error parsing form. Maybe upload cancelled. Error: ", err)
			return
		}
		form := r.MultipartForm
		files := form.File["uploadfile[]"]
		if files == nil {
			fmt.Println("files is empty")
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
				fmt.Println("Mime type check failed ", err)
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
				log.Println("Created file ", file.Filename)
				_, err = io.Copy(local_file, f)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer local_file.Close()
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<div>File saved successfully %s</div>", file.Filename)
		}
	}
}
