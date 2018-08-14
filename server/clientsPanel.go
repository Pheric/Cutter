package server

import (
	"bytes"
	"cutter/data"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func InitClientsPanel() Page {
	initCBPages()
	return Page{
		Path:      "clientsPanel",
		RegenFunc: genClientsPanel,
		PostFunc:  nil,
	}
}

var clientsPanelTemplate *template.Template

func genClientsPanel(initial bool) ([]byte, error) {
	var b bytes.Buffer
	var err error

	if initial {
		clientsPanelTemplate, err = template.ParseFiles("site/clientsPanel.html") // FIXME
		if err != nil {
			return nil, fmt.Errorf("error generating clientsPanel.html: %v\n", err) // FIXME
		}
	}

	clients, err := data.LoadAllClients()
	if err != nil {
		log.Printf("Error loading clients: %v\n", err)
		// Do not return
	}

	err = clientsPanelTemplate.Execute(&b, clients)
	return b.Bytes(), err
}

/* = = = = = Client Box = = = = = = */

var cbSummaryTemplate, cbLogTemplate, cbEditTemplate *template.Template

func initCBPages() {
	var err error
	cbSummaryTemplate, err = template.ParseFiles("site/cbSummary.html")
	if err != nil {
		log.Printf("error generating cbSummary.html: %v\n", err)
	}

	cbLogTemplate, err = template.ParseFiles("site/scripts/cbLog.js")
	if err != nil {
		log.Printf("error generating cbLog.js: %v\n", err)
	}

	cbEditTemplate, err = template.ParseFiles("site/cbEdit.html")
	if err != nil {
		log.Printf("error generating cbEdit.html: %v\n", err)
	}
}

func HandleGetCBSummary(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("cid")
	client, err := data.LoadClientWithUuid(clientId)
	if err != nil {
		log.Printf("Error loading client with uuid %s: %v\n", clientId, err)
		// Keep going..
	}

	err = cbSummaryTemplate.Execute(w, client)
	if err != nil {
		log.Printf("Error serving CB Summary: %v\n", err)
		ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
	}
}

func HandleGetCBLogJs(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("cid")
	client, err := data.LoadClientWithUuid(clientId)
	if err != nil {
		log.Printf("Error loading client with uuid %s: %v\n", clientId, err)
		// Keep going..
	}

	err = cbLogTemplate.Execute(w, client)
	if err != nil {
		log.Printf("Error serving CB Log (js): %v\n", err)
		ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
	}
}

func HandleCBEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		clientId := r.URL.Query().Get("cid")
		client, err := data.LoadClientWithUuid(clientId)
		if err != nil {
			log.Printf("Error loading client with uuid %s: %v\n", clientId, err)
			// Keep going..
		}

		err = cbEditTemplate.Execute(w, client)
		if err != nil {
			log.Printf("Error serving CB Edit: %v\n", err)
			ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
		}
	} else if r.Method == "POST" {
		handleCBEditPostRequest(w, r)
	} else {
		ServeCode(w, 405, "Method Not Allowed", "Your request is not acceptable.", 10)
	}
}

func handleCBEditPostRequest(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("cid")
	client, err := data.LoadClientWithUuid(clientId)
	if err != nil {
		log.Printf("Error loading client with uuid %s: %v\n", clientId, err)
		// Keep going..
	}

	client.Last = r.FormValue("last")
	client.First = r.FormValue("first")
	client.Phone = r.FormValue("phone")
	client.Address = r.FormValue("address")

	parsedQuote, err := strconv.ParseFloat(r.FormValue("quote"), 32)
	if err != nil {
		ServeCode(w, 400, "Improper Submission", "Your quote value is not a number.", 10)
		return
	}
	parsedTtc, err := strconv.Atoi(r.FormValue("ttc"))
	if err != nil {
		ServeCode(w, 400, "Improper Submission", "Your ttc value is not a number.", 10)
		return
	}
	parsedPeriod, err := strconv.Atoi(r.FormValue("period"))
	if err != nil {
		ServeCode(w, 400, "Improper Submission", "Your period value is not a number.\nStop messing with my website.\nBut if you do find anything interesting,\nEmail me. dimphoton@outlook.com. Thanks.", 10)
		return
	}
	client.Quote = float32(parsedQuote)
	client.Ttc = parsedTtc
	client.Period = parsedPeriod

	err = client.SaveShallow()
	if err != nil {
		log.Printf("Error saving client: %v\n", err)
		ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
		return
	}
	ServeCode(w, 200, "Success", "Operation successful.", 3)

	clientsPg := Pages["clientsPanel"]
	err = (&clientsPg).Recache(false)
	if err != nil {
		log.Printf("Error re-caching clients list after modification of client: %v\n", err)
		ServeCode(w, 500, "Internal Server Error", "Re-caching of your previous page has failed. The data may not show up in its latest version.", 3)
		// Can't really do anything about it at this point :(
	}
}
