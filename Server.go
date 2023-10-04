package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"-"`
	Email          string `json:"email"`
	IsAdmin        bool   `json:"is_admin"`
	HashedPassword string `json:"-"`
}

type Display struct {
	Diagonal   float32 `json:"diag"`
	Resolution string  `json:"resolution"`
	TypeMatrix string  `json:"type_matrix"`
	GSync      bool    `json:"gsync"`
}

type Monitor struct {
	VoltagePower    float32 `json:"voltage"`
	DisplayMonitor  Display `json:"display"`
	GSyncPrem       bool    `json:"gsync_prem"`
	Curved          bool    `json:"curved"`
	Type_Display_ID int     `json:"type_display_id"`
}

var mapOfUsers = make(map[string]User)

var db *sql.DB

func main() {
	connStr := "user=postgres password= dbname= sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Запуск сервера...")

	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/addDisplay", addDisplay)
	http.HandleFunc("/addMonitor", addMonitor)
	http.HandleFunc("/getAll", getAll)
	http.HandleFunc("/getMonitor", getMonitor)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		login(w, r)
	case http.MethodGet:
		w.Write([]byte("Форма авторизации"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		register(w, r)
	case http.MethodGet:
		w.Write([]byte("Форма регистрации"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	existingUser, err := getUserByUsername(user.Username)
	if err == nil && existingUser.Username == user.Username {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Пользователь с таким именем уже существует"))
		return
	}

	user.HashedPassword = hashPassword(user.Password, user.Username)
	_, err = db.Exec("INSERT INTO Type_Users (Name_Username, Name_Password, Name_email, Name_Is_Admin) VALUES ($1, $2, $3, $4)",
		user.Username, user.HashedPassword, user.Email, user.IsAdmin)

	if err != nil {
		log.Println("Ошибка при добавлении пользователя в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbUser, err := getUserByUsername(user.Username)
	if err != nil || dbUser.HashedPassword != hashPassword(user.Password, user.Username) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := generateToken(user.Username, dbUser.ID)
	mapOfUsers[token] = dbUser

	w.Write([]byte(token))
}

func getUserByUsername(username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT * FROM Type_Users WHERE Name_Username = $1", username).
		Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Email, &user.IsAdmin)

	return user, err
}

func hashPassword(password, username string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + username))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateToken(username string, userID int) string {
	token := sha256.Sum256([]byte(fmt.Sprintf("%s%d", username, userID)))
	return hex.EncodeToString(token[:])
}

func addDisplay(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !isAuthorized(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := getUserByToken(token)
	if err != nil || !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var display Display
	err = json.NewDecoder(r.Body).Decode(&display)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO Type_Display (Name_Diagonal, Name_Resolution, Type_Type, Type_Gsync) VALUES ($1, $2, $3, $4)",
		display.Diagonal, display.Resolution, display.TypeMatrix, display.GSync)
	if err != nil {
		log.Println("Ошибка при добавлении в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addMonitor(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !isAuthorized(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := getUserByToken(token)
	if err != nil || !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var monitor Monitor
	err = json.NewDecoder(r.Body).Decode(&monitor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO Type_Monitor (Name_Voltage, Name_Gsync_Prem, Name_Curved, Type_Display_ID) VALUES ($1, $2, $3, $4)",
		monitor.VoltagePower, monitor.GSyncPrem, monitor.Curved, monitor.Type_Display_ID)

	if err != nil {
		log.Println("Ошибка при добавлении в базу данных:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !isAuthorized(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
			&monitor.VoltagePower,
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getMonitor(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !isAuthorized(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
		&monitor.VoltagePower,
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	monitor.DisplayMonitor = display

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monitor)
}

func isAuthorized(token string) bool {
	_, authorized := mapOfUsers[token]
	return authorized
}

func getUserByToken(token string) (User, error) {
	user, exists := mapOfUsers[token]
	if !exists {
		return User{}, fmt.Errorf("пользователь не найден")
	}
	return user, nil
}
