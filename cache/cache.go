package cache

import (
	"bytes"
	"encoding/gob"
	"time"
)

// New new cache
func New(n string, s Store) *Cache {
	return &Cache{
		store: s,
	}
}

// Cache cache
type Cache struct {
	store Store
}

//Set cache item
func (p *Cache) Set(key string, val interface{}, ttl time.Duration) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}
	return p.store.Set(key, buf.Bytes(), ttl)
}

//Get get from cache
func (p *Cache) Get(key string, val interface{}) error {
	bys, err := p.store.Get(key)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(bys)
	return dec.Decode(val)
}
