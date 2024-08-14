package model

type OracleConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Service  string
}

type MattermostConfig struct {
	Url   string
	Token string
}

type Config struct {
	OracleDb   OracleConfig
	Mattermost MattermostConfig
}
