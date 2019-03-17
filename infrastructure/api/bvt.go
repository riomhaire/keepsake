package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var lock sync.Mutex

func (r *RestAPI) HandleBVT(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "plain/text")

	lock.Lock()
	token, err := bvtAuthenticate(r.Configuration.Test.BaseURL, r.Configuration.ClientCredentials[0].ClientID, r.Configuration.ClientCredentials[0].ClientSecret)
	if err != nil {
		log.Println("BVT > bvtAuthenticate failure", err.Error())
		handleSimpleError(w, http.StatusForbidden, err.Error())
		lock.Unlock()
		return
	}

	err = bvtVerify(r.Configuration.Test.BaseURL, token)
	if err != nil {
		log.Println("BVT > bvtVerify failure", err.Error())
		handleSimpleError(w, http.StatusUnauthorized, err.Error())
	} else {
		w.WriteHeader(http.StatusOK) // Success
	}
	lock.Unlock()
	return

}

func bvtAuthenticate(host, id, secret string) (accesstoken string, err error) {
	// Authenticate
	authenticateURI := fmt.Sprintf("%s/oauth/authorize", host)
	jsonStr := fmt.Sprintf("{\"grant_type\":\"client_credentials\",\"client_id\":\"%s\",\"client_secret\":\"%s\"}", id, secret)

	authenticateURIReq, err := http.NewRequest("POST", authenticateURI, bytes.NewBuffer([]byte(jsonStr)))
	authenticateURIReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(authenticateURIReq)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Check status code
	if resp.StatusCode != http.StatusOK {
		err = errors.New("Authenticate failed")
		return
	}
	authenticateResponse := make(map[string]interface{})
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &authenticateResponse)
	rawaccesstoken := authenticateResponse["access_token"]
	if rawaccesstoken != nil {
		accesstoken = (string(rawaccesstoken.(string)))
	}
	return
}

func bvtVerify(host, token string) (err error) {
	// Verify
	verifyURI := fmt.Sprintf("%s/oauth/verify?token=%s", host, token)
	verifyRespose, err := http.Get(verifyURI)
	if err != nil {
		return
	}
	defer verifyRespose.Body.Close()
	if verifyRespose.StatusCode != http.StatusOK {
		err = errors.New("Verify failed")
	}
	return
}
