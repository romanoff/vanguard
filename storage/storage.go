package storage

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	KeyPresent(key string) (bool, error)
	Delete(key string) error
}

var storage Storage

func GetStorage() Storage {
	if storage == nil {
		storage = NewEtcdStorage()
	}
	return storage
}
