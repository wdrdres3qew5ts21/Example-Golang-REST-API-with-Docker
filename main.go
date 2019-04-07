package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
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
	fmt.Println(redisClient.Get("name"))
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

func getVistedTime(w http.ResponseWriter, r *http.Request) {
	redisClient.Set("name", "supakorn", 0)
	name, _ := redisClient.Get("name").Result()
	totalVisiting, _ := redisClient.Get("totalVisiting").Result()
	test := map[string]string{
		"name": name,
		"totalVisiting": totalVisiting,
	}
	json.NewEncoder(w).Encode(test)
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

func connectToRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "redis-server:6379",
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}

var redisClient = connectToRedis()

func main() {
	loadProductServer()
	redisClient.Set("totalVisiting", 0, 0)
	router := mux.NewRouter()
	router.Use(loggingVistitor)
	router.HandleFunc("/product", addProduct).Methods("POST")
	router.HandleFunc("/products", getAllProductList).Methods("GET")
	router.HandleFunc("/product/{productId}", getProductById).Methods("GET")
	router.HandleFunc("/product/{productId}", deleteProductsById).Methods("DELETE")
	router.HandleFunc("/status", getVistedTime).Methods("GET")
	fmt.Println("----- Product Server X Golang ------")
	http.ListenAndServe(":3000", router)
}

func loggingVistitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		totalVisiting, _ := redisClient.Get("totalVisiting").Result()
		totalVisitingInt,_ := strconv.ParseInt(totalVisiting,10,64)
		totalVisitingInt = totalVisitingInt +1
		redisClient.Set("totalVisiting",totalVisitingInt,0)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
