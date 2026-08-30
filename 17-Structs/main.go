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
