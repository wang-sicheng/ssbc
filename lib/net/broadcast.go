package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	cfsslapi "github.com/cloudflare/cfssl/api"
	"github.com/cloudflare/cfssl/log"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"

)

type Client struct {

	httpClient *http.Client
	Url string
}

type Req struct {

}

var Urls []string = []string{"http://127.0.0.1:8000","http://127.0.0.1:8000","http://127.0.0.1:8000","http://127.0.0.1:8000"}

func Broadcast(s string,reqBody []byte)error{

	for _,k := range Urls{
		client := &Client{
			Url: k,
		}

		err :=client.initHTTPClient()
		if err != nil{
			log.Infof("broadcast fail:%s",k)
			return err
		}
		endPoint := k + "/" + s

		req,err := NewPost(endPoint, reqBody)
		if err != nil{
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8;")

		err = client.SendReq(req, nil)

	}
	return nil
}


func NewPost(endPoint string, reqBody []byte)(*http.Request, error){
	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed posting to %s", endPoint)
	}
	return req, nil

}

func (c *Client)SendReq(req *http.Request, result interface{}) (err error) {
	reqStr := "test"

	log.Info("Sending request\n")


	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "%s failure of request", req.Method)
	}
	var respBody []byte
	if resp.Body != nil {
		respBody, err = ioutil.ReadAll(resp.Body)
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Debugf("Failed to close the response body: %s", err.Error())
			}
		}()
		if err != nil {
			return errors.Wrapf(err, "Failed to read response of request: %s", reqStr)
		}
		log.Info("Received response\n")
	}
	log.Info(string(respBody))
	var body *cfsslapi.Response
	if respBody != nil && len(respBody) > 0 {
		body = new(cfsslapi.Response)
		err = json.Unmarshal(respBody, body)
		if err != nil {
			return errors.Wrapf(err, "Failed to parse response: %s", respBody)
		}
		if len(body.Errors) > 0 {
			var errorMsg string
			for _, err := range body.Errors {
				msg := fmt.Sprintf("Response from server: Error Code: %d - %s\n", err.Code, err.Message)
				if errorMsg == "" {
					errorMsg = msg
				} else {
					errorMsg = errorMsg + fmt.Sprintf("\n%s", msg)
				}
			}
			return errors.Errorf(errorMsg)
		}
	}
	scode := resp.StatusCode
	if scode >= 400 {
		return errors.Errorf("Failed with server status code %d for request:\n%s", scode, reqStr)
	}
	if body == nil {
		return errors.Errorf("Empty response body:\n%s", reqStr)
	}
	if !body.Success {
		return errors.Errorf("Server returned failure for request:\n%s", reqStr)
	}
	log.Debugf("Response body result: %+v", body.Result)
	if result != nil {
		return mapstructure.Decode(body.Result, result)
	}
	return nil
}

func (c *Client) getURL(endpoint string) (string, error) {
	nurl := "127.0.0.1:8000"

	rtn := fmt.Sprintf("%s/%s", nurl, endpoint)
	return rtn, nil
}

func (c *Client) initHTTPClient() error {
	tr := new(http.Transport)

	c.httpClient = &http.Client{Transport: tr}
	return nil
}
