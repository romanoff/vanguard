package container

import (
	"encoding/json"
	"errors"
	"github.com/romanoff/vanguard/storage"
	"strconv"
)

func ReserveIp() (string, error) {
	db := storage.GetStorage()
	counter := 0
	ip := ""
	for {
		counter++
		ip = "10.0.1." + strconv.Itoa(counter)
		ok, err := db.KeyPresent(ip)
		if err != nil {
			return "", err
		}
		if ok == false {
			err = db.Set(ip, "reserved")
			if err != nil {
				return "", err
			}
			break
		}
	}
	return ip, nil
}

func FreeIp(ip string) error {
	if ip == "" {
		return errors.New("no ip to free")
	}
	db := storage.GetStorage()
	container, err := GetByIp(ip)
	if err == nil {
		db.Delete(container.ContainerId)
	}
	return db.Delete(ip)
}

func Persist(c *Container) error {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if c.Ip == "" {
		return errors.New("Can't persist container that is not running -" + c.String())
	}
	db := storage.GetStorage()
	err = db.Set(c.Ip, string(jsonBytes))
	if err != nil {
		return err
	}
	err = db.Set(c.ContainerId, c.Ip)
	if err != nil {
		return err
	}
	return nil
}

func GetByIp(ip string) (*Container, error) {
	db := storage.GetStorage()
	jsonString, err := db.Get(ip)
	if err != nil {
		return nil, err
	}
	var container *Container
	err = json.Unmarshal([]byte(jsonString), &container)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func GetByContainerId(containerId string) (*Container, error) {
	db := storage.GetStorage()
	ip, err := db.Get(containerId)
	if err != nil {
		return nil, err
	}
	return GetByIp(ip)
}

func GetContainers() ([]*Container, error) {
	db := storage.GetStorage()
	ips, err := db.Keys("10.0.1.")
	if err != nil {
		return nil, err
	}
	containers := []*Container{}
	for _, ip := range ips {
		container, err := GetByIp(ip)
		if err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}
	return containers, nil
}
