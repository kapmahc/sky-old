package sky

import (
	"log/syslog"

	"golang.org/x/text/language"

	log "github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/facebookgo/inject"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type injectLogger struct {
}

func (p *injectLogger) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// IocAction ioc action
func IocAction(fn cli.ActionFunc) cli.ActionFunc {
	return CfgAction(func(c *cli.Context) error {
		inj := inject.Graph{Logger: &injectLogger{}}
		// -------
		var tags []language.Tag
		for _, l := range viper.GetStringSlice("languages") {
			if lng, err := language.Parse(l); err == nil {
				tags = append(tags, lng)
			} else {
				return err
			}
		}

		// ---------------
		if err := inj.Provide(
			&inject.Object{Value: language.NewMatcher(tags)},
		); err != nil {
			return err
		}
		// -----------------
		if err := Walk(func(en Engine) error {
			if err := en.Map(&inj); err != nil {
				return err
			}
			return inj.Provide(&inject.Object{Value: en})
		}); err != nil {
			return err
		}

		if err := inj.Populate(); err != nil {
			return err
		}

		return fn(c)
	})
}

// CfgAction cfg action
func CfgAction(f cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		log.Infof("read config from config.toml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		// -----------
		if IsProduction() {
			log.SetLevel(log.InfoLevel)
			if wrt, err := syslog.New(syslog.LOG_INFO, viper.GetString("app.name")); err == nil {
				log.AddHook(&logrus_syslog.SyslogHook{Writer: wrt})
			} else {
				log.Error(err)
			}
		} else {
			log.SetLevel(log.DebugLevel)
		}
		return f(c)
	}
}
