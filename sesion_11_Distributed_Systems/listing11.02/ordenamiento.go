/*
 * Archivo:     ordenamiento.go
 *
 * Propósito:   Ordenamiento distribuido en cadena (versión del curso, en texto
 *              plano). Cada nodo guarda el menor número recibido y reenvía el
 *              mayor al siguiente nodo. Cuando un nodo recibe `n` números,
 *              imprime el valor que se quedó: el nodo k termina con el k-ésimo
 *              menor.
 *
 * Compilación: go build ordenamiento.go
 * Ejecución:   Cuatro terminales (4 números => 4 nodos):
 *                go run ordenamiento.go -local 9080 -remoto 9081 -n 4
 *                go run ordenamiento.go -local 9081 -remoto 9082 -n 3
 *                go run ordenamiento.go -local 9082 -remoto 9083 -n 2
 *                go run ordenamiento.go -local 9083 -n 1      (último nodo)
 *              Luego: go run ../listing11.01/centinela.go
 *
 * Entrada:     Números enteros, uno por conexión TCP.
 * Salida:      "Número recibido = X" por cada número y, al completar `n`,
 *              "[NODO p] valor final = X". Con 9, 12, 10, 2 se obtiene
 *              2, 9, 10, 12 en los nodos 9080..9083.
 *
 * Notas:
 *    1. Revisión del ordenamiento-dist.go original:
 *       - Se usan flags en lugar de leer por consola, para lanzar los nodos
 *         rápido.
 *       - El original no tenía "último nodo": siempre llamaba a net.Dial y,
 *         si nadie escuchaba, conn era nil y el programa caía con panic.
 *         Aquí `-remoto ""` marca el último nodo.
 *       - Se revisan los errores de Listen, Accept y Dial.
 *    2. `chCont` es un canal con búfer de tamaño 1 usado como cerrojo: quien
 *       saca el valor entra a la sección crítica y lo devuelve al salir. Así
 *       se protegen `min` y el contador entre goroutines del mismo nodo.
 *    3. Memoria distribuida: los nodos NO comparten variables; solo se
 *       comunican por mensajes TCP.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	dirLocal  string
	dirRemoto string
	chCont    chan int // canal-cerrojo que guarda el número de valores recibidos
	num       int      // cuántos valores recibirá este nodo
	min       int      // menor valor guardado por este nodo
)

func main() {
	portLocal := flag.String("local", "9080", "puerto local")
	portRemoto := flag.String("remoto", "", "puerto del siguiente nodo (vacío = último nodo)")
	flag.IntVar(&num, "n", 1, "número de valores que recibirá este nodo")
	flag.Parse()

	dirLocal = "localhost:" + *portLocal
	if *portRemoto != "" {
		dirRemoto = "localhost:" + *portRemoto
	}

	chCont = make(chan int, 1)
	chCont <- 0

	ls, err := net.Listen("tcp", dirLocal)
	if err != nil {
		fmt.Println("Error al escuchar:", err)
		os.Exit(1)
	}
	defer ls.Close()
	fmt.Printf("[NODO %s] escuchando; siguiente = %q; espera %d valores\n",
		*portLocal, dirRemoto, num)

	for {
		conn, err := ls.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	dato, _ := bufio.NewReader(conn).ReadString('\n')
	dNum, err := strconv.Atoi(strings.TrimSpace(dato))
	if err != nil {
		fmt.Println("Dato inválido:", dato)
		return
	}
	fmt.Printf("Número recibido = %d\n", dNum)

	cont := <-chCont // entrar a la sección crítica

	if cont == 0 {
		min = dNum // primer valor: se guarda
	} else if dNum < min {
		enviar(min) // el antiguo mínimo pasa al siguiente nodo
		min = dNum
	} else {
		enviar(dNum) // el recibido es mayor: pasa al siguiente nodo
	}

	cont++
	if cont == num {
		fmt.Printf("[NODO %s] valor final = %d\n", dirLocal, min)
		cont = 0
	}

	chCont <- cont // salir de la sección crítica
}

func enviar(numero int) {
	if dirRemoto == "" {
		fmt.Printf("[NODO %s] soy el último nodo; descarto %d (¿-n mal configurado?)\n",
			dirLocal, numero)
		return
	}
	conn, err := net.Dial("tcp", dirRemoto)
	if err != nil {
		fmt.Printf("[NODO %s] no se pudo enviar %d a %s: %v\n", dirLocal, numero, dirRemoto, err)
		return
	}
	defer conn.Close()
	fmt.Fprintln(conn, numero)
}
