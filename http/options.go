package http

func Options(
	client Client,
	url string,
) (Result, error) {
	return Request(client, "OPTIONS", url, nil, nil)
}
