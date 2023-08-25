package fetch

import (
	"fmt"
	"github.com/jessjenkins/octoplay/octopus"
	"log"
	"time"
)

const (
	existingProd = "VAR-20-10-01"
	agileProd    = "AGILE-FLEX-22-11-25"
)

type Compare struct {
	API               *octopus.API
	ElectricityMPAN   string
	ElectricitySerial string
}

type Result struct {
	GSP           string
	Usage         Usage
	ExistingSCs   Charges
	ExistingUnits Charges
	AgileSCs      Charges
	AgileUnits    Charges
}

func (c *Compare) Run() (*Result, error) {
	mpan := c.ElectricityMPAN
	serial := c.ElectricitySerial

	log.Println("Getting gsp")
	gsp, err := c.API.GetGSP(mpan)
	if err != nil {
		return nil, err
	}

	//from, _ := time.Parse(time.RFC3339, "2023-02-07T00:00:00Z")
	from, _ := time.Parse(time.RFC3339, "2023-02-08T00:00:00Z")
	//from, _ := time.Parse(time.RFC3339, "2023-02-15T00:00:00Z")

	log.Println("Getting all usage")
	usage, err := c.GetAllUsage(mpan, serial, from)
	if err != nil {
		return nil, err
	}

	log.Println("Getting existing product")
	existingProduct, err := c.API.GetProduct(existingProd, &octopus.RequestOptions{PeriodFrom: &from})
	if err != nil {
		return nil, err
	}
	// TODO error checking
	existingTariff := existingProduct.SingleRegElecTariffs[gsp].DDMonthly.Code

	log.Println("Getting agile product")
	agileProduct, err := c.API.GetProduct(agileProd, &octopus.RequestOptions{PeriodFrom: &from})
	if err != nil {
		return nil, err
	}
	agileTariff := agileProduct.SingleRegElecTariffs[gsp].DDMonthly.Code

	log.Println("Getting existing standing charges")
	existingSCs, err := c.GetAllSCs(existingProd, existingTariff, from)
	if err != nil {
		return nil, err
	}

	log.Println("Getting existing unit costs")
	existingUnits, err := c.GetAllUnits(existingProd, existingTariff, from)
	if err != nil {
		return nil, err
	}

	log.Println("Getting agile standing charges")
	agileSCs, err := c.GetAllSCs(agileProd, agileTariff, from)
	if err != nil {
		return nil, err
	}

	log.Println("Getting agile unit costs")
	agileUnits, err := c.GetAllUnits(agileProd, agileTariff, from)
	if err != nil {
		return nil, err
	}

	days := 0
	avgSaving := 0.0
	for day := from.Truncate(24 * time.Hour); day.Before(time.Now().AddDate(0, 0, -1)); day = day.AddDate(0, 0, 1) {
		existingCost := existingSCs.ChargeAt(day)
		agileCost := agileSCs.ChargeAt(day)

		date := day.Format("2006-01-02")
		dayUsage, ok := usage[date]
		if ok {
			for halfHour := 0; halfHour < 48; halfHour++ {
				hourTime := day.Add(30 * time.Minute * time.Duration(halfHour))
				consumption := dayUsage[hourTime.Format("15:04")]

				existingUC := existingUnits.ChargeAt(hourTime)
				existingCost += consumption * existingUC

				agileUC := agileUnits.ChargeAt(hourTime)
				agileCost += consumption * agileUC
			}
		}

		saving := agileCost - existingCost
		avgSaving = (avgSaving*(float64(days)) + saving) / float64(days+1)
		days++
		fmt.Printf("%s (%s): £ %.2f  - £ %.2f   - £ %.2f\n", date, day.Format("Mon"), existingCost/100, agileCost/100, (saving)/100)
	}
	fmt.Printf("avg daily saving: £ %.2f\n", avgSaving/100)
	fmt.Printf("avg annual saving: £ %.2f\n", avgSaving*365/100)

	return &Result{
		GSP:           gsp,
		Usage:         usage,
		ExistingSCs:   existingSCs,
		ExistingUnits: existingUnits,
		AgileSCs:      agileSCs,
		AgileUnits:    agileUnits,
	}, nil
}

func (c *Compare) GetAllUsage(mpan, serial string, from time.Time) (Usage, error) {

	usage := Usage{}

	cn, err := c.API.GetConsumption(mpan, serial, &octopus.RequestOptions{PeriodFrom: &from})
	if err != nil {
		return nil, err
	}
	AddToUsage(usage, cn)

	for cn.GetNextPage() != nil {
		cn, err = c.API.GetConsumption(mpan, serial, &octopus.RequestOptions{PeriodFrom: &from, Page: cn.GetNextPage()})
		if err != nil {
			return nil, err
		}
		AddToUsage(usage, cn)
	}
	return usage, nil
}

type Usage map[string]map[string]float64

func AddToUsage(usage Usage, cn *octopus.Consumption) {
	for _, r := range cn.Results {
		date := r.IntervalStart.Format("2006-01-02")
		time := r.IntervalStart.Format("15:04")

		if usage[date] == nil {
			usage[date] = make(map[string]float64)
		}
		usage[date][time] = r.Consumption
	}
}

type Charges map[time.Time]float64

func (c *Compare) GetAllSCs(product, tariff string, from time.Time) (Charges, error) {
	charges := Charges{}

	cn, err := c.API.GetStandingCharges(product, tariff, &octopus.RequestOptions{PeriodFrom: &from})
	if err != nil {
		return nil, err
	}
	AddToChargesSC(charges, cn)

	for cn.GetNextPage() != nil {
		cn, err = c.API.GetStandingCharges(product, tariff, &octopus.RequestOptions{PeriodFrom: &from, Page: cn.GetNextPage()})
		if err != nil {
			return nil, err
		}
		AddToChargesSC(charges, cn)
	}
	return charges, nil
}

func (c *Compare) GetAllUnits(product, tariff string, from time.Time) (Charges, error) {
	charges := Charges{}

	cn, err := c.API.GetUnitRates(product, tariff, &octopus.RequestOptions{PeriodFrom: &from})
	if err != nil {
		return nil, err
	}
	AddToChargesUnit(charges, cn)

	for cn.GetNextPage() != nil {
		cn, err = c.API.GetUnitRates(product, tariff, &octopus.RequestOptions{PeriodFrom: &from, Page: cn.GetNextPage()})
		if err != nil {
			return nil, err
		}
		AddToChargesUnit(charges, cn)
	}
	return charges, nil
}

func AddToChargesSC(charges Charges, cn *octopus.StandingCharges) {
	for _, r := range cn.Results {
		time := r.ValidFrom
		charges[time] = r.ValueIncVAT
	}
}

func AddToChargesUnit(charges Charges, cn *octopus.UnitRates) {
	for _, r := range cn.Results {
		time := r.ValidFrom
		charges[time] = r.ValueIncVAT
	}
}

func (c *Charges) ChargeAt(t time.Time) float64 {
	var latestTime *time.Time
	val := 0.0 // TODO will return 'free' if no figures so handle error condition properly
	for k, v := range *c {
		if !k.After(t) && (latestTime == nil || k.After(*latestTime)) {
			betterTime := k
			latestTime = &betterTime
			val = v
		}
	}
	return val
}
