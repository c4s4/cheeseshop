package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	VERSION      = "UNKNOWN"
	LIST_HEAD    = "<html><head><title>Links for %s</title></head><body><h1>Links for %s</h1>"
	LIST_TAIL    = "</body></html>"
	LIST_ELEMENT = "<a href='%s'>%s</a><br/>"
)

var port = flag.Int("port", 8000, "The port CheeseShop is listening")
var path = flag.String("path", "simple", "The URL path")
var root = flag.String("root", ".", "The root directory for packages")
var shop = flag.String("shop", "http://pypi.python.org", "Redirection when not found")

func listRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("Listing root %s", *root)
	files, err := ioutil.ReadDir(*root)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing root directory %s", *root), 500)
		return
	}
	w.Write([]byte(fmt.Sprintf(LIST_HEAD, "root", "root")))
	for _, file := range files {
		if file.Mode().IsDir() {
			w.Write([]byte(fmt.Sprintf(LIST_ELEMENT, *path+file.Name(), file.Name())))
		}
	}
	w.Write([]byte(LIST_TAIL))
}

func listDirectory(dir string, w http.ResponseWriter, r *http.Request) {
	directory := filepath.Join(*root, dir)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		url := *shop + *path + dir
		log.Printf("Redirecting to %s", url)
		http.Redirect(w, r, url, 302)
		return
	}
	log.Printf("Listing directory %s", directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing directory %s", dir), 500)
		return
	}
	w.Write([]byte(fmt.Sprintf(LIST_HEAD, dir, dir)))
	for _, file := range files {
		w.Write([]byte(fmt.Sprintf(LIST_ELEMENT, *path+dir+"/"+file.Name(), file.Name())))
	}
	w.Write([]byte(LIST_TAIL))
}

func servePackage(dir, file string, w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(*root, dir, file)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		url := *shop + *path + dir + "/" + file
		log.Printf("Redirecting to %s", url)
		http.Redirect(w, r, url, 302)
		return
	}
	log.Printf("Serving file %s", filename)
	http.ServeFile(w, r, filename)
}

func copyFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m := r.MultipartForm
	files := m.File["content"]
	for _, file := range files {
		name := file.Filename
		pack := name[:strings.LastIndex(name, "-")]
		log.Printf("Writing file %s", name)
		f, err := file.Open()
		defer f.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dst, err := os.Create(filepath.Join(*root, pack, name))
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(dst, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parts := strings.Split(r.URL.Path[len(*path):], "/")
		if len(parts) > 2 {
			http.Error(w, fmt.Sprintf("%s is not a valid path", r.URL.Path), 404)
			return
		} else if len(parts) == 1 && parts[0] == "" {
			listRoot(w, r)
		} else if len(parts) == 1 {
			listDirectory(parts[0], w, r)
		} else {
			servePackage(parts[0], parts[1], w, r)
		}
	} else if r.Method == "POST" {
		copyFile(w, r)
	}
}

func parseCommandLine() {
	flag.Parse()
	absroot, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal("Error building root directory")
	}
	root = &absroot
	info, err := os.Stat(*root)
	if err != nil {
		log.Fatalf("Root directory %s not found", *root)
	}
	if !info.Mode().IsDir() {
		log.Fatalf("Root %s is not a directory", *root)
	}
	if !strings.HasPrefix(*path, "/") {
		p := "/" + *path
		path = &p
	}
	if !strings.HasSuffix(*path, "/") {
		p := *path + "/"
		path = &p
	}
	if *port > 65535 || *port < 0 {
		log.Fatalf("Bad port number %d", *port)
	}
}

func main() {
	parseCommandLine()
	http.HandleFunc(*path, handler)
	log.Print("Starting CheeseShop (version: ", VERSION, ", path: ", *path, ", port: ", *port, ", root: ", *root, ")")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	log.Print("Stopping CheeseShop")
}
