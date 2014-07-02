package main

import (
	"github.com/sysr-q/assets"
	"html/template"
	"net/http"
)

const Foo_1234 = []byte{0x01,0x02,0x03}

func tmpl() *template.Template {
	s := string(assets.MustRead("templates/layout.html"))
	t := template.New("layout")
	template.Must(t.Parse(s))
	return t
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	s := string(assets.MustRead("templates/index.html"))
	t := template.Must(tmpl().Parse(s))

	t.Execute(w, nil)
}

func main() {
	// /dev/null for you!
	_ = string(assets.MustRead("templates/layout.html"))

	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)
}
