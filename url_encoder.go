package main

import "net/url"

func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
