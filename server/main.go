package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("You need to specifiy a directory to run in")
	}

	dir := args[1]

	fs := http.FileServer(http.Dir(dir))
	log.Print("Serving " + dir + " on http://localhost:8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Println(req.URL)
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(resp, req)
	}))

	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal("An error has occured:", err)
		}
	}
}
