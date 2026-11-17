# Listado 11.12: Ricart-Agrawala con permisos

`ricart_agrawala.go` implementa el **algoritmo de Ricart-Agrawala** para la
exclusión mutua distribuida (Ben-Ari, Algoritmo 10.2) sobre TCP con mensajes
JSON. Un nodo entra a su **sección crítica** solo cuando **todos** los demás
nodos le dieron permiso.

El recurso compartido del listado 11.11 observa las entradas y salidas y
avisa si dos nodos están dentro a la vez.

**Regla: primero el recurso, luego los nodos.** Aquí no hay centinela: cada
nodo empieza a trabajar solo, 2 segundos después de arrancar.

## Topología: todos con todos

```
   nodo 1 (9101) ◄──────► nodo 2 (9102)
          ▲                    ▲
          └──► nodo 3 (9103) ◄─┘
                    │
       ENTRA/SALE   ▼   (los tres nodos)
             recurso (9100)
```

Cada nodo conoce la dirección de todos los demás (flag `-nodos`). Envía y
recibe mensajes de cada uno.

## El algoritmo

Cada nodo saca un **ticket** (`myNum = highestNum + 1`), como en el algoritmo
de la panadería, y pide permiso a todos:

| Momento | Qué hace el nodo | Ben-Ari |
|---|---|---|
| quiere entrar | `requestCS = true`, saca ticket, envía `REQUEST` a los N−1 demás | p2–p5 |
| espera | recibe N−1 `REPLY` | p6 |
| sección crítica | usa el recurso 500 ms | p7 |
| sale | `requestCS = false`, envía `REPLY` a los que difirió | p8–p11 |

Cuando llega un `REQUEST` (proceso *Receive*, p12–p16):

```go
responder := !requestCS || menor(m.Num, m.Origen, myNum, miID)
```

- Si el nodo **no quiere** entrar, o el ticket del otro es **menor**:
  responde `REPLY` ya.
- Si no: **difiere** la respuesta hasta salir de su sección crítica.
- **Empate** de tickets: gana el id menor (función `menor`, el `<<` del libro).

Costo: **2(N−1) mensajes** por cada entrada (N−1 `REQUEST` + N−1 `REPLY`).
Con 3 nodos, 4 mensajes.

