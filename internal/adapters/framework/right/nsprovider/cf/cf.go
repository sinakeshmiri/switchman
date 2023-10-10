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

func (cfa Adapter) extractZone(domin string) (string,string, error) {
	parts := strings.Split(domin, ".")
	n := len(parts)
	if n < 2 {
		return "", "",errors.New("wrong domin format")
	}
	tld := parts[n-1]

	sld := parts[n-2]

	zoneName := sld + "." + tld
	log.Println(zoneName)
	id, err := cfa.api.ZoneIDByName(zoneName) // Assuming zoneName exists in your Cloudflare account already
	if err != nil {
		log.Println(err)
		return "","", err
	}
	return zoneName,id, nil
}

func (cfa Adapter) GetRecords(domin string) ([]string, error) {
	_,id, err := cfa.extractZone(domin)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	records, _, err := cfa.api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(id), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	r := []string{}
	for _, rec := range records {
		if rec.Name == domin {
			r = append(r, rec.Content)
		}
	}
	return r, nil

}

func (cfa Adapter) DelRecord(domin string,value string) ( error) {
	_,zid, err := cfa.extractZone(domin)
	if err != nil {
		log.Println(err)
		return  err
	}
	records, _, err := cfa.api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zid), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Println(err)
		return  err
	}

	for  _,r:= range records {
		if(r.Name==domin&&r.Content==value){
			return  cfa.api.DeleteDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zid),r.ID)
		}
	}
	
	return errors.New("dns rcord not found")
}


func (cfa Adapter) SetRecord(domin string, value string) error {
	zname,id, err := cfa.extractZone(domin)
	if err != nil {
		log.Println(err)
		return err
	}
	f:=false
	_,err=cfa.api.CreateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(id),cloudflare.CreateDNSRecordParams{
		Type: "A",
		Tags: nil,
		Proxiable: true,
		ZoneName: zname,
		ZoneID: id,
		Data: nil,
		TTL:60,
		Proxied: &f,
		Name: domin,
		Content: value,
		Comment: "",
	})

	if err != nil {
		return err
	}
	return nil
}
