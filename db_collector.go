package main

import (
	"log"
	"time"
)

func getAllDB() {
	tokensMap()
}

func DBLoop() {
dbLoop:
	time.Sleep(10 * time.Second)
	getAllDB()
	goto dbLoop
}

func tokensMap() {
	authMapMutex.Lock()
	clearMapStringInt(tokenMap)
	clearMapIntBool(adminMap)
	var (
		user_id       int
		token         string
		admin         int
		validity_date string
	)
	rows, err := db.Query("SELECT user_id, token, is_admin ,validity_date  FROM tokens WHERE validity_date >= NOW()")
	if err != nil {
		log.Fatal("TOKENS_GET_SQL_QUERY_ERROR_" + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user_id, &token, &admin, &validity_date); err != nil {
			log.Fatal("TOKENS_GET_ROW_SCAN_ERROR_" + err.Error())
			return
		}
		tokenMap[token] = user_id
		if admin == 1 {
			adminMap[user_id] = true
		}

	}
	authMapMutex.Unlock()
	return
}
