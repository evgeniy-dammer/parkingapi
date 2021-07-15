package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"parkingapi/models"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

//CreateNewParkingIn create new parkingIn
func CreateNewParkingIn(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p models.ParkingIn
	err := dec.Decode(&p)

	if err != nil {

		log.Println(err)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id := models.CreateNewParkingIn(p)

	if id != "" {
		log.Println("ParkingIn inserted successfully!")
	} else {
		log.Println("Error inserting parkingIn!")
	}
}

//GetNewEvent returns new parkingIn event
func GetNewEvent(w http.ResponseWriter, r *http.Request) {
	rows := models.GetNewParkingInEventsList()
	c := models.NewParkingInEvent{}

	for rows.Next() {

		err := rows.Scan(&c.ID, &c.Tariff, &c.Entry, &c.Carplate, &c.IsShown)

		if err != nil {
			fmt.Fprint(w, err)
			continue
		}
	}

	jsonData, err := json.Marshal(c)

	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(jsonData))

	defer rows.Close()
}

//Test test
func Test(w http.ResponseWriter, r *http.Request) {
	log.Println("TEST")
}

//DeniedParkingIn denied new parkingIn
func DeniedParkingIn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	denied := models.DeniedNewParkingInEvent(id)

	count, err := denied.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count != 0 {
		log.Println("ParkingIn denied successfully!")
	} else {
		log.Println("Error deniing parkingIn!")
	}
}

//AcceptedParkingIn accepted new parkingIn
func AcceptedParkingIn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	accepted := models.AcceptedNewParkingInEvent(id)

	count, err := accepted.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count != 0 {
		log.Println("ParkingIn accepted successfully!")
	} else {
		log.Println("Error accepting parkingIn!")
	}
}

//GetAllEvents returns all parkingIn events
func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	var i = 0

	counter := models.GetCountAllEvents()
	rows := models.GetAllEventsList()
	object := make([]models.AllParkingInEvent, counter)

	for rows.Next() {
		c := models.AllParkingInEventForScan{}
		err := rows.Scan(&c.ID, &c.Tariff, &c.Entry, &c.Exit, &c.Sum, &c.Carplate, &c.IsShown, &c.IsAccepted, &c.IsDenied, &c.CarType, &c.NationalityType, &c.ParkingID, &c.IsClosed)

		if err != nil {
			fmt.Fprint(w, err)
			continue
		}

		if c.Exit.Valid {
			object[i].ID = c.ID
			object[i].Tariff = c.Tariff
			object[i].Entry = c.Entry
			object[i].Exit = c.Exit.String
			object[i].Sum = c.Sum
			object[i].Carplate = c.Carplate
			object[i].IsShown = c.IsShown
			object[i].IsAccepted = c.IsAccepted
			object[i].IsDenied = c.IsDenied
			object[i].CarType = c.CarType
			object[i].NationalityType = c.NationalityType
			object[i].ParkingID = c.ParkingID
			object[i].IsClosed = c.IsClosed
		} else {
			object[i].ID = c.ID
			object[i].Tariff = c.Tariff
			object[i].Entry = c.Entry
			object[i].Exit = ""
			object[i].Sum = c.Sum
			object[i].Carplate = c.Carplate
			object[i].IsShown = c.IsShown
			object[i].IsAccepted = c.IsAccepted
			object[i].IsDenied = c.IsDenied
			object[i].CarType = c.CarType
			object[i].NationalityType = c.NationalityType
			object[i].ParkingID = c.ParkingID
			object[i].IsClosed = c.IsClosed
		}

		i++
	}

	jsonData, err := json.MarshalIndent(object, " ", "  ")

	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(jsonData))

	defer rows.Close()

}

//CreateInvoice creates invoice
func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p models.ParkingInExitEvent
	err := dec.Decode(&p)

	if err != nil {

		log.Println(err)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id := models.CreateInvoice(p)

	if id != "" {
		log.Println("Invoice inserted successfully!")
		fmt.Fprint(w, string(id))
	} else {
		log.Println("Error inserting invoice!")
	}

}

//ParkingInExitEvent invoice payed
func ParkingInExitEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	invoicePayed := models.InvoicePayed(id)

	count, err := invoicePayed.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count != 0 {
		log.Println("Invoice payed successfully!")
	} else {
		log.Println("Error paying invoice!")
	}

	eventUpdated := models.UpdateEvent(id)

	count2, err2 := eventUpdated.RowsAffected()
	if err2 != nil {
		panic(err)
	}

	if count2 != 0 {
		log.Println("Event updated successfully!")
	} else {
		log.Println("Error event updating!")
	}
}

//GetBalance returns balance of account
func GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := models.GetBalance(id)
	c := models.Balance{}

	for row.Next() {
		err := row.Scan(&c.Balance)
		if err != nil {
			fmt.Fprint(w, err)
			continue
		}
	}

	jsonData, err := json.Marshal(c)

	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(jsonData))

	defer row.Close()
}


