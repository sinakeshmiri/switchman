package downstream

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type Healthcheck struct {
	Type string	`json:"type"`
	Uri  string	`json:"endpoint"`
}

type Adapter struct {
	IP             string
	HealthcheckUri Healthcheck
}

// NewAdapter creates a new Adapter
func NewAdapter(healthcheckUri Healthcheck, ip string) (*Adapter, error) {
	return &Adapter{HealthcheckUri: healthcheckUri, IP: ip}, nil
}

func (dsa Adapter) GetIP() string {
	return dsa.IP
}
func (dsa Adapter) CheckHealth() error {
	if dsa.HealthcheckUri.Type == "http" {
		req, err := http.NewRequest(http.MethodGet, dsa.HealthcheckUri.Uri, nil)
		if err != nil {
			return err
		}
		client := http.Client{
			Timeout: 300 * time.Second,
		}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return errors.New("Host is Down")
		}
		return nil
	} else if dsa.HealthcheckUri.Type == "tcp" {
		address := dsa.HealthcheckUri.Uri
		addr := strings.Split(address, ":")
		ip := addr[0]
		port := addr[1]

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
		if err != nil {
			return errors.New("Port not open:"+ ip+ port)
		}

		defer conn.Close()
		return nil
	}
	return errors.New("Healthcheck type not implemented")
}
