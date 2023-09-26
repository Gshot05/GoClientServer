package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Display struct {
	Diagonal   float32 `json:"diag"`
	Resolution string  `json:"resolution"`
	TypeMatrix string  `json:"type_matrix"`
	GSync      bool    `json:"gsync"`
}

type Monitor struct {
	Voltage        float32 `json:"voltage"`
	DisplayMonitor Display `json:"display"`
	GSyncPrem      bool    `json:"gsync_prem"`
	Curved         bool    `json:"curved"`
}

var db *sql.DB

func main() {
	connStr := "user=postgres password=0Shikhrik12$& dbname=GoDataBase sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO Type_Displays (Name_Diagonal, Name_Resolution, Type_Type, Type_Gsync) VALUES ($1, $2, $3, $4)",
		display.Diagonal, display.Resolution, display.TypeMatrix, display.GSync)
	if err != nil {
		log.Println("Ошибка при добавлении в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addMonitor(w http.ResponseWriter, r *http.Request) {
	var monitor Monitor
	err := json.NewDecoder(r.Body).Decode(&monitor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO Type_Monitor (Name_Voltage, Name_GsyncPrem, Name_Curved, Type_Displays_ID) VALUES ($1, $2, $3, $4)",
		monitor.Voltage, monitor.GSyncPrem, monitor.Curved, monitor.Type_Displays_ID)
	if err != nil {
		log.Println("Ошибка при добавлении в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT
			m.Name_Voltage,
			m.Name_GsyncPrem,
			m.Name_Curved,
			d.Name_Diagonal,
			d.Name_Resolution,
			d.Type_Type,
			d.Type_Gsync
		FROM
			Type_Monitor AS m
		INNER JOIN
			Type_Displays AS d ON m.Type_Displays_ID = d.ID_Type_Displays
	`)
	if err != nil {
		log.Println("Ошибка при запросе данных из базы данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var monitors []Monitor
	for rows.Next() {
		var monitor Monitor
		var display Display
		err := rows.Scan(
			&monitor.Voltage,
			&monitor.GSyncPrem,
			&monitor.Curved,
			&display.Diagonal,
			&display.Resolution,
			&display.TypeMatrix,
			&display.GSync,
		)
		if err != nil {
			log.Println("Ошибка при сканировании данных из базы данных:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		monitor.DisplayMonitor = display
		monitors = append(monitors, monitor)
	}

	data := struct {
		Monitors []Monitor `json:"monitors"`
	}{
		Monitors: monitors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
