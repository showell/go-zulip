package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
)

import (
	"go-zulip/html"
	"go-zulip/zulip"
)

//go:embed style.css
var styleCSS []byte

type StringWriterForBytes struct {
	w io.Writer
}

func (sw StringWriterForBytes) WriteString(s string) (int, error) {
	sw.w.Write([]byte(s))
	return 0, nil
}

func main() {
	db, err := zulip.BuildDatabase("config.json")
	if err != nil {
		fmt.Printf("Error building database: %v\n", err)
		return
	}

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(styleCSS)
	})

	http.HandleFunc("/{path...}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		stringWriter := StringWriterForBytes{w: w}
		stringWriter.WriteString("<html><head><link rel='stylesheet' href='/style.css'></head><body>\n")
		path := r.PathValue("path")
		html.Html(db, "/"+path, stringWriter)
		stringWriter.WriteString("</body></html>\n")
	})

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}
