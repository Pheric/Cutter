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
		// TODO 405
		http.Error(w, "method unsupported", http.StatusMethodNotAllowed)
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
		// TODO 500
		log.Printf("Error while loading client.html: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func servePostRequest(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("cid")
	log.Println(r.URL.Query())
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
			// TODO 500
			log.Printf("Error saving client: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	} else if quote := r.FormValue("quote"); quote != "" { // job information
		parsedQuote, err := strconv.ParseFloat(quote, 32)
		if err != nil {
			// TODO 400
			http.Error(w, "Quote value is NaN", http.StatusBadRequest)
			return
		}
		parsedTtc, err := strconv.Atoi(r.FormValue("ttc"))
		if err != nil {
			// TODO 400
			http.Error(w, "TTC value is NaN", http.StatusBadRequest)
			return
		}
		parsedPeriod, err := strconv.Atoi(r.FormValue("period"))
		if err != nil {
			// TODO 400
			http.Error(w, "Period value is NaN...\nStop messing with my website.\nBut if you do find anything interesting,\nEmail me. dimphoton@outlook.com. Thanks.", http.StatusBadRequest)
			return
		}
		client.Quote = float32(parsedQuote)
		client.Ttc = parsedTtc
		client.Period = parsedPeriod

		err = client.SaveShallow()
		if err != nil {
			// TODO 500
			log.Printf("Error saving client: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	} else if date := r.FormValue("date"); date != "" { // balance information (payment)
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			// TODO 400
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		parsedAmount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
		if err != nil {
			// TODO 400
			http.Error(w, "Amount value is NaN", http.StatusBadRequest)
			return
		}

		client.Balance += float32(parsedAmount)
		err = client.SaveShallow()
		if err != nil {
			// TODO 500
			log.Printf("Error saving client: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		err = data.Payment {
			Target: client.Uuid,
			Amount: float32(parsedAmount),
			Date: parsedDate,
		}.SavePayment()
		if err != nil {
			// TODO 500
			log.Printf("Error saving payment of %f to %s: %v\n", parsedAmount, client.Uuid, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		// What did they do...?
		// TODO 400
		http.Error(w, "Invalid form section...\nStop messing with my website.\nBut if you do find anything interesting,\nEmail me. dimphoton@outlook.com. Thanks.", http.StatusBadRequest)
		return
	}
}