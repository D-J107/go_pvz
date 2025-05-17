package rest

type PvzCreationRequest struct {
	Id               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

type PvzCreationResponse struct {
	Id               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

type GetAllFilterResponse struct {
	Pvzs []PvzsResponseStruct
}

type PvzsResponseStruct struct {
	PvzInfo       PvzInfoResponseStruct          `json:"pvz"`
	ReceptionInfo []ReceptionsInfoResponseStruct `json:"receptions"`
}

type PvzInfoResponseStruct struct {
	Id               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

type ReceptionsInfoResponseStruct struct {
	Reception ReceptionInfoResponseStruct  `json:"reception"`
	Products  []ProductsInfoResponseStruct `json:"products"`
}

type ReceptionInfoResponseStruct struct {
	Id       string `json:"id"`
	DateTime string `json:"dateTime"`
	PvzId    string `json:"pvzId"`
	Status   string `json:"status"`
}

type ProductsInfoResponseStruct struct {
	Id          string `json:"id"`
	DateTime    string `json:"dateTime"`
	Type        string `json:"type"`
	ReceptionId string `json:"receptionId"`
}
