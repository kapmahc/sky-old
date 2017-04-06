package sky

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	// Version version
	Version string
	// BuildTime build time
	BuildTime string
)

// New new app
func New() *cli.App {

	app := cli.NewApp()
	app.Name = os.Args[0]
	app.Version = fmt.Sprintf("%s(%s)", Version, BuildTime)
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{}

	for _, en := range engines {
		cmd := en.Shell()
		app.Commands = append(app.Commands, cmd...)
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	return app
}

func init() {
	viper.SetEnvPrefix("sky")
	viper.BindEnv("env")
	viper.SetDefault("env", "development")

}
