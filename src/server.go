package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	logger *log.Logger
)

func index(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("index.html")
	if err != nil {
		logger.Println(err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = io.Copy(w, file)
	if err != nil {
		logger.Println(err.Error())
		fmt.Fprintf(w, err.Error())
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	dot := r.FormValue("dot")
	cmd := exec.Command("dot", "-Tpng")
	cmd.Stdin = strings.NewReader(dot)

	w.Header().Set("Content-Type", "image/png")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Println(err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		logger.Println(err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
	_, err = io.Copy(w, stdout)
	if err != nil {
		logger.Println(err.Error())
		return
	}
	err = cmd.Wait()
	if err != nil {
		logger.Println(err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
	println("listen 0.0.0.0:8888")
	http.HandleFunc("/api/graphviz", update)
	http.HandleFunc("/", index)
	http.ListenAndServe("0.0.0.0:8888", nil)
}
