package main // main.go project main.go

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Questlog/LogSocket/server"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Hello World!!!")

	//someOption := flag.String("name", "defaultValue", "Description")
	//fmt.Println(*someOption)

	flag.Parse()

	fw := NewFileWatcher(flag.Args()[0:])
	defer fw.close()

	ws := server.NewServer("/fail2orgel")
	go ws.Listen()

	go func() {
		for {
			select {
			case event := <-fw.eventCh:
				ws.SendAll(server.NewMessage(event.Name))
			}
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Filewatcher struct {
	paths   []string
	watcher *fsnotify.Watcher
	eventCh chan *fsnotify.Event
}

func NewFileWatcher(paths []string) *Filewatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	eventCh := make(chan *fsnotify.Event)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					eventCh <- &event
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	for _, path := range paths {
		fmt.Println("Watching:", path)
		err = watcher.Add(path)

		if err != nil {
			log.Println(path)
			log.Fatal(err)
		}
	}

	return &Filewatcher{paths, watcher, eventCh}
}

func (fw *Filewatcher) close() {
	fw.watcher.Close()
}
