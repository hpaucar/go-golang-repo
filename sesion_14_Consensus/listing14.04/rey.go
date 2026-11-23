/*
 * Archivo:     rey.go
 *
 * Propósito:   Algoritmo del Rey (Berman y Garay; Ben-Ari, Algoritmo 12.4)
 *              sobre TCP con mensajes JSON. Tolera t traidores bizantinos con
 *              n >= 4t+1 generales y envía muchos menos mensajes que el de los
 *              Generales Bizantinos. Se ejecuta en dos fases; en cada fase un
 *              general distinto es el "rey".
 *
 * Compilación: go build rey.go
 * Ejecución:   Cinco terminales; Mike es el traidor y el primer rey
 *              (Ben-Ari, sección 12.10, segundo escenario):
 *                go run rey.go -yo Basil -plan A
 *                go run rey.go -yo John  -plan A
 *                go run rey.go -yo Leo   -plan R
 *                go run rey.go -yo Zoe   -plan R
 *                go run rey.go -yo Mike  -mentira Basil=R,John=A,Leo=A,Zoe=R
 *              Para el primer escenario (Zoe rey leal primero), agregue
 *              -reyes Zoe,Mike a los cinco.
 *
 * Entrada:     Ninguna (flags: -yo, -plan, -reyes, -t, -generales, -host,
 *              -espera, -caida, -mentira).
 * Salida:      En cada ronda, la fila del libro (planes, myMajority,
 *              votesMajority, kingPlan). Al final, todos los leales tienen
 *              el mismo plan.
 *
 * Notas:
 *    1. Correspondencia con Ben-Ari (Alg. 12.4):
 *         p2  do two times                     -> for fase := 1..2
 *         p3-p6 enviar y recibir planes        -> rondas 1 y 3
 *         p7  myMajority <- majority(plan)     -> miMayoria
 *         p8  votesMajority                    -> votos
 *         p9-p12 si soy rey: enviar myMajority -> rondas 2 y 4
 *         p13-p16 si no: aceptar el plan del rey solo si votos <= n/2 + t
 *         p17 finalPlan <- plan[myID]
 *    2. Mayoría "abrumadora": votos > n/2 + t (división entera). Con n = 5 y
 *       t = 1, más de 3 votos (4 o 5). Así un traidor no puede cambiarla.
 *    3. El rey de cada fase lo fija -reyes, igual en todos los nodos. Con dos
 *       fases y un solo traidor, al menos un rey es leal.
 *    4. Mensajes: (t+1)(n+1)(n-1). Con n = 5 y t = 1: 48. El de los Generales
 *       Bizantinos necesita 4 generales y 36 mensajes para t = 1, pero crece
 *       mucho más rápido (t = 2: 392 contra 240).
 *    5. El traidor (-mentira G=P,...) envía el plan P a G en TODOS sus
 *       mensajes, también cuando es rey.
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
	Plan   string `json:"plan"` // "A" o "R"
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
	textoReyes := flag.String("reyes", "Mike,Zoe", "rey de la fase 1 y de la fase 2")
	t := flag.Int("t", 1, "número de traidores tolerados")
	lista := flag.String("generales", "Basil,John,Leo,Mike,Zoe", "todos los generales, separados por coma")
	host := flag.String("host", "localhost", "host de los generales (si no se da host:puerto)")
	flag.DurationVar(&espera, "espera", 2*time.Second, "tiempo máximo de espera por ronda")
	flag.IntVar(&caida, "caida", -1, "TRAIDOR: deja de enviar después de k mensajes")
	textoMentira := flag.String("mentira", "", "TRAIDOR: plan falso por destino, ej. Basil=R,John=A")
	flag.Parse()

	generales = leerGenerales(*lista, *host)
	mentira = leerMapa(*textoMentira)
	reyes := strings.Split(*textoReyes, ",")
	n := len(generales)
	umbral := n/2 + *t
	if len(reyes) < 2 || n < 4**t+1 {
		fmt.Println("Se necesitan 2 reyes y al menos 4t+1 generales")
		os.Exit(1)
	}
	escuchar()
	esperarATodos()

	plan := map[string]string{yo: strings.ToUpper(*miPlan)}
	log("mi plan es %s; reyes: %v; mayoría abrumadora: más de %d votos", plan[yo], reyes[:2], umbral)
	encabezado()

	for fase := 1; fase <= 2; fase++ { // p2
		rImpar, rPar := 2*fase-1, 2*fase

		// p3-p6: todos envían su plan a todos
		for _, g := range otros() {
			enviar(g, Mensaje{Ronda: rImpar, Origen: yo, Plan: planPara(g, plan[yo])})
		}
		for g := range plan {
			if g != yo {
				delete(plan, g) // los planes de otros se vuelven a recibir
			}
		}
		for _, m := range recibir(rImpar, n-1) {
			plan[m.Origen] = m.Plan
		}
		// p7-p8
		var valores []string
		for _, p := range plan {
			valores = append(valores, p)
		}
		miMayoria := mayoria(valores)
		votos := 0
		for _, p := range plan {
			if p == miMayoria {
				votos++
			}
		}
		fila(rImpar, plan, miMayoria, votos, "")

		// p9-p16: ronda del rey
		rey := strings.TrimSpace(reyes[fase-1])
		planRey := ""
		if rey == yo {
			for _, g := range otros() { // p10-p11
				enviar(g, Mensaje{Ronda: rPar, Origen: yo, Plan: planPara(g, miMayoria)})
			}
			plan[yo] = miMayoria // p12
			planRey = miMayoria
		} else {
			for _, m := range recibir(rPar, 1) { // p13
				planRey = m.Plan
			}
			if votos > umbral || planRey == "" { // p14-p15
				plan[yo] = miMayoria
			} else { // p16
				plan[yo] = planRey
			}
		}
		fila(rPar, map[string]string{yo: plan[yo]}, "", 0, planRey)
		if rey == yo && esTraidor() {
			log("  (ronda %d: YO soy el rey TRAIDOR, envío planes distintos: %v)", rPar, mentira)
		} else if rey == yo {
			log("  (ronda %d: YO soy el rey, envío %s)", rPar, miMayoria)
		} else if votos > umbral {
			log("  (ronda %d: rey %s; mi mayoría %s es abrumadora (%d votos): la mantengo)", rPar, rey, miMayoria, votos)
		} else {
			log("  (ronda %d: rey %s; mi mayoría no es abrumadora (%d votos): adopto %s)", rPar, rey, votos, guion(planRey))
		}
	}

	decidir(plan[yo]) // p17
}

// encabezado y fila imprimen la tabla del libro para este general.
func encabezado() {
	fmt.Printf("\n  %-6s", "ronda")
	for _, g := range nombres() {
		fmt.Printf(" %-5s", g)
	}
	fmt.Printf(" | myMajority votesMajority kingPlan\n")
}

func fila(ronda int, plan map[string]string, miMayoria string, votos int, planRey string) {
	fmt.Printf("  %-6d", ronda)
	for _, g := range nombres() {
		p, ok := plan[g]
		if !ok {
			p = " "
		}
		fmt.Printf(" %-5s", p)
	}
	v := ""
	if votos > 0 {
		v = fmt.Sprint(votos)
	}
	fmt.Printf(" | %-10s %-13s %s\n", miMayoria, v, planRey)
}

func esTraidor() bool { return caida >= 0 || len(mentira) > 0 }

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
