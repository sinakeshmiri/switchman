package api

import (
	"errors"
	"log"

	"github.com/sinakeshmiri/switchman/internal/ports"
)

// Application implements the APIPort interface
type Application struct {
	downstreams []ports.DownStreamPort
	nsprovider  ports.NsproviderPort
	maindomin   string
}

// NewApplication creates a new Application
func NewApplication(
	maindomin string, nsprovider ports.NsproviderPort, downstreams []ports.DownStreamPort) *Application {
	return &Application{
		maindomin:   maindomin,
		nsprovider:  nsprovider,
		downstreams: downstreams,
	}
}

// update record will check all available services and find out wich service is up and running
// then updates the record if the preferd service is down
func (apia Application) Check() error {
	var ok bool
	ip, err := apia.nsprovider.GetRecord(apia.maindomin)
	if err != nil {
		return err
	}
	for _, s := range apia.downstreams {
		if s.GetIP() == ip {
			err := s.CheckHealth()
			if err == nil {
				ok = true
				break
			} else {
				ok = false
			}
		}
	}
	if !ok {
		for _, s := range apia.downstreams {
			err := s.CheckHealth()
			if err == nil {
				err = apia.nsprovider.SetRecord(apia.maindomin, s.GetIP())
				if err != nil {
					log.Println(err)
				} else {
					ok = true
					break
				}
			}
		}
		if !ok {
			log.Println("No Host is alive")
			return errors.New("No host is alvie")

		}

	}
	return nil
}
