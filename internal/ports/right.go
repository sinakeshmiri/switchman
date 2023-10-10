package ports

type NsproviderPort interface {
	SetRecord(domin string,value string)  error 
	GetRecords(domin string) ([]string, error) 
	DelRecord(domin string, value  string) (error) 
}

type DownStreamPort interface {
	IsPrimary()(bool)
	GetIP()(string)
	CheckHealth() (error)
}