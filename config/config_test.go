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
	serverExpose := config.Servers[0].Expose[0]
	if serverExpose != "3306:172.0.0.2:3306" {
		t.Errorf("Expected to get '%v' as server expose, but got %v", "3306:172.0.0.2:3306", serverExpose)
	}
	expose := config.Servers[0].Containers[0].Expose[0]
	if expose != "8500:8500" {
		t.Errorf("Expected first container to expose '8500:8500' ports, but got %v", expose)
	}
	variable := config.Servers[0].Containers[0].Variables[0]
	if variable != "ENV:production" {
		t.Errorf("Expected first container to have variable 'ENV:production', but got %v", variable)
	}
	volume := config.Servers[0].Containers[0].Volumes[0]
	if volume != "/data/volume:/data" {
		t.Errorf("Expected first container to have '/data/volume:/data' volume, but got %v", volume)
	}
	command := config.Servers[0].Containers[0].Command
	if command != "/bin/bash" {
		t.Errorf("Expected to get '/bin/bash' as command, but got %v", command)
	}
	dockerfile := config.Servers[0].Containers[0].Dockerfile
	expected := "/home/user/project/Dockerfile"
	if dockerfile != expected {
		t.Errorf("Expected to get '%v' as dockerfile, but got %v", expected, dockerfile)
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

func TestParseConfigRemote(t *testing.T) {
	config, err := ParseConfig("test_files/vanguard.yml")
	if err != nil {
		t.Errorf("Expected to not get error while parsing vanguard.yml, but got %v", err)
	}
	remote := config.Remote
	if remote == nil {
		t.Error("Expected to parse vanguard.yml remote, but got nil")
	}
	if remote.Type != "s3" {
		t.Errorf("Expected to get remote type 's3', but got '%v' ", remote.Type)
	}
	if remote.Bucket != "bucket" {
		t.Errorf("Expected to get remote bucket 'bucket', but got '%v' ", remote.Bucket)
	}
	if remote.Access_Key != "access" {
		t.Errorf("Expected to get remote access key 'access', but got '%v' ", remote.Access_Key)
	}
	if remote.Secret_Key != "secret" {
		t.Errorf("Expected to get remote access key 'secret', but got '%v' ", remote.Secret_Key)
	}
	if remote.Region != "us" {
		t.Errorf("Expected to get remote region 'us', but got '%v' ", remote.Region)
	}
	if len(remote.Files) != 1 {
		t.Errorf("Expected to get 1 remote file, but got '%v' ", len(remote.Files))
	}
	remoteFile := remote.Files[0]
	if remoteFile.Name != "sphinx.tar.bz2" {
		t.Errorf("Expected to get 'sphinx.tar.bz2' remote file name, but got %v", remoteFile.Name)
	}
	if remoteFile.Path != "sphinx" {
		t.Errorf("Expected to get 'sphinx' remote file path, but got '%v'", remoteFile.Path)
	}
}
