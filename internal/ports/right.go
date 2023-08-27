package ports

type NsproviderPort interface {
	SetRecord(domin string,value string)  error 
	GetRecord(domin string) (string, error) 
}

type DownStreamPort interface {
	GetIP()(string)
	CheckHealth() (error)
}