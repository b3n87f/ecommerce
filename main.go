package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	connStr             = postgredbinfo()
	db                  *sqlx.DB
	err                 error
	sq                  *sqlQuery
	orderPerPage        int = 30
	dsnFirst, dsnSecond     = migrationdbinfo()
	envDbName           string
	logmode             string

	//maps
	tokenMap               = make(map[string]int)
	adminMap               = make(map[int]bool)
	adminAcceptedStatusMap = make(map[int]bool)
	userAcceptedStatusMap  = make(map[int]bool)
	notChangableStatusMap  = make(map[int]bool)
	statusNamesMap         = make(map[int]string)

	//mutex
	authMapMutex sync.RWMutex
	orderMutex   sync.RWMutex

	//channel
	dataChannel = make(chan string, 100000)
)

type CustomErrors int

const (
	jsonCheckerStringError CustomErrors = iota
	jsonCheckerFloat64Error
	jsonCheckerBoolenError
	userNotFound
)

func (n CustomErrors) String() string {
	switch n {
	case jsonCheckerStringError:
		return "String Json Search Error"
	case jsonCheckerFloat64Error:
		return "Float or Integer Json Search Error"
	case jsonCheckerBoolenError:
		return "Boolean Json Search Error"
	case userNotFound:
		return "User Not Found"
	default:
		return ""
	}
}

func init() {
	getenvparameters()
	migration()
	dbConnection()
	getAllDB()
}

func main() {
	adminAcceptedStatusMap = map[int]bool{
		3: true,
		4: true,
		5: true,
		7: true,
	}
	userAcceptedStatusMap = map[int]bool{
		6: true,
	}
	notChangableStatusMap = map[int]bool{
		0: true,
		6: true,
		7: true,
	}
	statusNamesMap = map[int]string{
		0: "Error",
		1: "Waiting for payment",
		2: "Payment okey, waiting for prepare",
		3: "Preparing",
		4: "On the way",
		5: "Delivered",
		6: "Canceled By User",
		7: "Canceled By System",
	}

	go DBLoop()
	go saveDataFromChannelToDB(dataChannel)

	r := mux.NewRouter()
	r.HandleFunc("/customers", getCustomers).Methods("GET")
	r.HandleFunc("/customers", createCustomer).Methods("POST")
	r.HandleFunc("/customers", updateCustomer).Methods("PUT")
	///r.HandleFunc("/customers", deleteCustomer).Methods("DELETE")

	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products", updateProduct).Methods("PUT")
	r.HandleFunc("/products", deleteProduct).Methods("DELETE")

	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders", createOrders).Methods("POST")
	r.HandleFunc("/orders", updateOrders).Methods("PUT")

	r.HandleFunc("/myorders", getOrdersUser).Methods("GET")
	r.HandleFunc("/myorders", updateOrdersUser).Methods("PUT")

	r.HandleFunc("/methodx-payment", paymentListenerMethodX).Methods("POST")
	r.HandleFunc("/methody-payment", paymentListenerMethodY).Methods("POST")
	r.HandleFunc("/companya-shipping", shippingListenerShipperA).Methods("POST")
	r.HandleFunc("/companya-delivered", deliveredListenerShipperA).Methods("POST")

	http.ListenAndServe(":8080", r)
}

func dbConnection() {
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error while opening database connection :", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error while connecting database :", err.Error())
	}
	db.SetMaxOpenConns(3000)
	db.SetMaxIdleConns(200)
	db.SetConnMaxLifetime(time.Minute * 1)
	sq = newSQLQuery(db)
}
