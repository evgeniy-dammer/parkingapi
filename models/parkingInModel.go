package models

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

//Balance struct
type Balance struct {
	Balance string `json:"balance"`
}

//ParkingIn struct
type ParkingIn struct {
	ParkingID       string `json:"parkingId"`
	CarPlate        string `json:"carplate"`
	CarType         string `json:"carType"`
	NationalityType string `json:"nationalityType"`
}

//ParkingInExitEvent struct
type ParkingInExitEvent struct {
	Type         string `json:"type"`
	Carplate     string `json:"carplate"`
	Datetime     string `json:"datetime"`
	PayingMethod string `json:"payingmethod"`
}

//NewParkingInEvent struct
type NewParkingInEvent struct {
	ID       string `json:"id"`
	Tariff   string `json:"tariff"`
	Entry    string `json:"entry"`
	Carplate string `json:"carplate"`
	IsShown  bool   `json:"isshown"`
}

//AllParkingInEventForScan struct
type AllParkingInEventForScan struct {
	ID              string         `json:"id"`
	Tariff          string         `json:"tariff"`
	Entry           string         `json:"entry"`
	Exit            sql.NullString `json:"exit"`
	Sum             float32        `json:"sum"`
	Carplate        string         `json:"carplate"`
	IsShown         bool           `json:"isshown"`
	IsAccepted      bool           `json:"isaccepted"`
	IsDenied        bool           `json:"isdenied"`
	CarType         string         `json:"cartype"`
	NationalityType string         `json:"nationalitytype"`
	ParkingID       string         `json:"parkingid"`
	IsClosed        bool           `json:"isclosed"`
}

//ExitUpdateParkingIn struct
type ExitUpdateParkingIn struct {
	ID         string         `json:"id"`
	Tariff     string         `json:"tariff"`
	Entry      string         `json:"entry"`
	Exit       sql.NullString `json:"exit"`
	Sum        float32        `json:"sum"`
	Carplate   string         `json:"carplate"`
	IsShown    bool           `json:"isshown"`
	IsAccepted bool           `json:"isaccepted"`
	IsDenied   bool           `json:"isdenied"`
	PerEntry   float32        `json:"perentry"`
	PerHour    float32        `json:"perhour"`
	PerDay     float32        `json:"perday"`
	PerMonth   float32        `json:"permonth"`
}

//AllParkingInEvent struct
type AllParkingInEvent struct {
	ID              string  `json:"id"`
	Tariff          string  `json:"tariff"`
	Entry           string  `json:"entry"`
	Exit            string  `json:"exit"`
	Sum             float32 `json:"sum"`
	Carplate        string  `json:"carplate"`
	IsShown         bool    `json:"isshown"`
	IsAccepted      bool    `json:"isaccepted"`
	IsDenied        bool    `json:"isdenied"`
	CarType         string  `json:"cartype"`
	NationalityType string  `json:"nationalitytype"`
	ParkingID       string  `json:"parkingid"`
	IsClosed        bool    `json:"isclosed"`
}

//LocalParkingEvent struct
type LocalParkingEvent struct {
	ID              string  `json:"id"`
	Tariff          string  `json:"tariff"`
	Entry           string  `json:"entry"`
	Exit            string  `json:"exit"`
	Sum             float32 `json:"sum"`
	Carplate        string  `json:"carplate"`
	IsShown         string    `json:"isshown"`
	IsAccepted      string    `json:"isaccepted"`
	IsDenied        string    `json:"isdenied"`
	CarType         string  `json:"cartype"`
	NationalityType string  `json:"nationalitytype"`
	ParkingID       string  `json:"parkingid"`
	IsClosed        string    `json:"isclosed"`
}

//LocalParkingEvent struct
type LocalInvoice struct {
	ID              string  `json:"id"`
	DateTime          string  `json:"datetime"`
	CarPlate           string  `json:"carplate"`
	Sum             float32 `json:"sum"`
	EventID        string  `json:"eventid"`
	IsPayed         string    `json:"ispayed"`
}

