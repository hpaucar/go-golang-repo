/*
 * Archivo:     codeApi.go
 *
 * Propósito:   Microservicio REST mínimo en Go. Devuelve una lista de alumnos
 *              en JSON en GET /alumnos/. Sirve para practicar Dockerfile y
 *              Docker Compose.
 *
 * Compilación: go build codeApi.go
 * Ejecución:   Local:         go run codeApi.go
 *              Docker simple: docker build -t api-simple -f Dockerfile .
 *                             docker run --rm -p 9015:9015 api-simple
 *              Multi-stage:   docker build -t api-multi -f Dockerfile.multistage .
 *                             docker run --rm -p 9015:9015 api-multi
 *              Compose:       docker compose up -d --build
 *              Probar:        curl http://localhost:9015/alumnos/
 *
 * Entrada:     Peticiones HTTP en el puerto 9015 (o el de la variable PORT).
 * Salida:      JSON con tres alumnos; el servidor escribe "Invoca listar
 *              alumnos" en su log por cada petición.
 *
 * Notas:
 *    1. Código de la guía de Docker Compose, con un cambio: el puerto se lee
 *       de la variable de entorno PORT (la guía recomienda no fijarlo en el
 *       código). Si PORT no existe, se usa 9015.
 *    2. Comparar el tamaño de las imágenes con `docker images`: la imagen
 *       simple (golang:alpine) pesa cientos de MB y la multi-stage, pocos MB.
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Alumno struct {
	Codigo  string `json:"cod"`
	Nombres string `json:"nom"`
	Dni     int    `json:"dni"`
}

var alumnos = []Alumno{
	{"123", "Juan", 45678912},
	{"456", "Luis", 78945612},
	{"789", "Jorge", 12345678},
}

func listarAlumnos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alumnos)
	log.Println("Invoca listar alumnos")
}

func main() {
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "9015"
	}
	http.HandleFunc("/alumnos/", listarAlumnos)
	fmt.Printf("Servidor en http://localhost:%s/alumnos/\n", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, nil))
}
