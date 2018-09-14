package infrastructure

import (
	"encoding/json"
	"net/http"
)

//HTTPClient wras http requests
type HTTPClient struct {
}

//GetJSON execute an HTTP request and parse the JSON
func (c HTTPClient) GetJSON(uri string, data *interface{}) error {
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Deserialize the json
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&data); err != nil {
		return err
	}
	return nil
}