//SyncLocalInEvents inserts local parking in events
func SyncLocalInEvents(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p models.LocalParkingEvent
	err := dec.Decode(&p)

	if err != nil {

		log.Println(err)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id := models.SyncLocalInEvents(p)

	if id != "" {
		log.Println("Local ParkingIn event inserted successfully!")
		fmt.Fprint(w, "OK")
	} else {
		log.Println("Error inserting local ParkingIn event!")
		fmt.Fprint(w, "NO")
	}
}

//SyncLocalOutEvents inserts local parking out events
func SyncLocalOutEvents(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p models.LocalParkingEvent
	err := dec.Decode(&p)

	if err != nil {

		log.Println(err)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id := models.SyncLocalOutEvents(p)

	if id != "" {
		log.Println("Local ParkingOut event updating successfully!")
		fmt.Fprint(w, "OK")
	} else {
		log.Println("Error updating local ParkingOut event!")
		fmt.Fprint(w, "NO")
	}
}


//SyncLocalCreatedInvoices inserts local created invoices
func SyncLocalCreatedInvoices(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p models.LocalInvoice
	err := dec.Decode(&p)

	if err != nil {

		log.Println(err)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id := models.SyncLocalCreatedInvoices(p)

	if id != "" {
		log.Println("Local created invoice inserted successfully!")
		fmt.Fprint(w, "OK")
	} else {
		log.Println("Error inserting local created invoice!")
		fmt.Fprint(w, "NO")
	}
}

//SyncLocalPayedInvoices inserts local payed invoices
func SyncLocalPayedInvoices(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p models.LocalInvoice
	err := dec.Decode(&p)

	if err != nil {

		log.Println(err)

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id := models.SyncLocalPayedInvoices(p)

	if id != "" {
		log.Println("Local payed invoice updating successfully!")
		fmt.Fprint(w, "OK")
	} else {
		log.Println("Error updating local payed invoice!")
		fmt.Fprint(w, "NO")
	}
}


//SyncNotLocalInEvents returns not local parking in events
func SyncNotLocalInEvents(w http.ResponseWriter, r *http.Request) {
	var i = 0

	counter := models.GetCountNotLocalInEvents()
	rows := models.GetNotLocalInEventsList()
	object := make([]models.AllParkingInEvent, counter)

	for rows.Next() {
		c := models.AllParkingInEventForScan{}
		err := rows.Scan(&c.ID, &c.Tariff, &c.Entry, &c.Exit, &c.Sum, &c.Carplate, &c.IsShown, &c.IsAccepted, &c.IsDenied, &c.CarType, &c.NationalityType, &c.ParkingID, &c.IsClosed)

		if err != nil {
			fmt.Fprint(w, err)
			continue
		}

		if c.Exit.Valid {
			object[i].ID = c.ID
			object[i].Tariff = c.Tariff
			object[i].Entry = c.Entry
			object[i].Exit = c.Exit.String
			object[i].Sum = c.Sum
			object[i].Carplate = c.Carplate
			object[i].IsShown = c.IsShown
			object[i].IsAccepted = c.IsAccepted
			object[i].IsDenied = c.IsDenied
			object[i].CarType = c.CarType
			object[i].NationalityType = c.NationalityType
			object[i].ParkingID = c.ParkingID
			object[i].IsClosed = c.IsClosed
		} else {
			object[i].ID = c.ID
			object[i].Tariff = c.Tariff
			object[i].Entry = c.Entry
			object[i].Exit = ""
			object[i].Sum = c.Sum
			object[i].Carplate = c.Carplate
			object[i].IsShown = c.IsShown
			object[i].IsAccepted = c.IsAccepted
			object[i].IsDenied = c.IsDenied
			object[i].CarType = c.CarType
			object[i].NationalityType = c.NationalityType
			object[i].ParkingID = c.ParkingID
			object[i].IsClosed = c.IsClosed
		}

		i++
	}

	jsonData, err := json.MarshalIndent(object, " ", "  ")

	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(jsonData))

	defer rows.Close()

}

//SyncAcceptedNotLocalIn sync not local parking in event
func SyncNotLocalIn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	synced := models.SyncNotLocalIn(id)

	count, err := synced.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count != 0 {
		log.Println("Not local parking in event synced successfully!")
	} else {
		log.Println("Error syncing not local parking in event!")
	}
}

//SyncNotLocalOutEvents returns not local parking out events
func SyncNotLocalOutEvents(w http.ResponseWriter, r *http.Request) {
	var i = 0

	counter := models.GetCountNotLocalOutEvents()
	rows := models.GetNotLocalOutEventsList()
	object := make([]models.AllParkingInEvent, counter)

	for rows.Next() {
		c := models.AllParkingInEventForScan{}
		err := rows.Scan(&c.ID, &c.Tariff, &c.Entry, &c.Exit, &c.Sum, &c.Carplate, &c.IsShown, &c.IsAccepted, &c.IsDenied, &c.CarType, &c.NationalityType, &c.ParkingID, &c.IsClosed)

		if err != nil {
			fmt.Fprint(w, err)
			continue
		}

		if c.Exit.Valid {
			object[i].ID = c.ID
			object[i].Tariff = c.Tariff
			object[i].Entry = c.Entry
			object[i].Exit = c.Exit.String
			object[i].Sum = c.Sum
			object[i].Carplate = c.Carplate
			object[i].IsShown = c.IsShown
			object[i].IsAccepted = c.IsAccepted
			object[i].IsDenied = c.IsDenied
			object[i].CarType = c.CarType
			object[i].NationalityType = c.NationalityType
			object[i].ParkingID = c.ParkingID
			object[i].IsClosed = c.IsClosed
		} else {
			object[i].ID = c.ID
			object[i].Tariff = c.Tariff
			object[i].Entry = c.Entry
			object[i].Exit = ""
			object[i].Sum = c.Sum
			object[i].Carplate = c.Carplate
			object[i].IsShown = c.IsShown
			object[i].IsAccepted = c.IsAccepted
			object[i].IsDenied = c.IsDenied
			object[i].CarType = c.CarType
			object[i].NationalityType = c.NationalityType
			object[i].ParkingID = c.ParkingID
			object[i].IsClosed = c.IsClosed
		}

		i++
	}

	jsonData, err := json.MarshalIndent(object, " ", "  ")

	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(jsonData))

	defer rows.Close()

}

//SyncNotLocalOut sync not local parking Out event
func SyncNotLocalOut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	synced := models.SyncNotLocalOut(id)

	count, err := synced.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count != 0 {
		log.Println("Not local parking out event synced successfully!")
	} else {
		log.Println("Error syncing not local parking out event!")
	}
}

