package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
	negroniprometheus "github.com/zbindenren/negroni-prometheus"
)

func (a *RestAPI) defineRoutes() {
	negroni := negroni.Classic()
	prometheusService := negroniprometheus.NewMiddleware(a.Configuration.ApplicationName, 60, 300, 1200, 3600)
	// if you want to use other buckets than the default (300, 1200, 5000) you can run:
	// m := negroniprometheus.NewMiddleware("serviceName", 400, 1600, 700)

	mux := mux.NewRouter()
	negroni.Use(prometheusService)

	// Add handlers
	mux.HandleFunc("/api/v2/token/oauth/authorize", a.HandleAuthorize).Methods("POST")
	mux.HandleFunc("/api/v2/token/oauth/verify", a.HandleVerify).Methods("GET")
	mux.HandleFunc("/api/v2/token/jwt/sign", a.HandleSignJSONViaRSA).Methods("POST")
	mux.HandleFunc("/api/v2/token/jwt/verify", a.HandleVerifyJWTViaRSA).Methods("GET")
	mux.HandleFunc("/api/v2/token/health", a.HandleHealth).Methods("GET")
	mux.HandleFunc("/health", a.HandleHealth).Methods("GET")
	mux.Handle("/metrics", prometheus.Handler())
	negroni.UseFunc(a.AddWorkerHeader)  // Add which instance
	negroni.UseFunc(a.AddWorkerVersion) // Which version
	negroni.UseHandler(mux)

	a.Negroni = negroni
}
