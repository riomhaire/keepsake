package api

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/riomhaire/keepsake/infrastructure/facades/serviceregistry"
	"github.com/riomhaire/keepsake/models"
	"github.com/urfave/negroni"
)

var bearerPrefix = "bearer "

type RestAPI struct {
	Negroni                 *negroni.Negroni
	Configuration           *models.Configuration
	TokenEncoderDecoder     models.TokenEncoderDecoder
	JWTEncoderDecoder       models.JWTEncoderDecoder
	ClientStore             models.StorageInteractor
	ExternalServiceRegistry serviceregistry.ServiceRegistry
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

func (this *RestAPI) Start() {
	if this.Configuration.Port == 0 {
		this.Configuration.Port = 10101
	}
	this.ExternalServiceRegistry.Register()
	this.Negroni.Run(fmt.Sprintf(":%d", this.Configuration.Port))
}

func (this *RestAPI) Stop() {
	log.Println("Shutting Down REST API")
	this.ExternalServiceRegistry.Deregister()
}
