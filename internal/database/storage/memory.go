package storage

import "sync"

type MemoryStorage struct {
	mtx     *sync.RWMutex
	hashMap map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		mtx:     &sync.RWMutex{},
		hashMap: make(map[string]string),
	}
}

func (s *MemoryStorage) Set(key string, value string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hashMap[key] = value
	return nil
}

func (s *MemoryStorage) Get(key string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.hashMap[key], nil
}

func (s *MemoryStorage) Delete(key string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.hashMap, key)
	return nil
}
