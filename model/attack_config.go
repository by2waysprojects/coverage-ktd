package model

type AttackConfig struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Method      string            `json:"method"`
	Endpoint    string            `json:"endpoint"`
	Headers     map[string]string `json:"headers"`
	Parameters  map[string]string `json:"parameters"`
}
