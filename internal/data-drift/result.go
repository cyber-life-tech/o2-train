package data_drift

type Result struct {
	Detected bool    `json:"detected"`
	Value    float64 `json:"value"`
}
