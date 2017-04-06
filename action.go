package sky

import (
	"log/syslog"

	log "github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Action cfg action
func Action(f cli.ActionFunc) cli.ActionFunc {
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
