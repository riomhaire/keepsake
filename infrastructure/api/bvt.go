package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (r *RestAPI) HandleBVT(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "plain/text")

	// Get HOST + Build URLS
	fmt.Fprintf(w, (string(r.Configuration.Test.BaseURL)))

	// Authenticate
	authenticateURI := fmt.Sprintf("%s/oauth/authorize", r.Configuration.Test.BaseURL)
	jsonStr := fmt.Sprintf("{\"grant_type\":\"client_credentials\",\"client_id\":\"%s\",\"client_secret\":\"%s\"}", r.Configuration.ClientCredentials[0].ClientID, r.Configuration.ClientCredentials[0].ClientSecret)
	fmt.Fprintf(w, (jsonStr))

	authenticateURIReq, err := http.NewRequest("POST", authenticateURI, bytes.NewBuffer([]byte(jsonStr)))
	authenticateURIReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(authenticateURIReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Convert json to map
	authenticateResponse := make(map[string]interface{})
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &authenticateResponse)
	fmt.Fprintf(w, "\n"+(string(body)+"\n"))
	rawaccesstoken := authenticateResponse["access_token"]
	if rawaccesstoken != nil {
		accesstoken := (string(rawaccesstoken.(string)))
		fmt.Fprintf(w, "\n"+accesstoken+"\n")

		// Verify
		verifyURI := fmt.Sprintf("%s/oauth/verify?token=%s", r.Configuration.Test.BaseURL, accesstoken)
		verifyRespose, err := http.Get(verifyURI)
		if err != nil {
			fmt.Fprintf(w, "\n"+(err.Error()+"\n"))
		}
		defer verifyRespose.Body.Close()

	}

}
