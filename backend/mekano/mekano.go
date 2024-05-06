package mekano

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mozillazg/go-unidecode"
	"github.com/xuri/excelize/v2"
)

type MekanoInterface interface {
	ProcessPaymentFile(file string) (PaymentStatistics, error)
	ProcessBillFile(file string, extras string) (BillingStatistics, error)
}

type mekanoRepository struct {
	dr         DatabaseInterface
	statistics StatisticInterface
}

func NewMekano(dr DatabaseInterface) MekanoInterface {
	sta := NewStatistics(dr)
	return &mekanoRepository{
		dr,
		sta,
	}
}

func (mr *mekanoRepository) ProcessPaymentFile(file string) (PaymentStatistics, error) {
	mekanoSlice := make([]Mekano, 0)
	mr.statistics.SetFile(file)
	var consecutive = 0

	paymentFile, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return PaymentStatistics{}, err
	}

	rows, err := paymentFile.GetRows(paymentFile.GetSheetName(0))
	if err != nil {
		fmt.Println(err)
		return PaymentStatistics{}, err
	}

	c, err := mr.dr.GetPayment()
	if err != nil {
		return PaymentStatistics{}, err
	}

	for i, row := range rows[1:] {
		cashier, err := mr.dr.GetCashiers(row[9])
		if err != nil {
			return PaymentStatistics{}, err
		}
		consecutive = c.Consecutive + i + 1

		accountOne := Mekano{
			Tipo:          "RC",
			Prefijo:       "_",
			Numero:        strconv.Itoa(consecutive),
			Fecha:         row[4],
			Cuenta:        "13050501",
			Terceros:      row[1],
			CentroCostos:  "C1",
			Nota:          "RECAUDO POR VENTA SERVICIOS",
			Debito:        "0",
			Credito:       row[5],
			Base:          "0",
			Usuario:       "SUPERVISOR",
			NombreTercero: row[2],
			NombreCentro:  "CENTRO DE COSTOS GENERAL",
		}

		accountTwo := Mekano{
			Tipo:          "RC",
			Prefijo:       "_",
			Numero:        strconv.Itoa(consecutive),
			Fecha:         row[4],
			Cuenta:        cashier.Code,
			Terceros:      row[1],
			CentroCostos:  "C1",
			Nota:          "RECAUDO POR VENTA SERVICIOS",
			Debito:        row[5],
			Credito:       "0",
			Base:          "0",
			Usuario:       "SUPERVISOR",
			NombreTercero: row[2],
			NombreCentro:  "CENTRO DE COSTOS GENERAL",
		}
		mekanoSlice = append(mekanoSlice, accountOne, accountTwo)
	}

	exporterFile(mekanoSlice)
	data := mr.statistics.Payment(mekanoSlice, c.Consecutive, consecutive)
	return data, nil
}

