package config

import "github.com/spf13/viper"

type Cnf struct {
	Db          Db
	GithubToken string
}

type Db struct {
	Name string
	User string
	Pass string
	Host string
	Port string
}

func NewConfig() *Cnf {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	user := viper.GetString("SUPA_POSTGRE_USER")
	pass := viper.GetString("SUPA_POSTGRE_PASSWORD")
	host := viper.GetString("SUPA_POSTGRE_HOST")
	port := viper.GetString("SUPA_POSTGRE_PORT")
	name := viper.GetString("SUPA_POSTGRE_DB")
	gitToken := viper.GetString("GITHUB_TOKEN")

	return &Cnf{
		Db: Db{
			User: user,
			Name: name,
			Pass: pass,
			Host: host,
			Port: port,
		},
		GithubToken: gitToken,
	}
}
