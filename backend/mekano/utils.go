package mekano

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Round(d string, b string, i string) (float64, float64, float64, error) {
	var debitFinal, baseFinal, ivaFinal float64
	debit, err := strconv.ParseFloat(d, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error to parse debit: %v", err)
	}
	_, decimalDebito := math.Modf(debit)
	if decimalDebito >= 0.5 {
		debitFinal = math.Ceil(debit)
	} else {
		debitFinal = math.Round(debit)
	}

	base, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error to parse base: %v", err)
	}
	_, decimalBase := math.Modf(base)
	if decimalBase >= 0.5 {
		baseFinal = math.Ceil(base)
	} else {
		baseFinal = math.Round(base)
	}

	iva, err := strconv.ParseFloat(strings.TrimSpace(i), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error to parse iva: %v", err)
	}
	_, decimalIva := math.Modf(iva)

	if decimalIva >= 0.5 {
		ivaFinal = math.Ceil(iva)
	} else {
		ivaFinal = math.Round(iva)
	}

	return debitFinal, baseFinal, ivaFinal, nil
}
