package controllers

import (
	"encoding/json"
	"fmt"
	"gormTest/config"
	"gormTest/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type DeleteResponse struct {
	Id string `json:"id"`
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Go Api Server")
	fmt.Println("Root endpoint is hooked!")
}

func fetchAllItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	// modelの呼び出し
	models.GetAllItems(&items)
	responseBody, err := json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func fetchSingleItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item models.Item
	// modelの呼び出し
	models.GetSingleItem(&item, id)
	responseBody, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var item models.Item
	if err := json.Unmarshal(reqBody, &item); err != nil {
		log.Fatal(err)
	}
	// modelの呼び出し
	models.InsertItem(&item)
	responseBody, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// modelの呼び出し
	models.DeleteItem(id)
	responseBody, err := json.Marshal(DeleteResponse{Id: id})
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)

	var updateItem models.Item
	if err := json.Unmarshal(reqBody, &updateItem); err != nil {
		log.Fatal(err)
	}
	// modelの呼び出し
	models.UpdateItem(&updateItem, id)
	convertUintId, _ := strconv.ParseUint(id, 10, 64)
	updateItem.Model.ID = uint(convertUintId)
	responseBody, err := json.Marshal(updateItem)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func StartWebServer() error {
	fmt.Println("Rest API with Mux Routers")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/health", rootPage)
	router.HandleFunc("/items", fetchAllItems).Methods("GET")
	router.HandleFunc("/item/{id}", fetchSingleItem).Methods("GET")

	router.HandleFunc("/item", createItem).Methods("POST")
	router.HandleFunc("/item/{id}", deleteItem).Methods("DELETE")
	router.HandleFunc("/item/{id}", updateItem).Methods("PUT")

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.ServerPort), router)
}
