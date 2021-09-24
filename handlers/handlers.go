package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func Upload(w http.ResponseWriter, r *http.Request) {
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
			fmt.Println(file.Filename)

			f, err := file.Open()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()

			local_file, err := os.OpenFile(path.Join("./static/images", file.Filename), os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println("Can't create file:", err)
				return
			}
			io.Copy(local_file, f)
			defer local_file.Close()
			mime_type, err := mimetype.DetectFile(path.Join("./static/images", file.Filename))
			if err != nil {
				fmt.Println("Mime type check failed ", err)
				return
			}
			fmt.Println("	mime type", mime_type)
			if !strings.Contains(mime_type.String(), "image") {
				err = local_file.Close()
				if err != nil {
					log.Println("Can't close file: ", err)
				}
				log.Printf("File %s is not image. Removing...\n", file.Filename)
				err = os.Remove(path.Join("./static/images", file.Filename))
				if err != nil {
					log.Println("Error removing file: ", err)
				}
				log.Printf("Successfully removed: %s\n", file.Filename)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "File not saved because is not image %s. Type %s", file.Filename, mime_type)
				return
			}
			w.WriteHeader(http.StatusOK)
			// w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "File saved successfully %s", file.Filename)
		}
	}
}
