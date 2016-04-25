package webinterface

import (
	"html/template"
	"log"
	"net/http"

	"github.com/christer79/ccTray2Slack/cctray"
	"github.com/gorilla/mux"
)

//Statuses holds map of identifier to Projects mapping
type Statuses map[string][]cctray.Project

type WebInterface struct {
	ChStatus chan Statuses
	statuses Statuses
}

type statusPage struct {
	Projects     Statuses
	StatusFilter string
}

func (web *WebInterface) statusHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	status := statusPage{web.statuses, r.Form.Get("status")}
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Id: %s Status: %s ", id, status.StatusFilter)
	t, err := template.ParseFiles("html/status.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, status)
}

//Start start a http server to expose configuration adn status
func (w *WebInterface) Start(port string) {
	w.ChStatus = make(chan Statuses)
	log.Println("Starting web interface")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/status/", w.statusHandler)
	go func() {
		for {
			select {
			case s := <-w.ChStatus:
				log.Println("New statuses recieved")
				w.statuses = s
			}
		}
	}()
	log.Fatal(http.ListenAndServe(":"+port, router))
}
