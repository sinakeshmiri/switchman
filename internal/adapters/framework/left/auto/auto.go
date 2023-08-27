package auto

import (
	"log"
	"time"

	"github.com/sinakeshmiri/switchman/internal/ports"
)

type Adapter struct {
	api ports.APIPort
}

// NewAdapter creates a new Adapter
func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

func (autoa Adapter) checkHealth() error {
	err := autoa.api.Check()
	if err != nil {
		return err
	}
	return nil
}

func (autoa Adapter) Check() {
	for {
		err := autoa.checkHealth()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Minute)
	}
}
