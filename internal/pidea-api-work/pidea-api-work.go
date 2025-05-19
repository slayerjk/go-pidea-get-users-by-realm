package pideaapiwork

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Token struct {
	Result struct {
		Value struct {
			Token string `json:"token"`
		} `json:"value"`
	} `json:"result"`
}

type Users struct {
	Result struct {
		Value []User `json:"value"`
	} `json:"result"`
}

// user struct of Get user Pidea API response.
// excluded fields are "editable" & "userid" as no valuable info
type User struct {
	// Editable  bool     `json:"editable"`
	Email     string   `json:"email"`
	Givenname string   `json:"givenname"`
	MemberOf  []string `json:"memberOf"`
	Mobile    string   `json:"mobile"`
	Phone     string   `json:"phone"`
	Resolver  string   `json:"resolver"`
	Surname   string   `json:"surname"`
	// Userid    string   `json:"userid"`
	Username string `json:"username"`
}

// Get Pidea API Token for given user(POST)
func GetPideaApiToken(httpClient *http.Client, baseUrl, userName, UserPassword string) (string, error) {
	var tokenData Token

	// form URL query
	query := fmt.Sprintf("%s/auth?username=%s&password=%s", baseUrl, userName, UserPassword)

	// form request
	request, err := http.NewRequest(http.MethodPost, query, nil)
	if err != nil {
		return "", fmt.Errorf("failed to Form getToken POST request,\n\t%v", err)
	}

	// do request
	response, err := httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to Make getToken POST request,\n\t%v", err)
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return "", fmt.Errorf("wrong credentials for getToken POST request StatusCode,\n\t%s", response.Status)
		}
		return "", fmt.Errorf("check getToken POST request StatusCode,\n\t%s", response.Status)
	}

	// read response
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to Read getToken POST request,\n\t%v", err)
	}

	// unmarshall json response
	err = json.Unmarshal(responseBytes, &tokenData)
	if err != nil {
		return "", fmt.Errorf("failed to Unmarshall getToken POST request,\n\t%v", err)
	}

	if token := tokenData.Result.Value.Token; len(token) != 0 {
		return token, nil
	}

	return "", fmt.Errorf("token result of getToken POST request is Empty\n\t%+v", tokenData)
}

// Get Pidea users of given realm(GET)
func GetPideaUsersByRealm(httpClient *http.Client, baseUrl, realm, apiToken string) ([]User, error) {
	var usersData Users

	// form URL query
	query := fmt.Sprintf("%s/user?realm=%s", baseUrl, realm)

	// form request
	request, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to Form getUsersByRealm GET request,\n\t%v", err)
	}
	request.Header.Set("Authorization", apiToken)

	// do request
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to Make getUsersByRealm GET request,\n\t%v", err)
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return nil, fmt.Errorf("auth failure for getUsersByRealm GET request StatusCode,\n\t%s", response.Status)
		}
		return nil, fmt.Errorf("check getUsersByRealm GET request StatusCode,\n\t%s", response.Status)
	}

	// read response
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to Read getUsersByRealm GET request,\n\t%v", err)
	}

	// unmarshall json response
	err = json.Unmarshal(responseBytes, &usersData)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshall getUsersByRealm GET request,\n\t%v", err)
	}

	if users := usersData.Result.Value; len(users) != 0 {
		return users, nil
	}

	return nil, fmt.Errorf("token result of usersByRealm GET request is Empty\n\t%+v", usersData)
}
