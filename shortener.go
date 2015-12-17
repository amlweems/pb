package main

import (
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
)

func ftoa(f string) string {
	tmpBuf := make([]byte, 8)

	robo, _ := strconv.Atoi(f)
	n := binary.PutUvarint(tmpBuf, uint64(robo))
	user := base64.RawURLEncoding.EncodeToString(tmpBuf[:n])
	return user
}

func atof(f string) (string, error) {
	buf, err := base64.RawURLEncoding.DecodeString(f)
	if err != nil {
		return "", err
	}
	robo, _ := binary.Uvarint(buf)
	return strconv.FormatUint(robo, 10), nil
}

func usage(w http.ResponseWriter) {
	io.WriteString(w, `
    command line pastebin.

    ~$ echo Hello world. | curl -F 'f:1=<-' lf.lc
    http://lf.lc/fpW

    ~$ curl lf.lc/fpW
        Hello world.
`)
}

func fail(w http.ResponseWriter, err string) {
	w.WriteHeader(500)
	io.WriteString(w, err)
}

var directory string
var listeningPort string

func mux(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r.ParseMultipartForm(4096)
		if r.MultipartForm == nil {
			usage(w)
			return
		}
		if val, ok := r.MultipartForm.Value["f:1"]; ok {
			filePath, err := ioutil.TempFile(directory, "")
			if err != nil {
				fail(w, "Error: could not create file")
				return
			}
			defer filePath.Close()

			filePath.WriteString(val[0])

			info, _ := filePath.Stat()
			fmt.Fprintf(w, "%s\n", ftoa(info.Name()))
		} else {
			fail(w, "Error: invalid request")
			return
		}
	} else {
		urlPath, err := atof(r.URL.Path[1:])
		if err != nil {
			fail(w, "Error: invalid path")
			return
		}
		filePath := path.Join(directory, urlPath)
		slurp, err := ioutil.ReadFile(filePath)
		if err != nil {
			fail(w, "Error: invalid file")
			return
		}
		w.Write(slurp)
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.StringVar(&directory, "storage", "/srv/http/short", "storage directory")
	flag.StringVar(&listeningPort, "port", ":9000", "listening port")
	flag.Parse()

	http.HandleFunc("/", mux)
	err := http.ListenAndServe(listeningPort, Log(http.DefaultServeMux))
	if err != nil {
		log.Fatal("Unable to bind to address: ", listeningPort)
	}
}
