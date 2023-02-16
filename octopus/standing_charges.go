package octopus

import (
	"net/url"
	"strconv"
	"time"
)

type StandingCharges struct {
	Count    int                     `json:"count"`
	Next     string                  `json:"next"`
	Previous string                  `json:"previous"`
	Results  []StandingChargesResult `json:"results"`
}

type StandingChargesResult struct {
	ValueIncVAT float64    `json:"value_inc_vat"`
	ValidFrom   time.Time  `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
}

func (api *API) GetStandingCharges(product, tariff string, opt *RequestOptions) (*StandingCharges, error) {
	url := api.url.JoinPath("v1/products/", product, "electricity-tariffs/", tariff, "/standing-charges/")

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

	standingCharges := &StandingCharges{Results: make([]StandingChargesResult, 0)}
	err := api.get(url.String(), standingCharges)
	if err != nil {
		return nil, err
	}

	return standingCharges, nil
}

func (cn *StandingCharges) GetNextPage() *int {
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
