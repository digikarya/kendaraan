package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

func Curl(method, url string,r *http.Request, payload interface{}) (int, []byte, error) {
	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return 0, nil, errors.New("Service gagal")
	}
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest(method,url, body)
	if err != nil {
		return 0, nil, errors.New("Service gagal")
	}
	req.Header.Set("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	req.Header.Set("Authorization", reqToken)
	//req.Header.Set("CITILINK-Signature", secret[1])
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: time.Second * 5,Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal("Error reading response. ", err)
		return 0, nil, errors.New("Service gagal")
	}
	defer resp.Body.Close()
	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, errors.New("Service gagal")
	}
	return resp.StatusCode,bodyResponse,nil
}
