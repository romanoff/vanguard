package command

import (
	"github.com/romanoff/vanguard/config"
)

func getHostname(hostname string) string {
	if hostname == "" {
		hostname = "127.0.0.1"
	}
	cfg, _ := config.ParseConfig("vanguard.yml")
	if cfg != nil {
		for _, server := range cfg.Servers {
			if server.Label == hostname {
				hostname = server.Hostname
			}
		}
	}
	return hostname
}
