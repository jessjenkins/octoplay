package octopus

import (
	"net/url"
	"strconv"
	"time"
)

type Consumption struct {
	Count    int                 `json:"count"`
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Results  []ConsumptionResult `json:"results"`
}

type ConsumptionResult struct {
	Consumption   float64   `json:"consumption"`
	IntervalStart time.Time `json:"interval_start"`
	IntervalEnd   time.Time `json:"interval_end"`
}

func (api *API) GetConsumption(mpan, serial string, opt *RequestOptions) (*Consumption, error) {
	url := api.url.JoinPath("/v1/electricity-meter-points/", mpan, "/meters/", serial, "/consumption")

	if opt != nil && opt.PeriodFrom != nil {
		q := url.Query()
		q.Set("period_from", opt.PeriodFrom.Format(time.RFC3339))
		url.RawQuery = q.Encode()
	}
	if opt != nil && opt.Page != nil {
		q := url.Query()
		q.Set("page", strconv.Itoa(*opt.Page))
		url.RawQuery = q.Encode()
	}

	consumption := &Consumption{Results: make([]ConsumptionResult, 0, 100)}
	err := api.get(url.String(), consumption)
	if err != nil {
		return nil, err
	}

	return consumption, nil
}

func (cn *Consumption) GetNextPage() *int {
	u, err := url.Parse(string(cn.Next))
	if err != nil {
		return nil
	}
	page, err := strconv.Atoi(u.Query().Get("page"))
	if err != nil {
		return nil
	}
	return &page
}
