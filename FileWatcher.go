package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Filewatcher struct {
	watcher *fsnotify.Watcher
	eventCh chan *fsnotify.Event
}

func NewFileWatcher() *Filewatcher {
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

	return &Filewatcher{watcher, eventCh}
}

func (fw *Filewatcher) Close() {
	fw.watcher.Close()
}

func (fw *Filewatcher) WatchFiles(paths []string) {
	for _, path := range paths {
		log.Println("Watching:", path)
		err := fw.watcher.Add(path)

		if err != nil {
			log.Println(path)
			log.Fatal(err)
		}
	}
}
