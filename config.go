package sky

import (
	"fmt"

	"github.com/spf13/viper"
)

// Home home url
func Home() string {
	if IsProduction() {
		name := viper.GetString("server.name")
		if viper.GetBool("server.ssl") {
			return "https://" + name
		}
		return "http://" + name
	}
	return fmt.Sprintf("http://localhost:%d", viper.GetInt("server.port"))
}

// IsProduction production mode ?
func IsProduction() bool {
	return viper.GetString("env") == "production"
}

// DataSource database source url
func DataSource() string {
	//"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s"
	args := ""
	for k, v := range viper.GetStringMapString("database.args") {
		args += fmt.Sprintf(" %s=%s ", k, v)
	}
	return args

	//"postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	// return fmt.Sprintf(
	// 	"%s://%s:%s@%s:%d/%s?sslmode=%s",
	// 	viper.GetString("database.driver"),
	// 	viper.GetString("database.args.user"),
	// 	viper.GetString("database.args.password"),
	// 	viper.GetString("database.args.host"),
	// 	viper.GetInt("database.args.port"),
	// 	viper.GetString("database.args.dbname"),
	// 	viper.GetString("database.args.sslmode"),
	// )
}
