package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
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

func main() {
	var choice int
	fmt.Println("Выберите опцию:")
	fmt.Println("1. Авторизация")
	fmt.Println("2. Регистрация")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		loginMenu()
	case 2:
		registerMenu()
	default:
		fmt.Println("Неверный выбор")
	}
}

func loginMenu() {
	var user User
	fmt.Println("Введите имя пользователя:")
	fmt.Scan(&user.Username)
	fmt.Println("Введите пароль:")
	fmt.Scan(&user.Password)

	token := login(user)

	if token == "" {
		log.Println("Ошибка при входе в систему.")
		return
	}

	headers := map[string]string{"Authorization": token}

	var monitor Monitor
	fmt.Println("Введите напряжение (например, 220):")
	fmt.Scan(&monitor.VoltagePower)
	fmt.Println("Введите поддержку GSync Premium (true/false):")
	fmt.Scan(&monitor.GSyncPrem)
	fmt.Println("Введите кривизну (true/false):")
	fmt.Scan(&monitor.Curved)
	fmt.Println("Введите ID монитора:")
	fmt.Scan(&monitor.Type_Display_ID)

	monitor.DisplayMonitor = getDisplayInfo()
	addMonitor(monitor, headers)

	getAll(headers)

	var id string
	fmt.Println("Введите ID монитора:")
	fmt.Scan(&id)

	getMonitor(id, headers)
}

func registerMenu() {
	var user User
	fmt.Println("Введите имя пользователя:")
	fmt.Scan(&user.Username)
	fmt.Println("Введите пароль:")
	fmt.Scan(&user.Password)
	fmt.Println("Введите адрес электронной почты:")
	fmt.Scan(&user.Email)

	register(user)
}

func login(user User) string {
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return ""
	}

	resp, err := http.Post("http://127.0.0.1:8080/auth",
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return ""
	}

	defer resp.Body.Close()
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(token)
}

func register(user User) {
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return
	}

	resp, err := http.Post("http://127.0.0.1:8080/register",
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	defer resp.Body.Close()
}

func getDisplayInfo() Display {
	var display Display
	fmt.Println("Введите диагональ (например, 27):")
	fmt.Scan(&display.Diagonal)
	fmt.Println("Введите разрешение (например, 2560x1440):")
	fmt.Scan(&display.Resolution)
	fmt.Println("Введите тип матрицы (например, IPS):")
	fmt.Scan(&display.TypeMatrix)
	fmt.Println("Введите поддержку GSync (true/false):")
	fmt.Scan(&display.GSync)
	return display
}

func addMonitor(monitor Monitor, headers map[string]string) {
	data, err := json.Marshal(monitor)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/addMonitor", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Добавление заголовка Authorization
	req.Header.Add("Authorization", headers["Authorization"])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	defer resp.Body.Close()
}

func getAll(headers map[string]string) {
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/getAll", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	req.Header.Add("Authorization", headers["Authorization"])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}

func getMonitor(id string, headers map[string]string) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/getMonitor?id="+id, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	req.Header.Add("Authorization", headers["Authorization"])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var monitor Monitor
	err = json.Unmarshal(body, &monitor)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Monitor: %+v\n", monitor)
}
