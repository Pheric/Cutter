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
	mux.Handle("/scripts/", http.StripPrefix("/scripts", (http.FileServer(http.FileSystem(http.Dir("site/scripts/"))))))
	mux.Handle("/client", http.HandlerFunc(ServeClientPage))
	mux.Handle("/clients", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		p := Pages["clients"]
		(&p).ServePage(w, r)
	}))
	mux.Handle("/clientspanel", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		p := Pages["clientsPanel"]
		(&p).ServePage(w, r)
	}))
	mux.Handle("/getCbSummary", http.HandlerFunc(HandleGetCBSummary))
	mux.Handle("/getCbLogJs", http.HandlerFunc(HandleGetCBLogJs))
	mux.Handle("/getCbLogHtml", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "site/cbLog.html")
	}))
	mux.Handle("/cbEdit", http.HandlerFunc(HandleCBEdit))
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
	index.InitAutoRecache(time.Minute * 30)

	clients := InitClientsPage()
	Pages[clients.Path] = clients
	clients.InitAutoRecache(time.Hour * 1)

	cPanel := InitClientsPanel()
	Pages[cPanel.Path] = cPanel
	cPanel.InitAutoRecache(time.Minute * 30)
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
			ServeCode(w, 401, "Unauthorized", "Get back where you belong!", 5)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
