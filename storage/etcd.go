package storage

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"strings"
)

func NewEtcdStorage() Storage {
	return &EtcdStorage{etcd.NewClient(nil)}
}

type EtcdStorage struct {
	client *etcd.Client
}

func (self *EtcdStorage) Set(key, value string) error {
	_, err := self.client.Set(key, value, 0)
	return err
}

func (self *EtcdStorage) KeyPresent(key string) (bool, error) {
	resp, err := self.client.Get(key, false, false)
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

func (self *EtcdStorage) Get(key string) (string, error) {
	resp, err := self.client.Get(key, false, false)
	if err != nil {
		return "", err
	}
	if resp.Node != nil {
		return resp.Node.Value, nil
	}
	return "", errors.New("key missing")
}

func (self *EtcdStorage) Delete(key string) error {
	_, err := self.client.Delete(key, true)
	return err
}

func (self *EtcdStorage) Keys(prefix string) ([]string, error) {
	resp, err := self.client.Get("/", false, false)
	if err != nil {
		return nil, err
	}
	keys := []string{}
	if resp.Node == nil {
		return keys, nil
	}
	for _, node := range resp.Node.Nodes {
		if strings.HasPrefix(node.Key, "/" + prefix) {
			keys = append(keys, node.Key)
		}
	}
	return keys, nil
}