//SyncNotLocalCreatedInvoices returns not local created invoices
func SyncNotLocalCreatedInvoices(w http.ResponseWriter, r *http.Request) {
	var i = 0

	counter := models.GetCountNotLocalCreatedInvoices()
	rows := models.GetNotLocalCreatedInvoicesList()
	object := make([]models.NotLocalInvoice, counter)

	for rows.Next() {
		c := models.NotLocalInvoice{}
		err := rows.Scan(&c.ID, &c.DateTime, &c.CarPlate, &c.Sum, &c.EventID, &c.IsPayed, &c.IsLocalInvoice)

		if err != nil {
			fmt.Fprint(w, err)
			continue
		}

		object[i].ID = c.ID
		object[i].DateTime = c.DateTime
		object[i].CarPlate = c.CarPlate
		object[i].Sum = c.Sum
		object[i].EventID = c.EventID
		object[i].IsPayed = c.IsPayed
		object[i].IsLocalInvoice = c.IsLocalInvoice

		i++
	}

	jsonData, err := json.MarshalIndent(object, " ", "  ")

	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(jsonData))

	defer rows.Close()

}

//SyncNotLocalCreatedInvoice sync not local created invoice
func SyncNotLocalCreatedInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	synced := models.SyncNotLocalCreatedInvoices(id)

	count, err := synced.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count != 0 {
		log.Println("Not local created invoice synced successfully!")
	} else {
		log.Println("Error syncing not local creted invoice!")
	}
}

//SyncNotLocalPayedInvoices returns not local payed invoices
func SyncNotLocalPayedInvoices(w http.ResponseWriter, r *http.Request) {
	var i = 0

	counter := models.GetCountNotLocalPayedInvoices()
	rows := models.GetNotLocalPayedInvoicesList()
	object := make([]models.NotLocalInvoice, counter)

	for rows.Next() {
		c := models.NotLocalInvoice{}
		err := rows.Scan(&c.ID, &c.DateTime, &c.CarPlate, &c.Sum, &c.EventID, &c.IsPayed, &c.IsLocalInvoice)

		if err != nil {
			fmt.Fprint(w, err)
			continue
		}

		object[i].ID = c.ID
		object[i].DateTime = c.DateTime
		object[i].CarPlate = c.CarPlate
		object[i].Sum = c.Sum
		object[i].EventID = c.EventID
		object[i].IsPayed = c.IsPayed
		object[i].IsLocalInvoice = c.IsLocalInvoice

		i++
	}

	jsonData, err := json.MarshalIndent(object, " ", "  ")

	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, string(jsonData))

	defer rows.Close()

}

//SyncNotLocalPayedInvoice sync not local payed invoice
func SyncNotLocalPayedInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	synced := models.SyncNotLocalPayedInvoices(id)

	count, err := synced.RowsAffected()
	if err != nil {
		panic(err)
	}

	if count != 0 {
		log.Println("Not local payed invoice synced successfully!")
	} else {
		log.Println("Error syncing not local payed invoice!")
	}
}