package store

import (
	"fmt"
	"strings"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	indices map[string]map[string]map[string]any
}

func New() *Store {
	return &Store{indices: make(map[string]map[string]map[string]any)}
}

func (s *Store) CreateIndex(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.indices[name]; ok {
		return false
	}
	s.indices[name] = make(map[string]map[string]any)
	return true
}

func (s *Store) DeleteIndex(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.indices[name]; !ok {
		return false
	}
	delete(s.indices, name)
	return true
}

func (s *Store) IndexExists(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.indices[name]
	return ok
}

func (s *Store) ListIndices() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.indices))
	for k := range s.indices {
		out = append(out, k)
	}
	return out
}

func (s *Store) PutDoc(index, id string, doc map[string]any) (created bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.indices[index]; !ok {
		s.indices[index] = make(map[string]map[string]any)
	}
	_, exists := s.indices[index][id]
	s.indices[index][id] = doc
	return !exists
}

func (s *Store) GetDoc(index, id string) (map[string]any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if idx, ok := s.indices[index]; ok {
		if doc, ok := idx[id]; ok {
			return doc, true
		}
	}
	return nil, false
}

func (s *Store) DeleteDoc(index, id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if idx, ok := s.indices[index]; ok {
		if _, ok := idx[id]; ok {
			delete(s.indices[index], id)
			return true
		}
	}
	return false
}

func (s *Store) UpdateDoc(index, id string, partial map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx, ok := s.indices[index]
	if !ok {
		return fmt.Errorf("index '%s' tidak ditemukan", index)
	}
	doc, ok := idx[id]
	if !ok {
		return fmt.Errorf("dokumen '%s' tidak ditemukan", id)
	}
	for k, v := range partial {
		doc[k] = v
	}
	return nil
}

func (s *Store) Search(index, query string) []map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	query = strings.ToLower(query)
	var hits []map[string]any
	idx, ok := s.indices[index]
	if !ok {
		return hits
	}
	for id, doc := range idx {
		for _, v := range doc {
			if strings.Contains(strings.ToLower(fmt.Sprintf("%v", v)), query) {
				result := map[string]any{"_id": id, "_source": doc}
				hits = append(hits, result)
				break
			}
		}
	}
	return hits
}

func (s *Store) Count(index string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if idx, ok := s.indices[index]; ok {
		return len(idx)
	}
	return 0
}
