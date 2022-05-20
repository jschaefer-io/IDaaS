package utils

import (
	"net/url"
)

func AddQueryToUrl(baseUrl string, params map[string]string) (string, error) {
	pUrl, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return baseUrl, err
	}
	q := pUrl.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	pUrl.RawQuery = q.Encode()
	return pUrl.String(), nil
}

func GetResetConfirmUrl(baseUrl string, params map[string]string) string {
	link, _ := AddQueryToUrl(baseUrl+"/reset/confirm", params)
	return link
}

func GetUserConfirmUrl(baseUrl string, params map[string]string) string {
	link, _ := AddQueryToUrl(baseUrl+"/user/confirm", params)
	return link
}

func GetResetUrl(baseUrl string, params map[string]string) string {
	link, _ := AddQueryToUrl(baseUrl+"/reset", params)
	return link
}

func GetLoginUrl(baseUrl string, params map[string]string) string {
	link, _ := AddQueryToUrl(baseUrl+"/login", params)
	return link
}
