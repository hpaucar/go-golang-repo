package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	error := godotenv.Load(".env")
	if error != nil {
		log.Fatal("Error al cargar el archivo .env")
		log.Fatal("Agregar el archivo .env en la ruta /28-VariablesEntorno/.env")
	}

	appName := os.Getenv("APP_NAME")
	appEnv := os.Getenv("APP_ENV")
	appPort := os.Getenv("APP_PORT")

	fmt.Println("Nombre de la aplicación:", appName)
	fmt.Println("Entorno de la aplicación:", appEnv)
	fmt.Println("Puerto de la aplicación:", appPort)
}
