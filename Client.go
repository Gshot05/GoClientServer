// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
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
// 	Voltage        float32 `json:"voltage"`
// 	DisplayMonitor Display `json:"display"`
// 	GSyncPrem      bool    `json:"gsync_prem"`
// 	Curved         bool    `json:"curved"`
// 	TypeDisplayID  int     `json:"type_display_id"`
// }

// func main() {
// 	go addMonitor()
// 	go addDisplay()
// 	getAll()
// }

// func addDisplay() {
// 	var display Display

// 	fmt.Println("Введите данные для дисплея:")
// 	fmt.Print("Диагональ: ")
// 	_, err := fmt.Scanf("%f", &display.Diagonal)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Print("Разрешение: ")
// 	_, err = fmt.Scanf("%s", &display.Resolution)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Print("Тип матрицы: ")
// 	_, err = fmt.Scanf("%s", &display.TypeMatrix)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Print("Поддержка GSync (true/false): ")
// 	_, err = fmt.Scanf("%t", &display.GSync)
// 	if err != nil {
// 		log.Fatal(err)
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
// 	var display Display

// 	fmt.Println("Введите данные для монитора:")
// 	fmt.Print("Диагональ: ")
// 	_, err := fmt.Scanf("%f", &display.Diagonal)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Print("Разрешение: ")
// 	_, err = fmt.Scanf("%s", &display.Resolution)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Print("Тип матрицы: ")
// 	_, err = fmt.Scanf("%s", &display.TypeMatrix)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Print("Поддержка GSync (true/false): ")
// 	_, err = fmt.Scanf("%t", &display.GSync)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var voltage float32
// 	fmt.Print("Напряжение: ")
// 	_, err = fmt.Scanf("%f", &voltage)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var gsyncPrem bool
// 	fmt.Print("Премиум GSync (true/false): ")
// 	_, err = fmt.Scanf("%t", &gsyncPrem)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var curved bool
// 	fmt.Print("Изогнутый экран (true/false): ")
// 	_, err = fmt.Scanf("%t", &curved)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var typeDisplayID int
// 	fmt.Print("ID типа дисплея: ")
// 	_, err = fmt.Scanf("%d", &typeDisplayID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	monitor := Monitor{
// 		Voltage:        voltage,
// 		DisplayMonitor: display,
// 		GSyncPrem:      gsyncPrem,
// 		Curved:         curved,
// 		TypeDisplayID:  typeDisplayID,
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
	go addMonitor()
	go addDisplay()
	getAll()
}

func addDisplay() {
	display := Display{
		Diagonal:   27,
		Resolution: "2560x1440",
		TypeMatrix: "IPS",
		GSync:      true,
	}

	data, err := json.Marshal(display)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://127.0.0.1:8080/addDisplay",
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
}

func addMonitor() {
	display := Display{
		Diagonal:   27,
		Resolution: "2560x1440",
		TypeMatrix: "IPS",
		GSync:      true,
	}
	monitor := Monitor{
		VoltagePower:   220,
		DisplayMonitor: display,
		GSyncPrem:      true,
		Curved:         false,
	}

	data, err := json.Marshal(monitor)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://127.0.0.1:8080/addMonitor",
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		log.Fatal(err)
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
