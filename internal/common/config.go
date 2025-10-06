package common

import "os"

type Config interface {
	Parse()
}

type AuthBotConfig struct {
	AppName   string
	SecretKey string
	Token     string
}

func NewAuthBotConfig() *AuthBotConfig {
	config := new(AuthBotConfig)
	config.parse()
	return config
}

func (config *AuthBotConfig) parse() {
	config.AppName = "durak"
	config.Token = parseVar("TOKEN")
	config.SecretKey = parseVar("SECRET_KEY")
}

func parseVar(varName string) string {
	variable := os.Getenv(varName)
	if variable == "" {
		panic(varName + " not provided")
	}
	return variable
}
