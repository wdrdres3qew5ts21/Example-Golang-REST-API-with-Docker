package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"./model"
	"github.com/gorilla/mux"
)

type Product struct {
	ID            int    `json:"id"`
	ProductName   string `json:"productName"`
	ProductDetail string `json:productDetail`
}

var productList []Product

func addProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)
	productList = append(productList, product)
	json.NewEncoder(w).Encode(productList)
}

func getAllProductList(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	fmt.Println(r.URL.Path)
	json.NewEncoder(w).Encode(productList)
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	requestParam := mux.Vars(r)
	productId := requestParam["productId"]
	for _, product := range productList {
		if strconv.Itoa(product.ID) == productId {
			json.NewEncoder(w).Encode(product)
			return
		}
	}
}

func deleteProductsById(w http.ResponseWriter, r *http.Request) {
	pathVariable := mux.Vars(r)
	productId := pathVariable["productId"]

	for _, product := range productList {
		if strconv.Itoa(product.ID) == productId {
			fmt.Println("is bug ", product.ID, productId)
			productList = append(productList[:product.ID], productList[product.ID+1:]...)
			json.NewEncoder(w).Encode(productList)
			return
		}
	}
	json.NewEncoder(w).Encode(nil)

}

func loadProductServer() {
	productList = append(productList, Product{
		ID:            1,
		ProductName:   "Dell XPS 13",
		ProductDetail: "Charming white alpine white and compact",
	}, Product{
		ID:            2,
		ProductName:   "Thinkpad Carbon X1",
		ProductDetail: "Carbon Fiber Made real robust and military grade test",
	}, Product{
		ID:            3,
		ProductName:   "Dell Latitude 7340",
		ProductDetail: "Aluminium and enterprise grade",
	})
}

func main() {
	model.Init()
	food := model.Food{
		FoodName:   "ffff",
		FoodDetail: "ddd",
	}
	fmt.Println(food)

	loadProductServer()
	router := mux.NewRouter()
	router.HandleFunc("/product", addProduct).Methods("POST")
	router.HandleFunc("/products", getAllProductList).Methods("GET")
	router.HandleFunc("/product/{productId}", getProductById).Methods("GET")
	router.HandleFunc("/product/{productId}", deleteProductsById).Methods("DELETE")
	fmt.Println("----- Product Server X Golang ------")
	http.ListenAndServe(":3000", router)
}
