// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"
// )

// type Display struct {
// 	Diagonal   float32 `json:"diag"`
// 	Resolution string  `json:"resolution"`
// 	TypeMatrix string  `json:"type_matrix"`
// 	GSync      bool    `json:"gsync"`
// }

// type Monitor struct {
// 	VoltagePower   float32 `json:"voltage"`
// 	DisplayMonitor Display `json:"display"`
// 	GSyncPrem      bool    `json:"gsync_prem"`
// 	Curved         bool    `json:"curved"`
// }

// func main() {
// 	go addMonitor()
// 	go addDisplay()
// 	getAll()
// }

// func addDisplay() {
// 	display := Display{
// 		Diagonal:   27,
// 		Resolution: "2560x1440",
// 		TypeMatrix: "IPS",
// 		GSync:      true,
// 	}

// 	data, err := json.Marshal(display)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	resp, err := http.Post("http://127.0.0.1:8080/addDisplay",
// 		"application/json",
// 		bytes.NewBuffer(data))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer resp.Body.Close()
// }

// func addMonitor() {
// 	display := Display{
// 		Diagonal:   27,
// 		Resolution: "2560x1440",
// 		TypeMatrix: "IPS",
// 		GSync:      true,
// 	}
// 	monitor := Monitor{
// 		VoltagePower:   220,
// 		DisplayMonitor: display,
// 		GSyncPrem:      true,
// 		Curved:         false,
// 	}

// 	data, err := json.Marshal(monitor)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	resp, err := http.Post("http://127.0.0.1:8080/addMonitor",
// 		"application/json",
// 		bytes.NewBuffer(data))

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// }

// func getAll() {
// 	time.Sleep(1 * time.Second)
// 	resp, err := http.Get("http://127.0.0.1:8080/getAll")

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(string(body))
// }

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
	VoltagePower   float32 `json:"voltage"`
	DisplayMonitor Display `json:"display"`
	GSyncPrem      bool    `json:"gsync_prem"`
	Curved         bool    `json:"curved"`
}

func main() {
	var display Display
	fmt.Println("Введите диагональ (например, 27):")
	fmt.Scan(&display.Diagonal)
	fmt.Println("Введите разрешение (например, 2560x1440):")
	fmt.Scan(&display.Resolution)
	fmt.Println("Введите тип матрицы (например, IPS):")
	fmt.Scan(&display.TypeMatrix)
	fmt.Println("Введите поддержку GSync (true/false):")
	fmt.Scan(&display.GSync)

	addDisplay(display)

	var monitor Monitor
	fmt.Println("Введите напряжение (например, 220):")
	fmt.Scan(&monitor.VoltagePower)
	fmt.Println("Введите поддержку GSync Premium (true/false):")
	fmt.Scan(&monitor.GSyncPrem)
	fmt.Println("Введите кривизну (true/false):")
	fmt.Scan(&monitor.Curved)

	monitor.DisplayMonitor = display
	addMonitor(monitor)

	getAll()
}

func addDisplay(display Display) {
	data, err := json.Marshal(display)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return
	}

	resp, err := http.Post("http://127.0.0.1:8080/addDisplay",
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	defer resp.Body.Close()
}

func addMonitor(monitor Monitor) {
	data, err := json.Marshal(monitor)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return
	}

	resp, err := http.Post("http://127.0.0.1:8080/addMonitor",
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	defer resp.Body.Close()
}

func getAll() {
	time.Sleep(1 * time.Second)
	resp, err := http.Get("http://127.0.0.1:8080/getAll")

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
