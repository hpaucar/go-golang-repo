# Sesión 11: Algoritmos distribuidos con memoria distribuida en Go

Programación Concurrente y Distribuida (CC57), UPC.

En el modelo de **memoria distribuida** los nodos no comparten variables. Cada
nodo es un proceso con su propia memoria, y solo se comunica con los demás
**enviando mensajes** (aquí, JSON sobre sockets TCP). Dentro de un nodo sí puede
haber varias goroutines; esas se sincronizan con `sync.Mutex` o con canales.

## Diapositivas

[`Sesion_11_Algoritmos_Distribuidos.pptx`](Sesion_11_Algoritmos_Distribuidos.pptx):
29 diapositivas 16:9 con código, diagramas y notas del orador.

## Objetivos

1. Construir nodos que escuchan (`net.Listen`) y envían (`net.Dial`) mensajes
   JSON.
2. Reconocer los patrones clásicos: pipeline, anillo con token, fan-out y fan-in.
3. Implementar la exclusión mutua distribuida: Ricart-Agrawala con permisos y
   con token (Ben-Ari, cap. 10).
4. Empaquetar un nodo Go en Docker y levantar varios nodos con Docker Compose.

## Material base

- Guía de estudio "Semana 11: Algoritmos Distribuidos bajo el Modelo de Memoria
  Distribuida en Golang" (Ms. Carlos Jara García).
- M. Ben-Ari, *Principles of Concurrent and Distributed Programming*, 2.ª ed.,
  cap. 10 "Distributed Algorithms".
- Guías de Dockerfile y Docker Compose del curso.

Los programas de esta carpeta parten de ese material, con correcciones. Cada
archivo explica sus cambios en la sección `Notas` de su cabecera.

## Programas

| Listado | Archivo | Tema | Puertos |
|---|---|---|---|
| 11.01 | `centinela.go` | Centinela en texto plano | → 9080 |
| 11.02 | `ordenamiento.go` | Ordenamiento distribuido en cadena (versión de clase) | 9080–9083 |
| 11.03 | `centinela_json.go` | Centinela JSON con `-modo conway\|potato\|pipeline\|fanout` | → 9001 |
| 11.04 | `conway_compressor.go` | Conway, nodo 1: compresión RLE | 9001 |
| 11.05 | `conway_output.go` | Conway, nodo 2: salida en líneas de ancho K | 9002 |
| 11.06 | `hot_potato_node.go` | Papa caliente en anillo (token ring) | 9001–9003 |
| 11.07 | `pipeline_sort_node.go` | Ordenamiento en pipeline con mensaje FIN | 9001–9005 |
| 11.08 | `fanout_master.go` | Fan-out: reparte bloques de tareas | 9001 |
| 11.09 | `fanout_worker.go` | Worker (eleva al cuadrado) | 9002, 9003 |
| 11.10 | `fanin_aggregator.go` | Fan-in: suma los resultados | 9010 |
| 11.11 | `recurso_compartido.go` | "Observador" que detecta violaciones de exclusión mutua | 9100 |
| 11.12 | `ricart_agrawala.go` | Ricart-Agrawala con permisos (Ben-Ari, Alg. 10.2) | 9101–9103 |
| 11.13 | `ra_token.go` | Ricart-Agrawala con token (Ben-Ari, Alg. 10.3) | 9101–9103 |
| 11.14 | `codeApi.go`, `Dockerfile`, `Dockerfile.multistage`, `docker-compose.yml` | API REST en contenedor | 9015 |
| 11.15 | `Dockerfile`, `docker-compose.yml` | Anillo de papa caliente: un contenedor por nodo | 9001 |

## Cómo ejecutar

Regla general: **primero los nodos, al final el centinela.** Cada nodo va en
su propia terminal y se detiene con `Ctrl+C`. Las órdenes exactas están en la
cabecera de cada archivo. Ejemplos:

