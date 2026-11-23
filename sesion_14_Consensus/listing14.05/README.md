# Listado 14.05: Votación por mayoría en malla completa

`consenso_mesh.go` es un nodo de **votación por mayoría** en una **malla
completa** (full mesh): cada nodo conoce a todos los demás. Cuando el
centinela (listado 14.06) envía la orden `TRIGGER`, cada nodo registra su
voto, lo difunde a sus N−1 pares y, con los N votos, decide **GO A** o **GO B**.

Es el algoritmo de una ronda de Ben-Ari (14.01) **sin fallas**: como todos
reciben los mismos N votos y aplican la misma regla, todos deciden lo mismo.

Parte del código de la guía de laboratorio (Ms. Carlos Jara García), con las
correcciones que se listan en la cabecera del archivo.

**Regla: primero los nodos, al final el centinela.**

## Topología: malla completa de 5 nodos

```
            nodo1 (8001)
           ╱  │   │  ╲
     nodo5 ───┼───┼─── nodo2
     (8005) ╲ │   │ ╱  (8002)
          nodo4 ─── nodo3
          (8004)    (8003)

   centinela ──TRIGGER──► los 5 nodos
```

Cada nodo envía su voto a los otros 4: 5 × 4 = **20 mensajes** por ronda, más
5 `TRIGGER`.

## El mensaje

Las mismas constantes del código original (`143consenso.go`):

```go
const (
	cnum     = iota // 0: voto entre pares
	opa             // 1: Opción A
	opb             // 2: Opción B
	ctrigger        // 3: orden del centinela
)
```

```json
{"code":0,"addr":"localhost:8003","op":1}    // nodo3 vota A
{"code":3,"addr":"centinela","op":0}         // orden de votar
```

## El monitor de votos

Varias goroutines (una por conexión) quieren tocar el mapa de votos a la vez.
El mapa vive **dentro de un canal** de capacidad 1: quien lo saca es el único
que lo usa, y al terminar lo devuelve.

```go
info := <-chInfo        // entra al monitor
info[msg.Addr] = msg.Op
...
chInfo <- info          // sale del monitor
```

Es otra forma de exclusión mutua, sin `sync.Mutex`.

## Parámetros

| Flag | Variable | Por defecto | Significado |
|---|---|---|---|
| `-name` | `NODE_NAME` | `nodo1` | nombre del nodo (en Docker, también su hostname) |
| `-host` | `NODE_HOST` | el nombre | host de la dirección propia (`localhost` fuera de Docker) |
| `-port` | `NODE_PORT` | `8000` | puerto de escucha |
| `-peers` | `PEERS` | (vacío) | pares `host:puerto` separados por coma; puede incluir al propio nodo |
| `-vote` | `NODE_VOTE` | `0` | `1` = A, `2` = B, `0` = aleatorio |
| `-delay_ms` | `DELAY_MS` | `300` | retardo antes de difundir el voto |

## Paso 0: abrir 6 terminales

En **cada** terminal:

```bash
cd sesion_14_Consensus
P=localhost:8001,localhost:8002,localhost:8003,localhost:8004,localhost:8005
```

## Pasos 1 a 5: levantar los cinco nodos (terminales 1 a 5)

Votos de la guía: A, B, A, A, B.

```bash
go run listing14.05/consenso_mesh.go -name nodo1 -host localhost -port 8001 -vote 1 -peers $P
go run listing14.05/consenso_mesh.go -name nodo2 -host localhost -port 8002 -vote 2 -peers $P
go run listing14.05/consenso_mesh.go -name nodo3 -host localhost -port 8003 -vote 1 -peers $P
go run listing14.05/consenso_mesh.go -name nodo4 -host localhost -port 8004 -vote 1 -peers $P
go run listing14.05/consenso_mesh.go -name nodo5 -host localhost -port 8005 -vote 2 -peers $P
```

```
===============================================================
[nodo1] NODO LISTO | Addr: localhost:8001 | Voto: OPCIÓN A (1)
[nodo1] Malla completa (4 pares): [localhost:8002 localhost:8003 localhost:8004 localhost:8005]
===============================================================
```

## Paso 6: lanzar el centinela (terminal 6)

```bash
go run listing14.06/centinela_consenso.go -nodes $P
```

## Resultado esperado

En **cada** nodo (terminal 1, por ejemplo):

```
[nodo1] TRIGGER del centinela -> voto OPCIÓN A
[nodo1] Voto de [localhost:8002] -> OPCIÓN B
[nodo1] Voto de [localhost:8003] -> OPCIÓN A
[nodo1] Voto de [localhost:8004] -> OPCIÓN A
[nodo1] Voto de [localhost:8005] -> OPCIÓN B
---------------------------------------------------------------
[nodo1] RONDA 1, CONTEO FINAL (5 nodos): Opción A = 3 | Opción B = 2
>>> [nodo1] RESULTADO DEL CONSENSO: ¡GO A! <<<
---------------------------------------------------------------
```

El orden de los votos cambia en cada ejecución; el conteo y el resultado, no.
Lance el centinela otra vez: empieza la **ronda 2**, con el mismo resultado.

## Paso 7: detener todo

`Ctrl+C` en las terminales 1 a 5.

## Errores comunes

- **`No se pudo escuchar en 0.0.0.0:8001: ... address already in use`**: quedó
  un nodo de otra ejecución.
- **El nodo nunca muestra `CONTEO FINAL`**: falta un nodo o el centinela no lo
  contactó. Los demás esperan su voto **para siempre** (no hay timeout).
- **Malla con 5 pares en vez de 4**: olvidó `-host localhost`; la dirección
  propia quedó como `nodo1:8001` y no se descartó de la lista.
- **`[CENTINELA ERROR] No se pudo contactar localhost:8003`**: ese nodo no
  está levantado.

## Para pensar y ejercicios

1. **Empate.** Use 4 nodos con votos A, A, B, B. ¿Qué deciden? ¿Por qué la
   guía usa un número impar?
2. **Votos aleatorios.** Use `-vote 0` en los cinco y lance el centinela
   varias veces. ¿Alguna vez dos nodos deciden distinto?
3. **Quórum.** Detenga un nodo y modifique `consensus()` para decidir con
   `(N/2)+1` votos después de un tiempo límite.
4. **Traidor.** Compare con 14.01: ¿qué pasa si un nodo envía votos distintos
   a cada par?
