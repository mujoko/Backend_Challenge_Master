package controller

import (
	"encoding/json" // package untuk enkode dan mendekode json menjadi struct dan sebaliknya
	"fmt"
	"strconv" // package yang digunakan untuk mengubah string menjadi tipe int

	"net/http" // digunakan untuk mengakses objek permintaan dan respons dari api

	"go-inventory/config"
	"go-inventory/models" //models package dimana Stock didefinisikan
	lr "go-inventory/util/logger"

	"github.com/go-chi/chi" // digunakan untuk mendapatkan parameter dari router
	_ "github.com/lib/pq"   // postgres golang driver
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type Response struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []models.Stock `json:"data"`
}
 
// Add Stock
func AddStock(w http.ResponseWriter, r *http.Request) {
	appConf := config.AppConfig()
	logger := lr.New(appConf.Debug)

	// kita buat empty stock dengan tipe models.Stock
	var stock models.Stock
	// Prepare response objectnya
	res := response{}

	// decode data json request ke stock
	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		res.Message = "Invalid json format."
		logger.Info().Msgf("Invalid json format %v", err)
		json.NewEncoder(w).Encode(res)
		return
	}
	if stock.Price <= 0 {
		res.Message = "price should be bigger than 0."
		json.NewEncoder(w).Encode(res)
		return
	}
	if stock.Availability < 0 {
		res.Message = "availability can not negative"
		json.NewEncoder(w).Encode(res)
		return
	}

	// insert stock
	insertID := models.AddStock(stock)

	// format response objectnya
	res.ID = insertID
	res.Message = "A Stock has been added "

	// return response
	json.NewEncoder(w).Encode(res)
}

// Get Stock   single data with parameter id
func GetStock(w http.ResponseWriter, r *http.Request) {
	// set headernya
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// getId from request
	// params := mux.Vars(r)
	res := Response{Status: 0, Data: nil}

	// Convert string to int
	// id, err := strconv.Atoi(params["id"])
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)

	if err != nil {
		res.Message = "Param can not convert to int."
		json.NewEncoder(w).Encode(res)
		return
	}

	// Get Single Stock by ID
	stock, err := models.GetOneStock(int64(id))

	if err != nil {
		res.Message = "Stock not found."
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Status = 1
	res.Data = append(res.Data, stock)

	//  send response
	json.NewEncoder(w).Encode(res)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {

	// kita ambil request parameter idnya
	// params := mux.Vars(r)
	res := Response{Status: 0}

	// konversikan ke int yang sebelumnya adalah string
	// id, err := strconv.Atoi(params["id"])
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)

	if err != nil {
		res.Message = "Param can not convert to int."
		json.NewEncoder(w).Encode(res)
		return
	}

	// buat variable stock dengan type models.Stock
	var stock models.Stock

	// decode json request ke variable stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		res.Message = "Invalid json format."
		json.NewEncoder(w).Encode(res)
		return
	}

	// panggil updatestock untuk mengupdate data
	updatedRows := models.UpdateStock(int64(id), stock)

	// format message string
	msg := fmt.Sprintf("Stock succesfully update. Total %v rows/record", updatedRows)

	// ini adalah format response message
	res.Message = msg
	res.Status = 1
	res.Data = append(res.Data, stock)
	// kirim berupa response
	json.NewEncoder(w).Encode(res)
}
