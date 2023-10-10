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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func main() {
	var choice int
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1. Авторизация")
		fmt.Println("2. Регистрация")
		fmt.Println("3. Выход")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			token := login()
			if token != "" {
				performAuthenticatedActions(token)
			}
		case 2:
			register()
		case 3:
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Некорректный выбор. Пожалуйста, выберите 1, 2 или 3.")
		}
	}
}

func login() string {
	var username, password string
	fmt.Println("Введите имя пользователя:")
	fmt.Scan(&username)
	fmt.Println("Введите пароль:")
	fmt.Scan(&password)

	user := User{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return ""
	}

	resp, err := http.Post("http://127.0.0.1:8080/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var authResponse AuthResponse
		err := json.NewDecoder(resp.Body).Decode(&authResponse)
		if err != nil {
			fmt.Println("Ошибка при декодировании ответа:", err)
			return ""
		}
		fmt.Println("Авторизация успешна.")
		return authResponse.Token
	} else {
		fmt.Println("Ошибка при авторизации. Пожалуйста, проверьте имя пользователя и пароль.")
		return ""
	}
}

func register() {
	var username, password, email string
	fmt.Println("Введите имя пользователя:")
	fmt.Scan(&username)
	fmt.Println("Введите пароль:")
	fmt.Scan(&password)
	fmt.Println("Введите email:")
	fmt.Scan(&email)

	user := User{
		Username: username,
		Password: password,
		Email:    email,
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return
	}

	resp, err := http.Post("http://127.0.0.1:8080/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Регистрация успешно завершена.")
	} else {
		fmt.Println("Ошибка при регистрации. Пожалуйста, попробуйте снова.")
	}
}

func performAuthenticatedActions(token string) {
	for {
		var action int
		fmt.Println("Выберите действие:")
		fmt.Println("1. Добавить дисплей")
		fmt.Println("2. Добавить монитор")
		fmt.Println("3. Получить все мониторы")
		fmt.Println("4. Получить монитор по ID")
		fmt.Println("5. Выйти из аккаунта")
		fmt.Scan(&action)

		switch action {
		case 1:
			var display Display
			fmt.Println("Введите диагональ (например, 27):")
			fmt.Scan(&display.Diagonal)
			fmt.Println("Введите разрешение (например, 2560x1440):")
			fmt.Scan(&display.Resolution)
			fmt.Println("Введите тип матрицы (например, IPS):")
			fmt.Scan(&display.TypeMatrix)
			fmt.Println("Введите поддержку GSync (true/false):")
			fmt.Scan(&display.GSync)

			addDisplay(display, token)
		case 2:
			var monitor Monitor
			fmt.Println("Введите напряжение (например, 220):")
			fmt.Scan(&monitor.VoltagePower)
			fmt.Println("Введите поддержку GSync Premium (true/false):")
			fmt.Scan(&monitor.GSyncPrem)
			fmt.Println("Введите кривизну (true/false):")
			fmt.Scan(&monitor.Curved)
			fmt.Println("Введите ID")
			fmt.Scan(&monitor.Type_Display_ID)

			addMonitor(monitor, token)
		case 3:
			getAll(token)
		case 4:
			var id string
			fmt.Println("Введите ID монитора:")
			fmt.Scan(&id)

			getMonitor(id, token)
		case 5:
			fmt.Println("Выход из аккаунта.")
			return
		default:
			fmt.Println("Некорректный выбор. Пожалуйста, выберите 1, 2, 3, 4 или 5.")
		}
	}
}

func addDisplay(display Display, token string) {
	data, err := json.Marshal(display)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/addDisplay", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	defer resp.Body.Close()
}

func addMonitor(monitor Monitor, token string) {
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
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	defer resp.Body.Close()
}

func getAll(token string) {
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/getAll", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}
	log.Println(string(body))
}

func getMonitor(ID string, token string) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/getMonitor?id="+ID, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	var monitor Monitor
	err = json.Unmarshal(body, &monitor)
	if err != nil {
		fmt.Println("Ошибка при декодировании ответа:", err)
		return
	}

	log.Printf("Monitor: %+v\n", monitor)
}
