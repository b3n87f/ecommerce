package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Customers struct {
	First_name string    `gorm:"column:first_name;not null"`
	Last_name  string    `gorm:"column:last_name;not null"`
	Phone      string    `gorm:"column:phone;unique;not null"`
	Address    string    `gorm:"column:address;not null"`
	Email      string    `gorm:"column:email;unique;not null"`
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Created_at time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

type Loggers struct {
	Log string `gorm:"column:log;not null"`
	Id  int64  `gorm:"column:id;primaryKey;autoIncrement"`
}

type Order_details struct {
	Order_id           string `gorm:"column:order_id;not null"`
	Product_id         int64  `gorm:"column:product_id;not null"`
	Product_base_price int    `gorm:"column:product_base_price;not null"`
	Product_discount   int    `gorm:"column:product_discount;not null"`
	Product_pay_price  int    `gorm:"column:product_pay_price;not null"`
	Product_vat        int    `gorm:"column:product_vat;not null"`
	Product_name       string `gorm:"column:product_name;not null"`
	Quantity           int    `gorm:"column:quantity;not null"`
	Id                 int64  `gorm:"column:id;primaryKey;autoIncrement"`
}

type Orders struct {
	Customer           int64     `gorm:"column:customer;not null"`
	Inserted_date      time.Time `gorm:"column:inserted_date;not null;default:CURRENT_TIMESTAMP"`
	Updated_date       time.Time `gorm:"column:updated_date"`
	Status             int64     `gorm:"column:status;not null"`
	Order_id           string    `gorm:"column:order_id;unique;not null"`
	Transaction_method string    `gorm:"column:transaction_method;not null"`
	Total_price        int64     `gorm:"column:total_price;not null"`
	Id                 int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Products           string    `gorm:"column:products;not null"`
}

type Products struct {
	Name       string    `gorm:"column:name;not null"`
	Created_at time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	Is_valid   bool      `gorm:"column:is_valid;not null"`
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Base_price int       `gorm:"column:base_price;not null"`
	Discount   int       `gorm:"column:discount;not null;default:0"`
	Pay_price  int       `gorm:"column:pay_price;not null"`
	Vat        int       `gorm:"column:vat"`
	Category   string    `gorm:"column:category"`
}

type Tokens struct {
	Token         string    `gorm:"column:token;not null"`
	User_id       int64     `gorm:"column:user_id;not null"`
	Created_at    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	Validity_date time.Time `gorm:"column:validity_date;not null"`
	Id            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Is_admin      int       `gorm:"column:is_admin;not null"`
}

func migration() {
	var (
		db  *gorm.DB
		err error
	)
	db, err = gorm.Open(postgres.Open(dsnFirst), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var dbName string
	err = db.Raw("SELECT datname FROM pg_database WHERE datname = ?", envDbName).Scan(&dbName).Error
	if dbName == "" {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", envDbName)).Error
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database created successfully! : " + envDbName)
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Database already exists!")
	}

	newLogger := logger.Default.LogMode(logger.Silent)
	db, err = gorm.Open(postgres.Open(dsnSecond), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("failed to connect to the new database")
	}

	err = db.AutoMigrate(&Customers{}, &Loggers{}, &Order_details{}, &Orders{}, &Products{}, &Tokens{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	tokensData := []Tokens{
		{
			Token:         "admintoken",
			User_id:       1,
			Created_at:    time.Date(2023, 10, 30, 8, 11, 26, 853843000, time.FixedZone("+03", 3*3600)),
			Validity_date: time.Date(2023, 11, 30, 8, 11, 26, 853843000, time.UTC),
			Id:            1,
			Is_admin:      1,
		},
		{
			Token:         "usertoken1",
			User_id:       2,
			Created_at:    time.Date(2023, 10, 30, 8, 11, 26, 853843000, time.FixedZone("+03", 3*3600)),
			Validity_date: time.Date(2023, 11, 30, 8, 11, 26, 853843000, time.UTC),
			Id:            2,
			Is_admin:      0,
		},
		{
			Token:         "usertoken2",
			User_id:       3,
			Created_at:    time.Date(2023, 10, 30, 8, 11, 26, 853843000, time.FixedZone("+03", 3*3600)),
			Validity_date: time.Date(2023, 11, 30, 8, 11, 26, 853843000, time.UTC),
			Id:            3,
			Is_admin:      0,
		},
		{
			Token:         "usertoken3",
			User_id:       4,
			Created_at:    time.Date(2023, 10, 30, 8, 11, 26, 853843000, time.FixedZone("+03", 3*3600)),
			Validity_date: time.Date(2023, 11, 30, 8, 11, 26, 853843000, time.UTC),
			Id:            4,
			Is_admin:      0,
		},
	}

	for _, tokenData := range tokensData {
		_ = db.Create(&tokenData)
	}

	customersData := []Customers{
		{
			First_name: "Ali",
			Last_name:  "Yılmaz",
			Phone:      "90555111223311",
			Address:    "Çankaya Kuğulu, Park, 06690 Çankaya/Ankara",
			Email:      "aliyilmaz@ali.com",
			Id:         1,
			Created_at: time.Date(2023, 10, 30, 9, 43, 54, 354650000, time.FixedZone("+03", 3*3600)),
		},
		{
			First_name: "Mehmet",
			Last_name:  "Kaya",
			Phone:      "90555111223411",
			Address:    "Bahçelievler, İsmet İnönü Blv. No:4, 06490 Çankaya/Ankara",
			Email:      "mehmetkaya@mehmet.com",
			Id:         2,
			Created_at: time.Date(2023, 10, 30, 9, 44, 26, 55276000, time.FixedZone("+03", 3*3600)),
		},
		{
			First_name: "Ayşe",
			Last_name:  "Çaycı",
			Phone:      "90555111223511",
			Address:    "Çankaya, Şht. Ersan Cd. No: 14, 06690 Çankaya/Ankara",
			Email:      "aysecayci@ayse.com",
			Id:         3,
			Created_at: time.Date(2023, 10, 30, 9, 45, 16, 133103000, time.FixedZone("+03", 3*3600)),
		},
		{
			First_name: "Fatma",
			Last_name:  "Girik",
			Phone:      "123456789",
			Address:    "Çankaya, Merkez, 06690 Çankaya/Ankara",
			Email:      "fatma@girik.com",
			Id:         4,
			Created_at: time.Date(2023, 10, 30, 9, 46, 43, 71556000, time.FixedZone("+03", 3*3600)),
		},
	}

	for _, customerData := range customersData {
		_ = db.Create(&customerData)
	}

	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()
	return
}
