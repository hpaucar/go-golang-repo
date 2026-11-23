/*
 * Archivo:     inundacion.go
 *
 * Propósito:   Algoritmo de inundación (flooding) para consenso con fallas por
 *              CAÍDA (Ben-Ari, Algoritmo 12.3) sobre TCP con mensajes JSON.
 *              Durante t+1 rondas, cada general envía a todos el CONJUNTO de
 *              planes que conoce y lo une con los conjuntos que recibe. Al
 *              final decide por mayoría.
 *
 * Compilación: go build inundacion.go
 * Ejecución:   Cuatro terminales, dos traidores que caen (Ben-Ari, sec. 12.9):
 *                go run inundacion.go -yo Basil -plan A -caida 4
 *                go run inundacion.go -yo John  -plan A
 *                go run inundacion.go -yo Leo   -plan R -caida 1
 *                go run inundacion.go -yo Zoe   -plan R
 *              Con -t 2 (3 rondas, por defecto) John y Zoe deciden lo mismo.
 *              Repita con -t 1 en los cuatro (2 rondas): ya no coinciden.
 *
 * Entrada:     Ninguna (flags: -yo, -plan, -t, -generales, -host, -espera,
 *              -caida).
 * Salida:      En cada ronda, los planes que conoce el general; al final, la
 *              tabla y la decisión por mayoría.
 *
 * Notas:
 *    1. Correspondencia con Ben-Ari (Alg. 12.3):
 *         plan <- {chooseAttackOrRetreat}   -> planes[yo] = -plan
 *         p1  do t+1 times                  -> for ronda := 1..t+1
 *         p2-p3 send(G, plan)               -> enviar el mapa completo
 *         p4-p6 receive y plan <- plan ∪ r  -> unir(planes, m.Planes)
 *         p7  finalPlan <- majority(plan)   -> mayoria(...)
 *    2. El libro habla de un "conjunto de planes"; para votar hay que saber
 *       de quién es cada plan, así que aquí el conjunto es un mapa
 *       general -> plan.
 *    3. Escenario del libro: Leo cae tras enviar su plan solo a Basil
 *       (-caida 1); Basil envía la ronda 1 (3 mensajes) y en la ronda 2 solo
 *       alcanza a enviar a John (-caida 4). John reenvía el plan de Leo a Zoe
 *       en la ronda 3. Con solo t rondas, ese reenvío no ocurre.
 *    4. Solo tolera fallas por caída. Un traidor bizantino podría enviar
 *       planes falsos de otros generales.
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
	Ronda  int               `json:"ronda"`
	Origen string            `json:"origen"`
	Planes map[string]string `json:"planes"` // general -> plan que conozco
}

var (
	yo        string
	generales map[string]string // nombre -> dirección
	espera    time.Duration
	caida     int               // -1 = leal; k = cae tras enviar k mensajes
	mentira   map[string]string // no se usa aquí (solo fallas por caída)

	muEnvio  sync.Mutex
	enviados int
	caido    bool
)

func main() {
	flag.StringVar(&yo, "yo", "Leo", "nombre de este general")
	miPlan := flag.String("plan", "A", "plan inicial: A (atacar) o R (retirarse)")
	t := flag.Int("t", 2, "número máximo de traidores; se hacen t+1 rondas")
	lista := flag.String("generales", "Basil,John,Leo,Zoe", "todos los generales, separados por coma")
	host := flag.String("host", "localhost", "host de los generales (si no se da host:puerto)")
	flag.DurationVar(&espera, "espera", 2*time.Second, "tiempo máximo de espera por ronda")
	flag.IntVar(&caida, "caida", -1, "TRAIDOR: deja de enviar después de k mensajes")
	flag.Parse()

	generales = leerGenerales(*lista, *host)
	escuchar()
	esperarATodos()

	planes := map[string]string{yo: strings.ToUpper(*miPlan)}
	log("mi plan es %s; se harán %d rondas (t = %d)", planes[yo], *t+1, *t)

	for ronda := 1; ronda <= *t+1; ronda++ { // p1
		copia := map[string]string{}
		for g, p := range planes {
			copia[g] = p
		}
		for _, g := range otros() { // p2-p3
			enviar(g, Mensaje{Ronda: ronda, Origen: yo, Planes: copia})
		}
		for _, m := range recibir(ronda, len(otros())) { // p4-p6
			for g, p := range m.Planes {
				if _, ya := planes[g]; !ya {
					log("  ronda %d: por %s se entera de que %s eligió %s", ronda, m.Origen, g, p)
				}
				planes[g] = p
			}
		}
		log("fin de la ronda %d: conoce %s", ronda, resumen(planes))
	}

	decidir(mostrarTabla(planes)) // p7
}

// mostrarTabla imprime los planes conocidos y devuelve la mayoría.
func mostrarTabla(plan map[string]string) string {
	fmt.Printf("\n  %-8s %s\n", "general", "plan")
	var valores []string
	for _, g := range nombres() {
		valores = append(valores, plan[g])
		fmt.Printf("  %-8s %s\n", g, guion(plan[g]))
	}
	final := mayoria(valores)
	fmt.Printf("  %-8s %s\n\n", "mayoría", final)
	return final
}

func resumen(plan map[string]string) string {
	var partes []string
	for _, g := range nombres() {
		if p, ok := plan[g]; ok {
			partes = append(partes, g+"="+p)
		}
	}
	return "{" + strings.Join(partes, ", ") + "}"
}

func esTraidor() bool { return caida >= 0 }

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
