/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo calcular el hash SHA-256 de una cadena con el
 *              paquete `crypto/sha256`: se crea el hash, se escriben los bytes
 *              con `Write` y se obtiene el resultado con `Sum`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El hash SHA-256 de "Hola, mundo!" en hexadecimal (64
 *              caracteres), siempre el mismo.
 *
 * Notas:
 *    1. Un hash es determinista: la misma entrada siempre produce la misma
 *       salida, y un cambio mínimo en la entrada cambia todo el resultado.
 *    2. `hash.Sum(nil)` devuelve el resumen como `[]byte`, y el verbo `%x` lo
 *       imprime en hexadecimal. Para un solo bloque de datos basta con
 *       `sha256.Sum256(datos)`.
 *    3. SHA-256 no sirve para guardar contraseñas; para eso se usan funciones
 *       lentas como bcrypt o argon2.
 */

package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	cadenaBase := "Hola, mundo!"

	hash := sha256.New()           // Crear un nuevo hash SHA-256
	hash.Write([]byte(cadenaBase)) // Escribir la cadena en el hash

	// Obtener el resultado del hash
	resultado := hash.Sum(nil)

	fmt.Printf("Hash SHA-256 de la cadena '%s': %x\n", cadenaBase, resultado)
}
