package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

var (
	ErrInvalid  = errors.New("invalid investor")
	ErrNotFound = errors.New("investor not found")
)

// Investor is a value snapshot. Country is a two-letter label, not verified residency.
type Investor struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// Input deliberately excludes the server-owned ID.
type Input struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

// Store is process-local memory. It must not be copied after first use.
type Store struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]Investor
}

func newStore() *Store {
	return &Store{items: make(map[int64]Investor)}
}

func normalize(input Input) (Input, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Country = strings.ToUpper(strings.TrimSpace(input.Country))
	if !utf8.ValidString(input.Name) || input.Name == "" || utf8.RuneCountInString(input.Name) > 100 {
		return Input{}, fmt.Errorf("%w: name must contain 1–100 Unicode code points", ErrInvalid)
	}
	if len(input.Country) != 2 || input.Country[0] < 'A' || input.Country[0] > 'Z' || input.Country[1] < 'A' || input.Country[1] > 'Z' {
		return Input{}, fmt.Errorf("%w: country must contain two ASCII letters", ErrInvalid)
	}
	return input, nil
}

func (s *Store) Create(input Input) (Investor, error) {
	input, err := normalize(input)
	if err != nil {
		return Investor{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	item := Investor{ID: s.nextID, Name: input.Name, Country: input.Country}
	s.items[item.ID] = item
	return item, nil
}

func (s *Store) Get(id int64) (Investor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return Investor{}, ErrNotFound
	}
	return item, nil
}

// List returns an ID-ordered page. Callers validate limit in the range 1–100.
func (s *Store) List(after int64, limit int) []Investor {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Investor, 0)
	for _, item := range s.items {
		if item.ID > after {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	if len(items) > limit {
		items = items[:limit]
	}
	return items
}

func (s *Store) Update(id int64, input Input) (Investor, error) {
	input, err := normalize(input)
	if err != nil {
		return Investor{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return Investor{}, ErrNotFound
	}
	item := Investor{ID: id, Name: input.Name, Country: input.Country}
	s.items[id] = item
	return item, nil
}

func (s *Store) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
