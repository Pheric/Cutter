package server

import (
	"net/http"
	"time"
	"log"
	"fmt"
)

type Page struct {
	Path      string
	RegenFunc func(initial bool) ([]byte, error)
	PostFunc  http.HandlerFunc

	autoRecache bool
}

func (p *Page) InitAutoRecache(interval time.Duration) {
	if p.autoRecache {
		return
	}

	t := time.NewTicker(interval)
	go func() {
		for start := range t.C {
		log.Printf("Re-caching page \"%s\"...\n", p.Path)
		err := p.Recache(false)
		if err != nil {
			log.Printf("Encountered an error in auto-recache loop: %v\n", err)
		}
		log.Printf("Re-caching finished for page \"%s\" in %gms\n", p.Path, time.Since(start).Seconds()/1000)
	}
}()
}

func (p *Page) ServePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		if r.Method != "POST" || p.PostFunc == nil {
			// Todo 405
			http.Error(w, "method unsupported", http.StatusMethodNotAllowed)
			return
		}
		p.PostFunc(w, r)
		return
	}

	data, err := GlobalCache.Get(p.Path)
	if err != nil {
		err := p.Recache(true)
		if err != nil {
			log.Printf("Error while attempting live page load: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		p.ServePage(w, r)
		return
	}

	w.Write(data)
}

func (p *Page) Recache(initial bool) error {
	data, err := p.RegenFunc(initial)
	if err != nil {
		return fmt.Errorf("error re-caching page \"%s\": %v", p.Path, err)
	}

	GlobalCache.Save(p.Path, data)
	return nil
}