/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo cargar variables de entorno desde un archivo `.env`
 *              con el paquete externo `github.com/joho/godotenv` y leerlas con
 *              `os.Getenv`. Es el único ejemplo de `go-basics` que es un módulo
 *              de Go con su propio `go.mod` (módulo `curso-go`).
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run .)
 *
 * Entrada:     El archivo `.env` en el directorio actual, con las variables
 *              `APP_NAME`, `APP_ENV` y `APP_PORT`. Se crea a partir del
 *              ejemplo: `cp .env.example .env`.
 * Salida:      El nombre, el entorno y el puerto de la aplicación. Si no existe
 *              `.env`, imprime "Error al cargar el archivo .env" y termina con
 *              código de salida 1.
 *
 * Notas:
 *    1. Compile y ejecute desde la carpeta `28-VariablesEntorno`:
 *       `godotenv.Load(".env")` busca el archivo en el directorio actual, no
 *       junto al ejecutable.
 *    2. La dependencia está declarada en `go.mod` y `go.sum`; si falta, se
 *       descarga con `go mod tidy`.
 *    3. `log.Fatal` termina el programa de inmediato, así que el segundo
 *       `log.Fatal` nunca se ejecuta. Además, la variable `error` oculta el
 *       tipo predefinido `error`; es mejor llamarla `err`.
 *    4. No suba el archivo `.env` real al repositorio, solo `.env.example`.
 */

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
