package server

import (
	"html/template"
	"bytes"
	"fmt"
	"cutter/data"
)

func InitClientsPage() Page {
	return Page {
		Path: "clients",
		RegenFunc: genClientsPg,
		PostFunc: nil,
	}
}

var clientsTemplate *template.Template
func genClientsPg(initial bool) ([]byte, error) {
	var b bytes.Buffer
	var err error

	if initial {
		clientsTemplate, err = template.ParseFiles("site/clients.html")
		if err != nil {
			return nil, fmt.Errorf("error generating clients.html: %v\n", err)
		}
	}

	clients, err := data.LoadAllClients()
	if err != nil {
		// RIP
		return nil, fmt.Errorf("error loading clients list: %v", err)
	}

	err = clientsTemplate.Execute(&b, clients)
	return b.Bytes(), err
}