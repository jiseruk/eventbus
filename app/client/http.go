package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var HTTPClient *http.Client

func init() {
	//Check endpoint
	HTTPClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 3 * time.Second,
			}).Dial,
		},
		Timeout: 5 * time.Second,
	}
}

func CheckEndpoint(url *string) (bool, error) {
	if url == nil {
		return true, nil
	}
	message := struct {
		payload interface{}
		topic   string
	}{
		payload: map[string]string{},
		topic:   "test_endpoint",
	}
	data, _ := json.Marshal(&message)
	resp, err := HTTPClient.Post(*url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
	if resp.StatusCode >= 400 {
		return false, nil
	}

	// handling error and doing stuff with body that needs to be unit tested
	return true, nil
}
