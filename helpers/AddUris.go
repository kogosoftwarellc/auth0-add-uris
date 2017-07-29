package helpers

import (
	"github.com/kogosoftwarellc/auth0-add-uris/types"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"net/http"
	"os"
	"net/url"
)

func GetNewArray(oldCallbacks []string, callbacksToAdd []*url.URL, maxSize int) []string {
	initialArray := make([]string, maxSize)
	index := 0

	// copy the callbacks to add at the beginning of the new array
	for ; index < len(callbacksToAdd) && index < maxSize; index++ {
		initialArray[index] = callbacksToAdd[index].String()
	}

	i := 0

	AddingOld: for index < maxSize && i < len(oldCallbacks){
		proposed := oldCallbacks[i]
		i++
		for j := 0;j < len(callbacksToAdd);j++ {
			if proposed == callbacksToAdd[j].String() {
				continue AddingOld
			}
		}

		initialArray[index] = proposed
		index++
	}

	newArray := make([]string, index)

	for i = 0; i < index; i++ {
		newArray[i] = initialArray[i]
	}

	return newArray
}

func AddUris(config types.ApplicationConfig) {
	var errors []error
	var res *http.Response
	var body []byte
	request := gorequest.New()
	logger := config.Logger

	createTokenBodyRequest := types.RequestBodyCreateToken {
		ClientId: config.ClientId,
		ClientSecret: config.ClientSecret,
		GrantType: "client_credentials",
		Audience: fmt.Sprintf("https://%s/api/v2/", config.Domain),
	}

	var createTokenBodyResponse types.ResponseBodyCreateToken
	res, _, errors = request.
		Post(fmt.Sprintf("https://%s/oauth/token", config.Domain)).
		Send(&createTokenBodyRequest).
		EndStruct(&createTokenBodyResponse)

	if len(errors) > 0 {
		logger.Errorf("Unknown errors occurred when attempting to POST a new token %s", errors)
		os.Exit(1)
	}


	if res.StatusCode != 200 {
		logger.Errorf("Expected statusCode 200 when posting a new token: %d", res.StatusCode)
		os.Exit(1)
	}

	token := createTokenBodyResponse.AccessToken

	var clientDetailsResponse types.RequestBodyUpdateClient
	res, _, errors = request.
		Get(fmt.Sprintf(
			"https://%s/api/v2/clients/%s?fields=callbacks%2Callowed_logout_urls&include_fields=true",
			config.Domain, config.AppClientId)).
		Set("Authorization", fmt.Sprintf("Bearer %s", token)).
		EndStruct(&clientDetailsResponse)


	if len(errors) > 0 {
		logger.Errorf("Unknown errors occurred when attempting to GET the client details %s", errors)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		logger.Errorf("Expected statusCode 200 when fetching the client details: %d", res.StatusCode)
		os.Exit(1)
	}

	callbacks := GetNewArray(clientDetailsResponse.Callbacks, config.Callbacks, config.MaxLen)
	allowedLogoutUrls := GetNewArray(clientDetailsResponse.AllowedLogoutUrls, config.Logouts, config.MaxLen)

	//fmt.Println("Length of callbacks: ", len(callbacks))
	//for i := 0; i < len(callbacks); i++ {
	//	fmt.Println(callbacks[i]);
	//}
	//
	//fmt.Println("Length of allowedLogoutUrls: ", len(allowedLogoutUrls))
	//for i := 0; i < len(allowedLogoutUrls); i++ {
	//	fmt.Println(allowedLogoutUrls[i]);
	//}

	updateClientBodyRequest := types.RequestBodyUpdateClient {
		Callbacks: callbacks,
		AllowedLogoutUrls: allowedLogoutUrls,
	}

	var updateClientBodyResponse types.RequestBodyUpdateClient
	res, body, errors = request.
	Patch(fmt.Sprintf("https://%s/api/v2/clients/%s", config.Domain, config.AppClientId)).
		Set("Authorization", fmt.Sprintf("Bearer %s", token)).
		Send(&updateClientBodyRequest).
		EndStruct(&updateClientBodyResponse)

	if len(errors) > 0 {
		logger.Errorf("Unknown errors occurred when attempting to update the client %s", errors)
		os.Exit(1)
	}


	if res.StatusCode != 200 {
		logger.Errorf("Expected statusCode 200 when updating the client: %d", res.StatusCode)
		logger.Error(string(body))
		os.Exit(1)
	}

}
