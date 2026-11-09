/*
 * Archivo:     conway_output.go
 *
 * Propósito:   Nodo 2 del problema de Conway distribuido. Recibe la cadena
 *              comprimida en JSON y la imprime en líneas de ancho fijo K, como
 *              en el enunciado original de Conway (Ben-Ari, Alg. 8.2).
 *
 * Compilación: go build conway_output.go
 * Ejecución:   go run conway_output.go            (ancho por defecto: 4)
 *              go run conway_output.go -k 6
 *              Se lanza ANTES que el compresor (listing11.04).
 *
 * Entrada:     {"tipo":"RLE_OUTPUT","origen":"COMPRESSOR","contenido":"..."}
 *              en el puerto 9002.
 * Salida:      Un bloque con la hora, el tipo, el origen y la cadena RLE
 *              partida en líneas de K caracteres. Con K=4 y "5A3B2C1D5A":
 *                5A3B
 *                2C1D
 *                5A
 *
 * Notas:
 *    1. La guía original solo mostraba la cadena; el problema de Conway
 *       también pide formatear la salida en líneas de longitud fija. Por eso
 *       se añade el flag -k.
 *    2. Cada conexión se atiende en su propia goroutine; si llegan varios
 *       mensajes a la vez, sus bloques pueden mezclarse en pantalla.
 *       Ejercicio: evitarlo con un sync.Mutex.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

type Mensaje struct {
	Tipo      string `json:"tipo"`
	Origen    string `json:"origen,omitempty"`
	Contenido string `json:"contenido,omitempty"`
}

func main() {
	local := flag.String("local", "localhost:9002", "dirección de escucha")
	k := flag.Int("k", 4, "caracteres por línea")
	flag.Parse()

	ln, err := net.Listen("tcp", *local)
	if err != nil {
		fmt.Printf("[NODO OUTPUT ERROR] No se pudo escuchar en %s: %v\n", *local, err)
		return
	}
	defer ln.Close()
	fmt.Printf("[NODO 2 - OUTPUT] Escuchando en %s (K = %d)\n", *local, *k)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleOutput(conn, *k)
	}
}

func handleOutput(conn net.Conn, k int) {
	defer conn.Close()
	var msg Mensaje
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		fmt.Printf("[NODO OUTPUT ERROR] JSON inválido: %v\n", err)
		return
	}
	sep := strings.Repeat("-", 40)
	fmt.Println(sep)
	fmt.Printf("[%s] PAQUETE RECIBIDO\n", time.Now().Format("15:04:05.000"))
	fmt.Printf("  Tipo:   %s\n  Origen: %s\n  Salida en líneas de %d:\n",
		msg.Tipo, msg.Origen, k)
	for _, linea := range partir(msg.Contenido, k) {
		fmt.Println("    " + linea)
	}
	fmt.Println(sep)
}

// partir divide s en trozos de k runas.
func partir(s string, k int) []string {
	runas := []rune(s)
	var lineas []string
	for i := 0; i < len(runas); i += k {
		fin := i + k
		if fin > len(runas) {
			fin = len(runas)
		}
		lineas = append(lineas, string(runas[i:fin]))
	}
	return lineas
}
