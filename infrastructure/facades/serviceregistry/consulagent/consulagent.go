package consulagent

import (
	"fmt"
	"log"
	"os"

	"github.com/riomhaire/consul"
	"github.com/riomhaire/keepsake/models"
)

type ConsulServiceRegistry struct {
	baseEndpoint   string
	healthEndpoint string
	id             string
	configuration  *models.Configuration
	consulClient   *consul.ConsulClient // This registers this service with consul - may extract this into a separate use case

}

func NewConsulServiceRegistry(configuration *models.Configuration, baseEndpoint, healthEndpoint string) *ConsulServiceRegistry {
	r := ConsulServiceRegistry{}
	r.baseEndpoint = baseEndpoint
	r.healthEndpoint = healthEndpoint
	r.configuration = configuration

	return &r
}

func (a *ConsulServiceRegistry) Register() error {
	// Register with consol (if required)
	hostname, _ := os.Hostname()
	id := fmt.Sprintf("%v-%v-%v", a.configuration.ApplicationName, hostname, a.configuration.Port)
	a.configuration.Consul.ConsulId = id // remember id for other system
	a.id = id                            // This is our safe copy

	a.consulClient, _ = consul.NewConsulClient(a.configuration.Consul.Host)
	health := fmt.Sprintf("http://%v:%v%v", hostname, a.configuration.Port, a.healthEndpoint)
	log.Println(fmt.Sprintf("Registering with Consul at %v with %v %v", a.configuration.Consul.Host, a.baseEndpoint, health))

	a.consulClient.PeriodicRegister(id, a.configuration.ApplicationName, hostname, a.configuration.Port, a.baseEndpoint, health, 15)
	return nil

}

/*

 */
func (a *ConsulServiceRegistry) Deregister() error {
	log.Println(fmt.Sprintf("De Registering %v with Consul at %v with %v ", a.configuration.Consul.ConsulId, a.configuration.Consul.Host, a.baseEndpoint))
	a.consulClient.DeRegister(a.id)
	return nil
}
