package dnsAttack

import (
	"github.com/by2waysprojects/coverage-ktd/model"
	"github.com/miekg/dns"
)

type DnsAttack struct {
	config model.AttackConfig
}

func New(config model.AttackConfig) model.Attack {
	return &DnsAttack{config: config}
}

func (a *DnsAttack) Execute(target string) error {
	msg := new(dns.Msg)
	domain := a.config.Parameters["domain"]
	recordType := a.config.Parameters["type"]

	qtype, ok := dns.StringToType[recordType]
	if !ok {
		qtype = dns.TypeA
	}

	msg.SetQuestion(dns.Fqdn(domain), qtype)

	client := new(dns.Client)

	_, _, err := client.Exchange(msg, target)
	if err != nil {
		return err
	}

	return nil
}
