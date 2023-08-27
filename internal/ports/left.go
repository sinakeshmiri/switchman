package ports
type APIPort interface {
	Check() (error)
}