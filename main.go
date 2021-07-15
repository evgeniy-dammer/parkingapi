package main

import (
	"fmt"
	"net/http"
	"parkingapi/controllers"
	"parkingapi/models"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	models.DatabaseConnect()

	router := mux.NewRouter()

	router.HandleFunc("/api/parkingin", controllers.CreateNewParkingIn).Methods("POST")
	router.HandleFunc("/api/test", controllers.Test).Methods("GET")
	router.HandleFunc("/api/getnewevent", controllers.GetNewEvent).Methods("GET")
	router.HandleFunc("/api/denied/{id:[0-9a-zA\\-]+}", controllers.DeniedParkingIn).Methods("GET")
	router.HandleFunc("/api/accepted/{id:[0-9a-zA\\-]+}", controllers.AcceptedParkingIn).Methods("GET")
	router.HandleFunc("/api/getallevents", controllers.GetAllEvents).Methods("GET")
	router.HandleFunc("/api/createinvoice", controllers.CreateInvoice).Methods("POST")
	router.HandleFunc("/api/exitevent/{id:[0-9a-zA\\-]+}", controllers.ParkingInExitEvent).Methods("GET")
	router.HandleFunc("/api/balance/{id}", controllers.GetBalance).Methods("GET")
	router.HandleFunc("/api/synclocalinevents", controllers.SyncLocalInEvents).Methods("POST")
	router.HandleFunc("/api/synclocaloutevents", controllers.SyncLocalOutEvents).Methods("POST")
	router.HandleFunc("/api/synclocalcreatedinvoices", controllers.SyncLocalCreatedInvoices).Methods("POST")
	router.HandleFunc("/api/synclocalpayedinvoices", controllers.SyncLocalPayedInvoices).Methods("POST")
	router.HandleFunc("/api/syncnotlocalinevents", controllers.SyncNotLocalInEvents).Methods("GET")
	router.HandleFunc("/api/syncnotlocalin/{id:[0-9a-zA\\-]+}", controllers.SyncNotLocalIn).Methods("GET")
	router.HandleFunc("/api/syncnotlocaloutevents", controllers.SyncNotLocalOutEvents).Methods("GET")
	router.HandleFunc("/api/syncnotlocalout/{id:[0-9a-zA\\-]+}", controllers.SyncNotLocalOut).Methods("GET")
	router.HandleFunc("/api/syncnotlocalcreatedinvoices", controllers.SyncNotLocalCreatedInvoices).Methods("GET")
	router.HandleFunc("/api/syncnotlocalcreatedinvoice/{id:[0-9a-zA\\-]+}", controllers.SyncNotLocalCreatedInvoice).Methods("GET")
	router.HandleFunc("/api/syncnotlocalpayedinvoices", controllers.SyncNotLocalPayedInvoices).Methods("GET")
	router.HandleFunc("/api/syncnotlocalpayedinvoice/{id:[0-9a-zA\\-]+}", controllers.SyncNotLocalPayedInvoice).Methods("GET")

	router.HandleFunc("/billing/", controllers.Dashbord)
	router.HandleFunc("/billing/events/", controllers.Events)
	router.HandleFunc("/billing/invoices/", controllers.Events)
	router.HandleFunc("/billing/parkings/", controllers.Events)
	router.HandleFunc("/billing/cartypes/", controllers.Events)
	router.HandleFunc("/billing/nationalitytypes/", controllers.Events)
	router.HandleFunc("/billing/places/", controllers.Events)
	router.HandleFunc("/billing/tariffs/", controllers.Events)
	router.HandleFunc("/billing/accounts/", controllers.Events)

	//CSS
	router.PathPrefix("/dist/css/skins/").Handler(http.StripPrefix("/dist/css/skins/", http.FileServer(http.Dir("templates/dist/css/skins/"))))
	router.PathPrefix("/dist/css/").Handler(http.StripPrefix("/dist/css/", http.FileServer(http.Dir("templates/dist/css/"))))
	router.PathPrefix("/bootstrap/css/").Handler(http.StripPrefix("/bootstrap/css/", http.FileServer(http.Dir("templates/bootstrap/css/"))))
	router.PathPrefix("/plugins/iCheck/flat/").Handler(http.StripPrefix("/plugins/iCheck/flat/", http.FileServer(http.Dir("templates/plugins/iCheck/flat/"))))
	router.PathPrefix("/plugins/morris/").Handler(http.StripPrefix("/plugins/morris/", http.FileServer(http.Dir("templates/plugins/morris/"))))
	router.PathPrefix("/plugins/jvectormap/").Handler(http.StripPrefix("/plugins/jvectormap/", http.FileServer(http.Dir("templates/plugins/jvectormap/"))))
	router.PathPrefix("/plugins/datepicker/").Handler(http.StripPrefix("/plugins/datepicker/", http.FileServer(http.Dir("templates/plugins/datepicker/"))))
	router.PathPrefix("/plugins/datatables/").Handler(http.StripPrefix("/plugins/datatables/", http.FileServer(http.Dir("templates/plugins/datatables/"))))
	router.PathPrefix("/plugins/daterangepicker/").Handler(http.StripPrefix("/plugins/daterangepicker/", http.FileServer(http.Dir("templates/plugins/daterangepicker/"))))
	router.PathPrefix("/plugins/bootstrap-wysihtml5/").Handler(http.StripPrefix("/plugins/bootstrap-wysihtml5/", http.FileServer(http.Dir("templates/plugins/bootstrap-wysihtml5/"))))

	//JS
	router.PathPrefix("/plugins/jQuery/").Handler(http.StripPrefix("/plugins/jQuery/", http.FileServer(http.Dir("templates/plugins/jQuery/"))))
	router.PathPrefix("/bootstrap/js/").Handler(http.StripPrefix("/bootstrap/js/", http.FileServer(http.Dir("templates/bootstrap/js/"))))
	router.PathPrefix("/plugins/sparkline/").Handler(http.StripPrefix("/plugins/sparkline/", http.FileServer(http.Dir("templates/plugins/sparkline/"))))
	router.PathPrefix("/plugins/knob/").Handler(http.StripPrefix("/plugins/knob/", http.FileServer(http.Dir("templates/plugins/knob/"))))
	router.PathPrefix("/plugins/iCheck/").Handler(http.StripPrefix("/plugins/iCheck/", http.FileServer(http.Dir("templates/plugins/iCheck/"))))
	router.PathPrefix("/plugins/slimScroll/").Handler(http.StripPrefix("/plugins/slimScroll/", http.FileServer(http.Dir("templates/plugins/slimScroll/"))))
	router.PathPrefix("/plugins/slimScroll/").Handler(http.StripPrefix("/plugins/slimScroll/", http.FileServer(http.Dir("templates/plugins/slimScroll/"))))
	router.PathPrefix("/plugins/fastclick/").Handler(http.StripPrefix("/plugins/fastclick/", http.FileServer(http.Dir("templates/plugins/fastclick/"))))
	router.PathPrefix("/dist/js/").Handler(http.StripPrefix("/dist/js/", http.FileServer(http.Dir("templates/dist/js/"))))
	router.PathPrefix("/dist/js/pages/").Handler(http.StripPrefix("/dist/js/pages/", http.FileServer(http.Dir("templates/dist/js/pages/"))))

	fmt.Println("API is started...")
	http.ListenAndServe(":8181", router)
}
