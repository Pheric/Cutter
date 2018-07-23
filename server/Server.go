package server

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"
	"sync"
)

var Pages = make(map[string]Page)

func Init(port int) {
	log.Println("Begin registration and caching of all pages...")

	registerPages() // Add pages to map
	var wg sync.WaitGroup
	wg.Add(len(Pages))
	for _, p := range Pages {
		go func(p Page) {
			err := p.Recache(true)
			if err != nil {
				log.Printf("Initial cache failed: %v\n", err)
			}
			wg.Done()
		}(p)
	}
	wg.Wait()

	log.Println("Finished")

	mux := http.NewServeMux()
	mux.Handle("/styling/", http.StripPrefix("/styling", (http.FileServer(http.FileSystem(http.Dir("site/styling/"))))))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		p := Pages["index"]
		(&p).ServePage(w, r)
	}))

	log.Printf("Initializing HTTP server on port %d\n", port)
	go http.ListenAndServe(fmt.Sprintf("127.1:%d", port), mux)
	log.Println("Now listening for HTTP traffic.")

	select {
	// Lock
	}
}

func registerPages() {
	index := InitIndex()
	Pages[index.Path] = index
	index.InitAutoRecache(time.Second * 30)


}

// FIXME:
/*
I'm not convinced `checkRootPath` works correctly. To disable directory listings, I would implement an `http.FileSystem` that looks at `IsDir()` on the files it opens.
Or that simply returns an error for `Readdir` on the `http.File`s it returns.
Decorate the `http.FileSystem` implemented by `http.Dir` with the behavior you desire, and you'll get something that you can re-use and reason about easier.
 */

// Returns 401 if client requests a directory, thereby avoiding a directory listing.
func checkRootPath(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Is URL pointing to a directory?
		if b := path.Base(r.URL.Path); b == "." || b == "/" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
