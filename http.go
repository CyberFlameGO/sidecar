package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/memberlist"
	"github.com/newrelic/bosun/services_state"
)


func makeHandler(fn func (http.ResponseWriter, *http.Request, *memberlist.Memberlist, *services_state.ServicesState), list *memberlist.Memberlist, state *services_state.ServicesState) http.HandlerFunc {
	return func(response http.ResponseWriter, req *http.Request) {
		fn(response, req, list, state)
	}
}

func servicesHandler(response http.ResponseWriter, req *http.Request, list *memberlist.Memberlist, state *services_state.ServicesState) {
	params := mux.Vars(req)

	defer req.Body.Close()

	if params["extension"] == ".json" {
		response.Header().Set("Content-Type", "application/json")
		jsonStr, _ := json.MarshalIndent(state.ByService(), "", "  ")
		response.Write(jsonStr)
		return
	}

	response.Header().Set("Content-Type", "text/html")
	response.Write(
		[]byte(`
 			<head>
 			<meta http-equiv="refresh" content="4">
 			</head> 
	    	<pre>` + state.Format(list) + "</pre>"))
}

func serveHttp(list *memberlist.Memberlist, state *services_state.ServicesState) {
	router := mux.NewRouter()

	router.HandleFunc(
		"/services{extension}", makeHandler(servicesHandler, list, state),
	).Methods("GET")

	router.HandleFunc(
		"/services", makeHandler(servicesHandler, list, state),
	).Methods("GET")

	http.Handle("/", router)

	err := http.ListenAndServe("0.0.0.0:7777", nil)
	exitWithError(err, "Can't start HTTP server")
}
