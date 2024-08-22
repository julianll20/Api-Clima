package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

// Estructura para almacenar los datos de la API
type WeatherData struct {
	Name string `json:"name"` // Nombre de la ciudad
	Main struct {
		Temp float64 `json:"temp"` // Temperatura en Celsius
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"` // Descripción del clima
		Icon        string `json:"icon"`        // Icono del clima
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"` // Velocidad del viento en m/s
	} `json:"wind"`
}

func main() {
	// Manejo de la ruta principal
	http.HandleFunc("/", weatherHandler)

	// Servir archivos estáticos (CSS, imágenes, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Iniciar el servidor
	fmt.Println("Servidor escuchando en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}
}

// Handler para mostrar la página con el clima
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// API Key visible en el código
	apiKey := "8c9150bf4cc8eb2d53798608be82e53e" // Aquí puedes colocar tu API Key directamente

	// Ciudad a consultar
	city := "Mendoza" // Puedes cambiar esta ciudad si lo deseas

	// URL para la consulta a la API (idioma español y unidades métricas)
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=es&units=metric", city, apiKey)

	// Solicitud HTTP GET a la API de OpenWeather
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error al obtener datos del clima", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Leer la respuesta del cuerpo de la solicitud
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Error al leer la respuesta del servidor", http.StatusInternalServerError)
		return
	}

	// Imprimir la respuesta JSON en la consola para depuración
	fmt.Println(string(body))

	// Deserializar los datos JSON
	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		http.Error(w, "Error al deserializar los datos JSON", http.StatusInternalServerError)
		return
	}

	// Verificar si el clima tiene datos
	if len(weatherData.Weather) == 0 {
		http.Error(w, "Datos del clima no disponibles", http.StatusInternalServerError)
		return
	}

	// Renderizar la plantilla y pasar los datos del clima
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al cargar la plantilla: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, weatherData)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}
}
