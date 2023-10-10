package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	// application
	"github.com/sinakeshmiri/switchman/internal/application/api"

	// adapters
	"github.com/sinakeshmiri/switchman/internal/adapters/framework/left/auto"
	"github.com/sinakeshmiri/switchman/internal/adapters/framework/left/http"
	"github.com/sinakeshmiri/switchman/internal/adapters/framework/right/downstream"
	"github.com/sinakeshmiri/switchman/internal/adapters/framework/right/nsprovider/cf"

	//ports
	"github.com/sinakeshmiri/switchman/internal/ports"
)

type NSProvider struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type Upstream struct {
	Primary     bool                   `json:"primary"`
	IP          string                 `json:"ip"`
	HealthCheck downstream.Healthcheck `json:"healthcheck"`
}

type Config struct {
	Intrerface string     `json:"interface"`
	MainDomain string     `json:"maindomin"`
	NSProvider NSProvider `json:"nsprovider"`
	Upstreams  []Upstream `json:"upstreams"`
}

func main() {
	file, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
		return
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	if config.NSProvider.Type != "cloudflare" {
		log.Fatal("Error NSprovider is  not impelmented")
	}
	for _, u := range config.Upstreams {
		if u.HealthCheck.Type != "http" && u.HealthCheck.Type != "tcp" {
			log.Fatal("healthceck type is  not impelmented")

		}
	}
	cfadaptor, err := cf.NewAdapter(config.NSProvider.Key)
	if err != nil {
		log.Fatal(err)
	}
	var downstreams []ports.DownStreamPort
	for _, s := range config.Upstreams {
		downstreams = append(downstreams, downstream.Adapter{
			Primary:        s.Primary,
			IP:             s.IP,
			HealthcheckUri: s.HealthCheck,
		})
	}

	// NOTE: The application's right side port for driven
	// adapters, in this case, a db adapter.
	// Therefore the type for the dbAdapter parameter
	// that is to be injected into the NewApplication will
	// be of type DbPort
	applicationAPI := api.NewApplication(config.MainDomain, cfadaptor, downstreams)

	// NOTE: We use dependency injection to give the grpc
	// adapter access to the application, therefore
	// the location of the port is inverted. That is
	// the grpc adapter accesses the hexagon's driving port at the
	// application boundary via dependency injection,
	// therefore the type for the applicaitonAPI parameter
	// that is to be injected into the gRPC adapter will
	// be of type APIPort which is our hexagons left side
	// port for driving adapters
	if config.Intrerface == "http" {
		httpAdapter := http.NewAdapter(applicationAPI)
		httpAdapter.Check()
	} else if config.Intrerface == "auto" {
		autoAdapter := auto.NewAdapter(applicationAPI)
		autoAdapter.Check()
	}
}