```bash
# Conway (3 terminales, desde listing11.04)
go run ../listing11.05/conway_output.go
go run conway_compressor.go
go run ../listing11.03/centinela_json.go -modo conway
# Resultado: 5A3B2C1D5A en líneas de 4 caracteres

# Ricart-Agrawala (4 terminales, desde listing11.12)
go run ../listing11.11/recurso_compartido.go
go run ricart_agrawala.go -id 1
go run ricart_agrawala.go -id 2
go run ricart_agrawala.go -id 3
# Resultado: 0 violaciones. Repita con -sinra en los tres nodos: aparecen violaciones.

# Anillo en Docker (desde listing11.15)
docker compose up --build
go run ../listing11.03/centinela_json.go -modo potato -pasos 5
```

Resultados esperados (comprobados):

| Prueba | Resultado |
|---|---|
| 11.01 + 11.02 con 9, 12, 10, 2 | nodos 9080..9083 = 2, 9, 10, 12 |
| Conway `AAAAABBBCCDAAAAA` | `5A3B2C1D5A` |
| Papa caliente, 5 pasos | pierde el nodo 9003 |
| Pipeline 45, 12, 89, 3, 27 | nodos 9001..9005 = 3, 12, 27, 45, 89 |
| Fan-out/fan-in 1..10 | total = 385 |
| Ricart-Agrawala, 3 nodos × 4 rondas | 12 entradas, 0 violaciones |
| Lo mismo con `-sinra` | varias violaciones (8 en nuestra prueba) |
| RA con token | 12 entradas, 0 violaciones |

## Ideas clave para la clase

1. **El orden de los mensajes.** TCP conserva el orden dentro de UNA conexión,
   pero no entre conexiones distintas. Por eso el pipeline (11.07) usa una
   conexión persistente: si no, el mensaje FIN podría llegar antes que un dato.
2. **Dos niveles de sincronización.** Entre nodos: mensajes. Dentro de un nodo:
   `sync.Mutex`, porque cada conexión se atiende en una goroutine. Ben-Ari
   (sección 10.1) supone que dentro del nodo todo es atómico; en Go eso hay que
   programarlo.
3. **Errores de red.** Un nodo caído o que aún no arrancó hace fallar
   `net.Dial`. Si no se revisa `err`, `conn` es `nil` y el programa cae con
   panic. Este era el fallo de las versiones originales de 11.01 y 11.02.
4. **Costo en mensajes.** Ricart-Agrawala con permisos usa 2(N−1) mensajes por
   cada entrada a la sección crítica. Con token, un nodo que ya tiene el token
   entra sin enviar ningún mensaje.
5. **Nodos inactivos.** En RA, un nodo que terminó su trabajo debe seguir
   respondiendo a los demás. Si no, los bloquea para siempre.

## Ejercicios

**De la guía (retos de entorno real):**

1. Auditoría de transacciones fraudulentas (pipeline de 3 nodos).
2. Balanceador de carga con token en anillo.
3. MapReduce de logs de seguridad (fan-out + fan-in).
4. Rate limiting distribuido (token bucket).
5. Elección de líder en anillo (Chang-Roberts). Base: 11.06.
6. Broker publicador/suscriptor por temas.
7. Barrera distribuida (fan-in + fan-out). Base: 11.08–11.10.
8. Sharding por hash (`user_id % num_shards`).
9. Two-Phase Commit (2PC) entre Pago e Inventario.
10. Deduplicación con ventana deslizante de UUIDs.

**Sobre el código de esta sesión:**

11. Conway (11.04): aplique la regla de Ben-Ari: una letra sola no lleva "1" y
    las corridas se cortan en 9.
12. Pipeline (11.07): haga que funcione con más valores que nodos.
13. Ricart-Agrawala (11.12): cuente los mensajes enviados y compare con
    2(N−1) por entrada. Haga lo mismo en 11.13.
14. Ricart-Agrawala (11.12): ¿qué pasa si dos nodos usan el mismo `-id`?
    (Ben-Ari, ejercicio 3).
15. Implemente el algoritmo de Neilsen-Mizuno (Ben-Ari, Alg. 10.4) con la misma
    estructura que 11.13.
16. Levante los nodos de 11.12 con Docker Compose, como en 11.15.
