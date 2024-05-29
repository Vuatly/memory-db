package storage

type Engine interface {
	Get(string) (string, error)
	Set(string, string) error
	Delete(string) error
}
