package main

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"time"
)

const (
	// UPLOADDIR UPLOADDIR
	UPLOADDIR = "./upload"
	// KEY KEY
	KEY = "photoweb"
)

// FileInfo FileInfo
type FileInfo struct {
	Name string
	Size int64
}

func toMd5(msg string) (msgMd5 string) {
	md5Inst := md5.New()
	md5Inst.Write([]byte(msg))
	return hex.EncodeToString(md5Inst.Sum([]byte(nil)))
}

func render(rw http.ResponseWriter, tmpName string, locals interface{}) {
	tmpl, err := template.ParseFiles("./view/" + tmpName + ".html")
	if err != nil {
		return
	}
	tmpl.Execute(rw, locals)
}

//------------
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(rw, e.Error(), http.StatusInternalServerError)
				log.Printf("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()

		fn(rw, req)
	}
}
func listHandle(rw http.ResponseWriter, req *http.Request) {
	files, err := ioutil.ReadDir(UPLOADDIR)

	checkErr(err)

	locals := make(map[string]interface{})

	sFiles := []FileInfo{}

	for _, file := range files {
		sFiles = append(sFiles, FileInfo{file.Name(), file.Size()})
	}
	locals["files"] = sFiles

	render(rw, "list", locals)
	// tmpl, err := template.ParseFiles("./view/list.html")
	// tmpl.Execute(rw, locals)
	// var htmlList = "<ul>"
	// for _, file := range files {

	// 	if !file.IsDir() {
	// 		name := file.Name()
	// 		htmlList += "<li><a href='/view?id=" + name + "'>" + name + "(" + fmt.Sprintf("%d", file.Size()) + ")</a></li>"
	// 	}
	// }
	// htmlList += "</ul>"

	// rw.Header().Set("Content-type", "text/html; charset=utf-8")

	// io.WriteString(rw, htmlList)
}
func viewHandle(rw http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	filepath := path.Join(UPLOADDIR, id)
	_, err := os.Stat(filepath)
	// if err != nil {
	// 	http.NotFound(rw, req)
	// 	return
	// }

	checkErr(err)
	rw.Header().Set("Content-type", "image")
	http.ServeFile(rw, req, filepath)
}
func uploadHandle(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		render(rw, "upload", nil)
		return
		// file, err := os.Open("./view/upload.html")
		// if err == nil {
		// 	rw.Header().Set("Content-type", "text/html; charset=utf-8")
		// 	io.Copy(rw, file)
		// 	return
		// }
	} else if req.Method == "POST" {
		file, header, err := req.FormFile("file_upload")
		checkErr(err)

		defer file.Close()
		filename := header.Filename

		filenameNew := toMd5(time.Now().String()) + path.Ext(filename)

		os.MkdirAll(UPLOADDIR, os.ModeDir)
		target, err := os.Create(UPLOADDIR + "/" + filenameNew)
		checkErr(err)
		defer target.Close()

		_, err1 := io.Copy(target, file)
		checkErr(err1)
		http.Redirect(rw, req, "/view?id="+filenameNew, http.StatusFound)
	}
}
func main() {
	os.MkdirAll(UPLOADDIR, os.ModeDir)
	http.HandleFunc("/list", safeHandler(listHandle))
	http.HandleFunc("/view", safeHandler(viewHandle))
	http.HandleFunc("/upload", safeHandler(uploadHandle))
	http.HandleFunc("/assets/", func(prefix string) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, req *http.Request) {
			file := path.Join("./public/", req.URL.Path[len(prefix)-1:])
			_, err := os.Stat(file)
			if err != nil {
				http.NotFound(rw, req)
			}

			http.ServeFile(rw, req, file)
		}
	}("/assets/"))

	http.ListenAndServe(":8080", nil)
}
