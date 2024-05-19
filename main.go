package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/")

	if i <= 0 {
		return p[1:], "/"
	}

	return p[1:i], p[i:]
}

type Route struct {
	Logger  bool
	Tester  bool
	Handler http.Handler
}

type App struct {
	User *Route
}

type User struct{}

func (h *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var next *Route
	var head string

	head, r.URL.Path = shiftPath(r.URL.Path)

	fmt.Println(head, r.URL.Path)

	if len(head) == 0 {
		next = &Route{
			Logger: true,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("home page"))
			}),
		}
	} else if head == "user" {
		var i interface{} = User{}

		next = &Route{
			Logger:  true,
			Tester:  true,
			Handler: i.(http.Handler),
		}
	} else {
		http.Error(w, "not found", http.StatusNotFound)

		return
	}

	if next.Logger {
		next.Handler = h.log(next.Handler)
	}
}

func (h *App) log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(time.Since(start), r.Method, r.URL.Path)
	})
}

type key int

const (
	ctxTestKey key = 1
	ctxUserID      = 2
)
