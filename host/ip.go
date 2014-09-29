package host

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

var ipRegexp *regexp.Regexp = regexp.MustCompile("inet addr:(.*)\\s")

func GetIpAddress(interf string) (string, error) {
	interfaceData, err := exec.Command("ifconfig", interf).Output()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not find %v interface", interf))
	}
	lines := bytes.Split(interfaceData, []byte("\n"))
	ip := ""
	for _, line := range lines {
		foundIp := ipRegexp.FindAllSubmatch(line, -1)
		if foundIp != nil && foundIp[0] != nil && foundIp[0][1] != nil {
			ip = string(foundIp[0][1])
			break
		}
	}
	if ip == "" {
		return "", errors.New("Ip address not found")
	}
	return ip, nil
}
