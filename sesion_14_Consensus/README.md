# Sesión 14: Algoritmos de consenso en Go

Programación Concurrente y Distribuida (CC57), UPC.

En un sistema distribuido confiable, varios nodos independientes deben ponerse
de acuerdo en **un mismo valor**, aunque algunos fallen. Ese es el **problema
del consenso**. Ben-Ari lo cuenta con la historia de los **generales
bizantinos**: cada general elige atacar (`A`) o retirarse (`R`), algunos son
**traidores**, y los generales **leales** deben llegar a la misma decisión.

Como en la sesión 11, cada nodo es un proceso que solo se comunica **enviando
mensajes JSON por TCP**.

## Objetivos

1. Entender por qué la votación simple de una ronda **falla** con un traidor.
2. Distinguir dos tipos de fallas: **por caída** (el nodo deja de enviar) y
   **bizantinas** (el nodo envía mensajes falsos).
3. Implementar y ejecutar los algoritmos de Ben-Ari, cap. 12: Generales
   Bizantinos, inundación (flooding) y el algoritmo del Rey.
4. Comparar su costo: cuántos generales y cuántos mensajes necesita cada uno.
5. Desplegar una votación de mayoría en malla completa con Docker Compose.

## Material base

- M. Ben-Ari, *Principles of Concurrent and Distributed Programming*, 2.ª ed.,
  cap. 12 "Consensus".
- Guía de laboratorio "Algoritmo de Consenso Distribuido P2P en Topología
  Malla con Golang 1.26 y Docker Compose" (Ms. Carlos Jara García).

Los programas 14.05 a 14.07 parten de la guía, con correcciones. Cada archivo
explica sus cambios en la sección `Notas` de su cabecera.

## Programas

| Listado | Archivo | Tema | Puertos |
|---|---|---|---|
| 14.01 | `consenso_una_ronda.go` ([guía paso a paso](listing14.01/README.md)) | Algoritmo de una ronda (Ben-Ari, Alg. 12.1): falla con un traidor | 9201–9205 |
| 14.02 | `generales_bizantinos.go` ([guía paso a paso](listing14.02/README.md)) | Generales Bizantinos, dos rondas (Alg. 12.2) | 9201–9205 |
| 14.03 | `inundacion.go` ([guía paso a paso](listing14.03/README.md)) | Inundación con fallas por caída, t+1 rondas (Alg. 12.3) | 9201–9205 |
| 14.04 | `rey.go` ([guía paso a paso](listing14.04/README.md)) | Algoritmo del Rey, 5 generales y 1 traidor (Alg. 12.4) | 9201–9205 |
| 14.05 | `consenso_mesh.go` ([guía paso a paso](listing14.05/README.md)) | Nodo de votación por mayoría en malla completa (guía) | 8001–8005 / 8000 |
| 14.06 | `centinela_consenso.go` ([guía paso a paso](listing14.06/README.md)) | Centinela: envía `TRIGGER` a toda la malla | → 8001–8005 |
| 14.07 | `Dockerfile`, `Dockerfile.centinela`, `docker-compose.yml` ([guía paso a paso](listing14.07/README.md)) | 5 nodos en contenedores + centinela con `profiles: [test]` | 8000 |

En 14.01 a 14.04 los generales llevan los nombres del libro y tienen puertos
fijos: Basil 9201, John 9202, Leo 9203, Mike 9204, Zoe 9205. Así cada
ejecución reproduce un escenario del capítulo, con las mismas tablas.

## Cómo ejecutar

Regla general: **un general por terminal**, todos desde la carpeta de la
sesión. Los generales esperan a que todos estén escuchando y luego empiezan
juntos la ronda 1. Se detienen con `Ctrl+C`.

```bash
cd sesion_14_Consensus

# Generales Bizantinos con 3 generales (3 terminales): NO hay consenso
go run listing14.02/generales_bizantinos.go -yo Basil -plan A -invertir Leo
go run listing14.02/generales_bizantinos.go -yo Leo   -plan R
go run listing14.02/generales_bizantinos.go -yo Zoe   -plan A
# Resultado: Leo decide R y Zoe decide A

# Algoritmo del Rey (5 terminales): consenso
go run listing14.04/rey.go -yo Basil -plan A
go run listing14.04/rey.go -yo John  -plan A
go run listing14.04/rey.go -yo Leo   -plan R
go run listing14.04/rey.go -yo Zoe   -plan R
go run listing14.04/rey.go -yo Mike  -mentira Basil=R,John=A,Leo=A,Zoe=R
# Resultado: los 4 leales deciden R

# Malla en Docker (desde listing14.07)
docker compose up --build -d
docker compose --profile test run --rm centinela
docker compose logs | grep RESULTADO
```

