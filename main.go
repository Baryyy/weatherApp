package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	ApiKey string `json:"ApiKey"`
}

type weatherData struct {
	Name string `json: name`
	Main struct {
		Kelvin float64 `json: "temp"`
	} `json"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var x apiConfigData
	err = json.Unmarshal(bytes, &x)
	if err != nil {
		return apiConfigData{}, err
	}
	return x, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello there!\n"))
}
func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	response, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.ApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer response.Body.Close()

	var d weatherData
	if err := json.NewDecoder(response.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}
func main() {
	http.HandleFunc("/hello", sayHello)

	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, req *http.Request) {
			city := strings.SplitN(req.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=8")
			json.NewEncoder(w).Encode(data)
		})

}
