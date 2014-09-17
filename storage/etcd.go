package storage

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"strings"
)

func NewEtcdStorage() Storage {
	return &EtcdClient{etcd.NewClient(nil)}
}

type EtcdStorage struct {
	client *etcd.Client
}

func (self *EtcdStorage) Set(key, value string) error {
	_, err := self.client.set(key, value, 0)
	return err
}

func (self *EtcdStorage) KeyPresent(key string) (bool, error) {
	resp, err := client.Get(key, false, false)
	if err != nil {
		if strings.HasPrefix(err.Error(), "100:") {
			return false, nil
		}
		return false, err
	}
	if resp.Node != nil {
		return true, nil
	}
	return false, nil
}

func (self *EtcdStorage) Get(key) (string, error) {
	resp, err := client.Get(ip, false, false)
	if err != nil {
		return "", err
	}
	if resp.Node != nil {
		return resp.Node.Value, nil
	}
	return "", errors.New("key missing")
}
