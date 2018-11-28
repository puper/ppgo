package workerman_client

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (this *ResponseError) Error() string {
	return this.Message
}

type Response struct {
	*ResponseError
	Data interface{} `json:"data"`
}
