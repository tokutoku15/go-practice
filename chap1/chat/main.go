package main

import (
	"log"
	"sync"

	"net/http"
	"path/filepath"
	"text/template"
)

// temp1は1つのテンプレートを表す
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

// ServeHTTPはHTTPリクエストを処理
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w, nil)
}

func main() {
	r := newRoom()
	//ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	//チャットルームを開始
	go r.run()
	//Webサーバーの開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
