package main

import (
    // "fmt"
    "net/http"
    "log"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("Hola senor\n"))
    log.Print("== we have a visitour .. :\n" , req)
}

func main() {
    http.HandleFunc("/holasenor", HelloServer)
    err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
