package controllers

import (
	"html/template"
	"net/http"
	"parkingapi/models"
)

type ViewData struct {
	CountEvents   int
	CountInvoices int
	CountParkings int
	CountCarTypes int
	CountNatTypes int
	CountPlaces   int
	CountTariffs  int
	CountAccounts int
}

//Dashbord shows billing main page
func Dashbord(w http.ResponseWriter, r *http.Request) {
	var countEvents = 0
	var countInvoices = 0
	var countParkings = 0
	var countCarTypes = 0
	var countNatTypes = 0
	var countPlaces = 0
	var countTariffs = 0
	var countAccounts = 0

	countEvents = models.GetCountObjects("parking_event")
	countInvoices = models.GetCountObjects("parking_invoice")
	countParkings = models.GetCountObjects("parking_lot")
	countCarTypes = models.GetCountObjects("parking_cartype")
	countNatTypes = models.GetCountObjects("parking_nationalitytype")
	countPlaces = models.GetCountObjects("parking_place")
	countTariffs = models.GetCountObjects("parking_tariff")
	countAccounts = models.GetCountObjects("parking_account")

	data := ViewData{
		CountEvents:   countEvents,
		CountInvoices: countInvoices,
		CountParkings: countParkings,
		CountCarTypes: countCarTypes,
		CountNatTypes: countNatTypes,
		CountPlaces:   countPlaces,
		CountTariffs:  countTariffs,
		CountAccounts: countAccounts,
	}
	tmpl, _ := template.ParseFiles("views/index.html")
	tmpl.Execute(w, data)
}

//Events shows billing events page
func Events(w http.ResponseWriter, r *http.Request) {
	var countEvents = 0
	var countInvoices = 0
	var countParkings = 0
	var countCarTypes = 0
	var countNatTypes = 0
	var countPlaces = 0
	var countTariffs = 0
	var countAccounts = 0

	countEvents = models.GetCountObjects("parking_event")
	countInvoices = models.GetCountObjects("parking_invoice")
	countParkings = models.GetCountObjects("parking_lot")
	countCarTypes = models.GetCountObjects("parking_cartype")
	countNatTypes = models.GetCountObjects("parking_nationalitytype")
	countPlaces = models.GetCountObjects("parking_place")
	countTariffs = models.GetCountObjects("parking_tariff")
	countAccounts = models.GetCountObjects("parking_account")

	data := ViewData{
		CountEvents:   countEvents,
		CountInvoices: countInvoices,
		CountParkings: countParkings,
		CountCarTypes: countCarTypes,
		CountNatTypes: countNatTypes,
		CountPlaces:   countPlaces,
		CountTariffs:  countTariffs,
		CountAccounts: countAccounts,
	}
	tmpl, _ := template.ParseFiles("views/events/events_list.html")
	tmpl.Execute(w, data)
}
