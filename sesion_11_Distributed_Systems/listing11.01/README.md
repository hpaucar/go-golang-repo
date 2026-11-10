# Listado 11.01: Centinela + ordenamiento distribuido en cadena

El **centinela** (`centinela.go`) es el programa que inicia el trabajo: envía
los números 9, 12, 10 y 2 al puerto 9080, cada uno en una conexión TCP
distinta. Por sí solo no hace nada: necesita que los cuatro nodos de
ordenamiento del listado 11.02 (`../listing11.02/ordenamiento.go`) ya estén
escuchando.

**Regla: primero los nodos, al final el centinela.**

## Paso 0: abrir 5 terminales

Cada programa es un **proceso** distinto, con su propia memoria. Esa es la idea
central de la sesión: los nodos no comparten variables y solo se comunican con
mensajes TCP. En cada terminal:

```bash
cd sesion_11_Distributed_Systems/listing11.02
```

`go run` compila el archivo en una carpeta temporal y lo ejecuta. La primera
vez tarda un poco más porque Go compila.

## Pasos 1 a 4: levantar los nodos, del último al primero

Pueden levantarse en cualquier orden, siempre que todos estén activos antes del
centinela. Empezar por el último es un buen hábito: cuando un nodo reenvía un
número, el siguiente ya está escuchando.

```bash
go run ordenamiento.go -local 9083 -n 1                 # terminal 4 (último nodo)
go run ordenamiento.go -local 9082 -remoto 9083 -n 2    # terminal 3
go run ordenamiento.go -local 9081 -remoto 9082 -n 3    # terminal 2
go run ordenamiento.go -local 9080 -remoto 9081 -n 4    # terminal 1
```

| Flag | Significado |
|---|---|
| `-local` | puerto donde escucha este nodo (`net.Listen`) |
| `-remoto` | puerto del siguiente nodo (`net.Dial`). Si se omite, es el **último nodo** |
| `-n` | cuántos números recibirá este nodo antes de imprimir su resultado |

Cada nodo muestra una línea como:

```
[NODO 9080] escuchando; siguiente = "localhost:9081"; espera 4 valores
```

## Paso 5: ejecutar el centinela

```bash
go run ../listing11.01/centinela.go      # terminal 5
```

Envía 9, 12, 10 y 2 al puerto 9080 y termina.

## Qué ocurre por dentro

Regla de cada nodo: **guardar el menor número visto y reenviar el mayor al
siguiente nodo.** Traza de una ejecución real:

| Nodo | Recibe | Se queda con | Reenvía |
|---|---|---|---|
| 9080 | 9, 12, 10, 2 | **2** | 12, 10, 9 |
| 9081 | 12, 10, 9 | **9** | 12, 10 |
| 9082 | 12, 10 | **10** | 12 |
| 9083 | 12 | **12** | (nada, es el último) |

El nodo *k* termina con el *k*-ésimo menor, y el resultado es 2, 9, 10, 12.
Por eso `-n` va 4, 3, 2, 1: cada nodo reenvía un número menos que los que
recibe.

Salida final esperada:

```
[NODO localhost:9080] valor final = 2
[NODO localhost:9081] valor final = 9
[NODO localhost:9082] valor final = 10
[NODO localhost:9083] valor final = 12
```

Nota: en los nodos 9081 y 9082 el orden de llegada (12, 10, 9) no es el orden
en que los envió el centinela. Cada conexión se atiende en su propia goroutine,
así que el orden puede cambiar entre ejecuciones. El resultado final siempre es
el mismo.

## Paso 6: detener todo

Los nodos siguen escuchando después de imprimir su resultado. Presione
`Ctrl+C` en cada terminal.

## Errores comunes

- **`No se pudo conectar a localhost:9080`**: se ejecutó el centinela antes
  que el nodo 9080.
- **`address already in use`**: un nodo de una ejecución anterior sigue
  activo. Deténgalo con `Ctrl+C` o con `pkill -f ordenamiento`.
- **Un nodo nunca imprime `valor final`**: su `-n` es demasiado alto y espera
  números que nunca llegan.
- **Se cambian los números del centinela**: hay que ajustar también la
  cantidad de nodos y los valores de `-n` (con *m* números se necesitan *m*
  nodos, con `-n` = *m*, *m*−1, …, 1).
