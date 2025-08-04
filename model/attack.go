package model

type Attack interface {
	Execute(target string) error
}
