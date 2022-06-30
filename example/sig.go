package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	serve := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		fmt.Println(serve.ListenAndServe())
	}()
	fmt.Printf("%+v\n", "running .... ")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	c := <-sigs
	fmt.Println("Got signal:", c.String())

}
