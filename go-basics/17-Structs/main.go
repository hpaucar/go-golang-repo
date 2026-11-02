/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar las estructuras (structs): definición del tipo `Person`,
 *              creación con literales que nombran los campos, una función
 *              constructora `newPerson` que devuelve un puntero y modificación
 *              de campos directamente o a través de un puntero a un campo.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      `{Alice 30}`, `{Bob 0}` y punteros como `&{Charlie 25}`; luego
 *              la edad de Marcela cambia de 30 a 31 y después a 32.
 *
 * Notas:
 *    1. Los campos omitidos en un literal toman su valor cero (Bob tiene edad
 *       0).
 *    2. Devolver `&Person{...}` desde una función es seguro en Go: el
 *       compilador mueve el valor al heap si hace falta (escape analysis).
 *    3. Con un puntero a struct se accede a los campos sin `*`: `marcela.age`
 *       equivale a `(*marcela).age`.
 */

package main

import "fmt"

type Person struct {
	name string
	age  int
}

// Función para crear una nueva persona
// Esta función devuelve un puntero a una nueva instancia de Person
func newPerson(name string, age int) *Person {
	return &Person{name, age}
}

func main() {
	fmt.Println("Estructuras en Go")

	fmt.Println(Person{name: "Alice", age: 30})
	fmt.Println(Person{name: "Bob"})

	fmt.Println(newPerson("Charlie", 25))

	pepito := newPerson("Pepito", 40)
	fmt.Println(pepito)

	marcela := newPerson("Marcela", 30)
	fmt.Println(marcela)

	edadMarcela := &marcela.age // Obteniendo un puntero a la edad de Marcela
	*edadMarcela = 31           // Modificando la edad de Marcela a través del puntero
	fmt.Println(marcela)

	marcela.age = 32 // Modificando la edad de Marcela directamente
	fmt.Println(marcela)
}
