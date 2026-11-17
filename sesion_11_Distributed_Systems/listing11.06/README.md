# Listado 11.06: Papa caliente en anillo (token ring)

`hot_potato_node.go` es **un nodo** de un anillo. Por el anillo circula un
único **token** JSON con un contador. Cada nodo lo recibe, resta 1 y lo pasa
al siguiente. **Pierde el nodo que recibe el token con el contador en 0.**

El mismo programa se lanza tres veces con flags distintos. El centinela del
listado 11.03 (`-modo potato`) pone el token en el anillo.

**Regla: primero los nodos, al final el centinela.**

## Topología: un anillo

```
                 centinela
                     │ token (pasos = 5)
                     ▼
            ┌──► 9001 ──► 9002 ──┐
            │                    │
            └─────── 9003 ◄──────┘
```

| Proceso | Escucha | Envía a | Rol |
|---|---|---|---|
| centinela (11.03) | no | 9001 | solo cliente, una vez |
| nodo 9001 | sí | 9002 | servidor y cliente |
| nodo 9002 | sí | 9003 | servidor y cliente |
| nodo 9003 | sí | **9001** | servidor y cliente |

A diferencia de la cadena del listado 11.01, aquí **el último nodo se conecta
con el primero**: eso cierra el anillo. Todos los nodos son iguales; ninguno
es "el último".

## El token

```json
{"tipo":"TOKEN","origen":"CENTINELA","pasos":5}
```

Cada nodo hace:

```go
if msg.PasosRestantes <= 0 {
	// ¡PAPA CALIENTE! este nodo pierde y el token deja de circular
	return
}
msg.PasosRestantes--
msg.Origen = yo
enviar(siguiente, msg)
```

Hay **un solo token** en todo el sistema: en cada momento, solo su dueño
actúa. Es la idea base de la exclusión mutua por token (Ben-Ari, sección 10.6;
ver listado 11.13).

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-local` | `localhost:9001` | dirección donde escucha (host:puerto) |
| `-siguiente` | `localhost:9002` | siguiente nodo del anillo (host:puerto) |
| `-nombre` | el valor de `-local` | nombre que se muestra en pantalla |

Los flags reciben direcciones completas y no solo puertos. Así el mismo
programa funciona en Docker con nombres de servicio (`nodo2:9001`, listado
11.15).

## Paso 0: abrir 4 terminales

En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

## Pasos 1 a 3: levantar los tres nodos del anillo

```bash
go run listing11.06/hot_potato_node.go -local localhost:9001 -siguiente localhost:9002   # terminal 1
go run listing11.06/hot_potato_node.go -local localhost:9002 -siguiente localhost:9003   # terminal 2
go run listing11.06/hot_potato_node.go -local localhost:9003 -siguiente localhost:9001   # terminal 3
```

Cada nodo muestra:

```
[NODO localhost:9001] Escuchando en localhost:9001. Siguiente en el anillo: localhost:9002
```

## Paso 4: ejecutar el centinela (terminal 4)

```bash
go run listing11.03/centinela_json.go -modo potato -pasos 5
```

```
[CENTINELA SENT] TOKEN -> localhost:9001
```

## Resultado esperado

Cada pase tarda medio segundo (latencia simulada con `time.Sleep`), así que se
puede seguir el token a ojo.

Terminal 1 (9001):

```
[NODO localhost:9001] Token recibido de CENTINELA con pasos = 5
[NODO localhost:9001] Reenviando token con pasos = 4 a localhost:9002
[NODO localhost:9001] Token recibido de localhost:9003 con pasos = 2
[NODO localhost:9001] Reenviando token con pasos = 1 a localhost:9002
```

Terminal 3 (9003):

```
[NODO localhost:9003] Token recibido de localhost:9002 con pasos = 3
[NODO localhost:9003] Reenviando token con pasos = 2 a localhost:9001
[NODO localhost:9003] Token recibido de localhost:9002 con pasos = 0
[NODO localhost:9003] ¡PAPA CALIENTE! El contador llegó a 0: este nodo pierde.
```

Traza completa:

| Pase | Nodo | Recibe `pasos` | Envía `pasos` |
|---|---|---|---|
| 1 | 9001 | 5 | 4 → 9002 |
| 2 | 9002 | 4 | 3 → 9003 |
| 3 | 9003 | 3 | 2 → 9001 |
| 4 | 9001 | 2 | 1 → 9002 |
| 5 | 9002 | 1 | 0 → 9003 |
| 6 | 9003 | 0 | **pierde** |

Con `-pasos N` y 3 nodos, pierde el nodo en la posición `N mod 3` (contando
9001 como 0). Pruebe `-pasos 3` (pierde 9001) y `-pasos 4` (pierde 9002).

## Paso 5: detener todo

Presione `Ctrl+C` en las terminales 1 a 3. Si ejecuta el centinela otra vez
sin detenerlos, empieza una partida nueva.

## Errores comunes

- **`[CENTINELA ERROR] No se pudo conectar a localhost:9001`**: el nodo 9001
  no estaba levantado.
- **`¡Token perdido! No se pudo conectar a localhost:9003`**: un nodo del
  anillo no estaba activo. El token desaparece y la partida termina sin
  perdedor. Es el punto débil de un anillo: un solo nodo caído lo rompe.
- **`Error al escuchar: ... address already in use`**: otro programa usa ese
  puerto (por ejemplo, el compresor de Conway en 9001). Deténgalo.

## Para pensar y ejercicios

1. **Dos tokens.** Ejecute el centinela dos veces seguidas, muy rápido. ¿Qué
   pasa con la regla "un solo dueño"?
2. **Nodo caído.** Detenga el nodo 9002 a mitad de la partida. ¿Cómo podría el
   nodo 9001 saltarlo y enviar el token a 9003?
3. **Más nodos.** Agregue un cuarto nodo (9004) al anillo. ¿Qué flags cambian?
4. **Elección de líder.** Use este anillo como base para el algoritmo de
   Chang-Roberts (ejercicio 5 de la sesión).