//NotLocalInvoice struct
type NotLocalInvoice struct {
	ID              string  `json:"id"`
	DateTime          string  `json:"datetime"`
	CarPlate           string  `json:"carplate"`
	Sum             float32 `json:"sum"`
	EventID        string  `json:"eventid"`
	IsPayed         string    `json:"ispayed"`
	IsLocalInvoice         string    `json:"islocalinvoice"`
}

//CreateNewParkingIn inserts new parkingIn in database
func CreateNewParkingIn(parkingIn ParkingIn) string {
	var parkingInID = ""
	var tariffID = ""
	var parkingType = false

	checkParkinType := GetDB().QueryRow("SELECT isopen FROM parking_lot WHERE id = '" + parkingIn.ParkingID + "'")
	checkParkinType.Scan(&parkingType)

	checkTariff := GetDB().QueryRow("SELECT id FROM parking_tariff WHERE nationalitytype = '" + parkingIn.NationalityType + "' AND cartype = '" + parkingIn.CarType + "' AND isopen = " + strconv.FormatBool(parkingType))
	checkTariff.Scan(&tariffID)

	//log.Println(time.Now().Format("2006-02-01 15:04:05"));

	uuid := GenerateUUID()
	err := GetDB().QueryRow("INSERT INTO parking_event (id, tariff, entry, carplate, cartype, nationalitytype, parkingid) VALUES ('" +
		uuid + "', '" +
		tariffID + "', '" +
		time.Now().Format("2006-01-02 15:04:05") + "', '" +
		parkingIn.CarPlate + "', '" +
		parkingIn.CarType + "', '" +
		parkingIn.NationalityType + "', '" +
		parkingIn.ParkingID + "') returning id").Scan(&parkingInID)

	if err != nil {
		log.Println("Error parkingIn inserting!", err)
	}

	if parkingInID != "" {
		return parkingInID
	}
	return ""
}

//GenerateUUID generates new UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

//GetNewParkingInEventsList return list of new parkingIn events
func GetNewParkingInEventsList() *sql.Rows {
	rows, err := GetDB().Query("SELECT id, tariff, entry, carplate, isshown FROM parking_event WHERE isshown = false ORDER BY entry ASC LIMIT 1")

	if err != nil {
		panic(err)
	}

	return rows
}

//DeniedNewParkingInEvent denied new parkingIn event
func DeniedNewParkingInEvent(ID string) sql.Result {
	eventID, err := GetDB().Exec("UPDATE parking_event SET isdenied = true, isshown = true WHERE id = $1", ID)

	if err != nil {
		log.Println("Error event denied!", err)
	}

	return eventID
}

//AcceptedNewParkingInEvent accepted new parkingIn event
func AcceptedNewParkingInEvent(ID string) sql.Result {
	var notLocal = false
	var sum = 0.0
	var tariffID = ""

	checkIsNotLocalEvent := GetDB().QueryRow("SELECT islocalevent FROM parking_event WHERE id = '" + ID + "'")
	checkIsNotLocalEvent.Scan(&notLocal)

	checkTariff := GetDB().QueryRow("SELECT tariff FROM parking_event WHERE id = '" + ID + "'")
	checkTariff.Scan(&tariffID)

	checkSum := GetDB().QueryRow("SELECT entry FROM parking_tariff WHERE id = '" + tariffID + "'")
	checkSum.Scan(&sum)


	eventID, err := GetDB().Exec("UPDATE parking_event SET isaccepted = true, isshown = true, entry = '"+
		time.Now().Format("01-02-2006 15:04:05")+"', sum = "+fmt.Sprint(sum)+" WHERE id = $1", ID)

	if err != nil {
		log.Println("Error event accepted!", err)
	}

	count, err2 := eventID.RowsAffected()
	if err2 != nil {
		panic(err2)
	}

	if count != 0 {
		//_, err3 := GetDB().Exec("UPDATE parking_lot SET isbusy = true WHERE id = $1", lotUUID)

		/*if err3 != nil {
			log.Println("Error lot updating!", err3)
		}*/
	}

	return eventID
}

