package server

import (
	"net/http"
	"html/template"
	"log"
)

var codePageTemplate *template.Template
func ServeCode(w http.ResponseWriter, code int, description, resolution string, time int) {
	if codePageTemplate == nil {
		var err error
		codePageTemplate, err = template.ParseFiles("site/responseCode.html")
		if err != nil {
			log.Printf("error generating responseCode.html: %v\n", err)
		}
	}

	w.WriteHeader(code)
	err := codePageTemplate.Execute(w, struct {
		Code       int
		Desc       string
		Resolution string
		Delay      int
	}{
		Code:       code,
		Desc:       description,
		Resolution: resolution,
		Delay:      time,
	})

	if err != nil {
		// Um this is a bad spot to be in...
		log.Printf("Error executing the code page template: %v\n", err)

		// Not even going to bother checking for an error here. What a mess.
		w.Write([]byte("A 500 series error has occurred.\n" +
			"Please retry your last operation. If it continues to fail,\n" +
			"manual intervention may be necessary. In this case, please contact me (Eric).\n" +
			"Sorry about this."))
	}
}