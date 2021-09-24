package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"main/handlers"
	"net/http"
	"os"

	"github.com/speedata/optionparser"
)

func init() {
	Image_directory, Video_directory = InitApp()
}

func InitApp() (string, string) {
	return "static/images", "static/videos"
}

//go:embed "static/index/*"
var embeddedFS embed.FS

var Image_directory string

var Video_directory string

func main() {
	serverRoot, err := fs.Sub(embeddedFS, "static/index")
	if err != nil {
		log.Println("Root directory not found: ", err)
	}
	op := optionparser.NewOptionParser()
	op.On("-i", "--images VAL", "set images directory", &Image_directory)
	op.On("-v", "--videos VAL", "set videos directory", &Video_directory)

	err = op.Parse()
	if err != nil {
		log.Println("Error parsing options: ", err)
	}
	fmt.Printf("images directory is %s\n", Image_directory)
	fmt.Printf("videos directory is %s\n", Video_directory)
	err = os.MkdirAll(Image_directory, os.ModePerm)
	if err != nil {
		log.Println("Can't create directory ", Image_directory, err)
	}
	err = os.MkdirAll(Video_directory, os.ModePerm)
	if err != nil {
		log.Println("Can't create directory ", Video_directory, err)
	}

	http.Handle("/", http.FileServer(http.FS(serverRoot)))
	http.HandleFunc("/upload", handlers.Upload)
	err = http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// upload logic
// func upload(w http.ResponseWriter, r *http.Request) {
// 	log.Println("method:", r.Method)
// 	if r.Method == "POST" {
// 		err := r.ParseMultipartForm(32 << 20)
// 		if err != nil {
// 			log.Println("Error parsing form. Maybe upload cancelled. Error: ", err)
// 			return
// 		}
// 		form := r.MultipartForm
// 		files := form.File["uploadfile[]"]
// 		if files == nil {
// 			fmt.Println("files is empty")
// 		}

// 		for _, file := range files {
// 			fmt.Print(file.Filename)

// 			f, err := file.Open()
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			defer f.Close()

// 			local_file, err := os.OpenFile(path.Join(Image_directory, file.Filename), os.O_WRONLY|os.O_CREATE, 0666)
// 			if err != nil {
// 				fmt.Println("Can't create file:", err)
// 				return
// 			}
// 			io.Copy(local_file, f)
// 			defer local_file.Close()
// 			mime_type, err := mimetype.DetectFile(path.Join(Image_directory, file.Filename))
// 			if err != nil {
// 				fmt.Println("Mime type check failed ", err)
// 				return
// 			}
// 			fmt.Println("	mime type", mime_type)
// 			if !strings.Contains(mime_type.String(), "image") {
// 				err = local_file.Close()
// 				if err != nil {
// 					log.Println("Can't close file: ", err)
// 				}
// 				log.Printf("File %s is not image. Removing...\n", file.Filename)
// 				err = os.Remove(path.Join(Image_directory, file.Filename))
// 				if err != nil {
// 					log.Println("Error removing file: ", err)
// 				}
// 				log.Printf("Successfully removed: %s\n", file.Filename)
// 				w.WriteHeader(http.StatusBadRequest)
// 				fmt.Fprintf(w, "File not saved because is not image %s. Type %s", file.Filename, mime_type)
// 				return
// 			}
// 			w.WriteHeader(http.StatusOK)
// 			// w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 			fmt.Fprintf(w, "File saved successfully %s", file.Filename)
// 		}
// 	}
// }
