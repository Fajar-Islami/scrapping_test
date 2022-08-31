package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type (
	ResponJSON struct {
		Status StatusStruct `json:"status"`
		Data   []DataStruct `json:"data"`
	}

	StatusStruct struct {
		Code    string `json:"code"`
		Message string `json:"Message"`
	}

	DataStruct struct {
		ReceivedBy string            `json:"receivedBy"`
		Histories  []HistoriesStruct `json:"histories"`
	}

	HistoriesStruct struct {
		Description string `json:"description"`
		CreatedAt   string `json:"createdAt"`
		Formatted   struct {
			CreatedAt string `json:"createdAt"`
		} `json:"formatted"`
	}
)

const URI = "https://gist.githubusercontent.com/nubors/eecf5b8dc838d4e6cc9de9f7b5db236f/raw/d34e1823906d3ab36ccc2e687fcafedf3eacfac9/jne-awb.html"
const TABLE_TIMEFORMAT = "02-01-2006 15:04 MST"

const TIMEFORMAT = "02 January 2006, 15:04 MST" // 04 Februari 2021, 10:22 WIB"

func main() {
	res := ScrappingURI(URI)
	data, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(data))
}

func ScrappingURI(uri string) DataStruct {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dataTracking DataStruct

	doc.Find("table:nth-child(4) tbody").Each(func(i int, tablehtml *goquery.Selection) {
		// Get Data Consignee
		tablehtml.Find("tr:last-child td:nth-child(2)").Each(func(_ int, rowhtml *goquery.Selection) {
			text := rowhtml.Text()
			arrStr := strings.Split(text, "|")
			name := strings.Split(arrStr[0], "[")[1]
			dataTracking.ReceivedBy = name
		})

		// Get Data Hsitories
		tablehtml.Find("tr").Each(func(i int, rows *goquery.Selection) {
			rows.Find("td").Each(func(i int, data *goquery.Selection) {
				var dataHistory HistoriesStruct

				if i == 0 {
					createdAt := data.Text()
					createdAt += " WIB"
					fmt.Println("createdAt origin", createdAt)

					parseCreatedAt, err := time.Parse(TABLE_TIMEFORMAT, createdAt)
					if err != nil {
						log.Println(err)
					}

					createdAt = parseCreatedAt.Format(time.RFC3339)
					formattedCreatedAt := parseCreatedAt.Format(TIMEFORMAT)

					dataHistory.CreatedAt = createdAt
					dataHistory.Formatted.CreatedAt = formattedCreatedAt

				}
				if i == 1 {
					description := data.Text()
					dataHistory.Description = description
				}

				// prepend to histories tracking
				prependItem := []HistoriesStruct{dataHistory}
				dataTracking.Histories = append(prependItem, dataTracking.Histories...)
			})
		})
	})

	return dataTracking
}
