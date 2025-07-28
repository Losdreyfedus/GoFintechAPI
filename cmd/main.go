package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "GoFintechAPI'ye hoş geldin!")
    })

    fmt.Println("Sunucu 8080 portunda başlatıldı...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Sunucu başlatılırken hata oluştu:", err)
    }
}