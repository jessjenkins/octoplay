package octopus

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	url    *url.URL
	client http.Client
}

type rt string

func New(apiurl, apiKey string) (*API, error) {
	u, err := url.Parse(apiurl)
	if err != nil {
		return nil, err
	}

	return &API{
		url: u,
		client: http.Client{
			Transport: rt(apiKey),
			Timeout:   time.Minute,
		},
	}, nil
}

type RequestOptions struct {
	PeriodFrom *time.Time
	PeriodTo   *time.Time
	Page       *int
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
	mpUrl := api.url.JoinPath("v1/electricity-meter-points/", mpan)

	err := api.get(mpUrl.String(), mp)
	if err != nil {
		return nil, err
	}

	return mp, nil
}

func (r rt) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(string(r), "")
	return http.DefaultTransport.RoundTrip(req)
}

func (api *API) get(url string, r interface{}) error {
	res, err := api.client.Get(url)
	if err != nil {
		return err
	}

	//TODO check response status, assuming success for now

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	return nil
}
