package server

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	Url      string
	Token    tokenSettings
	Redirect redirectSettings
	Mail     mailSettings
	Static   staticSettings
}

type tokenSettings struct {
	Secret string
}

type redirectSettings struct {
	Default   string
	Whitelist []string
}

type mailSettings struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

type staticSettings struct {
	Dir string
}

func SettingsFromEnv() (*Settings, error) {
	mailPort, err := strconv.Atoi(os.Getenv("JIO_MAIL_PORT"))
	if err != nil {
		return nil, errors.New("JIO_MAIL_PORT must be an integer")
	}
	return &Settings{
		Url: os.Getenv("JIO_PUBLIC_URL"),
		Token: tokenSettings{
			Secret: os.Getenv("JIO_TOKEN_SECRET"),
		},
		Redirect: redirectSettings{
			Default:   os.Getenv("JIO_REDIRECT_DEFAULT"),
			Whitelist: strings.Split(os.Getenv("JIO_REDIRECT_WHITELIST"), ","),
		},
		Mail: mailSettings{
			From:     os.Getenv("JIO_MAIL_FROM"),
			Host:     os.Getenv("JIO_MAIL_HOST"),
			Port:     mailPort,
			Username: os.Getenv("JIO_MAIL_USERNAME"),
			Password: os.Getenv("JIO_MAIL_PASSWORD"),
		},
		Static: staticSettings{
			Dir: os.Getenv("JIO_STATIC_DIR"),
		},
	}, nil
}
