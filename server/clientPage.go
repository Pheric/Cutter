package server

import (
	"net/http"
	"cutter/data"
	"log"
	"html/template"
	"strconv"
	"time"
)

var clientTemplate *template.Template
var initial = true
func ServeClientPage(w http.ResponseWriter, r *http.Request) {
	if initial {
		var err error
		clientTemplate, err = template.ParseFiles("site/client.html")
		if err != nil {
			log.Printf("Error loading template for client.html: %v\n", err)
			return
		}
		initial = false
	}

	switch r.Method {
	case "GET":
		serveGetRequest(w, r)
	break
	case "POST":
		servePostRequest(w, r)
	break
	default:
		ServeCode(w, 405, "Method Not Allowed", "Your request is not acceptable.", 10)
	}
}

func serveGetRequest(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("cid")
	client, err := data.LoadClientWithUuid(clientId)
	if err != nil {
		log.Printf("Error loading client with uuid %s: %v\n", clientId, err)
		// Keep going..
	}

	err = clientTemplate.Execute(w, client)
	if err != nil {
		log.Printf("Error while loading client.html: %v\n", err)
		ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
		return
	}
}

func servePostRequest(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("cid")
	client, err := data.LoadClientWithUuid(clientId)
	if err != nil {
		log.Printf("Error loading client with uuid %s: %v\n", clientId, err)
		// Keep going..
	}

	// Which section of the page is being submitted? Each of these fields is present in only one section...
	if last := r.FormValue("last"); last != "" { // personal information
		client.Last = last
		client.First = r.FormValue("first")
		client.Phone = r.FormValue("phone")
		client.Address = r.FormValue("address")

		err := client.SaveShallow()
		if err != nil {
			log.Printf("Error saving client: %v\n", err)
			ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
			return
		}
	} else if quote := r.FormValue("quote"); quote != "" { // job information
		if clientId == "" {
			return
		}

		parsedQuote, err := strconv.ParseFloat(quote, 32)
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
			ServeCode(w, 400, "Improper Submission", "Your period value is not acceptable.", 10)
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
	} else if date := r.FormValue("date"); date != "" { // balance information (payment)
		if clientId == "" {
			return
		}

		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			ServeCode(w, 400, "Improper Submission", "Your date value is not acceptable.", 10)
			return
		}
		parsedAmount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
		if err != nil {
			ServeCode(w, 400, "Improper Submission", "Your amount value is not a number.", 10)
			return
		}

		client.Balance += float32(parsedAmount)
		err = client.SaveShallow()
		if err != nil {
			log.Printf("Error saving client: %v\n", err)
			ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
			return
		}

		err = data.Payment {
			Target: client.Uuid,
			Amount: float32(parsedAmount),
			Date: parsedDate,
		}.SavePayment()
		if err != nil {
			log.Printf("Error saving payment of %f to %s: %v\n", parsedAmount, client.Uuid, err)
			ServeCode(w, 500, "Internal Server Error", "Please retry your last action.", 10)
			return
		}
	} else {
		// What did they do...?
		ServeCode(w, 400, "Improper Submission", "You're out of line. Please either stop messing with my site or be responsible and report bugs. Thanks. dimphoton@outlook.com", 10)
		return
	}

	clientsPg := Pages["clients"]
	err = (&clientsPg).Recache(false)
	if err != nil {
		log.Printf("Error re-caching clients list after modification of client: %v\n", err)
		// Can't really do anything about it at this point :(
	}
}