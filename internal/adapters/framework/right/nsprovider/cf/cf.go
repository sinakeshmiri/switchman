package cf

import (
	"context"
	"errors"
	"log"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

type Adapter struct {
	api *cloudflare.API
}

// NewAdapter creates a new Adapter
func NewAdapter(key string) (*Adapter, error) {
	api, err := cloudflare.NewWithAPIToken(key)
	if err != nil {
		log.Println(err)
	}

	return &Adapter{api: api}, nil
}

func (cfa Adapter) extractZone(domin string) (string, error) {
	parts := strings.Split(domin, ".")
	n := len(parts)
	if n < 2 {
		return "", errors.New("wrong domin format")
	}
	tld := parts[n-1]

	sld := parts[n-2]

	zoneName := sld + "." + tld

	id, err := cfa.api.ZoneIDByName(zoneName) // Assuming zoneName exists in your Cloudflare account already
	if err != nil {
		log.Println(err)
		return "", err
	}
	return id, nil
}

func (cfa Adapter) GetRecord(domin string) (string, error) {
	id, err := cfa.extractZone(domin)
	if err != nil {
		log.Println(err)
		return "", err
	}
	records, _, err := cfa.api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(id), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Println(err)
		return "", err
	}
	r := cloudflare.DNSRecord{}
	for _, rec := range records {
		if rec.Name == domin {
			r = rec
		}
	}
	return r.Content, nil

}

func (cfa Adapter) SetRecord(domin string, value string) error {
	id, err := cfa.extractZone(domin)
	if err != nil {
		log.Println(err)
		return err
	}
	records, _, err := cfa.api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(id), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Println(err)
		return err
	}
	r := cloudflare.DNSRecord{}
	for _, rec := range records {
		if rec.Name == domin {
			r = rec
		}
	}

	_, err = cfa.api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(id), cloudflare.UpdateDNSRecordParams{
		ID:       r.ID,
		Content:  value,
		Type:     r.Type,
		Name:     r.Name,
		Data:     r.Data,
		Priority: r.Priority,
		TTL:      r.TTL,
		Proxied:  r.Proxied,
		Comment:  r.Comment,
		Tags:     r.Tags,
	})
	if err != nil {
		return err
	}
	return nil
}
