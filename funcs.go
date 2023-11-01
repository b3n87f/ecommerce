package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/joho/godotenv"
)

func postgredbinfo() (connStr string) {
	err := godotenv.Load(".env")
	if err != nil {
		go sendDataToChannel("Env File Cannot Opened !" + err.Error())
		log.Fatalf("ENV NOT OPENED: %s", err)
	}
	user := os.Getenv("POSTGRESUSER")
	dbname := os.Getenv("POSTGRESDBNAME")
	sslmode := os.Getenv("POSTRESSSLMODE")
	password := os.Getenv("POSTGRESPASSWORD")
	host := os.Getenv("POSTGRESHOST")
	port := os.Getenv("POSTGRESPORT")

	connStr = fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s", user, dbname, sslmode, password, host, port)
	return connStr
}

func getenvparameters() {
	err := godotenv.Load(".env")
	if err != nil {
		go sendDataToChannel("Env File Cannot Opened !" + err.Error())
		log.Fatalf("ENV NOT OPENED: %s", err)
	}
	envDbName = os.Getenv("POSTGRESDBNAME")
	logmode = os.Getenv("LOGMODE")

}

func migrationdbinfo() (string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		go sendDataToChannel("Env File Cannot Opened !" + err.Error())
		log.Fatalf("ENV NOT OPENED: %s", err)
	}
	user := os.Getenv("POSTGRESUSER")
	dbname := os.Getenv("POSTGRESDBNAME")
	sslmode := os.Getenv("POSTRESSSLMODE")
	password := os.Getenv("POSTGRESPASSWORD")
	host := os.Getenv("POSTGRESHOST")
	port := os.Getenv("POSTGRESPORT")

	dsnFirst := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=postgres sslmode=%s", host, user, password, port, sslmode)
	dsnSecond := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s", host, user, password, port, dbname, sslmode)
	return dsnFirst, dsnSecond
}

func clearerrorreturn(message string) string {
	HTTPResponse := gabs.New()
	HTTPResponse.Set("OK", "data", "type")
	HTTPResponse.Set(message, "data", "response")
	HTTPResponse.Set("ERROR", "data", "status")
	return HTTPResponse.String()
}

func clearokreturn(message string) string {
	HTTPResponse := gabs.New()
	HTTPResponse.Set("OK", "data", "type")
	HTTPResponse.Set("OK", "data", "status")
	HTTPResponse.Set(message, "data", "response")
	return HTTPResponse.String()
}

func jsonCheckerString(jsonParsed *gabs.Container, searchPath string) (string, error) {
	var returnString string

	if ok := jsonParsed.ExistsP(searchPath); !ok {
		return "", errors.New(jsonCheckerStringError.String())
	} else {
		temp := jsonParsed.Path(searchPath).Data()
		switch temp.(type) {
		case string:
			returnString = temp.(string)
		case float64:
			returnString = strconv.FormatFloat(temp.(float64), 'f', -1, 64)
		case bool:
			returnString = strconv.FormatBool(temp.(bool))
		default:
			return "", errors.New(jsonCheckerStringError.String())
		}
	}
	return returnString, nil
}

func jsonCheckerFloat64(jsonParsed *gabs.Container, searchPath string) (float64, error) {
	var returnFl float64
	if ok := jsonParsed.ExistsP(searchPath); !ok {
		return 0, errors.New(jsonCheckerFloat64Error.String())
	} else {
		temp := jsonParsed.Path(searchPath).Data()
		switch temp.(type) {
		case string:
			return 0, errors.New(jsonCheckerFloat64Error.String())
		case float64:
			returnFl = temp.(float64)
		case bool:
			return 0, errors.New(jsonCheckerFloat64Error.String())
		default:
			return 0, errors.New(jsonCheckerFloat64Error.String())
		}
	}
	return returnFl, nil
}

func jsonCheckerBoolen(jsonParsed *gabs.Container, searchPath string) (bool, error) {
	var returnBoolen bool
	if ok := jsonParsed.ExistsP(searchPath); !ok {
		return false, errors.New(jsonCheckerBoolenError.String())
	} else {
		temp := jsonParsed.Path(searchPath).Data()
		switch temp.(type) {
		case string:
			return false, errors.New(jsonCheckerBoolenError.String())
		case float64:
			return false, errors.New(jsonCheckerBoolenError.String())
		case bool:
			returnBoolen = temp.(bool)
		default:
			return false, errors.New(jsonCheckerBoolenError.String())
		}
	}
	return returnBoolen, nil
}

