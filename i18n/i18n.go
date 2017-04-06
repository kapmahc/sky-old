package i18n

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
	"golang.org/x/text/language"
)

// I18n i18n
type I18n struct {
	Store Store `inject:""`
}

// F format message
func (p *I18n) F(lang, code string, obj interface{}) (string, error) {
	msg, err := p.Store.Get(lang, code)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("").Parse(msg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, obj)
	return buf.String(), err
}

//E create an i18n error
func (p *I18n) E(lang, code string, args ...interface{}) error {
	msg, err := p.Store.Get(lang, code)
	if err != nil {
		return errors.New(code)
	}
	return fmt.Errorf(msg, args...)
}

//T translate by lang tag
func (p *I18n) T(lang, code string, args ...interface{}) string {
	msg, err := p.Store.Get(lang, code)
	if err != nil {
		return code
	}
	return fmt.Sprintf(msg, args...)
}

// All all items
func (p *I18n) All(lang string) (map[string]interface{}, error) {
	rt := make(map[string]interface{})

	items, err := p.Store.All(lang)
	if err != nil {
		return nil, err
	}
	for k, v := range items {
		codes := strings.Split(k, ".")
		tmp := rt
		for i, c := range codes {
			if i+1 == len(codes) {
				tmp[c] = v
			} else {
				if tmp[c] == nil {
					tmp[c] = make(map[string]interface{})
				}
				tmp = tmp[c].(map[string]interface{})
			}
		}

	}
	return rt, nil
}

// Load sync records
func (p *I18n) Load(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		const ext = ".ini"
		name := info.Name()
		if info.Mode().IsRegular() && filepath.Ext(name) == ext {
			log.Debugf("Find locale file %s", path)
			if err != nil {
				return err
			}
			lang := name[0 : len(name)-len(ext)]
			if _, err := language.Parse(lang); err != nil {
				return err
			}
			cfg, err := ini.Load(path)
			if err != nil {
				return err
			}

			items := cfg.Section(ini.DEFAULT_SECTION).KeysHash()
			for k, v := range items {
				if err := p.Store.Set(lang, k, v, false); err != nil {
					return err
				}
			}
			log.Infof("find %d items", len(items))
		}
		return nil
	})
}
