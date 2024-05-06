package mekano

type Mekano struct {
	Tipo          string `json:"tipo"`
	Prefijo       string `json:"prefijo"`
	Numero        string `json:"numero"`
	Fecha         string `json:"fecha"`
	Cuenta        string `json:"cuenta"`
	Terceros      string `json:"terceros"`
	CentroCostos  string `json:"centro_costos"`
	Nota          string `json:"nota"`
	Debito        string `json:"debito"`
	Credito       string `json:"credito"`
	Base          string `json:"base"`
	Usuario       string `json:"usuario"`
	NombreTercero string `json:"nombre_tercero"`
	NombreCentro  string `json:"nombre_centro"`
}

type PaymentStatistics struct {
	FileName    string `json:"archivo"`
	RangoRC     string `json:"rango_rc"`
	Bancolombia int    `json:"bancolombia"`
	Davivienda  int    `json:"davivienda"`
	Susuerte    int    `json:"susuerte"`
	PayU        int    `json:"payu"`
	Efectivo    int    `json:"efectivo"`
	Total       int    `json:"total"`
}
type BillingStatistics struct {
	File    string  `json:"file"`
	Debito  float64 `json:"debito"`
	Credito float64 `json:"credito"`
	Base    float64 `json:"base"`
}
