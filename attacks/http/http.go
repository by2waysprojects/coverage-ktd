package httpAttack

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/by2waysprojects/coverage-ktd/model"
)

type HttpAttack struct {
	config model.AttackConfig
}

func New(config model.AttackConfig) model.Attack {
	return &HttpAttack{config: config}
}

func (a *HttpAttack) Execute(target string) error {
	u, _ := url.Parse(target)
	u.Path = a.config.Endpoint

	q := u.Query()
	for k, v := range a.config.Parameters {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(a.config.Method, u.String(), nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}
	for k, v := range a.config.Headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error executing HTTP request:", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