//GetCountAllEvents return count of all parkingIn events
func GetCountAllEvents() int {
	var counter int
	err := GetDB().QueryRow("SELECT count(*) FROM parking_event e WHERE isdenied = false").Scan(&counter)

	if err != nil {
		panic(err)
	}

	return counter
}

//GetAllEventsList return list of all parkingIn events
func GetAllEventsList() *sql.Rows {
	//rows, err := GetDB().Query("SELECT e.id, t.name as tariff, e.entry, e.exit, e.sum, e.carplate, e.isshown, e.isaccepted, e.isdenied " +
	//"FROM parking_event e " +
	//"LEFT JOIN parking_lot l ON e.lot = l.id " +
	//"LEFT JOIN parking_tariff t ON e.tariff = t.id " +
	//"WHERE e.isdenied = false " +
	//"GROUP BY e.id, t.name " +
	//"ORDER BY entry DESC"
	//)

	rows, err := GetDB().Query("SELECT id, tariff, entry, exit, sum, carplate, isshown, isaccepted, isdenied, cartype, nationalitytype, parkingid, isclosed " +
		"FROM parking_event WHERE islocalevent = false;")

	if err != nil {
		panic(err)
	}

	return rows
}

//UpdateParkingInExitEvent updates parkingIn with exit event
func UpdateParkingInExitEvent(parkingInExitEvent ParkingInExitEvent) sql.Result {
	var sum = 0.0
	var exitDatetime = ""

	rows := GetThisEvent(parkingInExitEvent.Carplate)
	c := ExitUpdateParkingIn{}
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.Tariff, &c.Entry, &c.Exit, &c.Sum, &c.Carplate, &c.IsShown, &c.IsAccepted, &c.IsDenied, &c.PerEntry, &c.PerHour, &c.PerDay, &c.PerMonth)

		if err != nil {
			log.Println("Error!", err)
			continue
		}
	}

	//timeFormat = "2006-01-02 15:04 MST"

	//then, err := time.Parse(timeFormat, v)
	//if err != nil {
	// ...
	//return
	//}

	//delta := time.Now().Sub(a)
	//fmt.Println(delta.Hours())

	if parkingInExitEvent.Type == "P" {
		var currentDateTime = time.Now()
		exitDatetime = currentDateTime.Format("2006-01-02 15:04:05")
	} else {
		exitDatetime = parkingInExitEvent.Datetime
	}

	locationID, err := GetDB().Exec("UPDATE parking_event SET exit = '"+exitDatetime+"', sum = "+fmt.Sprint(sum)+" WHERE carplate = $1 AND exit IS NULL", parkingInExitEvent.Carplate)

	if err != nil {
		log.Println("Error updating parkingOut!", err)
	}

	if err != nil {
		log.Println("Error parkingIn inserting!", err)
	}

	//if parkingInID != "" {
	//return parkingInID
	//}
	return locationID

}

//GetThisEvent returns special event
func GetThisEvent(Carplate string) *sql.Rows {
	rows, err := GetDB().Query("SELECT e.id, e.tariff, e.entry, e.exit, e.sum, e.carplate, e.isshown, e.isaccepted, e.isdenied, t.entry AS perentry, t.perhour AS perhour, t.perday AS perday, t.permounth AS permounth FROM parking_event e " +
		"LEFT JOIN parking_tariff t ON e.tariff = t.id " +
		"WHERE e.carplate = '" + Carplate + "' AND e.isaccepted = true AND e.exit IS NULL " +
		"GROUP BY e.id, t.entry, t.perhour, t.perday, t.permounth " +
		"ORDER BY e.entry DESC LIMIT 1")

	if err != nil {
		panic(err)
	}

	return rows
}

