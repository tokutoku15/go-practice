package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"net/http"
	"path/filepath"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.temp1.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("uZ6dIJDB22124iZNVPAzEECmW7jaIveFnvUBgmoa9dIHF1FIBVNAAqD87dbJlCPX")
	gomniauth.WithProviders(
		github.New("ef0cf84ecb944baa8bad", "83ac464463cc40f953777c09e0461ef77be5c76f", "http://localhost:8080/auth/callback/github"),
		google.New("930271046878-b7v4g9r50dfqfeod84u67n32qk8kgtie.apps.googleusercontent.com", "8KYqDwj9B_NvWim4s373UR_H", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	//ルート
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//チャットルームを開始
	go r.run()
	//Webサーバーの開始
	log.Println("Webサーバーを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
