package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type sqlQuery struct {
	db *sqlx.DB
}

func newSQLQuery(db *sqlx.DB) *sqlQuery {
	return &sqlQuery{db: db}
}

func (sq *sqlQuery) delete(tableName string, condition string, args ...interface{}) error {
	query := "DELETE FROM " + tableName + " WHERE " + condition
	_, err := sq.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (sq *sqlQuery) update(tableName string, values map[string]interface{}, condition string, args ...interface{}) error {
	query := "UPDATE " + tableName + " SET "
	params := []interface{}{}

	i := 1
	for column, value := range values {
		query += fmt.Sprintf("%s = $%d, ", column, i)
		params = append(params, value)
		i++
	}

	query = query[:len(query)-2] + " WHERE " + condition

	for _, arg := range args {
		params = append(params, arg)
	}

	_, err := sq.db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}

func (sq *sqlQuery) insert(tableName string, values map[string]interface{}) (int64, error) {
	query := "INSERT INTO " + tableName + " ("
	valuePlaceholders := ""
	params := []interface{}{}

	i := 1
	for column, value := range values {
		query += column + ", "
		valuePlaceholders += fmt.Sprintf("$%d, ", i)
		params = append(params, value)
		i++
	}

	query = query[:len(query)-2] + ") VALUES (" + valuePlaceholders[:len(valuePlaceholders)-2] + ")"

	sqlQuery, err := sq.db.Exec(query, params...)
	if err != nil {
		return 0, err
	}
	aff, _ := sqlQuery.RowsAffected()

	return aff, nil
}

func (sq *sqlQuery) updateWithLock(tableName string, values map[string]interface{}, condition string, args ...interface{}) error {
	tx, err := sq.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	lockQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s FOR UPDATE", tableName, condition)
	_, err = tx.Exec(lockQuery, args...)
	if err != nil {
		return err
	}

	updateQuery := "UPDATE " + tableName + " SET "
	params := []interface{}{}

	i := len(args) + 1 
	for column, value := range values {
		updateQuery += fmt.Sprintf("%s = $%d, ", column, i)
		params = append(params, value)
		i++
	}

	updateQuery = updateQuery[:len(updateQuery)-2] + " WHERE " + condition
	params = append(args, params...) 

	_, err = tx.Exec(updateQuery, params...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
