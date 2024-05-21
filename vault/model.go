package vault

type Token struct {
	Token string `json:"-"`
	//TokenId string `json:"token_id"`
	Path string `json:"path"`
}