**Alternativa en una sola terminal**: compile una vez y lance todos en
segundo plano. La salida de los generales se mezcla, pero cada línea lleva el
nombre del general.

```bash
go build -o /tmp/bg listing14.02/generales_bizantinos.go
/tmp/bg -yo Basil -plan A -invertir Leo & /tmp/bg -yo Leo -plan R & /tmp/bg -yo Zoe -plan A &
sleep 5; kill %1 %2 %3
```

Resultados esperados (comprobados):

| Prueba | Resultado |
|---|---|
| 14.01, Basil cae tras 1 mensaje | Leo decide A, Zoe decide R: sin consenso |
| 14.01, Basil bizantino (`-mentira Leo=A,Zoe=R`) | Leo A, Zoe R: sin consenso |
| 14.02, 3 generales, Basil cae (`-caida 1` o `-caida 3`) | Leo y Zoe deciden A: consenso |
| 14.02, 3 generales, Basil bizantino (`-invertir Leo`) | Leo R, Zoe A: sin consenso |
| 14.02, 4 generales, Zoe bizantina | Basil, John y Leo deciden R: consenso |
| 14.03, Leo y Basil caen, `-t 2` (3 rondas) | John y Zoe deciden R |
| 14.03, lo mismo con `-t 1` (2 rondas) | John R, Zoe A: sin consenso |
| 14.04, rey Mike (traidor) y luego Zoe | los 4 leales deciden R |
| 14.04, rey Zoe y luego Mike (`-reyes Zoe,Mike`) | los 4 leales deciden R |
| 14.05–14.07, votos A, B, A, A, B | los 5 nodos: A = 3, B = 2, ¡GO A! |

## Ideas clave para la clase

1. **Sin fallas, el consenso es trivial**: cada nodo envía su valor a todos y
   todos votan sobre los mismos datos (14.05). El problema aparece cuando un
   nodo falla a la mitad del envío o miente.
2. **Caída vs. bizantina.** Una caída se detecta con un tiempo límite (el
   mensaje no llegó). Una falla bizantina no se detecta: el traidor envía
   mensajes con formato correcto pero contenido falso, y distinto a cada uno.
3. **3t+1.** Con fallas bizantinas se necesitan al menos 3t+1 generales para t
   traidores. Con 3 generales y 1 traidor es **imposible** (Ben-Ari, teorema
   12.3), y 14.02 lo muestra.
4. **El costo crece rápido.** Generales Bizantinos: 36 mensajes con t = 1,
   392 con t = 2, 5408 con t = 4. El Rey usa un general más por traidor
   (4t+1) pero muchos menos mensajes: 48, 240, 1440.
5. **Las caídas son más fáciles.** La inundación tolera t caídas con t+1
   rondas y cualquier número de generales. Con una ronda menos, falla (14.03).
6. **Rondas síncronas.** Los algoritmos del libro suponen rondas: todos envían
   y luego todos reciben. En TCP eso se programa con un buzón por ronda y un
   tiempo límite (`-espera`).

## Ejercicios

**Del libro (Ben-Ari, cap. 12):**

1. Ejercicio 1: en 14.02 con 4 generales, haga que Zoe envíe A a Basil en vez
   de R. Compruebe que la decisión pasa a ser A.
2. Ejercicio 2: ¿qué pasa si un general se envía su plan a sí mismo?
3. Ejercicio 5: generalice 14.02 a t traidores (t+1 rondas) y pruébelo con 7
   generales y 2 traidores.
4. Ejercicio 6: derive cuántos mensajes envía la inundación con n nodos y
   cuéntelos en 14.03.
5. Ejercicio 8: encuentre otro escenario en el que 14.03 falla con t rondas.
6. Ejercicio 9: escenario del Rey con tres leales que eligen A y uno R.
7. Ejercicio 10: en 14.04, use dos traidores con 5 generales
   (`-generales` con 5 y `-t 1`, pero dos nodos con `-mentira`). Muestre que
   no hay consenso.

**De la guía de laboratorio (14.05–14.07):**

8. Votación aleatoria multirronda: `-vote=0` en los 5 servicios; lance el
   centinela 3 veces. Los 5 nodos siempre coinciden.
9. Quórum con timeout: apague `nodo5` y modifique `consensus()` para decidir
   con `(N/2)+1` votos después de un tiempo límite.
10. Nodo traidor: haga que un nodo envíe A a `nodo1` y `nodo2` y B a `nodo3` y
    `nodo4`. ¿Siguen de acuerdo? Agregue la segunda ronda de 14.02.

**Sobre el código de esta sesión:**

11. Cuente los mensajes que envía cada general en 14.02 y 14.04 y compárelos
    con las fórmulas de la tabla del libro.
12. Levante los generales de 14.04 con Docker Compose, como en 14.07. Pista:
    `-generales Basil=basil:9201,John=john:9202,...`.
