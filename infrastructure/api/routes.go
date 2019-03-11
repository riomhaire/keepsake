package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func (a *RestAPI) defineRoutes() {
	mux := mux.NewRouter()
	negroni := negroni.Classic()

	// Add handlers
	mux.HandleFunc("/api/v2/oauth/authorize", a.HandleAuthorize).Methods("POST")
	mux.HandleFunc("/api/v2/oauth/verify", a.HandleVerify).Methods("GET")
	mux.HandleFunc("/api/v2/jwt/sign", a.HandleSignJSONViaRSA).Methods("POST")
	mux.HandleFunc("/api/v2/jwt/verify", a.HandleVerifyJWTViaRSA).Methods("GET")

	handler := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler(mux) // Add coors
	negroni.UseHandler(handler)

	a.Negroni = negroni
}
