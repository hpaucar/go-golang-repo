/*
 * Archivo:     conway_compressor.go
 *
 * Propósito:   Nodo 1 del problema de Conway distribuido (Ben-Ari, Alg. 8.2).
 *              Recibe una cadena en JSON, la comprime con Run-Length Encoding
 *              (RLE: "AAAAA" -> "5A") y envía el resultado en JSON al nodo de
 *              salida (listing11.05).
 *
 * Compilación: go build conway_compressor.go
 * Ejecución:   Terminal 1: go run ../listing11.05/conway_output.go
 *              Terminal 2: go run conway_compressor.go
 *              Terminal 3: go run ../listing11.03/centinela_json.go -modo conway
 *
 * Entrada:     {"tipo":"DATA","contenido":"AAAAABBBCCDAAAAA"} en el puerto 9001.
 * Salida:      "[NODO COMPRESS] Entrada: ... -> Salida RLE: 5A3B2C1D5A" y el
 *              envío del resultado al puerto 9002.
 *
 * Notas:
 *    1. Cambios respecto a la guía: se rellena `Origen` (antes el nodo de
 *       salida mostraba un origen vacío), se usa strings.Builder en vez de
 *       concatenar strings y se recorre la cadena por runas, para no romper
 *       caracteres como 'Ñ'.
 *    2. Diferencia con Ben-Ari: en el Alg. 8.2 una letra aislada se deja tal
 *       cual ("D", no "1D") y las corridas se cortan en 9. Aquí se sigue la
 *       guía ("1D"). Ejercicio: implementar la regla de Ben-Ari.
 *    3. Compresión y presentación quedan desacopladas en dos procesos: es un
 *       pipeline de dos etapas comunicado por paso de mensajes.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Mensaje struct {
	Tipo      string `json:"tipo"`
	Origen    string `json:"origen,omitempty"`
	Contenido string `json:"contenido,omitempty"`
}

func main() {
	local := flag.String("local", "localhost:9001", "dirección de escucha")
	salida := flag.String("salida", "localhost:9002", "dirección del nodo de salida")
	flag.Parse()

	ln, err := net.Listen("tcp", *local)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("[CONWAY COMPRESSOR] Escuchando en %s; salida -> %s\n", *local, *salida)

	for {
		conn, err := ln.Accept()
		if err == nil {
			go procesarYComprimir(conn, *salida)
		}
	}
}

func procesarYComprimir(conn net.Conn, salida string) {
	defer conn.Close()
	var msg Mensaje
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		fmt.Println("[NODO COMPRESS] JSON inválido:", err)
		return
	}
	resultado := rle(msg.Contenido)
	if resultado == "" {
		return
	}
	fmt.Printf("[NODO COMPRESS] Entrada: %s -> Salida RLE: %s\n", msg.Contenido, resultado)
	enviar(salida, Mensaje{Tipo: "RLE_OUTPUT", Origen: "COMPRESSOR", Contenido: resultado})
}

// rle comprime corridas de caracteres iguales: "AAAB" -> "3A1B". O(n).
func rle(s string) string {
	runas := []rune(s)
	if len(runas) == 0 {
		return ""
	}
	var b strings.Builder
	actual, cuenta := runas[0], 1
	for _, r := range runas[1:] {
		if r == actual {
			cuenta++
			continue
		}
		b.WriteString(strconv.Itoa(cuenta))
		b.WriteRune(actual)
		actual, cuenta = r, 1
	}
	b.WriteString(strconv.Itoa(cuenta))
	b.WriteRune(actual)
	return b.String()
}

func enviar(dir string, msg Mensaje) {
	conn, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Printf("[NODO COMPRESS] No se pudo enviar a %s: %v\n", dir, err)
		return
	}
	defer conn.Close()
	json.NewEncoder(conn).Encode(msg)
}