//CreateInvoice inserts new invoice in database
func CreateInvoice(parkingInExitEvent ParkingInExitEvent) string {
	var sum float32
	var invoiceID = ""
	var exitDatetime = ""
	var checkInvoice = ""

	rows := GetThisEvent(parkingInExitEvent.Carplate)
	c := ExitUpdateParkingIn{}
	for rows.Next() {
		err := rows.Scan(&c.ID, &c.Tariff, &c.Entry, &c.Exit, &c.Sum, &c.Carplate, &c.IsShown, &c.IsAccepted, &c.IsDenied, &c.PerEntry, &c.PerHour, &c.PerDay, &c.PerMonth)

		if err != nil {
			log.Println("Error!", err)
			continue
		}
	}

	checkUUID := GetDB().QueryRow("SELECT id FROM parking_invoice WHERE carplate = '" + parkingInExitEvent.Carplate + "' AND eventid = '" + c.ID + "' AND ispayed = false")
	checkUUID.Scan(&checkInvoice)

	if parkingInExitEvent.Type == "P" {
		var currentDateTime = time.Now()
		exitDatetime = currentDateTime.Format("2006-01-02 15:04:05")
	} else {
		exitDatetime = parkingInExitEvent.Datetime
	}

	sum = GetTotalSum(c.Entry, exitDatetime, c.PerEntry, c.PerHour, c.PerDay, c.PerMonth)

	if checkInvoice == "" {
		uuid := GenerateUUID()
		err := GetDB().QueryRow("INSERT INTO parking_invoice (id, datetime, carplate, sum, eventid) VALUES ('" +
			uuid + "', '" +
			time.Now().Format("2006-01-02 15:04:05") + "', '" +
			parkingInExitEvent.Carplate + "', " +
			fmt.Sprint(sum) + ", '" +
			c.ID + "') returning id").Scan(&invoiceID)

		if err != nil {
			log.Println("Error invoice inserting!", err)
		}

		if invoiceID != "" {

			if parkingInExitEvent.PayingMethod == "O" {

				var balance float32 = 0.0
				checkBalance := GetDB().QueryRow("SELECT balance FROM parking_account WHERE carplate = '" + parkingInExitEvent.Carplate + "'")
				checkBalance.Scan(&balance)

				balance = balance - sum

				accountBalance, err := GetDB().Exec("UPDATE parking_account SET balance = "+fmt.Sprint(balance)+" WHERE carplate = $1", parkingInExitEvent.Carplate)

				if err != nil {
					log.Println("Error balance updating!", err)
				}

				count, err2 := accountBalance.RowsAffected()
				if err2 != nil {
					panic(err2)
				}

				if count != 0 {
					return invoiceID + "/" + c.ID + "/" + fmt.Sprintf("%f", sum)
				}
			}

			return invoiceID + "/" + c.ID + "/" + fmt.Sprintf("%f", sum)
		}
		return ""
	} else {
		_, err := GetDB().Exec("UPDATE parking_invoice SET datetime = '"+time.Now().Format("2006-02-01 15:04:05")+
			"', sum = "+fmt.Sprint(sum)+
			" WHERE id = $1", checkInvoice)

		if err == nil {
			return checkInvoice + "/" + c.ID + "/" + fmt.Sprintf("%f", sum)
		}

		return ""
	}

}




//GetTotalSum count total sum
func GetTotalSum(EntryTime string, ExitTime string, PerEntry float32, PerHour float32, PerDay float32, PerMonth float32) float32 {
	var sum float32 = PerEntry

	EntryTime = strings.ReplaceAll(EntryTime, "Z", ".000Z")
	ExitTime = strings.ReplaceAll(ExitTime, " ", "T")
	ExitTime = ExitTime + ".000Z"

	time1 := convertToTimeObject(EntryTime)
	time2 := convertToTimeObject(ExitTime)

	var temp = math.Ceil(time2.Sub(time1).Hours())

	var duration = float32(temp)

	//log.Println("Duration ", duration)

	if duration > 7200 {
		var mod = duration / 7200
		sum += PerMonth * mod
		duration = duration - (7200 * mod)

		//log.Println("Duration per Month ", duration)
		//log.Println("summ per Month ", fmt.Sprint(sum))
	}

	if duration > 24 {
		var mod = duration / 24
		sum += PerDay * mod
		duration = duration - (24 * mod)
		//log.Println("Duration per day ", duration)
		//log.Println("summ per day ", fmt.Sprint(sum))
	}

	if duration > 0 {
		sum += PerHour * duration
		//log.Println("Duration per hour ", duration)
		//log.Println("summ per hour ", fmt.Sprint(sum))
	}

	return sum
}

