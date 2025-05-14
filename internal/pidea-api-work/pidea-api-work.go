package pideaapiwork

import "net/http"

type Token struct {
	Result struct {
		Authentication string `json:"authentication"`
		Status         bool   `json:"status"`
		Value          struct {
			Token string `json:"token"`
		} `json:"value"`
	} `json:"result"`
}

func getPideaApiToken(httpClient http.Client, baseUrl, userName, UserPassword string) (string, error) {
	var result Token

	return
}
