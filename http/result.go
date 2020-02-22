package http

type Result interface {
	Unmarshal(object interface{}) error
}
