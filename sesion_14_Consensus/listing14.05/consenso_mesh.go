/*
 * Archivo:     consenso_mesh.go
 *
 * Propósito:   Consenso por votación de mayoría en una malla completa (full
 *              mesh) de N nodos, sobre TCP con mensajes JSON. Al recibir la
 *              orden TRIGGER del centinela, cada nodo registra su voto, lo
 *              difunde a sus N-1 pares y, cuando tiene los N votos, decide
 *              GO A o GO B. Es el algoritmo de una ronda de Ben-Ari (Alg.
 *              12.1) en el caso SIN fallas.
 *
 * Compilación: go build consenso_mesh.go
 * Ejecución:   Seis terminales en el host (5 nodos + centinela, desde la
 *              carpeta de la sesión):
 *                P=localhost:8001,localhost:8002,localhost:8003,localhost:8004,localhost:8005
 *                go run listing14.05/consenso_mesh.go -name nodo1 -host localhost -port 8001 -vote 1 -peers $P
 *                go run listing14.05/consenso_mesh.go -name nodo2 -host localhost -port 8002 -vote 2 -peers $P
 *                go run listing14.05/consenso_mesh.go -name nodo3 -host localhost -port 8003 -vote 1 -peers $P
 *                go run listing14.05/consenso_mesh.go -name nodo4 -host localhost -port 8004 -vote 1 -peers $P
 *                go run listing14.05/consenso_mesh.go -name nodo5 -host localhost -port 8005 -vote 2 -peers $P
 *                go run listing14.06/centinela_consenso.go -nodes $P
 *              En Docker Compose: ver listing14.07.
 *
 * Entrada:     Ninguna (flags -name, -host, -port, -peers, -vote, -delay_ms,
 *              o las variables NODE_NAME, NODE_HOST, NODE_PORT, PEERS,
 *              NODE_VOTE, DELAY_MS).
 * Salida:      Cada nodo muestra los votos que recibe y el conteo final.
 *              Los 5 nodos imprimen el mismo resultado: A = 3, B = 2, GO A.
 *
 * Notas:
 *    1. Parte del código de la guía de laboratorio "Algoritmo de Consenso
 *       Distribuido P2P en Topología Malla" (Ms. Carlos Jara García), que a
 *       su vez corrige 143consenso.go: cada nodo cuenta también su propio
 *       voto y espera N votos, no N-1.
 *    2. Correcciones respecto a la guía:
 *         a) La guía dice que NODE_PORT, NODE_VOTE y DELAY_MS se pueden dar
 *            como variables de entorno, pero el código solo las leía para
 *            el nombre y los pares. Ahora se leen todas.
 *         b) La dirección propia era siempre nombre:puerto. Con -host se
 *            puede usar localhost en el host (sin Docker).
 *         c) La lista -peers puede incluir al propio nodo: se descarta. Así
 *            todos los nodos usan la misma lista.
 *         d) Si el puerto está ocupado, el programa termina con un mensaje
 *            claro (antes: log.Panicf).
 *    3. El mapa de votos se protege con un canal de capacidad 1 (patrón
 *       monitor), igual que en la guía: quien lo saca del canal es el único
 *       que lo puede tocar, y al terminar lo devuelve.
 *    4. Sin traidores ni caídas. Si un nodo cae, los demás esperan su voto
 *       para siempre (ejercicio 2 de la guía: quórum con timeout).
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Códigos de mensaje y opciones (las mismas constantes de 143consenso.go)
const (
	cnum     = iota // 0: voto entre pares
	opa             // 1: Opción A
	opb             // 2: Opción B
	ctrigger        // 3: orden de inicio enviada por el centinela
)

// tmsg es el mensaje JSON que viaja por TCP
type tmsg struct {
	Code int    `json:"code"`
	Addr string `json:"addr"`
	Op   int    `json:"op"`
}

var (
	localAddr string
	nodeName  string
	myVote    int
	delayMs   int
	addrs     []string            // los N-1 pares
	chInfo    chan map[string]int // monitor: votos recibidos (dirección -> opción)
	ronda     int
)

func main() {
	// 1. Parámetros: flag, con respaldo en variables de entorno
	name := flag.String("name", getEnv("NODE_NAME", "nodo1"), "nombre del nodo")
	host := flag.String("host", getEnv("NODE_HOST", ""), "host propio (vacío = el nombre)")
	port := flag.Int("port", getEnvInt("NODE_PORT", 8000), "puerto TCP de escucha")
	peers := flag.String("peers", getEnv("PEERS", ""), "pares separados por coma, ej. nodo2:8000,nodo3:8000")
	vote := flag.Int("vote", getEnvInt("NODE_VOTE", 0), "voto: 1 = Opción A, 2 = Opción B, 0 = aleatorio")
	delay := flag.Int("delay_ms", getEnvInt("DELAY_MS", 300), "retardo antes de difundir el voto (ms)")
	flag.Parse()

	nodeName = *name
	delayMs = *delay
	if *host == "" {
		*host = nodeName
	}
	localAddr = fmt.Sprintf("%s:%d", *host, *port)

	// 2. Pares de la malla (sin incluirme a mí)
	for _, p := range strings.Split(*peers, ",") {
		if p = strings.TrimSpace(p); p != "" && p != localAddr {
			addrs = append(addrs, p)
		}
	}

	// 3. Voto fijo o aleatorio
	if *vote == opa || *vote == opb {
		myVote = *vote
	} else {
		myVote = rand.Intn(2) + 1 // 1 (opa) o 2 (opb)
	}

	// 4. El monitor empieza con el mapa vacío dentro del canal
	chInfo = make(chan map[string]int, 1)
	chInfo <- map[string]int{}

	fmt.Println("===============================================================")
	fmt.Printf("[%s] NODO LISTO | Addr: %s | Voto: %s (%d)\n", nodeName, localAddr, nombreOpcion(myVote), myVote)
	fmt.Printf("[%s] Malla completa (%d pares): %v\n", nodeName, len(addrs), addrs)
	fmt.Println("===============================================================")

	server(fmt.Sprintf("0.0.0.0:%d", *port))
}

func server(listenAddr string) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("[%s ERROR] No se pudo escuchar en %s: %v\n", nodeName, listenAddr, err)
		os.Exit(1)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	var msg tmsg
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		return
	}

	switch msg.Code {
	case ctrigger:
		// El centinela ordena votar: registro mi voto y lo difundo a la malla
		fmt.Printf("[%s] TRIGGER del centinela -> voto %s\n", nodeName, nombreOpcion(myVote))
		propio := tmsg{Code: cnum, Addr: localAddr, Op: myVote}
		consensus(propio)

		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		for _, addr := range addrs {
			go send(addr, propio)
		}

	case cnum:
		fmt.Printf("[%s] Voto de [%s] -> %s\n", nodeName, msg.Addr, nombreOpcion(msg.Op))
		consensus(msg)
	}
}

// consensus guarda un voto. Cuando ya están los N votos, cuenta y decide.
func consensus(msg tmsg) {
	info := <-chInfo // entra al monitor
	info[msg.Addr] = msg.Op

	total := len(addrs) + 1 // N-1 pares + este nodo
	if len(info) == total {
		ronda++
		ca, cb := 0, 0
		for _, op := range info {
			if op == opa {
				ca++
			} else {
				cb++
			}
		}
		fmt.Println("---------------------------------------------------------------")
		fmt.Printf("[%s] RONDA %d, CONTEO FINAL (%d nodos): Opción A = %d | Opción B = %d\n",
			nodeName, ronda, total, ca, cb)
		switch {
		case ca > cb:
			fmt.Printf(">>> [%s] RESULTADO DEL CONSENSO: ¡GO A! <<<\n", nodeName)
		case cb > ca:
			fmt.Printf(">>> [%s] RESULTADO DEL CONSENSO: ¡GO B! <<<\n", nodeName)
		default:
			fmt.Printf(">>> [%s] EMPATE -> desempate fijo: ¡GO B! <<<\n", nodeName)
		}
		fmt.Println("---------------------------------------------------------------")
		info = map[string]int{} // listo para otra ronda
	}
	chInfo <- info // sale del monitor
}

// send reintenta 5 veces: el par quizá aún no arrancó.
func send(remoteAddr string, msg tmsg) {
	for intento := 1; intento <= 5; intento++ {
		conn, err := net.DialTimeout("tcp", remoteAddr, 2*time.Second)
		if err == nil {
			json.NewEncoder(conn).Encode(msg)
			conn.Close()
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Printf("[%s ERROR] No se pudo enviar el voto a %s\n", nodeName, remoteAddr)
}

func nombreOpcion(op int) string {
	if op == opa {
		return "OPCIÓN A"
	}
	return "OPCIÓN B"
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if n, err := strconv.Atoi(getEnv(key, "")); err == nil {
		return n
	}
	return def
}
