package models

import (
	"database/sql"
	"go-inventory/config"

	lr "go-inventory/util/logger"

	_ "github.com/lib/pq" // postgres golang driver
)

type Stock struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Availability int     `json:"availability"`
	IsActive     bool    `json:"is_active"`
}

func AddStock(stock Stock) int64 {
	logger := lr.New(true)
	// connect to db
	db := config.CreateConnection()

	// close connectiononce finish
	defer db.Close()

	// kita buat insert query
	// mengembalikan nilai id akan mengembalikan id dari stock yang dimasukkan ke db
	sqlStatement := `INSERT INTO stock (name, price, availability,is_active ) VALUES ($1, $2, $3, $4) RETURNING id`

	// id yang dimasukkan akan disimpan di id ini
	var id int64

	// Scan function akan menyimpan insert id didalam
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Availability, stock.IsActive).Scan(&id)

	if err != nil {
		logger.Info().Msgf("Can not execute query. %v", err)
	}
	logger.Info().Msgf("Insert data single record %v", id)

	// return insert id
	return id
}

// get one  stock
func GetOneStock(id int64) (Stock, error) {
	logger := lr.New(true)
	// connect to db
	db := config.CreateConnection()

	// close connectiononce finish
	defer db.Close()

	var stock Stock

	// buat sql query
	sqlStatement := `SELECT * FROM stock WHERE id=$1`

	// exequte sql
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Availability, &stock.IsActive)

	switch err {
	case sql.ErrNoRows:
		logger.Info().Msgf("Can not find any record!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		logger.Info().Msgf("Can not retrieve from table. %v", err)
	}

	return stock, err
}

// update user in the DB
func UpdateStock(id int64, stock Stock) int64 {
	logger := lr.New(true)

	// connect to db
	db := config.CreateConnection()

	// close connectiononce finish
	defer db.Close()

	// kita buat sql query create
	sqlStatement := `UPDATE stock SET name=$2, price=$3, availability=$4 , is_active=$5 WHERE id=$1`

	// exequte sql
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Availability, stock.IsActive)

	if err != nil {
		logger.Info().Msgf("can not execute query. %v", err)
	}

	// check how many record updated
	rowsAffected, err := res.RowsAffected()

	//cek
	if err != nil {
		logger.Info().Msgf("Error while checking rows/data. %v", err)
	}
	logger.Info().Msgf("Total rows/record updated %v\n", rowsAffected)

	return rowsAffected
}
