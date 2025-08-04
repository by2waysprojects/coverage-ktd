package attacks

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"

	"github.com/by2waysprojects/coverage-ktd/model"
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
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != dir {
			return e.loadAttackType(path)
		}
		return nil
	})
}

func (e *AttackExecutor) loadAttackType(dir string) error {
	attackType := filepath.Base(dir)
	configFile := filepath.Join(dir, "config.json")

	config, err := loadAttackConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config for %s: %v", attackType, err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".so" {
			pluginPath := filepath.Join(dir, f.Name())
			if err := e.loadPlugin(attackType, pluginPath, config); err != nil {
				log.Printf("Error loading plugin %s: %v", pluginPath, err)
			}
		}
	}
	return nil
}

func loadAttackConfig(path string) (model.AttackConfig, error) {
	var config model.AttackConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	return config, json.Unmarshal(data, &config)
}

func (e *AttackExecutor) loadPlugin(attackType, path string, config model.AttackConfig) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	sym, err := p.Lookup("New")
	if err != nil {
		return err
	}

	newFunc, ok := sym.(func(model.AttackConfig) model.Attack)
	if !ok {
		return fmt.Errorf("invalid constructor signature")
	}

	e.attacks[attackType] = newFunc(config)
	log.Printf("Loaded attack: %s (%s)", attackType, config.Name)
	return nil
}

func (e *AttackExecutor) RunAll() {
	for _, attack := range e.attacks {
		attack.Execute(e.target)
	}
}

func (e *AttackExecutor) GetAttacks() map[string]model.Attack {
	return e.attacks
}
