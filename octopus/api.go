package octopus

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const api_url = "https://api.octopus.energy/"

type API struct {
	OctopusAPIKey string
}

func (api *API) GetGSP(mpan string) (string, error) {
	mpResponse, err := api.GetMeterPoint(mpan)
	if err != nil {
		return "", err
	}
	return mpResponse.GSP, nil
}

type MeterPoint struct {
	GSP string `json:"gsp"`
}

func (api *API) GetMeterPoint(mpan string) (*MeterPoint, error) {
	mp := &MeterPoint{}

	mpUrl, err := url.JoinPath(api_url, "v1/electricity-meter-points/", mpan)
	if err != nil {
		return nil, err
	}

	client := api.getClient()
	res, err := client.Get(mpUrl)
	if err != nil {
		return nil, err
	}

	//TODO check response status, assuming success for now

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, mp)
	if err != nil {
		return nil, err
	}

	return mp, nil
}

func (api *API) getClient() http.Client {
	return http.Client{
		Transport: api,
		Timeout:   time.Minute,
	}
}

func (api *API) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(api.OctopusAPIKey, "")
	return http.DefaultTransport.RoundTrip(req)
}
