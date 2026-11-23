/*
 * Archivo:     consenso_una_ronda.go
 *
 * Propósito:   Algoritmo de consenso de UNA ronda (Ben-Ari, Algoritmo 12.1)
 *              sobre TCP con mensajes JSON. Cada general elige un plan
 *              (A = atacar, R = retirarse), lo envía a todos los demás, recibe
 *              los planes de ellos y decide por mayoría. Sirve para ver que
 *              este algoritmo FALLA si un general es traidor.
 *
 * Compilación: go build consenso_una_ronda.go
 * Ejecución:   Tres terminales (escenario de Ben-Ari, sección 12.3):
 *                go run consenso_una_ronda.go -yo Basil -plan A -caida 1
 *                go run consenso_una_ronda.go -yo Leo   -plan R
 *                go run consenso_una_ronda.go -yo Zoe   -plan A
 *              Basil es el traidor: cae después de enviar 1 mensaje (a Leo).
 *              Para un traidor bizantino, cambie Basil por:
 *                go run consenso_una_ronda.go -yo Basil -mentira Leo=A,Zoe=R
 *
 * Entrada:     Ninguna (flags: -yo, -plan, -generales, -host, -espera,
 *              -caida, -mentira).
 * Salida:      Cada general muestra los planes que conoce y su decisión.
 *              Con un traidor, Leo decide A y Zoe decide R: NO hay consenso.
 *
 * Notas:
 *    1. Correspondencia con Ben-Ari (Alg. 12.1):
 *         p1  plan[myID] <- chooseAttackOrRetreat -> flag -plan
 *         p2-p3 send(G, myID, plan[myID])         -> ronda 1, enviar a otros
 *         p4-p5 receive(G, plan[G])               -> recibir(1, N-1)
 *         p6  finalPlan <- majority(plan)         -> mayoria(plan)
 *    2. El libro supone rondas síncronas: si un mensaje no llega en un tiempo
 *       fijo (-espera), el emisor se considera caído y su plan queda como "-".
 *    3. Empates: la mayoría se rompe a favor de R (retirarse), como en el libro.
 *    4. Los generales tienen los nombres del libro y puertos fijos:
 *       Basil 9201, John 9202, Leo 9203, Mike 9204, Zoe 9205.
 *    5. Traidores: -caida k = el general deja de enviar después de k mensajes
 *       (falla por caída). -mentira G=P,... = envía el plan P al general G
 *       (falla bizantina). La decisión de un traidor no cuenta.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Mensaje struct {
	Ronda  int    `json:"ronda"`
	Origen string `json:"origen"`
	Plan   string `json:"plan"` // "A" (atacar) o "R" (retirarse)
}

var (
	yo        string
	generales map[string]string // nombre -> dirección
	espera    time.Duration
	caida     int               // -1 = leal; k = cae tras enviar k mensajes
	mentira   map[string]string // destino -> plan falso (traidor bizantino)

	muEnvio  sync.Mutex
	enviados int
	caido    bool
)

func main() {
	flag.StringVar(&yo, "yo", "Leo", "nombre de este general")
	miPlan := flag.String("plan", "A", "plan inicial: A (atacar) o R (retirarse)")
	lista := flag.String("generales", "Basil,Leo,Zoe", "todos los generales, separados por coma")
	host := flag.String("host", "localhost", "host de los generales (si no se da host:puerto)")
	flag.DurationVar(&espera, "espera", 2*time.Second, "tiempo máximo de espera por ronda")
	flag.IntVar(&caida, "caida", -1, "TRAIDOR: deja de enviar después de k mensajes")
	textoMentira := flag.String("mentira", "", "TRAIDOR: plan falso por destino, ej. Leo=A,Zoe=R")
	flag.Parse()

	generales = leerGenerales(*lista, *host)
	mentira = leerMapa(*textoMentira)
	escuchar()
	esperarATodos()

	// p1: cada general elige su plan
	plan := map[string]string{yo: strings.ToUpper(*miPlan)}
	log("mi plan es %s", plan[yo])

	// p2-p3: ronda única, enviar mi plan a los demás
	for _, g := range otros() {
		enviar(g, Mensaje{Ronda: 1, Origen: yo, Plan: planPara(g, plan[yo])})
	}
	// p4-p5: recibir el plan de cada uno de los demás
	for _, m := range recibir(1, len(otros())) {
		plan[m.Origen] = m.Plan
		log("  recibe de %s: %s", m.Origen, m.Plan)
	}

	// p6: decisión por mayoría
	decidir(mostrarTabla(plan))
}

// mostrarTabla imprime la tabla del libro (general | plan) y devuelve la mayoría.
func mostrarTabla(plan map[string]string) string {
	fmt.Printf("\n  %-8s %s\n", "general", "plan")
	var valores []string
	for _, g := range nombres() {
		p := plan[g]
		valores = append(valores, p)
		fmt.Printf("  %-8s %s\n", g, guion(p))
	}
	final := mayoria(valores)
	fmt.Printf("  %-8s %s\n\n", "mayoría", final)
	return final
}

// ---------------------------------------------------------------------
// Decisión, mayoría y traidores
// ---------------------------------------------------------------------

// mayoria cuenta A y R (ignora los "-") y rompe empates a favor de R.
func mayoria(planes []string) string {
	a, r := 0, 0
	for _, p := range planes {
		switch p {
		case "A":
			a++
		case "R":
			r++
		}
	}
	if a > r {
		return "A"
	}
	return "R"
}

func decidir(final string) {
	if esTraidor() {
		log("soy el TRAIDOR: mi decisión (%s) no cuenta", final)
	} else {
		log("DECISIÓN FINAL: %s", final)
	}
	log("terminado (Ctrl+C para salir)")
	select {} // sigue vivo: los demás aún pueden enviarle mensajes
}

func esTraidor() bool { return caida >= 0 || len(mentira) > 0 }

// planPara devuelve el plan que este general envía a `destino`.
// Un general leal envía siempre el verdadero; un bizantino puede mentir.
func planPara(destino, verdadero string) string {
	if falso, ok := mentira[destino]; ok {
		return falso
	}
	return verdadero
}

func guion(p string) string {
	if p == "" {
		return "-"
	}
	return p
}

// ---------------------------------------------------------------------
// Red: buzón por ronda, recepción con tiempo límite y envío
// ---------------------------------------------------------------------

var (
	muBuzon sync.Mutex
	buzon   = map[int]chan Mensaje{} // ronda -> mensajes de esa ronda
)

func canal(ronda int) chan Mensaje {
	muBuzon.Lock()
	defer muBuzon.Unlock()
	if buzon[ronda] == nil {
		buzon[ronda] = make(chan Mensaje, 256)
	}
	return buzon[ronda]
}

func escuchar() {
	ln, err := net.Listen("tcp", generales[yo])
	if err != nil {
		fmt.Println("Error al escuchar:", err)
		os.Exit(1)
	}
	log("escuchando en %s; otros generales: %v", generales[yo], otros())
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go func(c net.Conn) {
				defer c.Close()
				var m Mensaje
				if json.NewDecoder(c).Decode(&m) == nil {
					canal(m.Ronda) <- m
				}
			}(conn)
		}
	}()
}

// esperarATodos espera a que todos los generales estén escuchando, para que
// la ronda 1 empiece a la vez (rondas síncronas del libro).
func esperarATodos() {
	for _, g := range otros() {
		for {
			conn, err := net.DialTimeout("tcp", generales[g], time.Second)
			if err == nil {
				conn.Close()
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
	log("todos los generales están listos")
	time.Sleep(300 * time.Millisecond)
}

// recibir junta hasta `cuantos` mensajes de la ronda. Si pasa el tiempo
// -espera sin que lleguen todos, supone que los que faltan están caídos.
func recibir(ronda, cuantos int) []Mensaje {
	var lista []Mensaje
	limite := time.After(espera)
	for len(lista) < cuantos {
		select {
		case m := <-canal(ronda):
			lista = append(lista, m)
		case <-limite:
			log("  ronda %d: llegaron %d de %d mensajes (los demás, caídos)",
				ronda, len(lista), cuantos)
			return lista
		}
	}
	return lista
}

// enviar entrega un mensaje. Un traidor con -caida k deja de enviar tras k.
func enviar(destino string, m Mensaje) {
	muEnvio.Lock()
	if caida >= 0 && enviados >= caida {
		if !caido {
			log("*** CAÍDO: ya no envía más mensajes ***")
		}
		caido = true
		muEnvio.Unlock()
		return
	}
	enviados++
	muEnvio.Unlock()

	for intento := 0; intento < 50; intento++ {
		conn, err := net.DialTimeout("tcp", generales[destino], time.Second)
		if err == nil {
			json.NewEncoder(conn).Encode(m)
			conn.Close()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	log("ERROR: no se pudo enviar a %s", destino)
}

// ---------------------------------------------------------------------
// Utilidades
// ---------------------------------------------------------------------

var puertos = map[string]int{"Basil": 9201, "John": 9202, "Leo": 9203, "Mike": 9204, "Zoe": 9205}

// leerGenerales acepta "Basil,Leo,Zoe" (puertos fijos) o "Basil=host:puerto,...".
func leerGenerales(s, host string) map[string]string {
	m := map[string]string{}
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if nombre, dir, ok := strings.Cut(item, "="); ok {
			m[nombre] = dir
		} else if p, ok := puertos[item]; ok {
			m[item] = fmt.Sprintf("%s:%d", host, p)
		} else {
			fmt.Println("General desconocido (use nombre=host:puerto):", item)
			os.Exit(1)
		}
	}
	if _, ok := m[yo]; !ok {
		fmt.Printf("%s no está en -generales\n", yo)
		os.Exit(1)
	}
	return m
}

func leerMapa(s string) map[string]string {
	m := map[string]string{}
	for _, par := range strings.Split(s, ",") {
		if k, v, ok := strings.Cut(strings.TrimSpace(par), "="); ok {
			m[k] = strings.ToUpper(v)
		}
	}
	return m
}

func nombres() []string {
	var n []string
	for g := range generales {
		n = append(n, g)
	}
	sort.Strings(n)
	return n
}

func otros() []string {
	var n []string
	for _, g := range nombres() {
		if g != yo {
			n = append(n, g)
		}
	}
	return n
}

func log(formato string, args ...any) {
	fmt.Printf("%s [%s] %s\n", time.Now().Format("15:04:05.000"), yo,
		fmt.Sprintf(formato, args...))
}
