package main

import (
	"fmt"
	"maps"
)

// En go los mapas son una estructura de datos que permite almacenar pares clave-valor.
func main() {
	// 1. Uso de mapas

	// Declaración de un mapa vacío
	paises := make(map[string]string)
	paises["ES"] = "España"
	paises["CO"] = "Colombia"
	paises["AR"] = "Argentina"

	fmt.Println(paises)

	pais_CO := paises["CO"] // Acceso al valor asociado a la clave "CO"
	fmt.Println(pais_CO)

	// Cantidad de elementos en el mapa
	fmt.Println("Longitud del mapa:", len(paises))

	_, existePais := paises["PE"] // El guion bajo se usa para ignorar el valor de retorno
	fmt.Println("¿Existe el país PE?", existePais)

	_, existePais = paises["AR"]
	fmt.Println("¿Existe el país AR?", existePais)

	// Eliminar un elemento del mapa
	delete(paises, "AR")
	fmt.Println("Mapa después de eliminar AR:", paises)

	// Limpiar todo el mapa
	clear(paises)
	fmt.Println("Mapa después de limpiar:", paises)

	// 2. Comparaciones
	ciudades1 := map[string]string{"ES": "Madrid", "CO": "Bogotá"}
	ciudades2 := map[string]string{"ES": "Madrid", "CO": "Bogotá"}
	if maps.Equal(ciudades1, ciudades2) {
		fmt.Println("Los mapas son iguales")
	} else {
		fmt.Println("Los mapas son diferentes")
	}
}
