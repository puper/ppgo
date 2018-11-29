package workerman_client

func GetResponseError(err error) *ResponseError {
	respErr, _ := err.(*ResponseError)
	return respErr
}

func Params(params ...interface{}) []interface{} {
	return params
}
