package sesame

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs-client/crypto"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type Client struct {
	endpointURL *url.URL
	httpClient  *http.Client
}

func NewClient(address string) (*Client, error) {
	endpointURL, err := url.Parse("http://" + address)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := &Client{
		endpointURL: endpointURL,
		httpClient:  httpClient,
	}
	return client, nil
}

func (client *Client) getAPIMethod(path string) string {
	endpointURL, _ := url.Parse(client.endpointURL.String())
	endpointURL.Path = filepath.Join(endpointURL.Path, path)
	return endpointURL.String()
}

func getBody(resp *http.Response) string {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "--"
	}
	if len(data) == 0 {
		return "--"
	}
	text := strings.Trim(string(data), " \n")
	if len(text) == 0 {
		return "--"
	}
	return text
}

func (client *Client) SendKeyPhrase(phrase string) error {
	data, err := json.Marshal(phrase)
	if err != nil {
		return err
	}
	r := bytes.NewReader(data)
	resp, err := client.httpClient.Post(client.getAPIMethod("/key"), "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("sesameFS API: %s (%s)", resp.Status, getBody(resp)))
	}
	return nil
}

func (client *Client) SendCertificate(n uint8, kesFilePath string, vrfFilePath string, certFilePath string) error {
	// read operational certificate
	pkg, err := read(kesFilePath, vrfFilePath, certFilePath)
	if err != nil {
		return err
	}
	packageData, err := pkg.Bytes()
	if err != nil {
		return err
	}
	// generate encrypted pads
	keys, pads, err := crypto.GeneratePads(int(n), packageData)
	if err != nil {
		return err
	}
	for i := 0; i < len(keys); i++ {
		fmt.Printf("#%d: %s\n", i+1, keys[i])
	}
	// send the encrypted pads
	padsData, err := json.Marshal(pads)
	if err != nil {
		return err
	}
	r := bytes.NewReader(padsData)
	resp, err := client.httpClient.Post(client.getAPIMethod("/pads"), "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("sesameFS API: %s (%s)", resp.Status, getBody(resp)))
	}
	return nil
}
