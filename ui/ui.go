package ui

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
    "time"

	"../model"
	"../lnd"
)

type Config struct {
	Assets http.FileSystem
}

func Start(cfg Config, m *model.Model, l *lnd.LndLn, listener net.Listener) {

	server := &http.Server{
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16}

	http.Handle("/", indexHandler(m))
	http.Handle("/people", peopleHandler(m))
	http.Handle("/invoice", invoiceHandler(l))
	http.Handle("/check_invoice", checkInvoiceHandler(l))
	http.Handle("/js/", http.FileServer(cfg.Assets))

	go server.Serve(listener)
}

const (
	cdnReact           = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"
	cdnReactDom        = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"
	cdnBabelStandalone = "https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.24.0/babel.min.js"
	cdnAxios           = "https://cdnjs.cloudflare.com/ajax/libs/axios/0.16.1/axios.min.js"
)

const indexHTML = `
<!DOCTYPE HTML>
<html>
  <head>
    <meta charset="utf-8">
    <title>Simple Go Web App</title>
  </head>
  <body>
    <div id='root'></div>
    <script src="` + cdnReact + `"></script>
    <script src="` + cdnReactDom + `"></script>
    <script src="` + cdnBabelStandalone + `"></script>
    <script src="` + cdnAxios + `"></script>
    <script src="/js/app.jsx" type="text/babel"></script>
  </body>
</html>
`

func indexHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, indexHTML)
	})
}

func peopleHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		people, err := m.People()
		if err != nil {
			http.Error(w, "Could not get people.", http.StatusBadRequest)
			return
		}

		js, err := json.Marshal(people)
		if err != nil {
			http.Error(w, "Can't convert people into json.", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(js))
	})
}

func invoiceHandler(lnd *lnd.LndLn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, "Not GET.", r.Method)
            return
        }

        amt, err := strconv.Atoi(r.FormValue("amt"))
		if err != nil {
			amt = 100
		}
        int_amt := int64(amt)
        if int_amt == 0 {
            int_amt = 100
        }

        desc := r.FormValue("desc")
        if desc == "" {
            desc = "gowebapp payment"
        }

        invoice, err := lnd.MakeInvoice(int_amt, desc)
		if err != nil {
			http.Error(w, "Error generating invoice.", http.StatusBadRequest)
			return
		}

		js, err := json.Marshal(invoice)
		if err != nil {
			http.Error(w, "Error coverting invoice to json.", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(js))
	})
}

func checkInvoiceHandler(lnd *lnd.LndLn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, "Not GET.", r.Method)
            return
        }

        r_hash := r.FormValue("r_hash")
        if r_hash == "" {
			http.Error(w, "Empty payment hash.", http.StatusBadRequest)
			return
        }

        invoice, err := lnd.ViewInvoice(r_hash)
		if err != nil {
			http.Error(w, "Error checking invoice.", http.StatusBadRequest)
			return
		}

		js, err := json.Marshal(invoice)
		if err != nil {
			http.Error(w, "Error coverting invoice to json.", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(js))
	})
}
