package rest

type ReceptionCreationRequest struct {
	PvzId string `json:"pvzId"`
}

type ReceptionCreationResponse struct {
	ReceptionId string `json:"id"`
	DateTime    string `json:"dateTime"`
	PvzId       string `json:"pvzId"`
	Status      string `json:"status"`
}

type ReceptionCloseRequest struct {
	PvzId string `json:"pvzId"`
}

type ReceptionCloseResponse struct {
	ReceptionId string `json:"id"`
	DateTime    string `json:"dateTime"`
	PvzId       string `json:"pvzId"`
	Status      string `json:"status"`
}
