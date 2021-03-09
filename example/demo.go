package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/18211167516/gosession"
	"github.com/18211167516/gosession/store"
	"github.com/gogf/gf/util/gconv"
)

var GloabSession *gosession.Manager

func init() {
	var err error
	store.RegisterRedis("192.168.99.100:6379", "", 10, 10)
	GloabSession, err = gosession.NewDefaultIDManager("redis", "gosession", 84600)
	if err != nil {
		log.Panic("Session Manager init ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.PATH=%q\n Header=%q\n", r.URL.Path, r.Header)
}

func login(w http.ResponseWriter, r *http.Request) {
	s := GloabSession.GetSession(w, r)
	v, _ := s.Get("isLogin")
	if gconv.Int(v) == 1 {
		fmt.Fprintf(w, "已登录")
	} else {
		s.Set("isLogin", 1)
		fmt.Fprintf(w, "正在登录")
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	s := GloabSession.GetSession(w, r)
	s.Delete("isLogin")
	s.DestroySession(w, r)
	fmt.Fprintf(w, "已退出登录")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}
