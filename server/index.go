package server

import (
	"html/template"
	"fmt"
	"bytes"
	"log"
	"cutter/data"
)

func InitIndex() Page {
	return Page {
		Path: "index",
		RegenFunc: gen,
		PostFunc: nil,
	}
}

var indexTemplate *template.Template
func gen(initial bool) ([]byte, error) {
	var b bytes.Buffer
	var err error

	if initial {
		indexTemplate, err = template.ParseFiles("site/index.html")
		if err != nil {
			return nil, fmt.Errorf("error generating index.html: %v\n", err)
		}
	}

	employees, err := data.LoadAllEmployees()
	if err != nil {
		log.Printf("Error loading employees: %v\n", err)
		// Do not return
	}

	clients, err := data.LoadAllClients()
	if err != nil {
		log.Printf("Error loading clients: %v\n", err)
		// Do not return
	}

	err = indexTemplate.Execute(&b, struct {
		Employees []data.Employee
		Clients []data.Client
	}{
		Employees: employees,
		Clients: clients,
	})

	return b.Bytes(), err
}