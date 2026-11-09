/*
 * Archivo:     ra_token.go
 *
 * Propósito:   Algoritmo de Ricart-Agrawala con paso de token (Ben-Ari,
 *              Algoritmo 10.3; similar a Suzuki-Kasami). Solo el nodo que
 *              tiene el token puede entrar a la sección crítica. El token
 *              lleva el arreglo `granted`, y cada nodo guarda en `requested`
 *              los tickets que le pidieron.
 *
 * Compilación: go build ra_token.go
 * Ejecución:   Cuatro terminales:
 *                go run ../listing11.11/recurso_compartido.go
 *                go run ra_token.go -id 1      (empieza con el token)
 *                go run ra_token.go -id 2
 *                go run ra_token.go -id 3
 *
 * Entrada:     Ninguna (flags: -id, -nodos, -rondas, -recurso).
 * Salida:      Las solicitudes, los envíos del token y las entradas/salidas de
 *              la sección crítica. En el recurso no debe haber violaciones.
 *
 * Notas:
 *    1. Correspondencia con Ben-Ari (Alg. 10.3):
 *         p2-p7  si no tengo el token: pedirlo a todos y esperarlo
 *         p10    granted[miID] <- myNum
 *         p12    sendToken (al salir)
 *         p13-p16 Receive: requested[src] <- max(...); si tengo el token y
 *                no estoy en la SC, lo envío.
 *    2. Si el nodo ya tiene el token y nadie lo pidió, entra SIN enviar
 *       mensajes. Con RA de permisos (listing11.12) siempre cuesta 2(N-1).
 *    3. sendToken: el libro elige "algún N" con requested[N] > granted[N],
 *       y eso permite inanición (ejercicio 13 del libro). Aquí se busca en
 *       orden circular a partir de miID+1, lo que evita la inanición.
 *    4. Los nodos se numeran 1..N; los arreglos usan índices 1..N (el 0 no se
 *       usa).
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Mensaje struct {
	Tipo    string `json:"tipo"` // REQUEST, TOKEN (y ENTRA/SALE con el recurso)
	Origen  int    `json:"origen"`
	Num     int    `json:"num,omitempty"`
	Granted []int  `json:"granted,omitempty"` // viaja dentro del token
}

var (
	miID    int
	n       int // número de nodos
	nodos   map[int]string
	recurso string

	mu        sync.Mutex
	haveToken bool
	inCS      bool
	myNum     int
	requested []int
	granted   []int

	tokenCh = make(chan []int, 1)
)

func main() {
	flag.IntVar(&miID, "id", 1, "identificador de este nodo")
	lista := flag.String("nodos", "1=localhost:9101,2=localhost:9102,3=localhost:9103",
		"todos los nodos: id=host:puerto,...")
	rondas := flag.Int("rondas", 3, "veces que el nodo entra a la sección crítica")
	flag.StringVar(&recurso, "recurso", "localhost:9100", "recurso compartido (vacío = no usar)")
	flag.Parse()

	nodos = leerNodos(*lista)
	n = len(nodos)
	requested = make([]int, n+1)
	granted = make([]int, n+1)
	haveToken = miID == 1

	ln, err := net.Listen("tcp", nodos[miID])
	if err != nil {
		fmt.Println("Error al escuchar:", err)
		os.Exit(1)
	}
	go aceptar(ln)
	log("escuchando en %s; tengo el token: %v", nodos[miID], haveToken)

	time.Sleep(2 * time.Second) // margen para lanzar los demás nodos
	for r := 1; r <= *rondas; r++ {
		time.Sleep(time.Duration(200+rand.Intn(800)) * time.Millisecond) // sección no crítica

		// --- preprotocolo (p2-p8) ---
		mu.Lock()
		if haveToken {
			inCS = true
			mu.Unlock()
			log("ya tengo el token: entro sin enviar mensajes")
		} else {
			myNum++
			num := myNum
			mu.Unlock()
			log("pide el token con ticket %d (ronda %d)", num, r)
			for _, id := range otros() {
				enviar(nodos[id], Mensaje{Tipo: "REQUEST", Origen: miID, Num: num})
			}
			g := <-tokenCh // p6: receive(token, granted)
			mu.Lock()
			granted = g
			haveToken = true
			inCS = true
			mu.Unlock()
			log("  recibió el token; granted = %v", g[1:])
		}

		log(">>> ENTRA a la sección crítica (ronda %d)", r)
		usarRecurso(500 * time.Millisecond)
		log("<<< SALE de la sección crítica (ronda %d)", r)

		// --- postprotocolo (p10-p12) ---
		mu.Lock()
		granted[miID] = myNum
		inCS = false
		destino, g := elegirDestino()
		mu.Unlock()
		if destino != 0 {
			enviarToken(destino, g)
		}
	}
	log("terminó sus %d rondas; sigue atendiendo (Ctrl+C para salir)", *rondas)
	select {}
}

// manejarRequest es el proceso Receive (p13-p16).
func manejarRequest(m Mensaje) {
	mu.Lock()
	if m.Num > requested[m.Origen] {
		requested[m.Origen] = m.Num
	}
	destino := 0
	var g []int
	if haveToken && !inCS {
		destino, g = elegirDestino()
	}
	mu.Unlock()
	log("  REQUEST de %d (ticket %d)", m.Origen, m.Num)
	if destino != 0 {
		enviarToken(destino, g)
	}
}

// elegirDestino es sendToken: busca un nodo con requested > granted, en orden
// circular desde miID+1. Se llama con `mu` tomado. Si encuentra uno, suelta
// el token (haveToken = false) y devuelve una copia de granted para enviarla.
func elegirDestino() (int, []int) {
	for k := 1; k < n; k++ {
		id := (miID-1+k)%n + 1
		if requested[id] > granted[id] {
			haveToken = false
			g := make([]int, len(granted))
			copy(g, granted)
			return id, g
		}
	}
	return 0, nil
}

func enviarToken(destino int, g []int) {
	log("  envía el TOKEN al nodo %d", destino)
	enviar(nodos[destino], Mensaje{Tipo: "TOKEN", Origen: miID, Granted: g})
}

// ---------------------------------------------------------------------
// Red: recepción y envío de mensajes
// ---------------------------------------------------------------------

func aceptar(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			var m Mensaje
			if json.NewDecoder(c).Decode(&m) != nil {
				return
			}
			switch m.Tipo {
			case "REQUEST":
				manejarRequest(m)
			case "TOKEN":
				tokenCh <- m.Granted
			}
		}(conn)
	}
}

func enviar(dir string, m Mensaje) {
	for intento := 0; intento < 50; intento++ {
		conn, err := net.DialTimeout("tcp", dir, time.Second)
		if err == nil {
			json.NewEncoder(conn).Encode(m)
			conn.Close()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	log("ERROR: no se pudo enviar %s a %s", m.Tipo, dir)
}

func usarRecurso(d time.Duration) {
	if recurso == "" {
		time.Sleep(d)
		return
	}
	conn, err := net.Dial("tcp", recurso)
	if err != nil {
		log("recurso no disponible: %v", err)
		time.Sleep(d)
		return
	}
	defer conn.Close()
	enc, dec := json.NewEncoder(conn), json.NewDecoder(conn)
	var ack Mensaje
	enc.Encode(Mensaje{Tipo: "ENTRA", Origen: miID})
	dec.Decode(&ack)
	time.Sleep(d)
	enc.Encode(Mensaje{Tipo: "SALE", Origen: miID})
	dec.Decode(&ack)
}

// ---------------------------------------------------------------------
// Utilidades
// ---------------------------------------------------------------------

func leerNodos(s string) map[int]string {
	m := map[int]string{}
	for _, par := range strings.Split(s, ",") {
		partes := strings.SplitN(par, "=", 2)
		id, err := strconv.Atoi(partes[0])
		if err != nil || len(partes) != 2 {
			fmt.Println("Formato de -nodos inválido:", par)
			os.Exit(1)
		}
		m[id] = partes[1]
	}
	for id := 1; id <= len(m); id++ {
		if _, ok := m[id]; !ok {
			fmt.Println("Los ids de -nodos deben ser 1..N sin huecos")
			os.Exit(1)
		}
	}
	if _, ok := m[miID]; !ok {
		fmt.Printf("El id %d no está en -nodos\n", miID)
		os.Exit(1)
	}
	return m
}

func otros() []int {
	var ids []int
	for id := range nodos {
		if id != miID {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)
	return ids
}

func log(formato string, args ...any) {
	fmt.Printf("%s [NODO %d] %s\n", time.Now().Format("15:04:05.000"), miID,
		fmt.Sprintf(formato, args...))
}
