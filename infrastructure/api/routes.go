package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
	"github.com/zbindenren/negroni-prometheus"
)

func (this *RestAPI) defineRoutes() {
	rest := negroni.Classic()
	prometheusService := negroniprometheus.NewMiddleware(this.Configuration.ApplicationName, 60, 300, 1200, 3600)
	// if you want to use other buckets than the default (300, 1200, 5000) you can run:
	// m := negroniprometheus.NewMiddleware("serviceName", 400, 1600, 700)

	router := mux.NewRouter()
	rest.Use(prometheusService)

	// Add handlers
	router.HandleFunc("/api/v2/token/oauth/authorize", this.HandleAuthorize).Methods("POST")
	router.HandleFunc("/api/v2/token/oauth/verify", this.HandleVerify).Methods("GET")
	router.HandleFunc("/api/v2/token/jwt/sign", this.HandleSignJSONViaRSA).Methods("POST")
	router.HandleFunc("/api/v2/token/jwt/verify", this.HandleVerifyJWTViaRSA).Methods("GET")
	router.HandleFunc("/api/v2/token/health", this.HandleHealth).Methods("GET")
	router.HandleFunc("/health", this.HandleHealth).Methods("GET")
	router.HandleFunc("/api/v2/token/bvt", this.HandleBVT).Methods("GET")
	router.HandleFunc("/bvt", this.HandleBVT).Methods("GET")
	router.Handle("/metrics", prometheus.Handler())
	rest.UseFunc(this.AddWorkerHeader)  // Add which instance
	rest.UseFunc(this.AddWorkerVersion) // Which version
	rest.UseHandler(router)

	this.Negroni = rest
}
