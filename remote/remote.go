package remote

type Remote interface {
	Push(filepath, destination string) error
	Pull(filepath, destination string) error
	Remove(filepath string) error
	FilesList(prefix string) ([]string, error)
}
