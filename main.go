package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Set values ...
const (
	zipCode     = "1431EX"                               // Set Adress zip
	houseNumber = "72"                                   // Set House number
	startDate   = "2026-01-01"                           // Calender start date
	endDate     = "2028-01-01"                           // Calender end date
	companyCode = "800bf8d7-6dd1-4490-ba9d-b419d6dc8a45" // meerlanden code NL
)

type CalendarResponse struct {
	Status      bool `json:"status"`
	LocationId  any
	AddressZip  string
	HouseNumber string
	Street      any
	City        any
	Community   any
	Data        []Box `json:"dataList"`
}

type LocId struct {
	Status bool          `json:"status"`
	Data   []map[any]any `json:"dataList"`
}

type Box struct {
	Type  string   `json:"_pickupTypeText"`
	Dates []string `json:"pickupDates"`
}

func postRequestClient(payload string, url string) ([]byte, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("language", "nl_NL")
	req.Header.Set("X-Api-Key", "3f09dfa3-7933-485a-b7bd-893c28a8a0b1")
	req.Header.Set("X-Api-Hash", "338907df711a4d424d33a1b114d865eb")
	req.Header.Set("X-Instance-ID", `"90718_800bf8d7-6dd1-4490-ba9d-b419d6dc8a45"`)
	req.Header.Set("X-App-Version", "3.5.1")
	req.Header.Set("X-Module-Name", "calendarGeneric")
	req.Header.Set("Origin", "https://www.meerlanden.nl")
	req.Header.Set("Referer", "https://www.meerlanden.nl/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en,en-US;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return output, nil

}

func main() {
	// LocationId request
	loc_url := "https://wasteprod2api.ximmio.com/api/GetAddress"
	loc_payload := fmt.Sprintf(`{"companyCode": "%s", "PostCode": "%s", "HouseNumber": "%s"}`, companyCode, zipCode, houseNumber)
	locationBody, err := postRequestClient(loc_payload, loc_url)
	if err != nil {
		log.Println(err)
	}
	loc_id := LocId{}
	err = json.Unmarshal(locationBody, &loc_id)
	if err != nil {
		log.Println(err)
	}

	if loc_id.Status != true {
		log.Fatal("Location Not found. Set Correct address..!")
	}

	uniqueAddressID := loc_id.Data[0]["UniqueId"]

	// Calendar Request
	url := "https://wasteprod2api.ximmio.com/api/GetCalendar"
	//
	calenderPayloadStr := fmt.Sprintf(`{"companyCode": "%s",
 		"uniqueAddressID": "%s",
 		"startDate": "%s",
 		"endDate": "%s"}`, companyCode, uniqueAddressID, startDate, endDate)
	calendarBody, err := postRequestClient(calenderPayloadStr, url)
	if err != nil {
		log.Println(err)
	}
	output := CalendarResponse{}
	output = CalendarResponse{AddressZip: zipCode, HouseNumber: houseNumber, LocationId: uniqueAddressID, Street: loc_id.Data[0]["Street"], City: loc_id.Data[0]["City"], Community: loc_id.Data[0]["Community"]}

	err = json.Unmarshal(calendarBody, &output)
	if err != nil {
		panic(err)
	}

	prettyjson, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(prettyjson))

}
