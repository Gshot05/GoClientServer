package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type Display struct {
	Diagonal   float32 `json:"diag"`
	Resolution string  `json:"resolution"`
	TypeMatrix string  `json:"type_matrix"`
	GSync      bool    `json:"gsync"`
}

type Monitor struct {
	Voltage         float32 `json:"voltage"`
	DisplayMonitor  Display `json:"display"`
	GSyncPrem       bool    `json:"gsync_prem"`
	Curved          bool    `json:"curved"`
	Type_Display_ID int     `json:"type_display_id"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

var (
	db         *sql.DB
	userTokens map[string]User
	nc         *nats.Conn
)

func main() {
	connStr := "user=postgres password= dbname= sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userTokens = make(map[string]User)

	fmt.Println("Запуск сервера...")

	http.HandleFunc("/addDisplay", addDisplay)
	http.HandleFunc("/addMonitor", addMonitor)
	http.HandleFunc("/getAll", getAll)
	http.HandleFunc("/getMonitor", getMonitor)
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginUser)

	natsURL := "nats://95.165.107.100:4222"
	log.Println("Connecting to NATS...")
	nc, err = nats.Connect(natsURL)
	if err != nil {
		log.Println("NATS connection failed:")
		log.Fatal(err)
	}
	log.Println("NATS connection success")
	defer nc.Close()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func createToken(username, password string) string {
	tokenData := password + username
	hasher := sha256.New()
	hasher.Write([]byte(tokenData))
	return hex.EncodeToString(hasher.Sum(nil))
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword := hashPassword(user.Password)

	var storedPassword string
	err = db.QueryRow("SELECT Name_Password FROM Type_Users WHERE Name_Username = $1", user.Username).Scan(&storedPassword)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if hashedPassword == storedPassword {
		token := createToken(user.Username, storedPassword)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token":"` + token + `"}`))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword := hashPassword(user.Password)

	_, err = db.Exec("INSERT INTO Type_Users (Name_Username, Name_Password, Name_email, Name_Is_Admin) VALUES ($1, $2, $3, $4)",
		user.Username, hashedPassword, user.Email, user.IsAdmin)
	if err != nil {
		log.Println("Ошибка при добавлении пользователя в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isAdmin(username string) bool {
	var isAdmin bool
	err := db.QueryRow("SELECT Name_Is_Admin FROM Type_Users WHERE Name_Username = $1", username).Scan(&isAdmin)
	if err != nil {
		return false
	}
	return isAdmin
}

func getUsernameFromToken(token string) (string, error) {
	if len(token) < 8 || token[:7] != "Bearer " {
		return "", errors.New("Некорректный формат токена")
	}

	username := token[7:]
	return username, nil
}

func addDisplay(w http.ResponseWriter, r *http.Request) {
	var display Display
	err := json.NewDecoder(r.Body).Decode(&display)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")

	username, err := getUsernameFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !isAdmin(username) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("INSERT INTO Type_Display (Name_Diagonal, Name_Resolution, Type_Type, Type_Gsync) VALUES ($1, $2, $3, $4)",
		display.Diagonal, display.Resolution, display.TypeMatrix, display.GSync)
	if err != nil {
		log.Println("Ошибка при добавлении в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := r.FormValue("userId")

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("addDisplay allowed for User id: %s", userID)
	nc.Publish("log", []byte(message))

	if err != nil {
		errorMessage := fmt.Sprintf("Error processing addDisplay for User id: %s - %s", userID, err.Error())
		nc.Publish("error", []byte(errorMessage))

		log.Println("Error processing addDisplay:", err)
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

	token := r.Header.Get("Authorization")

	username, err := getUsernameFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !isAdmin(username) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("INSERT INTO Type_Monitor (Name_Voltage, Name_Gsync_Prem, Name_Curved, Type_Display_ID) VALUES ($1, $2, $3, $4)",
		monitor.Voltage, monitor.GSyncPrem, monitor.Curved, monitor.Type_Display_ID)

	if err != nil {
		log.Println("Ошибка при добавлении в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := r.FormValue("userId")

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("addMonitor allowed for User id: %s", userID)
	nc.Publish("log", []byte(message))

	if err != nil {
		errorMessage := fmt.Sprintf("Error processing addMonitor for User id: %s - %s", userID, err.Error())
		nc.Publish("error", []byte(errorMessage))

		log.Println("Error processing addMonitor:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	var monitors []Monitor
	rows, err := db.Query(`
        SELECT
            m.Name_Voltage,
            m.Name_Gsync_Prem,
            m.Name_Curved,
            d.Name_Diagonal,
            d.Name_Resolution,
            d.Type_Type,
            d.Type_Gsync,
			m.Type_Display_ID
        FROM
            Type_Monitor AS m
        INNER JOIN
            Type_Display AS d ON m.Type_Display_ID = d.ID_Type_Display
    `)

	if err != nil {
		log.Println("Ошибка при запросе данных из базы данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

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
			&monitor.Type_Display_ID,
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

	message := "getMonitors requested"
	nc.Publish("log", []byte(message))

	if err != nil {
		errorMessage := fmt.Sprintf("Error processing getAll - %s", err.Error())
		nc.Publish("error", []byte(errorMessage))

		log.Println("Error processing getAll:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getMonitor(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	var monitor Monitor
	var display Display
	err := db.QueryRow(`
        SELECT
            m.Name_Voltage,
            m.Name_Gsync_Prem,
            m.Name_Curved,
            d.Name_Diagonal,
            d.Name_Resolution,
            d.Type_Type,
            d.Type_Gsync,
            m.Type_Display_ID
        FROM
            Type_Monitor AS m
        INNER JOIN
            Type_Display AS d ON m.Type_Display_ID = d.ID_Type_Display
        WHERE
            m.Type_Display_ID = $1
    `, id).Scan(
		&monitor.Voltage,
		&monitor.GSyncPrem,
		&monitor.Curved,
		&display.Diagonal,
		&display.Resolution,
		&display.TypeMatrix,
		&display.GSync,
		&monitor.Type_Display_ID,
	)

	if err != nil {
		log.Println("Ошибка при запросе данных из базы данных:", err)
		errorMessage := fmt.Sprintf("Error processing getMonitor - %s", err.Error())
		nc.Publish("error", []byte(errorMessage))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("getMonitor requested for Type_Display_ID: %s", id)
	nc.Publish("log", []byte(message))

	monitor.DisplayMonitor = display

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monitor)
}
