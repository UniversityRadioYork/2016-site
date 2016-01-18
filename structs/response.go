package structs

/**
 * Represents what we get back from the URY API
 */
type Response struct {
	Status  string      `json:"status"`
	Payload interface{} `json:"payload"`
	Time    float32     `json:"time,string"`
}
