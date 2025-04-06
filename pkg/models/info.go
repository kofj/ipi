package models

type Info struct {
	IP       string    `json:"ip"`
	Locale   string    `json:"locale"`
	Location *Location `json:"location"`
	ASN      *ASN      `json:"asn"`
}

type ASN struct {
	Number       uint   `json:"number"`
	Organization string `json:"organization"`
}

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Timezone  string  `json:"timezone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
