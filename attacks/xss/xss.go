package xss

import (
	"net/http"
	"net/url"

	"github.com/by2waysprojects/coverage-ktd/model"
)

type XSSAttack struct {
	config model.AttackConfig
}

func New(config model.AttackConfig) model.Attack {
	return &XSSAttack{config: config}
}

func (a *XSSAttack) Execute(target string) error {
	u, _ := url.Parse(target)
	u.Path = a.config.Endpoint

	q := u.Query()
	for k, v := range a.config.Parameters {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(a.config.Method, u.String(), nil)
	for k, v := range a.config.Headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
