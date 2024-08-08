package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type KanyeRestSvc struct {
	baseUrl string
}

type KanyeRestApiResponseJSON struct {
	Quote string 
}

func NewKanyeRestSvc(baseUrl string) *KanyeRestSvc {
	return &KanyeRestSvc{baseUrl: baseUrl}
}

func (s *KanyeRestSvc) FetchQuote() (*KanyeRestApiResponseJSON, error) {
	resp, err := http.Get(s.baseUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var data KanyeRestApiResponseJSON

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &data, nil
}
