package workerman_client

func GetResponseError(err error) *ResponseError {
	respErr, _ := err.(*ResponseError)
	return respErr
}

func Params(params ...interface{}) []interface{} {
	result := make([]interface{}, 0, len(params))
	for _, param := range params {
		result = append(result, param)
	}
	return result
}
