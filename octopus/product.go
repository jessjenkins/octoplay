package octopus

import (
	"time"
)

type Product struct {
	Code                 string                   `json:"code"`
	DisplayName          string                   `json:"display_name"`
	SingleRegElecTariffs map[string]ProductTariff `json:"single_register_electricity_tariffs"`
}

type ProductTariff struct {
	DDMonthly ProductTariffDetails `json:"direct_debit_monthly"`
}

type ProductTariffDetails struct {
	Code string `json:"code"`
}

func (api *API) GetProduct(code string, opt *RequestOptions) (*Product, error) {
	url := api.url.JoinPath("/v1/products/", code)

	if opt != nil && opt.PeriodFrom != nil {
		q := url.Query()
		q.Set("period_from", opt.PeriodFrom.Format(time.RFC3339))
		url.RawQuery = q.Encode()
	}

	product := &Product{}
	err := api.get(url.String(), product)
	if err != nil {
		return nil, err
	}

	return product, nil
}
