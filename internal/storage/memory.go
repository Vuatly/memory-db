package storage

type MemoryStorage struct {
	hashMap map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		hashMap: make(map[string]string),
	}
}

func (s *MemoryStorage) Set(key string, value string) error {
	s.hashMap[key] = value
	return nil
}

func (s *MemoryStorage) Get(key string) (string, error) {
	return s.hashMap[key], nil
}

func (s *MemoryStorage) Delete(key string) error {
	delete(s.hashMap, key)
	return nil
}
