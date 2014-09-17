package container

import (
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
		ip = "10.0.1." + strconv.Itoa(counter) + "/24"
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
	return db.Delete(ip)
}
