package common

import "os"

var Conf *AuthBotConfig

type Config interface {
	Parse()
}

type AuthBotConfig struct {
	SecretKey string
	Token     string
}

func NewAuthBotConfig() {
	config := new(AuthBotConfig)
	config.parse()
	Conf = config

}

func (config *AuthBotConfig) parse() {
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