//convertToTimeObject convert string to time object
func convertToTimeObject(dateStr string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, dateStr)

	if err != nil {
		log.Fatalf("error while parsing time: %s\n", err)
	}
	return t
}

func roundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

//InvoicePayed makes invoice payed
func InvoicePayed(ID string) sql.Result {
	invoiceID, err := GetDB().Exec("UPDATE parking_invoice SET ispayed = true WHERE id = $1", ID)

	if err != nil {
		log.Println("Error invoice paying!", err)
	}

	return invoiceID
}

//UpdateEvent returns event by invoice
func UpdateEvent(invoiceID string) sql.Result {
	var eventID = ""
	var sum = 0.0

	checkEvent := GetDB().QueryRow("SELECT eventid FROM parking_invoice WHERE id = '" + invoiceID + "'")
	checkEvent.Scan(&eventID)

	checkSum := GetDB().QueryRow("SELECT sum FROM parking_invoice WHERE id = '" + invoiceID + "'")
	checkSum.Scan(&sum)

	//log.Println(time.Now().Format("2006-01-02 15:04:05"))
	countEvent, err := GetDB().Exec("UPDATE parking_event SET exit = '"+time.Now().Format("2006-01-02 15:04:05")+"', sum = "+fmt.Sprint(sum)+", isclosed = true  WHERE id = $1", eventID)

	if err != nil {
		log.Println("Error invoice paying!", err)
	}

	return countEvent
}

//GetBalance returns balance of account
func GetBalance(Carplate string) *sql.Rows {
	rows, err := GetDB().Query("SELECT balance FROM parking_account WHERE carplate = '" + Carplate + "'")

	if err != nil {
		panic(err)
	}

	return rows
}

//SyncLocalInEvents inserts local parking in event in database
func SyncLocalInEvents(localParkingInEvent LocalParkingEvent) string {
	var insertID = ""
	var exit = ""

	if localParkingInEvent.Exit != "" {
		exit = "'" + localParkingInEvent.Exit + "', "
	}else{
		exit = "NULL, "
	}

	var sql = "INSERT INTO parking_event (id, tariff, entry, exit, sum, carplate, isshown, isaccepted, isdenied, cartype, nationalitytype, parkingid, isclosed, islocalevent) VALUES ('" +
		localParkingInEvent.ID + "', '" +
		localParkingInEvent.Tariff + "', '" +
		localParkingInEvent.Entry + "', " +
		exit +
		fmt.Sprint(localParkingInEvent.Sum) + ", '" +
		localParkingInEvent.Carplate + "', " +
		localParkingInEvent.IsShown + ", " +
		localParkingInEvent.IsAccepted + ", " +
		localParkingInEvent.IsDenied + ", '" +
		localParkingInEvent.CarType + "', '" +
		localParkingInEvent.NationalityType + "', '" +
		localParkingInEvent.ParkingID + "', " +
		localParkingInEvent.IsClosed + ", true) returning id"

		
	err := GetDB().QueryRow(sql).Scan(&insertID)
	//log.Println(insertID)
	if err != nil {
		log.Println("Error local ParkingIn event inserting!", err)
	}
	return insertID
}

