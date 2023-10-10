package api

import (
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
	found :=false
	ips, err := apia.nsprovider.GetRecords(apia.maindomin)
	if err != nil {
		return err
	}
	
	for _, ds := range apia.downstreams {
		found =false
		for _, item := range ips {
			if item == ds.GetIP() {
				found = true
				err = ds.CheckHealth()
				if err != nil {
					log.Println(err)
					err = apia.nsprovider.DelRecord(apia.maindomin, ds.GetIP())
					if err != nil {
						return err
					}
				}
				break // Exit the loop early since we found the target
			}
		}
		if !found {
			err = ds.CheckHealth()
			if err != nil {
				log.Println(err)
			} else {
				err = apia.nsprovider.SetRecord(apia.maindomin, ds.GetIP())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
