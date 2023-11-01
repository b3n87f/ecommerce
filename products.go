package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jeffail/gabs/v2"
)

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.URL.Query().Get("token")
	_, err = checkUserIdFromToken(token)
	if err != nil {
		go sendDataToChannel("Get Products Need Valid Token : " + token)
		http.Error(w, clearerrorreturn(err.Error()), http.StatusUnauthorized)
		return
	}
	rows, err := db.Query("SELECT id, name,created_at,is_valid, category, base_price, discount, pay_price FROM products WHERE is_valid = $1", true)
	if err != nil {
		go sendDataToChannel("Get Products Sql Query Error : " + err.Error())
		http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	HTTPResponse := gabs.New()
	HTTPResponse.Set("OK", "data", "type")
	HTTPResponse.Set("OK", "data", "status")

	var (
		i, id, basePrice, discount, payPrice int
		isValid                              bool
		name, category, createdAt            string
	)
	for rows.Next() {
		err := rows.Scan(&id, &name, &createdAt, &isValid, &category, &basePrice, &discount, &payPrice)
		if err != nil {
			go sendDataToChannel("Get Products Sql Row Scan Error : " + err.Error())
			http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
			return
		}
		i++
		HTTPResponse.Set(id, "data", "response", "products", strconv.Itoa(id), "id")
		HTTPResponse.Set(name, "data", "response", "products", strconv.Itoa(id), "name")
		HTTPResponse.Set(createdAt, "data", "response", "products", strconv.Itoa(id), "createdAt")
		HTTPResponse.Set(isValid, "data", "response", "products", strconv.Itoa(id), "isValid")
		HTTPResponse.Set(category, "data", "response", "products", strconv.Itoa(id), "category")
		HTTPResponse.Set(basePrice, "data", "response", "products", strconv.Itoa(id), "basePrice")
		HTTPResponse.Set(discount, "data", "response", "products", strconv.Itoa(id), "discount")
		HTTPResponse.Set(payPrice, "data", "response", "products", strconv.Itoa(id), "payPrice")
	}
	if i == 0 {
		HTTPResponse.Set("", "data", "response", "products")
	}
	HTTPResponse.Set(i, "data", "response", "totalCount")

	go sendDataToChannel("Get Products Success Token : " + token)

	fmt.Fprintf(w, HTTPResponse.String())
	return
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var (
		name, category, token    string
		isValid                  bool
		err                      error
		basePrice, discount, vat float64
		payPrice                 int
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Create Product Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Create Product token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}
	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Create Product Require An Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}

	name, err = jsonCheckerString(jsonParsed, "data.request.product.name")
	if err != nil {
		go sendDataToChannel("Create Product name Parse Error. Expected data->request->product->name as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	category, err = jsonCheckerString(jsonParsed, "data.request.product.category")
	if err != nil {
		go sendDataToChannel("Create Product category Parse Error. Expected data->request->product->category as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	isValid, err = jsonCheckerBoolen(jsonParsed, "data.request.product.isValid")
	if err != nil {
		go sendDataToChannel("Create Product isValid Parse Error. Expected data->request->product->isValid as boolen")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	basePrice, err = jsonCheckerFloat64(jsonParsed, "data.request.product.basePrice")
	if err != nil {
		go sendDataToChannel("Create Product basePrice Parse Error. Expected data->request->product->basePrice as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	discount, err = jsonCheckerFloat64(jsonParsed, "data.request.product.discount")
	if err != nil {
		go sendDataToChannel("Create Product discount Parse Error. Expected data->request->product->discount as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	vat, err = jsonCheckerFloat64(jsonParsed, "data.request.product.vat")
	if err != nil {
		go sendDataToChannel("Create Product vat Parse Error. Expected data->request->product->vat as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	payPriceTemp := int(basePrice) - ((int(basePrice) * int(discount)) / 100)
	payPrice = int(payPriceTemp) + ((int(payPriceTemp) * int(vat)) / 100)

	values := map[string]interface{}{
		"name":       name,
		"is_valid":   isValid,
		"category":   category,
		"base_price": basePrice,
		"discount":   discount,
		"vat":        vat,
		"pay_price":  payPrice,
	}

	_, errInsert := sq.insert("products", values)

	if errInsert != nil {
		go sendDataToChannel("Create Product Insert Error " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Create Product Insert Success Name : " + name + " Base Price : " + strconv.Itoa(int(basePrice)) + " Vat : " + strconv.Itoa(int(vat)) + " Discount : " + strconv.Itoa(int(discount)) + " Pay Price " + strconv.Itoa(int(payPrice)) + " Category : " + category)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("NEW PRODUCT INSERTED"))
		return
	}

}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	var (
		name, category, token        string
		isValid                      bool
		err                          error
		basePrice, discount, vat, id float64
		payPrice                     int
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Update Product Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Update Product token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}
	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Update Product Require An Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}

	id, err = jsonCheckerFloat64(jsonParsed, "data.request.product.id")
	if err != nil {
		go sendDataToChannel("Update Product id Parse Error. Expected data->request->product->id as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	name, err = jsonCheckerString(jsonParsed, "data.request.product.name")
	if err != nil {
		go sendDataToChannel("Update Product name Parse Error. Expected data->request->product->name as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	category, err = jsonCheckerString(jsonParsed, "data.request.product.category")
	if err != nil {
		go sendDataToChannel("Update Product category Parse Error. Expected data->request->product->category as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	isValid, err = jsonCheckerBoolen(jsonParsed, "data.request.product.isValid")
	if err != nil {
		go sendDataToChannel("Update Product isValid Parse Error. Expected data->request->product->isValid as boolen")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	basePrice, err = jsonCheckerFloat64(jsonParsed, "data.request.product.basePrice")
	if err != nil {
		go sendDataToChannel("Update Product basePrice Parse Error. Expected data->request->product->basePrice as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	discount, err = jsonCheckerFloat64(jsonParsed, "data.request.product.discount")
	if err != nil {
		go sendDataToChannel("Update Product discount Parse Error. Expected data->request->product->discount as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	vat, err = jsonCheckerFloat64(jsonParsed, "data.request.product.vat")
	if err != nil {
		go sendDataToChannel("Update Product vat Parse Error. Expected data->request->product->vat as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	payPriceTemp := int(basePrice) - ((int(basePrice) * int(discount)) / 100)
	payPrice = int(payPriceTemp) + ((int(payPriceTemp) * int(vat)) / 100)

	values := map[string]interface{}{
		"name":       name,
		"is_valid":   isValid,
		"category":   category,
		"base_price": basePrice,
		"discount":   discount,
		"vat":        vat,
		"pay_price":  payPrice,
	}
	errInsert := sq.updateWithLock("products", values, " id = "+strconv.Itoa(int(id)))

	if errInsert != nil {
		go sendDataToChannel("Update Product Insert Error " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Update Product Success Name : " + name + " Base Price : " + strconv.Itoa(int(basePrice)) + " Vat : " + strconv.Itoa(int(vat)) + " Discount : " + strconv.Itoa(int(discount)) + " Pay Price " + strconv.Itoa(int(payPrice)) + " Category : " + category)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("PRODUCT UPDATED"))
		return
	}

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    float64
		token string
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Delete Product Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Delete Product token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}
	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Delete Product Require An Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}

	id, err = jsonCheckerFloat64(jsonParsed, "data.request.product.id")
	if err != nil {
		go sendDataToChannel("Delete Product id Parse Error. Expected data->request->product->id as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	values := map[string]interface{}{
		"is_valid": false,
	}
	errInsert := sq.updateWithLock("products", values, " id = "+strconv.Itoa(int(id)))

	if errInsert != nil {
		go sendDataToChannel("Delete Product Error " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Delete Product Success ID : " + strconv.Itoa(int(id)))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("PRODUCT DELETED"))
		return
	}

}
