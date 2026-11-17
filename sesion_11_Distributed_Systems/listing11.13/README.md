# Listado 11.13: Ricart-Agrawala con token

`ra_token.go` implementa el **algoritmo de Ricart-Agrawala con paso de token**
(Ben-Ari, Algoritmo 10.3; parecido a Suzuki-Kasami). Hay un **único token**
en el sistema y solo el nodo que lo tiene puede entrar a su sección crítica.

Se ejecuta igual que el listado 11.12 (mismos puertos, mismo recurso del
listado 11.11), así que se pueden comparar directamente.

**Regla: primero el recurso, luego los nodos.** El nodo 1 empieza con el
token.

## Topología: todos con todos, más un token

```
   nodo 1 (9101) ◄──────► nodo 2 (9102)      REQUEST: a todos
          ▲                    ▲             TOKEN:   a uno solo
          └──► nodo 3 (9103) ◄─┘
                    │
       ENTRA/SALE   ▼   (los tres nodos)
             recurso (9100)
```

## Permisos (11.12) vs. token (11.13)

| | 11.12 permisos | 11.13 token |
|---|---|---|
| Para entrar hace falta | `REPLY` de los N−1 demás | tener el token |
| Mensajes por entrada | siempre 2(N−1) | N (N−1 `REQUEST` + 1 `TOKEN`), o **0** si ya tiene el token |
| Estado que viaja | nada | el arreglo `granted` dentro del token |

## El algoritmo

Estado de cada nodo:

- `requested[j]`: el ticket más alto que pidió el nodo *j* (lo que este nodo
  sabe).
- `granted[j]`: el ticket con el que *j* entró por última vez. **Viaja dentro
  del token**, así que el dueño del token siempre tiene la versión correcta.

| Momento | Qué hace el nodo | Ben-Ari |
|---|---|---|
| quiere entrar y **tiene** el token | entra sin enviar mensajes | p2 |
| quiere entrar y **no** lo tiene | `myNum++`, envía `REQUEST` a todos, espera el `TOKEN` | p3–p7 |
| sale | `granted[miID] = myNum`; si alguien tiene `requested[j] > granted[j]`, le envía el token | p10–p12 |
| llega un `REQUEST` | actualiza `requested`; si tiene el token y no está dentro, lo envía | p13–p16 |

**Cambio respecto al libro:** el libro envía el token a "algún" nodo que lo
pidió, y eso permite inanición. Aquí `elegirDestino` busca en **orden
circular** desde `miID+1`, así todos los que piden reciben el token.

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-id` | `1` | identificador; el nodo 1 empieza con el token |
| `-nodos` | `1=localhost:9101,2=localhost:9102,3=localhost:9103` | todos los nodos; los ids deben ser 1..N sin huecos |
| `-rondas` | `3` | cuántas veces entra a la sección crítica |
| `-recurso` | `localhost:9100` | recurso compartido (vacío = no usar) |

## Paso 0: abrir 4 terminales

En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

## Paso 1: levantar el recurso (terminal 1)

```bash
go run listing11.11/recurso_compartido.go
```

## Pasos 2 a 4: levantar los tres nodos (terminales 2 a 4)

```bash
go run listing11.13/ra_token.go -id 1   # terminal 2 (empieza con el token)
go run listing11.13/ra_token.go -id 2   # terminal 3
go run listing11.13/ra_token.go -id 3   # terminal 4
```

```
18:33:56.366 [NODO 1] escuchando en localhost:9101; tengo el token: true
18:33:56.670 [NODO 2] escuchando en localhost:9102; tengo el token: false
```

## Resultado esperado

Terminal 2 (nodo 1), fragmento:

```
[NODO 1] ya tengo el token: entro sin enviar mensajes
[NODO 1] >>> ENTRA a la sección crítica (ronda 1)
[NODO 1]   REQUEST de 2 (ticket 1)
[NODO 1] <<< SALE de la sección crítica (ronda 1)
[NODO 1]   envía el TOKEN al nodo 2
[NODO 1] pide el token con ticket 1 (ronda 2)
[NODO 1]   recibió el token; granted = [0 1 0]
[NODO 1] >>> ENTRA a la sección crítica (ronda 2)
```

Cómo leerlo:

1. El nodo 1 tiene el token y entra **sin enviar mensajes**.
2. Mientras está dentro, llega el `REQUEST` del nodo 2. No le envía el token
   todavía: está dentro.
3. Al salir, ve que `requested[2] = 1 > granted[2] = 0` y le envía el token.
4. Cuando el token vuelve al nodo 1, `granted = [0 1 0]` muestra que el nodo 2
   ya entró una vez con el ticket 1 (`granted[1..3]`).

Busque también en las terminales 3 y 4 la línea `ya tengo el token: entro sin
enviar mensajes`: un nodo que conserva el token entra gratis.

Terminal 1 (recurso): 9 pares ENTRA/SALE, siempre `dentro = 1`, **ninguna**
`VIOLACIÓN`.

## Paso 5: detener todo

Los nodos siguen atendiendo al terminar (puede que tengan el token y otro lo
necesite). Presione `Ctrl+C` en las terminales 2 a 4 y luego en la 1.

## Errores comunes

- **`Los ids de -nodos deben ser 1..N sin huecos`**: por ejemplo,
  `1=...,2=...,4=...`. Los arreglos usan el id como índice.
- **Nadie entra nunca**: no se levantó el nodo 1, que es el único con el
  token al principio. Los demás esperan un token que no existe.
- **Un nodo se queda en `pide el token`**: el nodo que tenía el token se
  detuvo con `Ctrl+C` y el token se perdió con él. Reinicie todos.
- **`Error al escuchar: ... address already in use`**: quedan nodos de 11.12
  activos (mismos puertos).

## Para pensar y ejercicios

1. **Contar mensajes.** Cuente los mensajes de cada nodo y compare con 11.12.
   ¿En qué caso el token es más barato y en cuál no?
2. **Inanición.** Cambie `elegirDestino` para que siempre busque desde el id 1
   (como "algún N" del libro). Con muchas rondas, ¿algún nodo espera mucho
   más? (Ben-Ari, ejercicio 13.)
3. **Token perdido.** ¿Cómo podría el sistema detectar que el token se perdió
   y generar uno nuevo?
4. **Neilsen-Mizuno.** Implemente el Alg. 10.4 de Ben-Ari con la misma
   estructura de este archivo.
