package main // main.go project main.go

import (
	"flag"
	"log"
	"net/http"

	"github.com/Questlog/LogSocket/server"
)

func main() {
	bindPath := flag.String("bindPath", "/logSocket", "Url to bind to")
	bindPort := flag.String("bindPort", ":8080", "Port to bind to")

	flag.Parse()

	fw := NewFileWatcher()
	defer fw.Close()
	fw.WatchFiles(flag.Args())

	ws := server.NewServer(*bindPath)
	go ws.Listen()

	go func() {
		for {
			select {
			case event := <-fw.eventCh:
				ws.SendAll(server.NewMessage(event.Name))
			}
		}
	}()

	//http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Fatal(http.ListenAndServe(*bindPort, nil))
}
