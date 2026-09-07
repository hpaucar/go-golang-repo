package main

func main() {
	nombre := "Pepito"
	edad := 32

	if nombre == "Pepito" {
		println("Hola Pepito")
	} else {
		println("Hola desconocido")
	}

	if edad >= 18 {
		println("Eres mayor de edad")
	} else {
		println("Eres menor de edad")
	}

	if 8%2 == 0 {
		println("El número 8 es par")
	} else {
		println("El número 8 es impar")
	}

	if numero := 10; numero < 0 {
		println("El número es negativo")
	} else if numero == 0 {
		println("El número es cero")
	} else {
		println("El número es positivo")
	}
}
