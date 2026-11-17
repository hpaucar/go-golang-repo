# Listado 11.02: Nodo de ordenamiento distribuido en cadena

`ordenamiento.go` es **un nodo** del ordenamiento distribuido en cadena (versión
de clase, en texto plano). Se lanza varias veces, una por nodo, con flags
distintos. Los datos los inyecta el centinela del listado 11.01
(`../listing11.01/centinela.go`).

**Regla: primero los nodos, al final el centinela.**

## Topología: una cadena

```
centinela ──► 9080 ──► 9081 ──► 9082 ──► 9083
 (11.01)      (servidor y cliente)       (solo servidor)
```

Con 4 números se necesitan 4 nodos. El último nodo no tiene `-remoto`.

## Qué hace cada nodo

Regla: **guardar el menor número visto y reenviar el mayor al siguiente nodo.**

```go
cont := <-chCont      // entrar a la sección crítica

if cont == 0 {
	min = dNum        // primer valor: se guarda
} else if dNum < min {
	enviar(min)       // el antiguo mínimo pasa al siguiente nodo
	min = dNum
} else {
	enviar(dNum)      // el recibido es mayor: pasa al siguiente nodo
}

cont++
...
chCont <- cont        // salir de la sección crítica
```

Dos detalles de Go:

- **Una goroutine por conexión.** `main` acepta conexiones en un bucle y
  lanza `go handle(conn)`. Varias goroutines pueden leer y cambiar `min` a la
  vez.
- **Un canal como cerrojo.** `chCont` es un canal con búfer de tamaño 1 que
  guarda el contador. Quien saca el valor entra a la sección crítica; al
  devolverlo, sale. Hace el mismo trabajo que un `sync.Mutex`. La
  sincronización es **dentro** del nodo; **entre** nodos solo hay mensajes TCP.

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-local` | `9080` | puerto donde escucha este nodo |
| `-remoto` | (vacío) | puerto del siguiente nodo. Vacío = **último nodo** |
| `-n` | `1` | cuántos números recibirá este nodo antes de imprimir su resultado |

## Paso 0: abrir 5 terminales

En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

## Pasos 1 a 4: levantar los nodos, del último al primero

```bash
go run listing11.02/ordenamiento.go -local 9083 -n 1                 # terminal 1 (último nodo)
go run listing11.02/ordenamiento.go -local 9082 -remoto 9083 -n 2    # terminal 2
go run listing11.02/ordenamiento.go -local 9081 -remoto 9082 -n 3    # terminal 3
go run listing11.02/ordenamiento.go -local 9080 -remoto 9081 -n 4    # terminal 4
```

Cada nodo muestra una línea como:

```
[NODO 9080] escuchando; espera 4 valores; siguiente = "localhost:9081".
```

## Paso 5: ejecutar el centinela (terminal 5)

```bash
go run listing11.01/centinela.go
```

```
[CENTINELA] Enviado 9 a localhost:9080
[CENTINELA] Enviado 12 a localhost:9080
[CENTINELA] Enviado 10 a localhost:9080
[CENTINELA] Enviado 2 a localhost:9080
```

## Resultado esperado

Terminal 4 (nodo 9080):

```
Número recibido = 9
Número recibido = 12
Número recibido = 10
Número recibido = 2
[NODO localhost:9080] valor que se queda = 2
```

Última línea de cada nodo:

| Nodo | Recibe | Se queda con | Reenvía |
|---|---|---|---|
| 9080 | 9, 12, 10, 2 | **2** | 12, 10, 9 |
| 9081 | 12, 10, 9 | **9** | 12, 10 |
| 9082 | 12, 10 | **10** | 12 |
| 9083 | 12 | **12** | (nada) |

El nodo *k* se queda con el *k*-ésimo menor: 2, 9, 10, 12. El orden de llegada
puede cambiar entre ejecuciones (una goroutine por conexión), pero el
resultado final es siempre el mismo.

## Paso 6: detener todo

Los nodos siguen escuchando. Presione `Ctrl+C` en las terminales 1 a 4.

## Errores comunes

- **`Error al escuchar: ... address already in use`**: un nodo de una
  ejecución anterior sigue activo. Deténgalo con `Ctrl+C`.
- **`no se pudo enviar X a localhost:9081: ... connection refused`**: el
  siguiente nodo no estaba levantado y el número se perdió.
- **`soy el último nodo; descarto X (¿-n mal configurado?)`**: hay menos
  nodos que números, o el `-n` del último nodo es menor que 1.
- **Un nodo nunca imprime `valor que se queda`**: su `-n` es demasiado alto y
  espera números que nunca llegan.

## Para pensar y ejercicios

1. **Canal o mutex.** Reemplace `chCont` por un `sync.Mutex` y un contador
   entero. ¿Cambia el comportamiento?
2. **Más números.** Cambie el centinela para enviar 6 números. ¿Cuántos nodos
   hacen falta y qué valores de `-n` lleva cada uno?
3. **El problema del orden.** Cada número viaja en una conexión distinta. Con
   un mensaje de fin (como el `FIN` del listado 11.07), ¿podría adelantarse a
   un número? ¿Por qué el listado 11.07 usa una sola conexión?
