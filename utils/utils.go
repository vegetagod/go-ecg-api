package utils

import (
	"net/http"
	"time"
)

func GetHttpClient() *http.Client {
	ht := &http.Transport{
		MaxIdleConns:          300,
		MaxConnsPerHost:       300,
		IdleConnTimeout:       15 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		ExpectContinueTimeout: 15 * time.Second,
	}

	return &http.Client{
		Transport: ht,
		Timeout:   5 * time.Minute,
	}

}
