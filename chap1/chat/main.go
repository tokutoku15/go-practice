package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"net/http"
	"path/filepath"
	"text/template"

	"github.com/tokutoku15/go-practice/chap1/trace"
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
	t.temp1.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	//チャットルームを開始
	go r.run()
	//Webサーバーの開始
	log.Println("Webサーバーを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
