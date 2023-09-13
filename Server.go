package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Display struct {
	Diagonal   float32 `json:"diag"`
	Resolution string  `json:"resolution"`
	TypeMatrix string  `json:"type_matrix"`
	GSync      bool    `json:"gsync"`
}

type Monitor struct {
	Voltage   float32 `json:"voltage"`
	DisplayMonitor Display `json:"display"`
	GSyncPrem      bool    `json:"gsync_prem"`
	Curved         bool    `json:"curved"`
}

var monitors []Monitor
var displays []Display

func main() {
	fmt.Println("Запуск сервера...")

	http.HandleFunc("/addDisplay", addDisplay)
	http.HandleFunc("/addMonitor", addMonitor)
	http.HandleFunc("/getAll", getAll)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addDisplay(w http.ResponseWriter, r *http.Request) {
	var display Display
	err := json.NewDecoder(r.Body).Decode(&display)
	if err != nil {
		w.WriteHeader(200)
		return
	}
	displays = append(displays, display)
}

func addMonitor(w http.ResponseWriter, r *http.Request) {
	var monitor Monitor
	w.WriteHeader(200)
	err := json.NewDecoder(r.Body).Decode(&monitor)
	if err != nil {

		w.WriteHeader(200)

		return

	}
	monitors = append(monitors, monitor)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Monitors []Monitor `json:"monitors"`
		Displays []Display `json:"displays"`
	}{
		Monitors: monitors,
		Displays: displays,
	}
	json.NewEncoder(w).Encode(data)
}
