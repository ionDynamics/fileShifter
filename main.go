package main

import (
	"flag"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

var folderPtr = flag.String("folder", "./tmp/", "")
var assets *rice.Box

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	assets = rice.MustFindBox("assets")

	http.Handle("/css/", http.FileServer(assets.HTTPBox()))
	http.Handle("/js/", http.FileServer(assets.HTTPBox()))
	http.Handle("/_/", http.StripPrefix("/_/", http.FileServer(http.Dir(*folderPtr))))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3001", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templateString, err := assets.String("index.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	// parse and execute the template
	tmplIndex, err := template.New("index").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}
	tmplIndex.Execute(w, map[string]string{"Message": "Hello, world!"})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload")
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	out, err := os.Create(*folderPtr + header.Filename)
	if err != nil {
		fmt.Println("Unable to create the file for writing. Check your write access privilege")
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("File uploaded successfully : ")
	fmt.Println(header.Filename)
}
