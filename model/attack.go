package model

type Attack interface {
	Execute(target string) error
}

const (
	HTTP = "http"
	DNS  = "dns"
)
