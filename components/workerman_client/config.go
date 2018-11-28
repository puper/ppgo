package workerman_client

type Config struct {
	Addr        string `json:"addr"`
	Timeout     int    `json:"timeout"`
	SuccessCode int    `json:"success_code"`
}