func clearMapStringInt(m map[string]int) {
	for key := range m {
		delete(m, key)
	}
	return
}

func clearMapIntBool(m map[int]bool) {
	for key := range m {
		delete(m, key)
	}
	return
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`

	match, _ := regexp.MatchString(pattern, email)
	return match
}

func isValidPhoneNumber(phoneNumber string) bool {
	match, _ := regexp.MatchString("^[0-9]+$", phoneNumber)
	return match
}

func checkTokensAdminStatus(token string) bool {
	var adminstate bool
	authMapMutex.RLock()
	if adminMap[tokenMap[token]] {
		adminstate = true
		authMapMutex.RUnlock()
	} else {
		authMapMutex.RUnlock()
		tokensMap()
		authMapMutex.RLock()
		if adminMap[tokenMap[token]] {
			adminstate = true
		}
		authMapMutex.RUnlock()
	}
	return adminstate
}

func checkUserIdFromToken(token string) (int, error) {
	authMapMutex.RLock()
	userId := tokenMap[token]
	authMapMutex.RUnlock()
	if userId == 0 {
		tokensMap()
		authMapMutex.RLock()
		userId = tokenMap[token]
		authMapMutex.RUnlock()
		if userId == 0 {
			return userId, errors.New(userNotFound.String())
		}
	}
	return userId, nil
}

type Product struct {
	ID        int
	Name      string
	IsValid   bool
	BasePrice int
	Discount  int
	PayPrice  int
	Vat       int
	Category  string
}

func itemsMap() map[int]Product {
	productMap := make(map[int]Product)

	rows, err := db.Query("SELECT id, name,is_valid, category, base_price, discount, pay_price, vat FROM products")
	if err != nil {
		go sendDataToChannel("Products Get Sql Query Error " + err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.Name, &item.IsValid, &item.Category, &item.BasePrice, &item.Discount, &item.PayPrice, &item.Vat); err != nil {
			go sendDataToChannel("Products Get Row Scan Error " + err.Error())
			return nil
		}
		productMap[item.ID] = item
	}
	return productMap
}

func randomCreatorText(digit int) string {
	seed := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(seed), func(i, j int) { seed[i], seed[j] = seed[j], seed[i] })

	hrand := make([]rune, digit)
	for i := 0; i < digit; i++ {
		hrand[i] = seed[rand.Intn(len(seed))]
	}

	return fmt.Sprintf("%x", string(hrand))
}

func currentOrderStatus(id int) int {
	sqlStatement := `SELECT status FROM orders WHERE id=$1;`
	var status int
	row := db.QueryRow(sqlStatement, id)
	if err := row.Scan(&status); err != nil {
		go sendDataToChannel("Order Status Get Error OrderId : " + strconv.Itoa(id) + " " + err.Error())
		return 0
	}
	return status
}

func paymentListenerChecker(transactionId string) (int, string, int) {
	sqlStatement := `SELECT id, transaction_method,status FROM orders WHERE order_id=$1;`
	var id, status int
	var transactionMethod string
	row := db.QueryRow(sqlStatement, transactionId)
	if err := row.Scan(&id, &transactionMethod, &status); err != nil {
		go sendDataToChannel("paymentListenerChecker Get Error transactionId : " + transactionId + " " + err.Error())
		return 0, "", 0
	}
	return id, transactionMethod, status
}

func randomString(n int) (string, error) {
	bytes := make([]byte, n/2) 
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sendDataToChannel(data string) {
	dataChannel <- data
	return
}

func saveDataFromChannelToDB(dataChannel <-chan string) {
	go func() {
		for {
			select {
			case data := <-dataChannel:
				if logmode == "true" {
					go insertLog(data)
				}
			case <-time.After(5 * time.Second):
				fmt.Println("No log for 5s, checking again...")
			}
		}
	}()
}

func insertLog(data string) {
	values := map[string]interface{}{
		"log": data,
	}
	_, errInsert := sq.insert("loggers", values)
	if errInsert != nil {
		fmt.Println("Log insert error ! ", data, " at ", time.Now(), " err : ", errInsert.Error())
		return
	}
	return

}
