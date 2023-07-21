package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	httpClient = http.Client{Timeout: 10 * time.Second}
)

const (
	QuotaNone      = "quota-non"
	QuotaFull      = "quota-r"
	QuotaAvailable = "quota-g"
)

func main() {
	fmt.Println(getSituation())
}

type Situation struct {
	Date     string `json:"date"`
	QuotaR   string `json:"quotaR"`
	OfficeId string `json:"officeId"`
	QuotaK   string `json:"quotaK"`
	OfficeDetail
}

type Office struct {
	OfficeId      string       `json:"officeId"`
	ChineseSimple OfficeDetail `json:"chs"`
}

type OfficeDetail struct {
	OfficeName string `json:"officeName"`
	Region     string `json:"region"`
	District   string `json:"district"`
}

func getSituation() ([]Situation, error) {
	req, err := http.NewRequest("GET", "https://eservices.es2.immd.gov.hk/surgecontrolgate/ticket/getSituation", nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var Resp struct {
		Data   []Situation `json:"data"`
		Office []Office    `json:"office"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&Resp); err != nil {
		return nil, err
	}
	offices := make(map[string]Office)
	for _, office := range Resp.Office {
		offices[office.OfficeId] = office
	}
	avaiable := make([]Situation, 0)
	for _, situation := range Resp.Data {
		if situation.QuotaK == QuotaAvailable {
			situation.OfficeDetail = offices[situation.OfficeId].ChineseSimple
			avaiable = append(avaiable, situation)
		}
	}
	return avaiable, nil
}
