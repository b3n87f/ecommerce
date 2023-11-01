package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jeffail/gabs/v2"
)

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.URL.Query().Get("token")
	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Get Customers Need Valid Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}

	rows, err := db.Query("SELECT id, first_name, last_name, phone, email, address FROM customers")
	if err != nil {
		go sendDataToChannel("Get Customers Sql Query Error : " + err.Error())
		http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	HTTPResponse := gabs.New()
	HTTPResponse.Set("OK", "data", "type")
	HTTPResponse.Set("OK", "data", "status")
	var (
		i, id                                      int
		firstName, lastName, phone, email, address string
	)
	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &phone, &email, &address)
		if err != nil {
			go sendDataToChannel("Get Customers Sql Row Scan Error : " + err.Error())
			http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
			return
		}
		i++
		HTTPResponse.Set(id, "data", "response", "customers", strconv.Itoa(id), "id")
		HTTPResponse.Set(firstName, "data", "response", "customers", strconv.Itoa(id), "firstName")
		HTTPResponse.Set(lastName, "data", "response", "customers", strconv.Itoa(id), "lastName")
		HTTPResponse.Set(phone, "data", "response", "customers", strconv.Itoa(id), "phone")
		HTTPResponse.Set(email, "data", "response", "customers", strconv.Itoa(id), "email")
		HTTPResponse.Set(address, "data", "response", "customers", strconv.Itoa(id), "address")
	}
	if i == 0 {
		HTTPResponse.Set("", "data", "response", "customers")
	}
	HTTPResponse.Set(i, "data", "response", "totalCount")

	go sendDataToChannel("Get Customers Success ")
	fmt.Fprintf(w, HTTPResponse.String())
	return
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		firstName, lastName, phone, email, address, token string
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Create Customer Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Create Customer token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}
	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Create Customer Require An Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}

	firstName, err = jsonCheckerString(jsonParsed, "data.request.customer.firstName")
	if err != nil {
		go sendDataToChannel("Create Customer firstName Parse Error. Expected data->request->customer->firstName as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	lastName, err = jsonCheckerString(jsonParsed, "data.request.customer.lastName")
	if err != nil {
		go sendDataToChannel("Create Customer lastName Parse Error. Expected data->request->customer->lastName as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	phone, err = jsonCheckerString(jsonParsed, "data.request.customer.phone")
	if err != nil {
		go sendDataToChannel("Create Customer phone Parse Error. Expected data->request->customer->phone as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	email, err = jsonCheckerString(jsonParsed, "data.request.customer.email")
	if err != nil {
		go sendDataToChannel("Create Customer email Parse Error. Expected data->request->customer->email as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if !isValidEmail(email) {
		go sendDataToChannel("Create Customer Need Valid Email Address : " + email)
		http.Error(w, clearerrorreturn("Need valid email address."), http.StatusBadRequest)
		return
	}

	if !isValidPhoneNumber(phone) {
		go sendDataToChannel("Create Customer Need Valid Phone Number : " + phone)
		http.Error(w, clearerrorreturn("Need valid phone number."), http.StatusBadRequest)
		return
	}

	address, err = jsonCheckerString(jsonParsed, "data.request.customer.address")
	if err != nil {
		go sendDataToChannel("Create Customer address Parse Error. Expected data->request->customer->address as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	values := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"phone":      phone,
		"email":      email,
		"address":    address,
	}

	_, errInsert := sq.insert("customers", values)
	if errInsert != nil {
		go sendDataToChannel("Create Customer Insert Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Create Customer Insert Success First Name : " + firstName + " Last Name : " + lastName + " phone : " + phone + " email : " + email + " address : " + address)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("NEW CUSTOMER INSERTED"))
		return
	}

}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		firstName, lastName, phone, email, address, token string
		id                                                float64
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Update Customer Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Update Customer token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}
	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Update Customer Require An Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}

	id, err = jsonCheckerFloat64(jsonParsed, "data.request.customer.id")
	if err != nil {
		go sendDataToChannel("Update Customer id Parse Error. Expected data->request->customer->id as float/int")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	firstName, err = jsonCheckerString(jsonParsed, "data.request.customer.firstName")
	if err != nil {
		go sendDataToChannel("Update Customer firstName Parse Error. Expected data->request->customer->firstName as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	lastName, err = jsonCheckerString(jsonParsed, "data.request.customer.lastName")
	if err != nil {
		go sendDataToChannel("Update Customer lastName Parse Error. Expected data->request->customer->lastName as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	phone, err = jsonCheckerString(jsonParsed, "data.request.customer.phone")
	if err != nil {
		go sendDataToChannel("Update Customer phone Parse Error. Expected data->request->customer->phone as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	email, err = jsonCheckerString(jsonParsed, "data.request.customer.email")
	if err != nil {
		go sendDataToChannel("Update Customer email Parse Error. Expected data->request->customer->email as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if !isValidEmail(email) {
		go sendDataToChannel("Update Customer Need Valid Email Address : " + email)
		http.Error(w, clearerrorreturn("Need valid email address."), http.StatusBadRequest)
		return
	}

	if !isValidPhoneNumber(phone) {
		go sendDataToChannel("Update Customer Need Valid Phone Number : " + phone)
		http.Error(w, clearerrorreturn("Need valid phone number."), http.StatusBadRequest)
		return
	}

	address, err = jsonCheckerString(jsonParsed, "data.request.customer.address")
	if err != nil {
		go sendDataToChannel("Update Customer address Parse Error. Expected data->request->customer->address as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	values := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"phone":      phone,
		"email":      email,
		"address":    address,
	}
	errInsert := sq.updateWithLock("customers", values, " id = "+strconv.Itoa(int(id)))
	if errInsert != nil {
		go sendDataToChannel("Update Customer Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Update Customer Success First Name : " + firstName + " Last Name : " + lastName + " phone : " + phone + " email : " + email + " address : " + address)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("CUSTOMER UPDATED"))
		return
	}

}
