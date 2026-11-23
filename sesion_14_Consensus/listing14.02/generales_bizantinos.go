/*
 * Archivo:     generales_bizantinos.go
 *
 * Propósito:   Algoritmo de los Generales Bizantinos de dos rondas (Lamport,
 *              Shostak y Pease; Ben-Ari, Algoritmo 12.2) sobre TCP con
 *              mensajes JSON. En la ronda 1 cada general envía su plan; en la
 *              ronda 2 informa a los demás lo que recibió de cada uno. Luego
 *              vota dos veces: primero el plan de cada general, después el
 *              plan final.
 *
 * Compilación: go build generales_bizantinos.go
 * Ejecución:   Tres generales, traidor bizantino (Ben-Ari, sección 12.7):
 *                go run generales_bizantinos.go -yo Basil -plan A -invertir Leo
 *                go run generales_bizantinos.go -yo Leo   -plan R
 *                go run generales_bizantinos.go -yo Zoe   -plan A
 *              Cuatro generales, Zoe traidora (Ben-Ari, sección 12.8):
 *                G=Basil,John,Leo,Zoe
 *                go run generales_bizantinos.go -generales $G -yo Basil -plan A
 *                go run generales_bizantinos.go -generales $G -yo John  -plan A
 *                go run generales_bizantinos.go -generales $G -yo Leo   -plan R
 *                go run generales_bizantinos.go -generales $G -yo Zoe \
 *                    -mentira Basil=R,Leo=R,John=A -invertir Basil,John,Leo
 *
 * Entrada:     Ninguna (flags: -yo, -plan, -generales, -host, -espera,
 *              -caida, -mentira, -invertir).
 * Salida:      Cada general muestra la tabla del libro (plan, "reported by",
 *              mayoría) y su decisión. Con 3 generales y un traidor
 *              bizantino NO hay consenso; con 4 generales y un traidor, sí.
 *
 * Notas:
 *    1. Correspondencia con Ben-Ari (Alg. 12.2):
 *         p1     plan[myID] <- chooseAttackOrRetreat   -> flag -plan
 *         p2-p5  primera ronda                         -> ronda 1
 *         p6-p11 segunda ronda: send(G', myID, G, ...) -> ronda 2, campo Sobre
 *         p12-p13 primera votación (por general)       -> mayoriaPlan[G]
 *         p14-p15 segunda votación (plan final)        -> mayoria(mayoriaPlan)
 *    2. En la ronda 2 el libro envía un mensaje por cada par (G', G). Aquí
 *       también, para poder contar los mensajes: (n-1) + (n-1)(n-2) por
 *       general. Con 4 generales, 3 + 6 = 9 por general y 36 en total.
 *    3. Si un general no recibió el plan de G (porque G cayó), igual envía
 *       el informe con plan vacío ("-"). Así nadie espera un mensaje que no
 *       llegará.
 *    4. Traidores: -caida k (cae tras k mensajes), -mentira G=P (en la ronda
 *       1 envía el plan P a G), -invertir G1,G2 (en la ronda 2 informa a G1 y
 *       G2 lo contrario de lo que recibió).
 *    5. Se necesitan al menos 3t+1 generales para t traidores bizantinos.
 *       Con 3 generales (t = 1) no alcanza: ese es el primer ejemplo.
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
	Sobre  string `json:"sobre,omitempty"` // ronda 2: general del que se informa
	Plan   string `json:"plan"`            // "A", "R" o "" (desconocido)
}

var (
	yo        string
	generales map[string]string // nombre -> dirección
	espera    time.Duration
	caida     int               // -1 = leal; k = cae tras enviar k mensajes
	mentira   map[string]string // destino -> plan falso (traidor bizantino)
	invertir  map[string]bool   // destinos a los que el traidor miente en la ronda 2

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
	textoMentira := flag.String("mentira", "", "TRAIDOR: plan falso en la ronda 1, ej. Leo=A,Zoe=R")
	textoInvertir := flag.String("invertir", "", "TRAIDOR: en la ronda 2 miente a estos generales, ej. Leo")
	flag.Parse()

	generales = leerGenerales(*lista, *host)
	mentira = leerMapa(*textoMentira)
	invertir = map[string]bool{}
	for _, g := range strings.Split(*textoInvertir, ",") {
		if g = strings.TrimSpace(g); g != "" {
			invertir[g] = true
		}
	}
	escuchar()
	esperarATodos()

	// p1
	plan := map[string]string{yo: strings.ToUpper(*miPlan)}
	log("mi plan es %s", plan[yo])

	// p2-p5: primera ronda
	for _, g := range otros() {
		enviar(g, Mensaje{Ronda: 1, Origen: yo, Plan: planPara(g, plan[yo])})
	}
	for _, m := range recibir(1, len(otros())) {
		plan[m.Origen] = m.Plan
		log("  ronda 1: %s dice que su plan es %s", m.Origen, m.Plan)
	}

	// p6-p8: segunda ronda. A cada general G' le informo el plan de cada
	// otro general G (sin incluirme a mí ni al propio G').
	for _, destino := range otros() {
		for _, g := range otros() {
			if g == destino {
				continue
			}
			informe := plan[g]
			if invertir[destino] {
				informe = contrario(informe)
			}
			enviar(destino, Mensaje{Ronda: 2, Origen: yo, Sobre: g, Plan: informe})
		}
	}
	// p9-p11: reportedPlan[G, G'] = lo que G dice que recibió de G'
	n := len(generales)
	reportado := map[string]map[string]string{}
	for _, m := range recibir(2, (n-1)*(n-2)) {
		if reportado[m.Sobre] == nil {
			reportado[m.Sobre] = map[string]string{}
		}
		reportado[m.Sobre][m.Origen] = m.Plan
		log("  ronda 2: %s informa que el plan de %s es %s", m.Origen, m.Sobre, guion(m.Plan))
	}

	// p12-p15: dos votaciones
	decidir(mostrarTabla(plan, reportado))
}

// mostrarTabla imprime la tabla del libro y calcula las dos votaciones:
// mayoriaPlan[G] = mayoría(plan[G] ∪ reportedPlan[*, G]) y la final.
func mostrarTabla(plan map[string]string, reportado map[string]map[string]string) string {
	fmt.Printf("\n  %-8s %-5s", "general", "plan")
	for _, g := range otros() {
		fmt.Printf(" %-6s", "<"+g)
	}
	fmt.Printf(" mayoría\n")

	var mayoriaPlan []string
	for _, g := range nombres() {
		fila := []string{plan[g]}
		fmt.Printf("  %-8s %-5s", g, guion(plan[g]))
		for _, informante := range otros() {
			v := ""
			if g != yo && informante != g {
				v = reportado[g][informante]
				fila = append(fila, v)
			}
			if g == yo || informante == g {
				fmt.Printf(" %-6s", "")
			} else {
				fmt.Printf(" %-6s", guion(v))
			}
		}
		m := mayoria(fila) // p13 (para mí mismo, p14: mi propio plan)
		if g == yo {
			m = plan[yo]
		}
		mayoriaPlan = append(mayoriaPlan, m)
		fmt.Printf(" %s\n", m)
	}
	final := mayoria(mayoriaPlan) // p15
	fmt.Printf("  %-8s %-5s%s %s\n", "mayoría", "", strings.Repeat(" ", 7*len(otros())), final)
	fmt.Println("  (<X = informado por X en la ronda 2)")
	fmt.Println()
	return final
}

func contrario(p string) string {
	switch p {
	case "A":
		return "R"
	case "R":
		return "A"
	}
	return p
}

func esTraidor() bool { return caida >= 0 || len(mentira) > 0 || len(invertir) > 0 }

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
