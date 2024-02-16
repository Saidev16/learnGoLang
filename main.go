package main
import (
	"encoding/json" 
	"log" 
	"net/http" 
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
)


// Init Cars var as a slice Car struct
var cars []Car

// Car Struct (Model)
type Car struct {
	ID string `json:"id"`
	Model string `json:"model"`
	Registration string `json:"registration"`
	Mileage float64 `json:"mileage"`
	Condition string `json:"condition"`
}

// Get All Cars
func getCars(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func getCar(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// Loop through cars and find with id
	for _, item := range cars {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Car{})
}

func createCar(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	var newCar Car
	_ = json.NewDecoder(r.Body).Decode(&newCar)

	for _, car := range cars {
		if car.Registration == newCar.Registration {
			json.NewEncoder(w).Encode("the registration number already exists")
			return
		}
	}


	newCar.ID = strconv.Itoa( rand.Intn(10000000) ) 
	cars = append(cars, newCar)
	json.NewEncoder(w).Encode(newCar.ID)

}


func rentCar(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	// Get the params
	params := mux.Vars(r) 
	registrationNumber := params["registration"]
	// Loop through cars and find with registration
	for index, car := range cars {
		if car.Registration == registrationNumber {
			// if car already rendted throw error and exit
			if car.Condition == "rented" {
				w.WriteHeader(http.StatusConflict)
                json.NewEncoder(w).Encode("Car is already rented")
                return
			}
			cars[index].Condition = "rented"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Car is rented")
			return
	
		}
	}
	// if car not found throw error
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode("Car with the given registration number does not exist")

}


func main()  {
	// Init Router
	r := mux.NewRouter()

	//MOCK Data
	cars = append(cars, Car{ID: "1",Model:"ford" , Registration: "AA222", Mileage: 1000, Condition: "available"})
	cars = append(cars, Car{ID: "2",Model:"mercedes", Registration: "AA413", Mileage: 1320, Condition: "available"})
	cars = append(cars, Car{ID: "3",Model:"tesla", Registration: "AA484", Mileage: 1470, Condition: "rented"})

	// Route Handlers / Endpoints
	r.HandleFunc("/cars",getCars).Methods("GET")
	r.HandleFunc("/cars/{id}",getCar).Methods("GET")
	r.HandleFunc("/cars",createCar).Methods("POST")
	r.HandleFunc("/cars/{registration}/rentals",rentCar).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000",r))

}