//SyncLocalOutEvents updates local parking out event in database
func SyncLocalOutEvents(localParkingOutEvent LocalParkingEvent) string {
	updateCount, err := GetDB().Exec("UPDATE parking_event SET " +
		"tariff = '" +localParkingOutEvent.Tariff+ "', " + 
		"entry = '" +localParkingOutEvent.Entry+ "', " + 
		"exit = '" +localParkingOutEvent.Exit+ "', " + 
		"sum = '" +fmt.Sprint(localParkingOutEvent.Sum)+ "', " + 
		"carplate = '" +localParkingOutEvent.Carplate+ "', " + 
		"isshown = '" +localParkingOutEvent.IsShown+ "', " + 
		"isaccepted = '" +localParkingOutEvent.IsAccepted+ "', " + 
		"isdenied = '" +localParkingOutEvent.IsDenied+ "', " + 
		"cartype = '" +localParkingOutEvent.CarType+ "', " + 
		"nationalitytype = '" +localParkingOutEvent.NationalityType+ "', " + 
		"parkingid = '" +localParkingOutEvent.ParkingID+ "', " + 
		"isclosed = '" +localParkingOutEvent.IsClosed+ "', " + 
		"islocalevent = true " + 
		"WHERE id = $1", localParkingOutEvent.ID)

	if err != nil {
		log.Println("Error local ParkingOut event updating!", err)
	}

	count, err2 := updateCount.RowsAffected()
	if err2 != nil {
		panic(err2)
	}

	if count != 0 {
		return fmt.Sprint(count)
	}

	return ""
}

//SyncLocalCreatedInvoices inserts local created invoices in database
func SyncLocalCreatedInvoices(localCreatedInvoice LocalInvoice) string {
	var insertID = ""

	var sql = "INSERT INTO parking_invoice (id, datetime, carplate, sum, eventid, ispayed, islocalinvoice) VALUES ('" +
		localCreatedInvoice.ID + "', '" +
		localCreatedInvoice.DateTime + "', '" +
		localCreatedInvoice.CarPlate + "', " +
		fmt.Sprint(localCreatedInvoice.Sum) + ", '" +
		localCreatedInvoice.EventID + "', " +
		localCreatedInvoice.IsPayed + ", true) returning id"

		
	err := GetDB().QueryRow(sql).Scan(&insertID)

	if err != nil {
		log.Println("Error local created invoice inserting!", err)
	}
	return insertID
}

//SyncLocalPayedInvoices updates local payed invoice in database
func SyncLocalPayedInvoices(localPayedInvoice LocalInvoice) string {
	updateCount, err := GetDB().Exec("UPDATE parking_invoice SET " +
		"datetime = '" +localPayedInvoice.DateTime+ "', " + 
		"carplate = '" +localPayedInvoice.CarPlate+ "', " + 
		"sum = '" +fmt.Sprint(localPayedInvoice.Sum)+ "', " + 
		"eventid = '" +localPayedInvoice.EventID+ "', " + 
		"ispayed = '" +localPayedInvoice.IsPayed+ "', " + 
		"islocalinvoice = true " + 
		"WHERE id = $1", localPayedInvoice.ID)

	if err != nil {
		log.Println("Error local payed invoice updating!", err)
	}

	count, err2 := updateCount.RowsAffected()
	if err2 != nil {
		panic(err2)
	}

	if count != 0 {
		return fmt.Sprint(count)
	}

	return ""
}

//GetCountAllEvents return count of not local parking in events
func GetCountNotLocalInEvents() int {
	var counter int
	err := GetDB().QueryRow("SELECT count(*) FROM parking_event e WHERE isdenied = false AND islocalevent = false AND issyncedin = false").Scan(&counter) //

	if err != nil {
		panic(err)
	}

	return counter
}

//GetNotLocalInEventsList return list of not local parking in events
func GetNotLocalInEventsList() *sql.Rows {
	rows, err := GetDB().Query("SELECT id, tariff, entry, exit, sum, carplate, isshown, isaccepted, isdenied, cartype, nationalitytype, parkingid, isclosed " +
		"FROM parking_event WHERE isdenied = false AND islocalevent = false AND issyncedin = false;") //

	if err != nil {
		panic(err)
	}

	return rows
}

