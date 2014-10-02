package config

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	config, err := ParseConfig("test_files/vanguard.yml")
	if err != nil {
		t.Errorf("Expected to not get error while parsing vanguard.yml, but got %v", err)
	}
	if config.Application != "SampleApp" {
		t.Errorf("Expected to get 'SampleApp' application name, but got '%v'", config.Application)
	}
	if len(config.Servers) != 1 {
		t.Errorf("Expected to get 1 application server, but got '%v'", len(config.Servers))
	}
	if len(config.Servers[0].Containers) != 2 {
		t.Errorf("Expected to get 2 container, but got '%v'", len(config.Servers[0].Containers))
	}
	expose := config.Servers[0].Containers[0].Expose[0]
	if expose != "8500:8500" {
		t.Errorf("Expected first container to expose '8500:8500' ports, but got %v", expose)
	}
	variable := config.Servers[0].Containers[0].Variables[0]
	if variable != "ENV:production" {
		t.Errorf("Expected first container to have variable 'ENV:production', but got %v", variable)
	}
}

func TestGetTiers(t *testing.T) {
	config, err := ParseConfig("test_files/tiers.yml")
	if err != nil {
		t.Errorf("Expected to not get error while parsing tiers.yml, but got %v", err)
	}
	tiers, err := config.GetTiers()
	if err != nil {
		t.Errorf("Expected to not get error while splitting config into tiers, but got %v", err)
	}
	if len(tiers) != 3 {
		t.Errorf("Expected to get 3 tiers, but got %v", len(tiers))
	}
}
