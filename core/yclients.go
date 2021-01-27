package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"lyag/data"
	"lyag/data/models"
	"net/url"
	"time"
)

type YClients interface {
	AvailabilityFor(staffIds []string, daysDepth int) []Seance
	Book(staffId int, seance time.Time) error
	Staff(staffId int) (data.Staff, error)
}

type yClients struct {
	token string
	proxy *url.URL
	name  string
	phone string
	email string
}

func NewYClients(token string, proxy string, settings *Settings) YClients {
	var proxyURL *url.URL

	if proxy == "" {
		proxyURL = nil
	} else {
		proxyURL, _ = url.Parse(proxy)
	}

	api := yClients{
		token: token,
		proxy: proxyURL,
		name:  settings.Name,
		phone: settings.Phone,
		email: settings.Email,
	}

	return &api
}

func (client *yClients) AvailabilityFor(staffIds []string, daysDepth int) []Seance {
	loc, _ := time.LoadLocation("Europe/Moscow")
	today := time.Now().In(loc)

	seances := make([]Seance, 0)

	for _, staffId := range staffIds {
		endpoint := "https://n18056.yclients.com/api/v1/book_dates/104300?service_ids%5B%5D=1959975&staff_id=" + staffId

		availableDatesJson, err := data.FetchJson(client.token, endpoint, client.proxy)

		if err != nil {
			log.Printf("Error on receiving book dates: %s", err.Error())
			continue
		}

		var availableSeances models.BookDates
		err = json.Unmarshal(availableDatesJson, &availableSeances)

		if err != nil {
			log.Printf("Error on Unmarshal json: %s", err.Error())
			continue
		}

		log.Printf("Available dates for %s:", staffId)

		for i := 0; i < daysDepth; i++ {
			checkDate := today.AddDate(0, 0, i)
			for _, date := range availableSeances.BookingDates {
				if checkDate.Format("2006-01-02") == date {
					dateSeances := fetchSeancesOn(client.token, staffId, date)
					seances = append(seances, dateSeances...)
				}
			}
		}
	}

	return seances
}

func (client *yClients) Book(staffId int, seance time.Time) error {
	urlCheck := "https://n18056.yclients.com/api/v1/book_check/104300"
	urlRecord := "https://n18056.yclients.com/api/v1/book_record/104300"

	form := models.BookForm{
		Phone:       client.phone,
		Fullname:    client.name,
		Email:       client.email,
		Code:        nil,
		NotifyBySms: 0,
		Comment:     "Прошу напомнить за день (ну или за пару часов) :)",
		Appointments: models.Appointments{
			{
				ID: 0,
				Services: []int{
					1959975,
				},
				StaffID:      staffId,
				Events:       []interface{}{},
				Datetime:     seance.Format("2006-01-02T15:04:00"),
				ChargeStatus: "",
			},
		},
		IsMobile:            false,
		Referrer:            "https://core.ru",
		AppointmentsCharges: nil,
		IsSupportCharge:     false,
		RedirectPrefix:      "https://n18056.yclients.com/company:104300",
		BookformID:          18056,
	}

	var payload, err = json.Marshal(form)

	if err != nil {
		return err
	}

	checkHttpCode, checkResponse, err := data.PostJson(client.token, urlCheck, payload, client.proxy)

	if err != nil {
		return err
	}

	if checkHttpCode != 201 {
		return errors.New(fmt.Sprintf("error on check; %s", string(checkResponse)))
	}

	bookHttpCode, bookResponse, err := data.PostJson(client.token, urlRecord, payload, client.proxy)
	if bookHttpCode != 201 {
		return errors.New(fmt.Sprintf("error on book; %s", string(bookResponse)))
	}

	return nil
}

func (client *yClients) Staff(staffId int) (data.Staff, error) {
	panic("implement me")
}

func fetchSeancesOn(token string, staffId string, date string) []Seance {
	times := getDateTimes(token, staffId, date)

	dateSeances := make([]Seance, len(times))

	for j, seanceTime := range times {
		seance := Seance{
			StaffId: staffId,
			Service: 1959975,
			Time:    seanceTime,
		}
		dateSeances[j] = seance
	}

	return dateSeances
}

func getDateTimes(token string, staffId string, date string) []time.Time {
	endpoint := "https://n18056.yclients.com/api/v1/book_times/104300/" + staffId + "/" + date + "?service_ids%5B%5D=1959975"

	var p *url.URL
	if true {
		p = nil
	} else {
		p, _ = url.Parse("http://127.0.0.1:9090")
	}

	availableTimesJson, err := data.FetchJson(token, endpoint, p)

	var seances []time.Time

	if err != nil {
		log.Printf("Error on receiving book time: %s", err.Error())
		return seances
	}

	var availableTimes models.BookTimes

	jsonErr := json.Unmarshal(availableTimesJson, &availableTimes)

	if jsonErr != nil {
		println("Error while parsing " + date)
	} else {
		parsedDate, _ := time.Parse("2006-01-02", date)
		dayOfWeek := data.DaysOfWeek[parsedDate.Weekday()]

		println(fmt.Sprintf("%s (%s)", parsedDate.Format("02.01.2006"), dayOfWeek))

		seances = make([]time.Time, len(availableTimes))

		for i, seance := range availableTimes {
			parsed, err := time.Parse("2006-01-02T15:04:05-0700", seance.Datetime)
			if err != nil {
				log.Panicf("Invalid RFC3339 format date: %s", seance.Datetime)
			}

			seances[i] = parsed
		}
	}

	return seances
}
