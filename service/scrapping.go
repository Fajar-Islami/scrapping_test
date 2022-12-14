package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Fajar-Islami/scrapping_test/model"
	"github.com/Fajar-Islami/scrapping_test/repository"
	"github.com/PuerkitoBio/goquery"
)

const TABLE_TIMEFORMAT = "02-01-2006 15:04 MST"
const TIMEFORMAT = "02 January 2006, 15:04 MST" // 04 Februari 2021, 10:22 WIB"

type TrackingService interface {
	GetDataTracking(uri string) (model.DataStruct, error)
}

type trackingServiceImpl struct {
	context       context.Context
	trackingRedis repository.RedisTrackingRepository
}

func NewTrackingService(context context.Context, trackingRedis repository.RedisTrackingRepository) TrackingService {
	return &trackingServiceImpl{
		context:       context,
		trackingRedis: trackingRedis,
	}
}

const basePrefix = "tracking:"
const timeExpire = 1

func (ts *trackingServiceImpl) GetDataTracking(uri string) (dataTracking model.DataStruct, err error) {
	timeStr := fmt.Sprint(timeExpire, "m")
	strKeys := ts.keyWithPrefix("GetDataTracking", timeStr)

	// Check if data already in redis
	data, err := ts.trackingRedis.GetTrackingByQueryCtx(ts.context, strKeys)
	if err != nil {
		log.Println(err)
	}

	// If data exist
	if data != nil {
		dataTracking = *data
		return dataTracking, nil
	}

	fmt.Println("Data not exist")

	res, err := http.Get(uri)
	if err != nil {
		log.Println("GetDataTracking err===", err)
		return
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		// log.Fatal(err)
		return
	}

	doc.Find("table:nth-child(4) tbody").Each(func(i int, tablehtml *goquery.Selection) {
		// Get Data Consignee
		tablehtml.Find("tr:last-child td:nth-child(2)").Each(func(_ int, rowhtml *goquery.Selection) {
			text := rowhtml.Text()
			arrStr := strings.Split(text, "|")
			name := strings.Split(arrStr[0], "[")[1]
			name = strings.TrimSpace(name)
			dataTracking.ReceivedBy = name
		})

		// Get Data Hsitories
		tablehtml.Find("tr").Each(func(i int, rows *goquery.Selection) {
			rows.Find("td").Each(func(i int, data *goquery.Selection) {
				var dataHistory model.HistoriesStruct

				if i == 0 {
					createdAt := data.Text()
					createdAt += " WIB"

					parseCreatedAt, err := time.Parse(TABLE_TIMEFORMAT, createdAt)
					if err != nil {
						// log.Fatal(err)
						return
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
				prependItem := []model.HistoriesStruct{dataHistory}
				dataTracking.Histories = append(prependItem, dataTracking.Histories...)
			})
		})
	})

	if err := ts.trackingRedis.SetTrackingCtx(ts.context, strKeys, timeExpire, &dataTracking); err != nil {
		log.Println(err)
	}

	return
}

func (ts *trackingServiceImpl) keyWithPrefix(subprefix, time string) string {
	return fmt.Sprintf("%s:%s:%s", basePrefix, subprefix, time)
}
