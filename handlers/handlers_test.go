package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

func TestUpload(t *testing.T) {
	// file_to_test, err := ioutil.ReadFile("test_image.png")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(file_to_test)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	b, w := createMultipartFormData(t, "uploadfile[]", "./test_image.png")

	req, err := http.NewRequest("POST", "/upload", &b)
	if err != nil {
		fmt.Println("FILE NOT FOUND ", "./test_image.png")
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Upload)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "File saved successfully test_image.png"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	original_file, err := ioutil.ReadFile("./test_image.png")
	if err != nil {
		t.Log(err)
	}
	t.Log("Checking equality")
	uploaded_file, err := ioutil.ReadFile("./static/images/test_image.png")
	if err != nil {
		t.Log(err)
	}
	if !reflect.DeepEqual(original_file, uploaded_file) {
		t.Error("Error. Files are different. Upload is corrupted")
	} else {
		t.Log("Files are equal")
	}

}
