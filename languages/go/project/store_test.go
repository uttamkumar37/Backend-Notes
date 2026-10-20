package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestLifecycleAndSnapshots(t *testing.T) {
	store := newStore()
	first, err := store.Create(Input{Name: " Asha ", Country: "in"})
	if err != nil || first.Name != "Asha" || first.Country != "IN" {
		t.Fatalf("create = %+v, %v", first, err)
	}
	first.Name = "changed outside store"
	got, err := store.Get(first.ID)
	if err != nil || got.Name != "Asha" {
		t.Fatalf("stored snapshot = %+v, %v", got, err)
	}
	updated, err := store.Update(first.ID, Input{Name: "Mira", Country: "gb"})
	if err != nil || updated.ID != first.ID || updated.Name != "Mira" || updated.Country != "GB" {
		t.Fatalf("update = %+v, %v", updated, err)
	}
	if err := store.Delete(first.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(first.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get deleted investor = %v", err)
	}
	if err := store.Delete(first.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing investor = %v", err)
	}
	if _, err := store.Update(first.ID, Input{Name: "Asha", Country: "IN"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("update missing investor = %v", err)
	}
}

func TestValidationDoesNotMutateStore(t *testing.T) {
	cases := []Input{
		{Name: " ", Country: "IN"},
		{Name: strings.Repeat("界", 101), Country: "IN"},
		{Name: "Asha", Country: "IND"},
		{Name: "Asha", Country: "12"},
		{Name: "Asha", Country: ""},
	}
	for _, input := range cases {
		t.Run(fmt.Sprintf("%d-%s", len(input.Name), input.Country), func(t *testing.T) {
			store := newStore()
			if _, err := store.Create(input); !errors.Is(err, ErrInvalid) {
				t.Fatalf("create invalid input = %v", err)
			}
			item, err := store.Create(Input{Name: "Asha", Country: "IN"})
			if err != nil || item.ID != 1 {
				t.Fatalf("invalid input consumed an ID: %+v, %v", item, err)
			}
			if _, err := store.Update(item.ID, input); !errors.Is(err, ErrInvalid) {
				t.Fatalf("update invalid input = %v", err)
			}
			got, err := store.Get(item.ID)
			if err != nil || got.Name != "Asha" {
				t.Fatalf("invalid update changed data: %+v, %v", got, err)
			}
		})
	}
}

func TestPaginationAndConcurrentAccess(t *testing.T) {
	store := newStore()
	var workers sync.WaitGroup
	for i := 0; i < 80; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			item, err := store.Create(Input{Name: "Fictional", Country: "IN"})
			if err != nil {
				t.Error(err)
				return
			}
			if _, err := store.Get(item.ID); err != nil {
				t.Error(err)
			}
			store.List(0, 10)
			if _, err := store.Update(item.ID, Input{Name: "Updated", Country: "GB"}); err != nil {
				t.Error(err)
			}
		}()
	}
	workers.Wait()
	page := store.List(40, 10)
	if len(page) != 10 || page[0].ID != 41 || page[9].ID != 50 {
		t.Fatalf("unexpected page: %+v", page)
	}
	page[0].Name = "external mutation"
	got, err := store.Get(41)
	if err != nil || got.Name != "Updated" {
		t.Fatalf("list exposed stored values: %+v, %v", got, err)
	}
	if got := store.List(80, 10); got == nil || len(got) != 0 {
		t.Fatalf("empty page = %#v", got)
	}
}
