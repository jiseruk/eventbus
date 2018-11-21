package client

import (
	"bytes"
	"errors"
	"fmt"
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
	message := []byte(`{"payload":{"test":true}, "topic":"test_topic"}`)

	resp, err := HTTPClient.Post(*url, "application/json", bytes.NewBuffer(message))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("BODY endpoint %s", string(body))
	if err != nil {
		return false, err
	}
	//io.Copy(ioutil.Discard, resp.Body)
	if resp.StatusCode >= 400 {
		return false, errors.New(string(body))
	}

	// handling error and doing stuff with body that needs to be unit tested
	return true, nil
}
