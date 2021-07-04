package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/plgd-dev/cloud/grpc-gateway/pb"
	"github.com/plgd-dev/kit/codec/json"
)

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
type authTokenResponse struct {
	Code string `json:"code"`
}

func getAccessToken(cloudConfiguration *pb.ClientConfigurationResponse) string {
	resp, err := http.Get(cloudConfiguration.GetAccessTokenUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var accessTokenStruct accessTokenResponse
	err = json.Decode([]byte(body), &accessTokenStruct)
	if err != nil {
		log.Fatal(err)
	}
	return accessTokenStruct.AccessToken
}

func getAuthToken(cloudConfiguration *pb.ClientConfigurationResponse) string {
	resp, err := http.Get(cloudConfiguration.GetAuthCodeUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var authTokenStruct authTokenResponse
	err = json.Decode([]byte(body), &authTokenStruct)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(authTokenStruct)
	return authTokenStruct.Code
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // Disable check of certificate for the whole app. This is because it's local and certificates are auto generated

	var cloudConfiguration pb.ClientConfigurationResponse
	resp, err := http.Get("https://localhost/.well-known/ocfcloud-configuration")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Decode(body, &cloudConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	accessToken := getAccessToken(&cloudConfiguration)
	authToken := getAuthToken(&cloudConfiguration)
	log.Print(authToken)
	ocfClient := Ocfclient{}

	err = ocfClient.Initialize(accessToken, string(body))
	if err != nil {
		log.Fatal(err)
	}
	discover, err := ocfClient.Discover(10)
	log.Print(discover)

	if err != nil {
		log.Fatal(err)
	}
	if len(discover) != 1 {
		log.Fatalf("Number of devices found is not correct, %d", len(discover))
	}
	deviceId := discover[0].ID

	log.Println("Owning device...")

	result, err := ocfClient.OwnDevice(deviceId, accessToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
	log.Println("Setting ACLs...")
	err = ocfClient.SetAccessForCloud(deviceId)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Onboarding...")
	err = ocfClient.OnboardDevice(deviceId, authToken)
	if err != nil {
		log.Fatal(err)
	}
}
