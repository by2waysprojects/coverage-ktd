package attacks

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	DNSAttack "github.com/by2waysprojects/coverage-ktd/attacks/dns"
	HTTPAttack "github.com/by2waysprojects/coverage-ktd/attacks/http"
	"github.com/by2waysprojects/coverage-ktd/model"
)

const (
	configDirectory = "config"
)

type AttackExecutor struct {
	attacks map[string]model.Attack
	target  string
}

func NewAttackExecutor(target string) *AttackExecutor {
	return &AttackExecutor{
		attacks: make(map[string]model.Attack),
		target:  target,
	}
}

func (e *AttackExecutor) LoadAttacks(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			subdirPath := filepath.Join(dir, file.Name())
			if err := e.loadAttackType(subdirPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *AttackExecutor) loadAttackType(dir string) error {
	attackType := filepath.Base(dir)
	configDirectory := filepath.Join(dir, configDirectory)

	files, err := os.ReadDir(configDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			configPath := filepath.Join(configDirectory, file.Name())
			config, err := loadAttackConfig(configPath)
			if err != nil {
				log.Printf("Error loading attack config from %s: %v", configPath, err)
				continue
			}
			for _, cfg := range config {
				switch attackType {
				case model.HTTP:
					e.attacks[cfg.Name] = HTTPAttack.New(cfg)
				case model.DNS:
					e.attacks[cfg.Name] = DNSAttack.New(cfg)
				default:
					log.Printf("Unknown attack type: %s", attackType)
				}
			}
		}
	}

	return nil
}

func loadAttackConfig(path string) ([]model.AttackConfig, error) {
	var config []model.AttackConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

func (e *AttackExecutor) RunAll() {
	for _, attack := range e.attacks {
		attack.Execute(e.target)
	}
}

func (e *AttackExecutor) GetAttacks() map[string]model.Attack {
	return e.attacks
}