Dentro de cada nodo, el bucle principal y los *Receive* son goroutines
concurrentes, así que el estado local (`myNum`, `requestCS`, `diferidos`) se
protege con un `sync.Mutex`. Ben-Ari supone que eso es atómico; en Go hay que
programarlo.

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-id` | `1` | identificador de este nodo |
| `-nodos` | `1=localhost:9101,2=localhost:9102,3=localhost:9103` | todos los nodos, `id=host:puerto` separados por coma |
| `-rondas` | `3` | cuántas veces entra a la sección crítica |
| `-recurso` | `localhost:9100` | recurso compartido (vacío = no usar) |
| `-sinra` | `false` | entra **sin pedir permiso** (para ver el fallo) |

## Paso 0: abrir 4 terminales

En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

## Paso 1: levantar el recurso (terminal 1)

```bash
go run listing11.11/recurso_compartido.go
```

```
[RECURSO] Escuchando en localhost:9100
```

## Pasos 2 a 4: levantar los tres nodos (terminales 2 a 4)

Lance los tres en menos de 2 segundos, o al menos sin demora. Si un nodo aún no
arrancó, los demás reintentan el envío durante unos 10 segundos.

```bash
go run listing11.12/ricart_agrawala.go -id 1   # terminal 2
go run listing11.12/ricart_agrawala.go -id 2   # terminal 3
go run listing11.12/ricart_agrawala.go -id 3   # terminal 4
```

```
18:33:36.827 [NODO 1] escuchando en localhost:9101; otros nodos: [2 3]
```

## Resultado esperado

Terminal 2 (nodo 1), fragmento:

```
[NODO 1] quiere entrar (ronda 1) con ticket 1
[NODO 1]   permiso recibido de 2
[NODO 1]   permiso recibido de 3
[NODO 1] >>> ENTRA a la sección crítica (ronda 1)
[NODO 1] <<< SALE de la sección crítica (ronda 1)
...
[NODO 1] quiere entrar (ronda 2) con ticket 4
[NODO 1]   permiso recibido de 2
[NODO 1]   permiso recibido de 3
[NODO 1] >>> ENTRA a la sección crítica (ronda 2)
[NODO 1]   REQUEST de 3 (ticket 5): lo difiere (mi ticket 4)
[NODO 1]   REQUEST de 2 (ticket 6): lo difiere (mi ticket 4)
[NODO 1] <<< SALE de la sección crítica (ronda 2)
[NODO 1]   responde al diferido 3
[NODO 1]   responde al diferido 2
...
[NODO 1] terminó sus 3 rondas; sigue respondiendo (Ctrl+C para salir)
```

Fíjese en las líneas **`lo difiere`**: el nodo 1 está dentro con el ticket 4,
así que los tickets 5 y 6 esperan. Al salir, responde a los diferidos.

Terminal 1 (recurso): 9 pares ENTRA/SALE (3 nodos × 3 rondas), siempre con
`dentro = 1` y **ninguna** `VIOLACIÓN`.

```
18:33:39.114 ENTRA nodo 1 (dentro = 1)
18:33:39.615 SALE  nodo 1
18:33:39.735 ENTRA nodo 2 (dentro = 1)
18:33:40.236 SALE  nodo 2
...
```

Los tiempos de la sección no crítica son aleatorios, así que el orden de
entrada cambia en cada ejecución. La ausencia de violaciones no cambia.

## Paso 5: detener todo

Los nodos siguen vivos al terminar sus rondas: un nodo inactivo debe seguir
respondiendo a los demás, o los bloquearía para siempre (Ben-Ari, *quiescent
nodes*). Presione `Ctrl+C` en las terminales 2 a 4 y luego en la 1.

## Paso 6: repetir SIN el algoritmo

Reinicie el recurso y lance los tres nodos con `-sinra`:

```bash
go run listing11.11/recurso_compartido.go              # terminal 1
go run listing11.12/ricart_agrawala.go -id 1 -sinra    # terminal 2
go run listing11.12/ricart_agrawala.go -id 2 -sinra    # terminal 3
go run listing11.12/ricart_agrawala.go -id 3 -sinra    # terminal 4
```

Ahora cada nodo entra sin pedir permiso. Terminal 1:

```
18:33:50.746 ENTRA nodo 1 (dentro = 1)
18:33:51.090 ENTRA nodo 2 (dentro = 2)
18:33:51.090 *** VIOLACIÓN DE EXCLUSIÓN MUTUA #1 ***
...
18:33:51.978 ENTRA nodo 1 (dentro = 3)
18:33:51.978 *** VIOLACIÓN DE EXCLUSIÓN MUTUA #4 ***
```

En nuestra prueba hubo 6 violaciones; el número cambia en cada ejecución.

## Errores comunes

- **`El id 4 no está en -nodos`**: el `-id` no aparece en la lista de nodos.
- **`ERROR: no se pudo enviar REQUEST a localhost:9103`**: un nodo no arrancó
  en unos 10 segundos. Los demás se quedan esperando su `REPLY` para siempre.
- **Un nodo se queda en `quiere entrar` y no avanza**: otro nodo se detuvo con
  `Ctrl+C` antes de responderle. Reinicie todos.
- **`recurso no disponible`**: el recurso no estaba levantado. El algoritmo
  funciona igual, pero nadie verifica la exclusión mutua.
- **`Error al escuchar: ... address already in use`**: quedó un nodo activo de
  una ejecución anterior (por ejemplo, de 11.13, que usa los mismos puertos).

## Para pensar y ejercicios

1. **Contar mensajes.** Cuente los mensajes que envía cada nodo y compare con
   2(N−1) por entrada.
2. **Mismo id.** ¿Qué pasa si dos nodos usan el mismo `-id`? (Ben-Ari,
   ejercicio 3.)
3. **Más nodos.** Agregue un cuarto nodo con
   `-nodos 1=localhost:9101,2=localhost:9102,3=localhost:9103,4=localhost:9104`
   en **los cuatro** nodos.
4. **Comparar con el token.** Ejecute 11.13 con los mismos flags. ¿Qué
   algoritmo envía menos mensajes?
5. **En Docker.** Levante estos nodos con Docker Compose, como en 11.15.
