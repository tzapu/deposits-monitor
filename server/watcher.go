package server

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/time/rate"
)

type Watcher struct {
	hub *Hub
}

func (w *Watcher) Watch(paths ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, p := range paths {
		err = watcher.Add(p)
		if err != nil {
			return err
		}
	}

	r := rate.Every(time.Second * 2)
	limiter := rate.NewLimiter(r, 1)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("watcher events not ok")
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				if limiter.Allow() {
					log.Debugf("modified file %s", event.Name)
					w.hub.Broadcast([]byte(`{"type":"build_complete"}`))
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return fmt.Errorf("watcher err not ok")
			}
			log.Println("error:", err)
		}
	}
}

func NewWatcher(hub *Hub) *Watcher {
	return &Watcher{
		hub: hub,
	}
}
