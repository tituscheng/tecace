package response

const (
	SOMETHING_WENT_WRONG = "something went wrong"
	KEY_NOT_FOUND        = "key not found"
	OK                   = "OK"
)

type Response struct {
	Result      int    `json:"result"`
	Description string `json:"description"`
}

type ResponseData struct {
	Result int               `json:"result"`
	Data   map[string]string `json:"data"`
}
