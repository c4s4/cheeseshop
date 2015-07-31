package main

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	VERSION      = "UNKNOWN"
	LIST_HEAD    = "<html><head><title>Links for %s</title></head><body><h1>Links for %s</h1>"
	LIST_TAIL    = "</body></html>"
	LIST_ELEMENT = "<a href='%s'>%s</a><br/>"
)

var DEFAULT_CONFIG = []string{"~/.cheeseshop.yml", "/etc/cheeseshop.yml"}

type Config struct {
	Http      int
	Https     int
	Path      string
	Root      string
	Shop      string
	Cert      string
	Key       string
	Overwrite bool
	Auth      map[string]string
}

var config Config

func listRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("Listing root %s", config.Root)
	files, err := ioutil.ReadDir(config.Root)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing root directory %s", config.Root), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf(LIST_HEAD, "root", "root")))
	for _, file := range files {
		if file.Mode().IsDir() {
			w.Write([]byte(fmt.Sprintf(LIST_ELEMENT, config.Path+file.Name(), file.Name())))
		}
	}
	w.Write([]byte(LIST_TAIL))
}

func listDirectory(dir string, w http.ResponseWriter, r *http.Request) {
	directory := filepath.Join(config.Root, dir)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		url := config.Shop + "/" + dir
		log.Printf("Redirecting to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	log.Printf("Listing directory %s", directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing directory %s", dir), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf(LIST_HEAD, dir, dir)))
	for _, file := range files {
		w.Write([]byte(fmt.Sprintf(LIST_ELEMENT, config.Path+dir+"/"+file.Name(), file.Name())))
	}
	w.Write([]byte(LIST_TAIL))
}

func servePackage(dir, file string, w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(config.Root, dir, file)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		url := config.Shop + "/" + dir + "/" + file
		log.Printf("Redirecting to %s", url)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	log.Printf("Serving file %s", filename)
	http.ServeFile(w, r, filename)
}

func copyFile(w http.ResponseWriter, r *http.Request) {
	if len(config.Auth) > 0 {
		username, password, ok := r.BasicAuth()
		sum := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		if !ok || config.Auth[username] != sum {
			log.Printf("Unauthorized access from %s", username)
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Granted access for user %s", username)
	}
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
		dir := filepath.Join(config.Root, pack)
		path := filepath.Join(config.Root, pack, name)
		if _, err := os.Stat(path); err == nil && !config.Overwrite {
			log.Printf("Failed attempt to overwrite package %s", name)
			http.Error(w, fmt.Sprintf("Package %s already exists", name), http.StatusBadRequest)
			return
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
			if err != nil {
				log.Printf("Error creating directory for package %s", pack)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			log.Printf("Created directory for package %s", pack)
		}
		log.Printf("Writing file %s", name)
		f, err := file.Open()
		defer f.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dst, err := os.Create(path)
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
		parts := strings.Split(r.URL.Path[len(config.Path):], "/")
		if len(parts) > 2 {
			http.Error(w, fmt.Sprintf("%s is not a valid path", r.URL.Path), http.StatusNotFound)
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

func normalizeFile(file string) string {
	if strings.HasPrefix(file, "~") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		file = filepath.Join(dir, file[1:])
	}
	absfile, err := filepath.Abs(file)
	if err != nil {
		log.Fatalf("Error getting absolute path for file %s", file)
	}
	return absfile
}

func loadConfig() {
	var file = ""
	if len(os.Args) > 1 {
		file = os.Args[1]
	} else {
		for _, path := range DEFAULT_CONFIG {
			path = normalizeFile(path)
			if _, err := os.Stat(path); err == nil {
				file = path
				break
			}
		}
	}
	if file == "" {
		log.Fatal("No configuration file found")
	}
	log.Printf("Loading %s configuration file", file)
	source, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error loading config file %s", file)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatalf("Error parsing config file %s: %s", file, err)
	}
}

func checkConfig() {
	config.Root = normalizeFile(config.Root)
	info, err := os.Stat(config.Root)
	if err != nil {
		log.Fatalf("Root directory %s not found", config.Root)
	}
	if !info.Mode().IsDir() {
		log.Fatalf("Root %s is not a directory", config.Root)
	}
	if !strings.HasPrefix(config.Path, "/") {
		config.Path = "/" + config.Path
	}
	if !strings.HasSuffix(config.Path, "/") {
		config.Path = config.Path + "/"
	}
	if config.Http > 65535 || config.Http < 0 {
		log.Fatalf("Bad HTTP port number %d", config.Http)
	}
	if config.Https > 65535 || config.Https < 0 {
		log.Fatalf("Bad HTTPS port number %d", config.Https)
	}
	if config.Http == 0 && config.Https == 0 {
		log.Fatal("At least one of HTTP or HTTPS must be enabled")
	}
}

func main() {
	loadConfig()
	checkConfig()
	log.Printf("Starting CheeseShop (ports: %d & %d, path: %s, root: %s, shop: %s, cert: %s, key: %s, overwrite: %t)",
		config.Http, config.Https, config.Path, config.Root, config.Shop, config.Cert, config.Key, config.Overwrite)
	http.HandleFunc(config.Path, handler)
	if config.Http != 0 {
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Http), nil))
		}()
	}
	if config.Https != 0 {
		go func() {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Https),
				normalizeFile(config.Cert), normalizeFile(config.Key), nil))
		}()
	}
	wait := make(chan bool, 1)
	<-wait
}
