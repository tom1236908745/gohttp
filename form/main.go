package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
	//注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。
	fmt.Println(r.Form) //これらのデータはサーバのプリント情報に出力されます
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //ここでwに書き込まれたものがクライアントに出力されます。
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //リクエストを取得するメソッド
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.gtpl")
		t.Execute(w, nil)
	} else {
		//ログインデータがリクエストされ、ログインのロジック判断が実行されます。
		r.ParseForm()
		// scrip文など外部からのコードを許容してしまうので、よろしくない。
		//fmt.Println("username:", r.Form["username"])
		//fmt.Println("password:", r.Form["password"])

		// 外部からのコードの対処法
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //サーバ側に出力されます。
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //クライアントに出力されます。
	}
}
func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //リクエストを受け取る方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("template/signup.gtpl")
		t.Execute(w, token)
	} else {
		//リクエストはログインデータです。ログインのロジックを実行して判断します。
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			//tokenの合法性を検証します。
			fmt.Println("token", token)
		} else {
			//tokenが存在しなければエラーを出します
			fmt.Println("token nothing")
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //サーバ側に出力します。
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //クライアントに出力します。
		if token != "" {
			template.HTMLEscape(w, []byte("\n"))
			template.HTMLEscape(w, []byte("your token is \""+token+"\""))
		}
	}
}
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //リクエストを受け取るメソッド
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("template/upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
func main() {
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