//SyncNotLocalIn sync not local parking in event
func SyncNotLocalIn(ID string) sql.Result {
	eventID, err := GetDB().Exec("UPDATE parking_event SET issyncedin = true WHERE id = $1", ID)

	if err != nil {
		log.Println("Error not local parking in event syncing!", err)
	}

	return eventID
}

//GetCountNotLocalOutEvents return count of not local parking out events
func GetCountNotLocalOutEvents() int {
	var counter int
	err := GetDB().QueryRow("SELECT count(*) FROM parking_event e WHERE isclosed = true AND islocalevent = false AND issyncedin = true AND issyncedout = false").Scan(&counter)

	if err != nil {
		panic(err)
	}

	return counter
}

//GetNotLocalOutEventsList return list of not local parking out events
func GetNotLocalOutEventsList() *sql.Rows {
	rows, err := GetDB().Query("SELECT id, tariff, entry, exit, sum, carplate, isshown, isaccepted, isdenied, cartype, nationalitytype, parkingid, isclosed " +
		"FROM parking_event WHERE isclosed = true AND islocalevent = false AND issyncedin = true AND issyncedout = false;")

	if err != nil {
		panic(err)
	}

	return rows
}

//SyncNotLocalOut sync not local parking out event
func SyncNotLocalOut(ID string) sql.Result {
	eventID, err := GetDB().Exec("UPDATE parking_event SET issyncedout = true WHERE id = $1", ID)

	if err != nil {
		log.Println("Error not local parking out event syncing!", err)
	}

	return eventID
}

//GetCountNotLocalCreatedInvoices return count of not local creted invoices
func GetCountNotLocalCreatedInvoices() int {
	var counter int
	err := GetDB().QueryRow("SELECT count(*) FROM parking_invoice WHERE islocalinvoice = false AND ispayed = false AND issyncedcreation = false").Scan(&counter)

	if err != nil {
		panic(err)
	}

	return counter
}

//GetNotLocalCreatedInvoicesList return list of not local creted invoices
func GetNotLocalCreatedInvoicesList() *sql.Rows {
	rows, err := GetDB().Query("SELECT id, datetime, carplate, sum, eventid, ispayed, islocalinvoice " +
		"FROM parking_invoice WHERE islocalinvoice = false AND ispayed = false AND issyncedcreation = false;")

	if err != nil {
		panic(err)
	}

	return rows
}

//SyncNotLocalCreatedInvoices sync not local creted invoice
func SyncNotLocalCreatedInvoices(ID string) sql.Result {
	eventID, err := GetDB().Exec("UPDATE parking_invoice SET issyncedcreation = true WHERE id = $1", ID)

	if err != nil {
		log.Println("Error not local creted invoice syncing!", err)
	}

	return eventID
}

//GetCountNotLocalPayedInvoices return count of not local payed invoices
func GetCountNotLocalPayedInvoices() int {
	var counter int
	err := GetDB().QueryRow("SELECT count(*) FROM parking_invoice WHERE islocalinvoice = false AND ispayed = true AND issyncedcreation = true AND issyncedpayed = false").Scan(&counter)

	if err != nil {
		panic(err)
	}

	return counter
}

//GetNotLocalPayedInvoicesList return list of not local payed invoices
func GetNotLocalPayedInvoicesList() *sql.Rows {
	rows, err := GetDB().Query("SELECT id, datetime, carplate, sum, eventid, ispayed, islocalinvoice " +
		"FROM parking_invoice WHERE islocalinvoice = false AND ispayed = true AND issyncedcreation = true AND issyncedpayed = false;")

	if err != nil {
		panic(err)
	}

	return rows
}

//SyncNotLocalPayedInvoices sync not local payed invoice
func SyncNotLocalPayedInvoices(ID string) sql.Result {
	eventID, err := GetDB().Exec("UPDATE parking_invoice SET issyncedpayed = true WHERE id = $1", ID)

	if err != nil {
		log.Println("Error not local payed invoice syncing!", err)
	}

	return eventID
}