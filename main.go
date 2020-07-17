package main

import(
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
)

// Page ...
type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// URL.Pathは例えば /view/testの部分を取ってくる
	title := r.URL.Path[len("/view/"):] // /view/以降の文字列をとる
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s<h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	p1 := &Page{Title: "test", Body: []byte("This is a simple Page")}
	p1.save()
	http.HandleFunc("/view/", viewHandler) // "/view/"の場合
	log.Fatal(http.ListenAndServe(":8080", nil)) // "/"でアクセスしたときにこれが呼ばれる nilの部分はdefault handler となり、page not foundをかえせる
}