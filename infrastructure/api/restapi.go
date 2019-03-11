package api

import (
	"fmt"

	"github.com/riomhaire/keepsake/models"
	"github.com/urfave/negroni"
)

var bearerPrefix = "bearer "

type RestAPI struct {
	Negroni             *negroni.Negroni
	Configuration       *models.Configuration
	TokenEncoderDecoder models.TokenEncoderDecoder
	JWTEncoderDecoder   models.JWTEncoderDecoder
	ClientStore         models.StorageInteractor
}

func NewRestAPI(configuration *models.Configuration, tokenizer models.TokenEncoderDecoder, jwtizer models.JWTEncoderDecoder, storageInteractor models.StorageInteractor) RestAPI {
	api := RestAPI{}
	api.Configuration = configuration
	api.TokenEncoderDecoder = tokenizer
	api.JWTEncoderDecoder = jwtizer
	api.ClientStore = storageInteractor
	api.defineRoutes()

	return api
}

func (a *RestAPI) Start() {
	if a.Configuration.Port == 0 {
		a.Configuration.Port = 10101
	}
	a.Negroni.Run(fmt.Sprintf(":%d", a.Configuration.Port))
}
