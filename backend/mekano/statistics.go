package mekano

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type StatisticInterface interface {
	SetFile(filename string)
	Payment(mekanoData []Mekano, initRange int, lastRange int) PaymentStatistics
	Billing(mekanoData []Mekano) BillingStatistics
}

type statistics struct {
	db   DatabaseInterface
	file string
}

func NewStatistics(db DatabaseInterface) StatisticInterface {
	return &statistics{
		db: db,
	}
}

func (s *statistics) SetFile(filename string) {
	s.file = filename
}

func (s statistics) Payment(mekanoData []Mekano, initRange int, lastRange int) PaymentStatistics {
	var efectivo, bancolombia, davivienda, susuerte, payU, total int = 0, 0, 0, 0, 0, 0

	for _, d := range mekanoData {
		debito, err := strconv.Atoi(d.Debito)
		total += debito
		if err != nil {
			log.Println(err)
		}
		switch d.Cuenta {
		case "11050501": //Efectivo
			efectivo += debito
		case "11200501": //Bancolombia
			bancolombia += debito
		case "11200510": //Davivienda
			davivienda += debito
		case "13452505": //Susuerte
			susuerte += debito
		case "13452501": //Pay U
			payU += debito
		}
	}

	statistic := PaymentStatistics{
		FileName:    s.file,
		RangoRC:     fmt.Sprintf("%d-%d", initRange+1, lastRange),
		Efectivo:    efectivo,
		Bancolombia: bancolombia,
		Davivienda:  davivienda,
		PayU:        payU,
		Susuerte:    susuerte,
		Total:       total,
	}
	_, err := s.db.SavePayment(Payment{Consecutive: lastRange, CreateAt: time.Now(), File: s.file})
	if err != nil {
		log.Println(err)
	}
	return statistic
}

func (s statistics) Billing(mekanoData []Mekano) BillingStatistics {
	var d, c, b float64 = 0, 0, 0
	for _, row := range mekanoData {
		debito, _ := strconv.ParseFloat(row.Debito, 64)
		d += debito
		credito, _ := strconv.ParseFloat(row.Credito, 64)
		c += credito
		base, _ := strconv.ParseFloat(row.Base, 64)
		b += base
	}

	bs := BillingStatistics{
		File:    s.file,
		Debito:  d,
		Credito: c,
		Base:    b,
	}

	_, err := s.db.SaveBilling(Billing{File: s.file, Base: int(b), Debit: int(d), Credit: int(c), CreateAt: time.Now()})
	if err != nil {
		log.Println(err)
	}

	return bs
}
