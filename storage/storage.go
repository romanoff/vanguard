package storage

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	KeyPresent(key string) (bool, error)
}

var storage Storage

func GetStorage() Stoage {
	if storage == nil {
		storage = NewEtcdStorage()
	}
	return storage
}
