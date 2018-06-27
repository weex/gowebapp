package ui

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
        "time"

	"github.com/weex/gowebapp/model"
	"github.com/weex/gowebapp/lnd"
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
	http.Handle("/invoice", invoiceHandler(l))
	http.Handle("/check_invoice", checkInvoiceHandler(l))
	http.Handle("/js/", http.FileServer(cfg.Assets))
	http.Handle("/basic.html", http.FileServer(cfg.Assets))
	http.Handle("/demo.html", http.FileServer(cfg.Assets))

	go server.Serve(listener)
}

const (
	cdnReact           = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"
	cdnReactDom        = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"
	cdnBabelStandalone = "https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.24.0/babel.min.js"
	cdnAxios           = "https://cdnjs.cloudflare.com/ajax/libs/axios/0.16.1/axios.min.js"
)

const indexHTML = `Please go to /demo.html`

func indexHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, indexHTML)
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

		js, err := json.MarshalIndent(invoice, "", "  ")
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

        payment_hash := r.FormValue("payment_hash")
        if payment_hash == "" {
			http.Error(w, "Empty payment hash.", http.StatusBadRequest)
			return
        }

        invoice, err := lnd.ViewInvoice(payment_hash)
		if err != nil {
			http.Error(w, "Error checking invoice.", http.StatusBadRequest)
			return
		}

		js, err := json.MarshalIndent(invoice, "", "  ")
		if err != nil {
			http.Error(w, "Error coverting invoice to json.", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(js))
	})
}
