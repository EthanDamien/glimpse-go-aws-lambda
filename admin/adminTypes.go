package admin

// response format for any admin calls
type AdminResponse struct {
	DESC string `json:"body"`
	OK   bool   `json:"ok"`
}
