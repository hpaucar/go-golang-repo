/*
 * Archivo:     centinela.go
 *
 * Propósito:   Nodo "centinela" (disparador) en texto plano. Envía cuatro
 *              números, uno por conexión TCP, al primer nodo de la cadena de
 *              ordenamiento distribuido (listing11.02).
 *
 * Compilación: go build centinela.go
 * Ejecución:   ./centinela -destino localhost:9080
 *              (Alternativamente: go run centinela.go)
 *
 * Entrada:     Ninguna (los números 9, 12, 10 y 2 están fijos en el código).
 * Salida:      Una línea por número enviado, o un mensaje de error si el
 *              nodo destino no está escuchando.
 *
 * Notas:
 *    1. Versión revisada del centinela.go original del curso. El original
 *       ignoraba el error de net.Dial: si el nodo no estaba levantado, `conn`
 *       era nil y `defer conn.Close()` provocaba un panic.
 *    2. Cada número viaja en una conexión TCP distinta: el protocolo es una
 *       línea de texto terminada en '\n' (fmt.Fprintln).
 *    3. Primero se levantan los nodos (listing11.02) y al final el centinela.
 */

package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	destino := flag.String("destino", "localhost:9080", "dirección del primer nodo")
	flag.Parse()

	for _, n := range []int{9, 12, 10, 2} {
		enviar(*destino, n)
	}
}

// enviar abre una conexión, escribe el número como texto y la cierra.
func enviar(dir string, numero int) {
	conn, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Printf("[CENTINELA] No se pudo conectar a %s: %v\n", dir, err)
		return
	}
	defer conn.Close()
	fmt.Fprintln(conn, numero)
	fmt.Printf("[CENTINELA] Enviado %d a %s\n", numero, dir)
}
