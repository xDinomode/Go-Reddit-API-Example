package main

import "net/http"
import "golang.org/x/oauth2"
import "fmt"
import "bytes"

// Replace with your own clientid
var ClientID = "Your Client ID"

// Replace with your own clientsecret
var ClientSecret = "Your Client Secret"
var RedirectURL = "http://localhost:3000/redirect"
var authURL = "https://www.reddit.com/api/v1/authorize"
var tokenURL = "https://www.reddit.com/api/v1/access_token"
var redditConf *oauth2.Config

func redirect(res http.ResponseWriter, req *http.Request) {
	code := req.FormValue("code")

	if len(code) != 0 {
		accessToken, err := redditConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println(err)
			http.NotFound(res, req)
			return
		}

		if accessToken.Valid() {
			// Create the *http.Client
			httpClient := redditConf.Client(oauth2.NoContext, accessToken)

			// Add the query parameters to submit content to /r/test
			query := []byte("sr=test&title=Test&kind=self&text=My First Bot Post")

			// POST request to the API
			apiRequest, err := http.NewRequest("POST", "https://oauth.reddit.com/api/submit?", bytes.NewBuffer(query))
			// Reddit requires a custom User-Agent with every request
			apiRequest.Header.Add("User-Agent", "MyAwesomeApp5:v1.0 (by /u/YOUR_USERNAME)")
			if err != nil {
				fmt.Println(err)
				http.NotFound(res, req)
				return
			}

			// Do the apiRequest
			apiResponse, err := httpClient.Do(apiRequest)
			if err != nil {
				fmt.Println(err)
				http.NotFound(res, req)
				return
			}

			defer apiResponse.Body.Close()
			fmt.Println("Response Status:", apiResponse.Status)
			fmt.Println("Response Headers:", apiResponse.Header)
			res.Write([]byte("Post Submitted to /r/test!"))
		} else {
			http.NotFound(res, req)
		}
	} else {
		http.NotFound(res, req)
	}
}

func homePage(res http.ResponseWriter, req *http.Request) {
	url := redditConf.AuthCodeURL("CSRF")
	res.Write([]byte(fmt.Sprintf("<a href='%s'>Allow this app access</a>", url)))
}

func main() {

	redditConf = &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: RedirectURL,
		Scopes:      []string{"submit"},
	}

	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":3000", nil)
}
