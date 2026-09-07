package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. Usando switch para evaluar una variable
	fmt.Println("1. Usando switch para evaluar una variable:")

	number := 3
	fmt.Println("🔢 El número:", number, ", se escribe cómo:")

	switch number {
	case 1:
		fmt.Println("Uno")
	case 2:
		fmt.Println("Dos")
	case 3:
		fmt.Println("Tres")
	default:
		fmt.Println("❓ Número desconocido")
	}

	// 2. Usando switch con múltiples casos
	fmt.Println("\n2. Usando switch con múltiples casos:")

	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("🏖️ Es fin de semana")
	default:
		fmt.Println("💼 Es un día de semana")
	}

	// 3. Usando switch sin una expresión
	fmt.Println("\n3. Usando switch sin una expresión:")

	t := time.Now()
	fmt.Println("La hora actual es:", t)

	switch {
	case t.Hour() < 12:
		fmt.Println("☀️ Buenos días")
	default:
		fmt.Println("🌇 Buenas tardes")
	}
}
