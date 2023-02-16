package octopus

import (
	"net/url"
	"strconv"
	"time"
)

type UnitRates struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []UnitRatesResult `json:"results"`
}

type UnitRatesResult struct {
	ValueIncVAT float64    `json:"value_inc_vat"`
	ValidFrom   time.Time  `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
}

func (api *API) GetUnitRates(product, tariff string, opt *RequestOptions) (*UnitRates, error) {
	url := api.url.JoinPath("v1/products/", product, "electricity-tariffs/", tariff, "/standard-unit-rates/")

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

	unitRates := &UnitRates{Results: make([]UnitRatesResult, 0, 100)}
	err := api.get(url.String(), unitRates)
	if err != nil {
		return nil, err
	}

	return unitRates, nil
}

func (cn *UnitRates) GetNextPage() *int {
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
