package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.URL.Query().Get("token")
	page := r.URL.Query().Get("p")
	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Get Orders Need Valid Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}
	var (
		startIndex int
	)
	if page == "1" {
		startIndex = 0
	} else {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			go sendDataToChannel("Get Orders Need Valid Page : " + page)
			http.Error(w, clearerrorreturn(err.Error()), http.StatusForbidden)
			return
		}
		startIndex = pageInt * orderPerPage
	}
	rows, err := db.Query(`SELECT o.id, o.order_id, o.inserted_date, o.status, o.transaction_method,o.total_price,c.email, c.phone, c.first_name, c.last_name
							FROM orders o
							JOIN customers c ON o.customer = c.id
							ORDER BY o.id DESC
							LIMIT $1 OFFSET $2`, orderPerPage, startIndex)
	if err != nil {
		go sendDataToChannel("Get Orders Sql Query Error : " + err.Error())
		http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	HTTPResponse := gabs.New()
	HTTPResponse.Set("OK", "data", "type")
	HTTPResponse.Set("OK", "data", "status")

	var (
		id, status, i, totalPrice, totalOrderCount                                                int
		userLastName, userFirstName, userPhone, userEmail, transactionMethod, insertDate, orderId string
		orderIdInList, virg                                                                       string
	)
	for rows.Next() {
		err := rows.Scan(&id, &orderId, &insertDate, &status, &transactionMethod, &totalPrice, &userEmail, &userPhone, &userFirstName, &userLastName)
		if err != nil {
			go sendDataToChannel("Get Orders Sql Row Scan Error : " + err.Error())
			http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
			return
		}
		i++
		HTTPResponse.Set(id, "data", "response", "orders", orderId, "id")
		HTTPResponse.Set(orderId, "data", "response", "orders", orderId, "orderId")
		HTTPResponse.Set(insertDate, "data", "response", "orders", orderId, "insertDate")
		HTTPResponse.Set(statusNamesMap[status], "data", "response", "orders", orderId, "status")
		HTTPResponse.Set(transactionMethod, "data", "response", "orders", orderId, "transactionMethod")
		HTTPResponse.Set(userEmail, "data", "response", "orders", orderId, "customer", "email")
		HTTPResponse.Set(userPhone, "data", "response", "orders", orderId, "customer", "phone")
		HTTPResponse.Set(userFirstName, "data", "response", "orders", orderId, "customer", "firstName")
		HTTPResponse.Set(userLastName, "data", "response", "orders", orderId, "customer", "lastName")
		HTTPResponse.Set(totalPrice, "data", "response", "orders", orderId, "totalPrice")
		orderIdInList = orderIdInList + virg + "'" + orderId + "'"
		virg = ","
	}
	if i == 0 {
		HTTPResponse.Set("", "data", "response", "orders")
	}
	sqlStatement := `SELECT count(id) FROM orders ;`
	row := db.QueryRow(sqlStatement)
	if err := row.Scan(&totalOrderCount); err != nil {
		go sendDataToChannel("Get Orders Count QueryRow Scan Error : " + err.Error())
		fmt.Println("ITEMS_GET_ROW_SCAN_ERROR_" + err.Error())
		totalOrderCount = 90
	}

	HTTPResponse.Set(totalOrderCount, "data", "response", "totalCount")
	HTTPResponse.Set(math.Ceil(float64(totalOrderCount)/float64(orderPerPage)), "data", "response", "maxPage")
	HTTPResponse.Set(i, "data", "response", "responseCount")

	if i > 0 {
		rows, err = db.Query("SELECT id,order_id,product_id, product_base_price, product_discount, product_pay_price,product_vat, product_name, quantity FROM order_details WHERE order_id IN (" + orderIdInList + ")")
		if err != nil {
			go sendDataToChannel("Get Orders Details Sql Query Error : " + err.Error())
			http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var (
			productId, productBasePrice, productDiscount, productPayPrice, productVat, quantity int
			productName                                                                         string
		)
		for rows.Next() {
			err := rows.Scan(&id, &orderId, &productId, &productBasePrice, &productDiscount, &productPayPrice, &productVat, &productName, &quantity)
			if err != nil {
				go sendDataToChannel("Get Orders Details Sql Row Scan Error : " + err.Error())
				http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
				return
			}
			HTTPResponse.Set(productId, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productId")
			HTTPResponse.Set(productBasePrice, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productBasePrice")
			HTTPResponse.Set(productDiscount, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productDiscount")
			HTTPResponse.Set(productPayPrice, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productPayPrice")
			HTTPResponse.Set(productVat, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productVat")
			HTTPResponse.Set(productName, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productName")
			HTTPResponse.Set(quantity, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "quantity")
		}
	}

	go sendDataToChannel("Get Orders Success Token : " + token + " Page : " + page + " Total Count : " + strconv.Itoa(totalOrderCount))

	fmt.Fprintf(w, HTTPResponse.String())
	return
}

func createOrders(w http.ResponseWriter, r *http.Request) {

	var (
		token, transactionMethod string
		productlist              string 
		userId, i                int
		itemMap                  = make(map[int]Product)
		totalPrice               float64
		itemCountMap             = make(map[int]int)
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Create Orders Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Create Orders token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	transactionMethod, err = jsonCheckerString(jsonParsed, "data.request.order.transactionMethod")
	if err != nil {
		go sendDataToChannel("Create Orders transactionMethod Parse Error. Expected data->request->order->transactionMethod as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	userId, err = checkUserIdFromToken(token)
	if err != nil {
		go sendDataToChannel("Create Orders User ID Check Error. Token : " + token)
		http.Error(w, clearerrorreturn(err.Error()), http.StatusUnauthorized)
		return
	}

	productlist, err = jsonCheckerString(jsonParsed, "data.request.order.productlist")
	if err != nil {
		go sendDataToChannel("Create Orders productlist Parse Error. Expected data->request->order->productlist as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	itemMap = itemsMap()

	items := strings.Split(productlist, ",")

	for _, item := range items {
		i++
		intValue, _ := strconv.Atoi(item)
		if itemMap[intValue].IsValid {
			itemCountMap[intValue]++
			totalPrice += float64(itemMap[intValue].PayPrice)
		}
	}
	if i == 0 {
		go sendDataToChannel("Create Orders productlist Empty. ")
		fmt.Fprintf(w, clearerrorreturn("Item list is empty."))
		return
	}

	orderId, err := randomString(20)
	if err != nil {
		go sendDataToChannel("Create Orders Order ID create error. ")
		fmt.Fprintf(w, clearerrorreturn("Order ID create error."))
		return
	}

	values := map[string]interface{}{
		"customer":           userId,
		"products":           productlist,
		"status":             1,
		"order_id":           orderId,
		"transaction_method": transactionMethod,
		"total_price":        totalPrice,
	}

	_, errInsert := sq.insert("orders", values)

	for k, v := range itemCountMap {
		values := map[string]interface{}{
			"order_id":           orderId,
			"product_id":         k,
			"product_base_price": itemMap[k].BasePrice,
			"product_discount":   itemMap[k].Discount,
			"product_pay_price":  itemMap[k].PayPrice,
			"product_vat":        itemMap[k].Vat,
			"product_name":       itemMap[k].Name,
			"quantity":           v,
		}
		_, err = sq.insert("order_details", values)
		if err != nil {
			go sendDataToChannel("Create Orders order_details insert error. " + err.Error())
		}
	}

	if errInsert != nil {
		go sendDataToChannel("Create Orders orders insert error. " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Create Orders orders insert success . Order ID : " + orderId + " Products : " + productlist + " Transaction Method : " + transactionMethod + " Total Price : " + strconv.Itoa(int(totalPrice)))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("NEW ORDER INSERTED"))
		return
	}
}

func updateOrders(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		id, status float64
		token      string
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Update Order Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Update Order token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	adminState := checkTokensAdminStatus(token)
	if !adminState {
		go sendDataToChannel("Update Order Require An Admin Token : " + token)
		http.Error(w, clearerrorreturn("Need Valid Token"), http.StatusForbidden)
		return
	}

	status, err = jsonCheckerFloat64(jsonParsed, "data.request.order.status")
	if err != nil {
		go sendDataToChannel("Update Order status Parse Error. Expected data->request->order->status as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	id, err = jsonCheckerFloat64(jsonParsed, "data.request.order.id")
	if err != nil {
		go sendDataToChannel("Update Order id Parse Error. Expected data->request->order->id as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if !adminAcceptedStatusMap[int(status)] {
		go sendDataToChannel("Update Order Status Value Not In Accepted Status List. Status : " + strconv.Itoa(int(status)))
		http.Error(w, clearerrorreturn("Not valid status."), http.StatusBadRequest)
		return
	}

	currentStatus := currentOrderStatus(int(id))

	if notChangableStatusMap[currentStatus] {
		go sendDataToChannel("Update Order Status Value Currently Not Changable. Current Status : " + strconv.Itoa(currentStatus) + " Requested Status : " + strconv.Itoa(int(status)))
		http.Error(w, clearerrorreturn("Not Changable Status."), http.StatusBadRequest)
		return
	}
	values := map[string]interface{}{
		"status": status,
	}
	errInsert := sq.updateWithLock("orders", values, " id = "+strconv.Itoa(int(id)))

	if errInsert != nil {
		go sendDataToChannel("Update Order Status Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Update Order Status Success : " + strconv.Itoa(int(status)) + " Order Id : " + strconv.Itoa(int(id)))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("ORDER UPDATED"))
		return
	}
}

func updateOrdersUser(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		id, status float64
		token      string
	)
	w.Header().Set("Content-Type", "application/json")
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("Update Orders User Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	token, err = jsonCheckerString(jsonParsed, "data.token")
	if err != nil {
		go sendDataToChannel("Update Orders User token Parse Error. Expected data->token as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	_, err = checkUserIdFromToken(token)
	if err != nil {
		go sendDataToChannel("Update Orders User Need Valid Token : " + token)
		http.Error(w, clearerrorreturn(err.Error()), http.StatusUnauthorized)
		return
	}

	status, err = jsonCheckerFloat64(jsonParsed, "data.request.order.status")
	if err != nil {
		go sendDataToChannel("Update Orders User status Parse Error. Expected data->request->order->status as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	id, err = jsonCheckerFloat64(jsonParsed, "data.request.order.id")
	if err != nil {
		go sendDataToChannel("Update Orders User id Parse Error. Expected data->request->order->id as float/integer")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if !userAcceptedStatusMap[int(status)] {
		go sendDataToChannel("Update Order User Status Value Not In Accepted Status List. Status : " + strconv.Itoa(int(status)))
		http.Error(w, clearerrorreturn("Not valid status."), http.StatusBadRequest)
		return
	}

	currentStatus := currentOrderStatus(int(id))

	if notChangableStatusMap[currentStatus] {
		go sendDataToChannel("Update Order User Status Value Currently Not Changable. Current Status : " + strconv.Itoa(currentStatus) + " Requested Status : " + strconv.Itoa(int(status)))
		http.Error(w, clearerrorreturn("Not Changable Status."), http.StatusBadRequest)
		return
	}
	values := map[string]interface{}{
		"status": status,
	}
	errInsert := sq.updateWithLock("orders", values, " id = "+strconv.Itoa(int(id)))

	if errInsert != nil {
		go sendDataToChannel("Update Order User Status Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, clearerrorreturn(errInsert.Error()))
		return
	} else {
		go sendDataToChannel("Update Order User Status Success : " + strconv.Itoa(int(status)) + " Order Id : " + strconv.Itoa(int(id)))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, clearokreturn("ORDER UPDATED"))
		return
	}
}

func getOrdersUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.URL.Query().Get("token")
	page := r.URL.Query().Get("p")
	userId, err := checkUserIdFromToken(token)
	if err != nil {
		go sendDataToChannel("Get Users Orders Need Valid Token : " + token)
		http.Error(w, clearerrorreturn(err.Error()), http.StatusUnauthorized)
		return
	}

	var (
		startIndex int
	)

	if page == "1" {
		startIndex = 0
	} else {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			go sendDataToChannel("Get Orders Need Valid Page : " + page)
			http.Error(w, clearerrorreturn(err.Error()), http.StatusForbidden)
			return
		}
		startIndex = pageInt * orderPerPage
	}

	rows, err := db.Query(`SELECT o.id, o.order_id, o.inserted_date, o.status, o.transaction_method,o.total_price,c.email, c.phone, c.first_name, c.last_name
							FROM orders o
							JOIN customers c ON o.customer = c.id
							WHERE o.customer = $1
							ORDER BY o.id DESC
							LIMIT $2 OFFSET $3`, userId, orderPerPage, startIndex)
	if err != nil {
		go sendDataToChannel("Get Users Orders Sql Query Error : " + err.Error())
		http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	HTTPResponse := gabs.New()
	HTTPResponse.Set("OK", "data", "type")
	HTTPResponse.Set("OK", "data", "status")

	var (
		id, status, i, totalPrice, totalOrderCount                                                                     int
		userLastName, userFirstName, userPhone, userEmail, transactionMethod, insertDate, orderId, orderIdInList, virg string
	)
	for rows.Next() {
		err := rows.Scan(&id, &orderId, &insertDate, &status, &transactionMethod, &totalPrice, &userEmail, &userPhone, &userFirstName, &userLastName)
		if err != nil {
			go sendDataToChannel("Get Users Orders Sql Row Scan Error : " + err.Error())
			http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
			return
		}
		i++
		HTTPResponse.Set(id, "data", "response", "orders", orderId, "id")
		HTTPResponse.Set(orderId, "data", "response", "orders", orderId, "orderId")
		HTTPResponse.Set(insertDate, "data", "response", "orders", orderId, "insertDate")
		HTTPResponse.Set(statusNamesMap[status], "data", "response", "orders", orderId, "status")
		HTTPResponse.Set(transactionMethod, "data", "response", "orders", orderId, "transactionMethod")
		HTTPResponse.Set(userEmail, "data", "response", "orders", orderId, "customer", "email")
		HTTPResponse.Set(userPhone, "data", "response", "orders", orderId, "customer", "phone")
		HTTPResponse.Set(userFirstName, "data", "response", "orders", orderId, "customer", "firstName")
		HTTPResponse.Set(userLastName, "data", "response", "orders", orderId, "customer", "lastName")
		HTTPResponse.Set(userLastName, "data", "response", "orders", orderId, "totalPrice")
		orderIdInList = orderIdInList + virg + "'" + orderId + "'"
		virg = ","
	}
	if i == 0 {
		HTTPResponse.Set("", "data", "response", "orders")
	}
	HTTPResponse.Set(i, "data", "response", "totalCount")

	sqlStatement := `SELECT count(id) FROM orders WHERE customer = $1 ;`
	row := db.QueryRow(sqlStatement, userId)
	if err := row.Scan(&totalOrderCount); err != nil {
		go sendDataToChannel("Get Orders Count QueryRow Scan Error : " + err.Error())
		fmt.Println("ITEMS_GET_ROW_SCAN_ERROR_" + err.Error())
		totalOrderCount = 90
	}

	HTTPResponse.Set(totalOrderCount, "data", "response", "totalCount")
	HTTPResponse.Set(math.Ceil(float64(totalOrderCount)/float64(orderPerPage)), "data", "response", "maxPage")
	if i > 0 {

		rows, err = db.Query("SELECT id,order_id,product_id, product_base_price, product_discount, product_pay_price,product_vat, product_name, quantity FROM order_details  WHERE order_id IN (" + orderIdInList + ")")
		if err != nil {
			go sendDataToChannel("Get Users Orders Details Sql Query Error : " + err.Error())
			http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var (
			productId, productBasePrice, productDiscount, productPayPrice, productVat, quantity int
			productName                                                                         string
		)
		for rows.Next() {
			err := rows.Scan(&id, &orderId, &productId, &productBasePrice, &productDiscount, &productPayPrice, &productVat, &productName, &quantity)
			if err != nil {
				go sendDataToChannel("Get Users Orders Details Sql Row Scan Error : " + err.Error())
				http.Error(w, clearerrorreturn(err.Error()), http.StatusInternalServerError)
				return
			}
			HTTPResponse.Set(productId, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productId")
			HTTPResponse.Set(productBasePrice, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productBasePrice")
			HTTPResponse.Set(productDiscount, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productDiscount")
			HTTPResponse.Set(productPayPrice, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productPayPrice")
			HTTPResponse.Set(productVat, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productVat")
			HTTPResponse.Set(productName, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "productName")
			HTTPResponse.Set(quantity, "data", "response", "orders", orderId, "details", strconv.Itoa(id), "quantity")
		}
	}
	go sendDataToChannel("Get Users Orders Success User : " + strconv.Itoa(userId))
	fmt.Fprintf(w, HTTPResponse.String())
	return
}

func paymentListenerMethodX(w http.ResponseWriter, r *http.Request) {
	//hashing
	var (
		err                   error
		status, transactionid string
	)
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("paymentListenerMethodX Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	status, err = jsonCheckerString(jsonParsed, "data.status")
	if err != nil {
		go sendDataToChannel("paymentListenerMethodX status Parse Error. Expected data->status as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	transactionid, err = jsonCheckerString(jsonParsed, "data.transactionid")
	if err != nil {
		go sendDataToChannel("paymentListenerMethodX transactionid Parse Error. Expected data->transactionid as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if status != "OK" {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}

	orderId, paymentMethod, currentStatus := paymentListenerChecker(transactionid)

	if currentStatus != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	}

	if paymentMethod != "MethodX" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	}

	values := map[string]interface{}{
		"status": 2,
	}
	errInsert := sq.updateWithLock("orders", values, " id = "+strconv.Itoa(orderId))

	if errInsert != nil {
		go sendDataToChannel("paymentListenerMethodX Payment Update Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	} else {
		go sendDataToChannel("paymentListenerMethodX Payment Update Success Transaction ID  : " + transactionid)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}
}

func paymentListenerMethodY(w http.ResponseWriter, r *http.Request) {
	//hashing
	var (
		err                   error
		status, transactionid string
	)
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("paymentListenerMethodY Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	status, err = jsonCheckerString(jsonParsed, "data.status")
	if err != nil {
		go sendDataToChannel("paymentListenerMethodY status Parse Error. Expected data->status as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	transactionid, err = jsonCheckerString(jsonParsed, "data.transactionid")
	if err != nil {
		go sendDataToChannel("paymentListenerMethodY transactionid Parse Error. Expected data->transactionid as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if status != "OK" {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}

	orderId, paymentMethod, currentStatus := paymentListenerChecker(transactionid)

	if currentStatus != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	}

	if paymentMethod != "MethodY" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	}

	values := map[string]interface{}{
		"status": 2,
	}
	errInsert := sq.updateWithLock("orders", values, " id = "+strconv.Itoa(orderId))

	if errInsert != nil {
		go sendDataToChannel("paymentListenerMethodY Payment Update Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	} else {
		go sendDataToChannel("paymentListenerMethodY Payment Update Success Transaction ID  : " + transactionid)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}
}

func shippingListenerShipperA(w http.ResponseWriter, r *http.Request) {
	var (
		err                   error
		status, transactionid string
	)
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("shippingListenerShipperA Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	status, err = jsonCheckerString(jsonParsed, "data.status")
	if err != nil {
		go sendDataToChannel("shippingListenerShipperA status Parse Error. Expected data->status as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	transactionid, err = jsonCheckerString(jsonParsed, "data.order")
	if err != nil {
		go sendDataToChannel("shippingListenerShipperA order Parse Error. Expected data->order as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if status != "OK" {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}

	orderId, _, currentStatus := paymentListenerChecker(transactionid)

	if currentStatus != 3 {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}

	values := map[string]interface{}{
		"status": 4,
	}
	errInsert := sq.updateWithLock("orders", values, " id = "+strconv.Itoa(orderId))

	if errInsert != nil {
		go sendDataToChannel("shippingListenerShipperA Update Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	} else {
		go sendDataToChannel("shippingListenerShipperA Update Success : " + transactionid)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}
}

func deliveredListenerShipperA(w http.ResponseWriter, r *http.Request) {
	var (
		err                   error
		status, transactionid string
	)
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		go sendDataToChannel("deliveredListenerShipperA Json Parse Error.")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	status, err = jsonCheckerString(jsonParsed, "data.status")
	if err != nil {
		go sendDataToChannel("deliveredListenerShipperA status Parse Error. Expected data->status as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	transactionid, err = jsonCheckerString(jsonParsed, "data.order")
	if err != nil {
		go sendDataToChannel("deliveredListenerShipperA order Parse Error. Expected data->order as string")
		http.Error(w, clearerrorreturn(err.Error()), http.StatusBadRequest)
		return
	}

	if status != "OK" {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}

	orderId, _, currentStatus := paymentListenerChecker(transactionid)

	if currentStatus != 4 {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}

	values := map[string]interface{}{
		"status": 5,
	}
	errInsert := sq.updateWithLock("orders", values, " id = "+strconv.Itoa(orderId))

	if errInsert != nil {
		go sendDataToChannel("deliveredListenerShipperA Update Error : " + errInsert.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ERROR")
		return
	} else {
		go sendDataToChannel("deliveredListenerShipperA Update Success Transaction ID : " + transactionid)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "OK")
		return
	}
}
