package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/eguevara/go-realtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

func main() {

	oauthClient := getOAuthClient()

	client := realtime.NewClient(realtime.WithHTTPClient(oauthClient))

	opts := &realtime.Options{
		IDs:     "ga:56576851",
		Metrics: "rt:activeUsers",
	}

	resp, err := client.GetRealTime(opts)
	if err != nil {
		log.Fatalf("error in calling resource %v", err)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
}

func getOAuthClient() *http.Client {

	key, err := ioutil.ReadFile("./sample.pem")
	if err != nil {
		log.Fatal("Error reading GA Service Account PEM key -", err)
	}

	conf := &jwt.Config{
		Email:      "",
		PrivateKey: key,
		Scopes: []string{
			"https://www.googleapis.com/auth/analytics.readonly",
		},
		TokenURL: google.JWTTokenURL,
	}
	client := conf.Client(oauth2.NoContext)

	return client
}
