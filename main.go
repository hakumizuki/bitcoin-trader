package main

import(
	// "fmt"
	"html/template"
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// URL.Pathは例えば /view/testの部分を取ってくる
	title := r.URL.Path[len("/view/"):] // /view/以降の文字列をとる
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	// fmt.Fprintf(w, "<h1>%s<h1><div>%s</div>", p.Title, p.Body)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	// URL.Pathは例えば /edit/testの部分を取ってくる
	title := r.URL.Path[len("/edit/"):] // /edit/以降の文字列をとる
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	// fmt.Fprintf(w, "<h1>%s<h1><div>%s</div>", p.Title, p.Body)
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body") // htmlの方でtextareaのname属性を"body"にして、submitで送信したため、"body"で中身をとることができる
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	p1 := &Page{Title: "test", Body: []byte("This is a simple Page")}
	p1.save()
	http.HandleFunc("/view/", viewHandler) // "/view/"の場合
	http.HandleFunc("/edit/", editHandler) // "/edit/"の場合
	http.HandleFunc("/save/", saveHandler) // "/save/"の場合
	log.Fatal(http.ListenAndServe(":8080", nil)) // "/"でアクセスしたときにこれが呼ばれる nilの部分はdefault handler となり、page not foundを返す
}