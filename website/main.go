package hello

import (
    "net/http"
)

func init() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
}
