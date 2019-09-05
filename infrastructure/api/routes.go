package api

import (
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	negroniprometheus "github.com/zbindenren/negroni-prometheus"
)

func (this *RestAPI) defineRoutes() {
	// Includes some default middlewares to all routes
	rest := negroni.New()
	rest.Use(negroni.NewRecovery())

	// add logrus
	rest.Use(negronilogrus.NewMiddlewareFromLogger(log.StandardLogger(), this.Configuration.ApplicationName))

	prometheusService := negroniprometheus.NewMiddleware(this.Configuration.ApplicationName, 60, 300, 1200, 3600)
	// if you want to use other buckets than the default (300, 1200, 5000) you can run:
	// m := negroniprometheus.NewMiddleware("serviceName", 400, 1600, 700)

	router := mux.NewRouter()
	rest.Use(prometheusService)

	// Add handlers
	router.HandleFunc("/api/v2/token/oauth/authorize", this.HandleAuthorize).Methods("POST")
	router.HandleFunc("/api/v2/token/oauth/verify", this.HandleVerify).Methods("POST")
	router.HandleFunc("/api/v2/token/jwt/sign", this.HandleSignJSONViaRSA).Methods("POST")
	router.HandleFunc("/api/v2/token/jwt/verify", this.HandleVerifyJWTViaRSA).Methods("POST")
	router.HandleFunc("/api/v2/token/health", this.HandleHealth).Methods("GET")
	router.HandleFunc("/health", this.HandleHealth).Methods("GET")
	router.HandleFunc("/api/v2/token/.wellknown/jwks.json", this.HandleJWKPublicGet).Methods("GET")
	router.HandleFunc("/.wellknown/jwks.json", this.HandleJWKPublicGet).Methods("GET")
	router.HandleFunc("/api/v2/token/bvt", this.HandleBVT).Methods("GET")
	router.HandleFunc("/bvt", this.HandleBVT).Methods("GET")
	router.Handle("/metrics", prometheus.Handler())
	rest.UseFunc(this.AddWorkerHeader)  // Add which instance
	rest.UseFunc(this.AddWorkerVersion) // Which version
	rest.UseHandler(router)

	this.Negroni = rest
}
