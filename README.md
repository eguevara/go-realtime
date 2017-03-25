# go-realtime

Simple go client to query Google Analytics API.

 https://developers.google.com/apis-explorer/#p/analytics/v3/analytics.data.realtime.get


## Example

```

client := realtime.NewClient(realtime.WithHTTPClient(oauthClient))

opts := &realtime.Options{
    IDs:     "ga:56576851",
    Metrics: "rt:activeUsers",
}

resp, err := client.GetRealTime(opts)
if err != nil {
    log.Fatalf("error in calling resource %v", err)
}

```

## OAuth Client

Requires an OAUTH http client to make the client requests. 

```
key, err := ioutil.ReadFile("./real.pem")
if err != nil {
    log.Fatal("Error reading GA Service Account PEM key -", err)
}

conf := &jwt.Config{
    Email:      "readonly@mm.iam.gserviceaccount.com",
    PrivateKey: key,
    Scopes: []string{
        "https://www.googleapis.com/auth/analytics.readonly",
    },
    TokenURL: google.JWTTokenURL,
}
client := conf.Client(oauth2.NoContext)

```
