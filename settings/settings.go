package settings

import (
	"bytes"
	"encoding/gob"

	"github.com/kapmahc/sky/security"
)

// New new Settings
func New(s Store, f *security.Factory) *Settings {
	return &Settings{
		store:   s,
		factory: f,
	}
}

// Settings setting helper
type Settings struct {
	store   Store
	factory *security.Factory
}

//Set save setting
func (p *Settings) Set(k string, v interface{}, f bool) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return err
	}
	var val []byte
	if f {
		if val, err = p.factory.To(buf.Bytes()); err != nil {
			return err
		}
	} else {
		val = buf.Bytes()
	}
	return p.store.Set(k, val, f)

}

//Get get setting value by key
func (p *Settings) Get(k string, v interface{}) error {
	val, enc, err := p.store.Get(k)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)

	if enc {
		vl, er := p.factory.From(val)
		if er != nil {
			return er
		}
		buf.Write(vl)
	} else {
		buf.Write(val)
	}

	return dec.Decode(v)
}
