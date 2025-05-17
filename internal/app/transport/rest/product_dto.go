package rest

type CreateProductRequest struct {
	Type  string `json:"type"`
	PvzId string `json:"pvzId"`
}

type CreateProductResponse struct {
	Id          string `json:"id"`
	DateTime    string `json:"dateTime"`
	Type        string `json:"type"`
	ReceptionId string `json:"receptionId"`
}
