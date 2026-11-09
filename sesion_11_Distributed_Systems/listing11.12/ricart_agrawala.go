/*
 * Archivo:     ricart_agrawala.go
 *
 * Propósito:   Algoritmo de Ricart-Agrawala para exclusión mutua distribuida
 *              (Ben-Ari, Algoritmo 10.2) sobre TCP con mensajes JSON. Un nodo
 *              entra a su sección crítica solo cuando TODOS los demás nodos le
 *              dieron permiso (mensaje REPLY).
 *
 * Compilación: go build ricart_agrawala.go
 * Ejecución:   Cuatro terminales:
 *                go run ../listing11.11/recurso_compartido.go
 *                go run ricart_agrawala.go -id 1
 *                go run ricart_agrawala.go -id 2
 *                go run ricart_agrawala.go -id 3
 *              Para ver qué pasa SIN el algoritmo, añada -sinra a los tres.
 *
 * Entrada:     Ninguna (flags: -id, -nodos, -rondas, -recurso, -sinra).
 * Salida:      Cada nodo muestra su ticket, las solicitudes y respuestas, y
 *              cuándo entra y sale de la sección crítica. En el recurso no
 *              debe aparecer ninguna VIOLACIÓN. Con -sinra sí aparecen.
 *
 * Notas:
 *    1. Correspondencia con Ben-Ari (Alg. 10.2):
 *         p2  requestCS <- true          -> requestCS = true
 *         p3  myNum <- highestNum + 1    -> myNum = highestNum + 1
 *         p4-p5 send(request) a todos    -> enviar REQUEST
 *         p6  await replies              -> <-replies, N-1 veces
 *         p9-p11 responder a diferidos   -> enviar REPLY a `diferidos`
 *         p12-p16 proceso Receive        -> manejarRequest
 *    2. Ben-Ari supone que cada proceso es atómico dentro del nodo. En Go el
 *       Main y los Receive son goroutines concurrentes, así que el estado
 *       local se protege con un sync.Mutex.
 *    3. Empates: si dos tickets son iguales, gana el id menor (el operador
 *       "<<" del libro). Ver la función `menor`.
 *    4. Mensajes por entrada a la SC: 2(N-1) (N-1 REQUEST + N-1 REPLY).
 *    5. Al terminar sus rondas, el nodo sigue vivo para responder a los
 *       demás (sección "Quiescent nodes"). Detenerlo con Ctrl+C.
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
	Tipo   string `json:"tipo"` // REQUEST, REPLY (y ENTRA/SALE con el recurso)
	Origen int    `json:"origen"`
	Num    int    `json:"num,omitempty"`
}

var (
	miID    int
	nodos   map[int]string // id -> dirección
	recurso string

	mu         sync.Mutex
	myNum      int
	highestNum int
	requestCS  bool
	diferidos  []int

	replies = make(chan int, 16)
)

func main() {
	flag.IntVar(&miID, "id", 1, "identificador de este nodo")
	lista := flag.String("nodos", "1=localhost:9101,2=localhost:9102,3=localhost:9103",
		"todos los nodos: id=host:puerto,...")
	rondas := flag.Int("rondas", 3, "veces que el nodo entra a la sección crítica")
	flag.StringVar(&recurso, "recurso", "localhost:9100", "recurso compartido (vacío = no usar)")
	sinRA := flag.Bool("sinra", false, "entrar SIN pedir permiso (para ver el fallo)")
	flag.Parse()

	nodos = leerNodos(*lista)
	ln, err := net.Listen("tcp", nodos[miID])
	if err != nil {
		fmt.Println("Error al escuchar:", err)
		os.Exit(1)
	}
	go aceptar(ln)
	log("escuchando en %s; otros nodos: %v", nodos[miID], otros())

	time.Sleep(2 * time.Second) // margen para lanzar los demás nodos
	for r := 1; r <= *rondas; r++ {
		time.Sleep(time.Duration(200+rand.Intn(800)) * time.Millisecond) // sección no crítica
		if *sinRA {
			seccionCritica(r)
			continue
		}

		// --- preprotocolo (p2-p6) ---
		mu.Lock()
		requestCS = true
		myNum = highestNum + 1
		num := myNum
		mu.Unlock()
		log("quiere entrar (ronda %d) con ticket %d", r, num)
		for _, id := range otros() {
			enviar(nodos[id], Mensaje{Tipo: "REQUEST", Origen: miID, Num: num})
		}
		for range otros() {
			de := <-replies
			log("  permiso recibido de %d", de)
		}

		seccionCritica(r)

		// --- postprotocolo (p8-p11) ---
		mu.Lock()
		requestCS = false
		pendientes := diferidos
		diferidos = nil
		mu.Unlock()
		for _, id := range pendientes {
			log("  responde al diferido %d", id)
			enviar(nodos[id], Mensaje{Tipo: "REPLY", Origen: miID})
		}
	}
	log("terminó sus %d rondas; sigue respondiendo (Ctrl+C para salir)", *rondas)
	select {}
}

func seccionCritica(r int) {
	log(">>> ENTRA a la sección crítica (ronda %d)", r)
	usarRecurso(500 * time.Millisecond)
	log("<<< SALE de la sección crítica (ronda %d)", r)
}

// manejarRequest es el proceso Receive del libro (p12-p16).
func manejarRequest(m Mensaje) {
	mu.Lock()
	if m.Num > highestNum {
		highestNum = m.Num
	}
	responder := !requestCS || menor(m.Num, m.Origen, myNum, miID)
	if !responder {
		diferidos = append(diferidos, m.Origen)
	}
	yo := myNum
	mu.Unlock()

	if responder {
		log("  REQUEST de %d (ticket %d): responde ya", m.Origen, m.Num)
		enviar(nodos[m.Origen], Mensaje{Tipo: "REPLY", Origen: miID})
	} else {
		log("  REQUEST de %d (ticket %d): lo difiere (mi ticket %d)", m.Origen, m.Num, yo)
	}
}

// menor implementa el "<<" de Ben-Ari: compara tickets y desempata por id.
func menor(num1, id1, num2, id2 int) bool {
	return num1 < num2 || (num1 == num2 && id1 < id2)
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
			case "REPLY":
				replies <- m.Origen
			}
		}(conn)
	}
}

// enviar reintenta durante unos segundos: el otro nodo quizá aún no arranca.
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

// usarRecurso avisa al recurso compartido al entrar y al salir.
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