func (mr *mekanoRepository) ProcessBillFile(file string, extras string) (BillingStatistics, error) {
	var itemIvaBaseFinal float64

	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return BillingStatistics{}, err
	}

	billingFile, err := xlsx.GetRows(xlsx.GetSheetName(0))
	if err != nil {
		fmt.Println(err)
		return BillingStatistics{}, err
	}

	ivaXlsx, err := excelize.OpenFile(extras)
	if err != nil {
		fmt.Println(err)
		return BillingStatistics{}, err
	}

	itemsIvaFile, err := ivaXlsx.GetRows(ivaXlsx.GetSheetName(0))
	if err != nil {
		log.Println(err, "itemsIvaFile")
		return BillingStatistics{}, err
	}

	var BillingDataSheet []Mekano

	for _, bRow := range billingFile[1:] {

		debit, base, iva, err := Round(bRow[14], bRow[12], bRow[13])
		if err != nil {
			return BillingStatistics{}, fmt.Errorf("error to parse values: %v", err)
		}

		if !strings.Contains(bRow[21], ",") {
			account, err := mr.dr.GetAccounts(bRow[21])
			if err != nil {
				return BillingStatistics{}, fmt.Errorf("error to get account %s", bRow[21])
			}

			costCenter, err := mr.dr.GetCostCenter(unidecode.Unidecode(bRow[17]))
			if err != nil {
				return BillingStatistics{}, fmt.Errorf("error to get cost center %s", bRow[17])
			}
			billingNormal := Mekano{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Fecha:         bRow[9],
				Cuenta:        account.Code,
				Terceros:      bRow[1],
				CentroCostos:  costCenter.Code,
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        "0",
				Credito:       fmt.Sprintf("%f", base),
				Base:          "0",
				Usuario:       "SUPERVISOR",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
			}

			billingIva := Mekano{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Fecha:         bRow[9],
				Cuenta:        "24080505",
				Terceros:      bRow[1],
				CentroCostos:  costCenter.Code,
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        "0",
				Credito:       fmt.Sprintf("%f", iva),
				Base:          fmt.Sprintf("%f", base),
				Usuario:       "SUPERVISOR",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
			}

			billingCxC := Mekano{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Fecha:         bRow[9],
				Cuenta:        "13050501",
				Terceros:      bRow[1],
				CentroCostos:  costCenter.Code,
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        fmt.Sprintf("%f", debit),
				Credito:       "0",
				Base:          "0",
				Usuario:       "SUPERVISOR",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
			}

			BillingDataSheet = append(BillingDataSheet, billingNormal, billingIva, billingCxC)
		} else {
			costCenter, err := mr.dr.GetCostCenter(unidecode.Unidecode(bRow[17]))
			if err != nil {
				return BillingStatistics{}, fmt.Errorf("error to get cost center %s", bRow[17])
			}
			splitBillingItems := strings.Split(bRow[21], ",")
			for _, item := range splitBillingItems {
				for _, itemIva := range itemsIvaFile[1:] {

					if itemIva[1] == strings.TrimSpace(item) && itemIva[0] == bRow[0] {
						itemIvaBase, _ := strconv.ParseFloat(itemIva[2], 64)
						_, decimalIvaBase := math.Modf(itemIvaBase)

						if decimalIvaBase >= 0.5 {
							itemIvaBaseFinal = math.Ceil(itemIvaBase)
						} else {
							itemIvaBaseFinal = math.Round(itemIvaBase)
						}

						account, err := mr.dr.GetAccounts(unidecode.Unidecode(strings.TrimSpace(item)))
						if err != nil {
							return BillingStatistics{}, fmt.Errorf("error to get account %s", strings.TrimSpace(item))

						}

						costCenter, err := mr.dr.GetCostCenter(unidecode.Unidecode(bRow[17]))
						if err != nil {
							return BillingStatistics{}, fmt.Errorf("error to get account %s", unidecode.Unidecode(bRow[17]))
						}

						billingNormalPlus := Mekano{
							Tipo:          "FVE",
							Prefijo:       "_",
							Numero:        bRow[8],
							Fecha:         bRow[9],
							Cuenta:        account.Code,
							Terceros:      bRow[1],
							CentroCostos:  costCenter.Code,
							Nota:          "FACTURA ELECTRÓNICA DE VENTA",
							Debito:        "0",
							Credito:       fmt.Sprintf("%f", itemIvaBaseFinal),
							Base:          "0",
							Usuario:       "SUPERVISOR",
							NombreTercero: bRow[2],
							NombreCentro:  bRow[17],
						}
						BillingDataSheet = append(BillingDataSheet, billingNormalPlus)
					}
				}
			}
			billingIvaPlus := Mekano{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Fecha:         bRow[9],
				Cuenta:        "24080505",
				Terceros:      bRow[1],
				CentroCostos:  costCenter.Code,
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        "0",
				Credito:       fmt.Sprintf("%f", iva),
				Base:          fmt.Sprintf("%f", base),
				Usuario:       "SUPERVISOR",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
			}

			billingCxCPlus := Mekano{
				Tipo:          "FVE",
				Prefijo:       "_",
				Numero:        bRow[8],
				Fecha:         bRow[9],
				Cuenta:        "13050501",
				Terceros:      bRow[1],
				CentroCostos:  costCenter.Code,
				Nota:          "FACTURA ELECTRÓNICA DE VENTA",
				Debito:        fmt.Sprintf("%f", debit),
				Credito:       "0",
				Base:          "0",
				Usuario:       "SUPERVISOR",
				NombreTercero: bRow[2],
				NombreCentro:  bRow[17],
			}

			BillingDataSheet = append(BillingDataSheet, billingIvaPlus, billingCxCPlus)
		}
	}

	exporterFile(BillingDataSheet)
	data := mr.statistics.Billing(BillingDataSheet)
	return data, nil
}

func exporterFile(mekanoData []Mekano) {
	f := excelize.NewFile()
	// Crea un nuevo sheet.
	index, _ := f.NewSheet("Sheet1")

	for i, m := range mekanoData {
		row := i + 2 // Comenzar en la fila 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), m.Tipo)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), m.Prefijo)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), m.Numero)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), m.Fecha)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), m.Cuenta)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), m.Terceros)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), m.CentroCostos)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), m.Nota)
		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), m.Debito)
		f.SetCellValue("Sheet1", fmt.Sprintf("J%d", row), m.Credito)
		f.SetCellValue("Sheet1", fmt.Sprintf("K%d", row), m.Base)
		f.SetCellValue("Sheet1", fmt.Sprintf("L%d", row), m.Usuario)
		f.SetCellValue("Sheet1", fmt.Sprintf("M%d", row), m.NombreTercero)
		f.SetCellValue("Sheet1", fmt.Sprintf("N%d", row), m.NombreCentro)
	}

	// Establece el sheet activo al primero
	f.SetActiveSheet(index)

	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	// Guarda el archivo de Excel
	if err := f.SaveAs(filepath.Join(dir, "CONTABLE.xlsx")); err != nil {
		fmt.Println(err)
	}
}
