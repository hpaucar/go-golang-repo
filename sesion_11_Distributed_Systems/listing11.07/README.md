# Listado 11.07: Ordenamiento en pipeline con mensaje FIN

`pipeline_sort_node.go` es **un nodo** de un ordenamiento en pipeline
("extrusión de mínimos"). Es la versión JSON del ordenamiento en cadena de los
listados 11.01 y 11.02, con dos mejoras:

1. Un mensaje **`FIN`** avisa que ya no llegan más datos. Así los nodos no
   necesitan un flag `-n`.
2. Cada nodo usa **una sola conexión TCP** con el siguiente, y no una por
   número. Así los mensajes llegan en orden.

**Regla: primero los nodos, al final el centinela.**

## Topología: un pipeline de 5 nodos

```
centinela ══► 9001 ══► 9002 ══► 9003 ══► 9004 ══► 9005
 (11.03)                                         (último)
   ═══ una conexión TCP abierta todo el tiempo
```

Con 5 números se necesitan 5 nodos. El último nodo no tiene `-siguiente`.

## La regla de cada nodo

```go
switch msg.Tipo {
case "FIN":
	// imprime su resultado y propaga el FIN
case "DATA":
	switch {
	case !tieneVal:            // primer valor: se guarda
	case nuevo < minVal:       // llegó uno menor: expulsa el antiguo mínimo
	default:                   // llegó uno mayor: lo expulsa
	}
}
```

Igual que en 11.02: **guardar el menor y expulsar el mayor**. Al final, el
nodo *k* guarda el *k*-ésimo menor.

## ¿Por qué una sola conexión?

TCP conserva el orden **dentro** de una conexión, pero **no entre**
conexiones distintas. Si cada número viajara en su propia conexión, el `FIN`
podría llegar antes que un dato, y el nodo imprimiría su resultado
demasiado pronto. Por eso:

- El centinela envía los 5 números y el `FIN` por la **misma** conexión
  (`enviarFlujo`).
- Cada nodo abre la conexión con el siguiente **una vez** y la guarda en la
  variable `salida`.
- `reenviar` se llama con el `sync.Mutex` tomado: los mensajes salen en el
  mismo orden en que se decidieron.

## Los mensajes

```json
{"tipo":"DATA","origen":"CENTINELA","valores":[45]}
{"tipo":"FIN","origen":"CENTINELA"}
```

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-local` | `localhost:9001` | dirección donde escucha |
| `-siguiente` | (vacío) | siguiente nodo. Vacío = **último nodo** |

## Paso 0: abrir 6 terminales

En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

## Pasos 1 a 5: levantar los nodos, del último al primero

```bash
go run listing11.07/pipeline_sort_node.go -local localhost:9005                             # terminal 1
go run listing11.07/pipeline_sort_node.go -local localhost:9004 -siguiente localhost:9005   # terminal 2
go run listing11.07/pipeline_sort_node.go -local localhost:9003 -siguiente localhost:9004   # terminal 3
go run listing11.07/pipeline_sort_node.go -local localhost:9002 -siguiente localhost:9003   # terminal 4
go run listing11.07/pipeline_sort_node.go -local localhost:9001 -siguiente localhost:9002   # terminal 5
```

Cada nodo muestra:

```
[PIPELINE localhost:9001] Listo. Siguiente: "localhost:9002"
```

## Paso 6: ejecutar el centinela (terminal 6)

```bash
go run listing11.03/centinela_json.go -modo pipeline
```

```
[CENTINELA SENT] DATA 45 -> localhost:9001
[CENTINELA SENT] DATA 12 -> localhost:9001
[CENTINELA SENT] DATA 89 -> localhost:9001
[CENTINELA SENT] DATA 3 -> localhost:9001
[CENTINELA SENT] DATA 27 -> localhost:9001
[CENTINELA SENT] FIN -> localhost:9001
```

## Resultado esperado

Terminal 5 (9001, el primer nodo):

```
[PIPELINE localhost:9001] Primer valor guardado: 45
[PIPELINE localhost:9001] Nuevo MIN 12; expulsa 45
[PIPELINE localhost:9001] Mantiene MIN 12; expulsa 89
[PIPELINE localhost:9001] Nuevo MIN 3; expulsa 12
[PIPELINE localhost:9001] Mantiene MIN 3; expulsa 27
[PIPELINE localhost:9001] RESULTADO: 3
```

Última línea de cada nodo:

| Nodo | Recibe | `RESULTADO` | Expulsa |
|---|---|---|---|
| 9001 | 45, 12, 89, 3, 27 | **3** | 45, 89, 12, 27 |
| 9002 | 45, 89, 12, 27 | **12** | 89, 45, 27 |
| 9003 | 89, 45, 27 | **27** | 89, 45 |
| 9004 | 89, 45 | **45** | 89 |
| 9005 | 89 | **89** | (nada) |

Ordenado: 3, 12, 27, 45, 89. Los `RESULTADO` aparecen del primer nodo al
último, porque el `FIN` recorre el pipeline en ese orden.

## Paso 7: detener todo

Presione `Ctrl+C` en las terminales 1 a 5. Antes de repetir la prueba hay que
reiniciar los nodos: cada uno guarda su `minVal` de la ejecución anterior.

## Errores comunes

- **`[CENTINELA ERROR] No se pudo conectar a localhost:9001`**: el nodo 9001
  no estaba levantado.
- **`No se pudo conectar a localhost:9003`**: un nodo intermedio no estaba
  levantado. Los valores que debía recibir se pierden.
- **`Último nodo: se descarta [X] (faltan nodos)`**: hay menos nodos que
  números.
- **`RESULTADO: (sin valor)`**: hay más nodos que números; a ese nodo no le
  llegó ningún dato.
- **`panic: listen tcp 127.0.0.1:9001: bind: address already in use`**: otro
  patrón (Conway, papa caliente, fan-out) sigue usando el puerto 9001.

## Para pensar y ejercicios

1. **Más valores que nodos.** Haga que el último nodo guarde una lista
   ordenada en lugar de un solo valor.
2. **Romper el orden.** Cambie `reenviar` para que abra una conexión nueva por
   mensaje (como en 11.02). Ejecute varias veces: ¿aparece algún `RESULTADO`
   antes de tiempo?
3. **Comparar con 11.02.** ¿Qué ventaja tiene `FIN` frente al flag `-n`?
