package http

func Get(
	client Client,
	url string,
	headers *map[string]string,
) (Result, error) {
	return Request(client, "GET", url, headers, nil)
}
