package main

import (
	"flag"
	"fmt"
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
var path = flag.String("path", "/simple/", "The URL path")
var root = flag.String("root", ".", "The root directory for packages")
var shop = flag.String("shop", "http://pypi.python.org", "Shop to redirect to when not found")

func listDirectory(dir string, w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing directory %s", dir), 500)
	}
	pkg := dir[strings.LastIndex(dir, "/")+1:]
	w.Write([]byte(fmt.Sprintf(LIST_HEAD, pkg, pkg)))
	for _, file := range files {
		w.Write([]byte(fmt.Sprintf(LIST_ELEMENT, *path+pkg+"/"+file.Name(), file.Name())))
	}
	w.Write([]byte(LIST_TAIL))
}

func servePackage(filename string, w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filename)
}

func handler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(*root, r.URL.Path[1:])
	if info, err := os.Stat(filename); err != nil {
		url := *shop + r.URL.Path
		log.Print("Redirecting to ", url)
		http.Redirect(w, r, url, 302)
	} else {
		switch mode := info.Mode(); {
		case mode.IsDir():
			log.Print("Listing directory ", filename)
			listDirectory(filename, w, r)
		case mode.IsRegular():
			log.Print("Serving package ", filename)
			servePackage(filename, w, r)
		}
	}
}

func parseCommandLine() {
	flag.Parse()
	absroot, err := filepath.Abs(*root)
	if err != nil {
		panic("Error building root directory")
	}
	root = &absroot
}

func main() {
	parseCommandLine()
	http.HandleFunc(*path, handler)
	log.Print("Starting CheeseShop version ", VERSION)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	log.Print("Stopping CheeseShop")
}
