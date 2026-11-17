# Listado 11.14: API REST en Docker y Docker Compose

`codeApi.go` es un **microservicio REST** mínimo en Go: responde
`GET /alumnos/` con una lista de alumnos en JSON. El código es sencillo a
propósito; el objetivo es practicar tres formas de ejecutarlo:

1. **Local**, con `go run`.
2. **En un contenedor**, con dos Dockerfiles: uno simple y uno multi-stage.
3. **Con Docker Compose**, que construye y levanta el contenedor con un solo
   comando.

Aquí no hay centinela: el cliente es `curl` (o el navegador).

## Archivos

| Archivo | Para qué sirve |
|---|---|
| `codeApi.go` | el servidor HTTP (puerto de la variable `PORT`, por defecto 9015) |
| `Dockerfile` | imagen simple: SDK de Go completo + `go run` al arrancar |
| `Dockerfile.multistage` | compila en una etapa y copia solo el binario a una imagen Alpine |
| `docker-compose.yml` | un servicio `nodo-api` construido con el Dockerfile multi-stage |

## Requisitos

- Go 1.22 o superior (para el paso local).
- Docker con el plugin Compose (`docker compose version` debe responder).

## Paso 0: ir a la carpeta del listado

Todos los comandos de esta guía se ejecutan **desde `listing11.14`**, porque
los Dockerfiles copian `codeApi.go` de la carpeta actual:

```bash
cd sesion_11_Distributed_Systems/listing11.14
```

## Paso 1: ejecutar en local

Terminal 1:

```bash
go run codeApi.go
```

```
Servidor en http://localhost:9015/alumnos/
```

Terminal 2:

```bash
curl http://localhost:9015/alumnos/
```

```json
[{"cod":"123","nom":"Juan","dni":45678912},{"cod":"456","nom":"Luis","dni":78945612},{"cod":"789","nom":"Jorge","dni":12345678}]
```

La terminal 1 muestra `Invoca listar alumnos` por cada petición. Detenga el
servidor con `Ctrl+C` **antes del paso 2**: si no, el puerto 9015 queda
ocupado.

El puerto se lee de la variable de entorno `PORT`. Pruebe:

```bash
PORT=8080 go run codeApi.go
curl http://localhost:8080/alumnos/
```

## Paso 2: Dockerfile simple

```bash
docker build -t api-simple -f Dockerfile .
docker run --rm -p 9015:9015 api-simple
```

En otra terminal: `curl http://localhost:9015/alumnos/`. Detenga con `Ctrl+C`.

Esta imagen ejecuta `go run` **cada vez que arranca el contenedor**: tarda
unos segundos en responder la primera vez, porque compila. Si `curl` falla
con `Connection reset` o `refused`, espere un poco y repita.

`-p 9015:9015` publica el puerto: `HOST:CONTENEDOR`. Sin `-p`, el servidor
corre pero no es accesible desde fuera del contenedor (`EXPOSE` solo
documenta el puerto).

## Paso 3: Dockerfile multi-stage

```bash
docker build -t api-multi -f Dockerfile.multistage .
docker run --rm -p 9015:9015 api-multi
```

Pruebe con `curl` y detenga con `Ctrl+C`. Esta vez responde de inmediato: el
binario ya está compilado.

Compare el tamaño de las dos imágenes:

```bash
docker images | grep api-
```

```
api-multi    latest   ...   16.2MB
api-simple   latest   ...   262MB
```

| | `Dockerfile` | `Dockerfile.multistage` |
|---|---|---|
| Imagen final | `golang:1.24-alpine` (con el SDK) | `alpine:3.20` + un binario |
| Tamaño | ~260 MB | ~16 MB |
| Arranque | compila cada vez | inmediato |
| Uso | laboratorio | producción |

`CGO_ENABLED=0` genera un binario **estático** (no depende de la libc), por
eso puede correr en una imagen mínima.

## Paso 4: Docker Compose

```bash
docker compose up -d --build
```

```
 Container node-api  Started
```

| Comando | Qué hace |
|---|---|
| `docker compose up -d --build` | construye la imagen y levanta el contenedor en segundo plano (`-d`) |
| `docker compose ps` | muestra el estado y los puertos |
| `docker compose logs -f nodo-api` | muestra los logs en vivo (`Ctrl+C` para salir de los logs) |
| `docker compose down` | detiene y borra el contenedor y su red |

Pruebe:

```bash
curl http://localhost:9015/alumnos/
docker compose ps
docker compose logs nodo-api
```

```
node-api  Up 2 seconds  0.0.0.0:9015->9015/tcp
node-api  | Servidor en http://localhost:9015/alumnos/
node-api  | 2026/11/16 23:35:38 Invoca listar alumnos
```

## Paso 5: detener todo

```bash
docker compose down
```

Para borrar también las imágenes de práctica:

```bash
docker rmi api-simple api-multi
```

## Errores comunes

- **`failed to bind host port 0.0.0.0:9015/tcp: address already in use`**:
  el servidor local del paso 1 u otro contenedor sigue usando el puerto 9015.
  Deténgalo con `Ctrl+C`, `docker ps` + `docker stop <id>`, o
  `docker compose down`.
- **`COPY failed: file not found ... codeApi.go`**: se ejecutó `docker build`
  desde otra carpeta. Vaya a `listing11.14`.
- **`curl: (56) Recv failure: Connection reset by peer`** con `api-simple`:
  el contenedor todavía está compilando. Espere unos segundos.
- **`permission denied while trying to connect to the Docker daemon`**
  (Linux): agregue su usuario al grupo `docker` o use `sudo`.

## Para pensar y ejercicios

1. **Otro puerto.** Cambie `PORT` y `ports:` en `docker-compose.yml` para
   publicar el servicio en el puerto 8080 del host, sin tocar el código.
2. **Dos réplicas.** Agregue un segundo servicio `nodo-api-2` en el puerto
   9016 del host. ¿Hace falta otro Dockerfile?
3. **Nuevo endpoint.** Agregue `GET /alumnos/{cod}` que devuelva un solo
   alumno (Go 1.22 acepta patrones con `{cod}` en `http.HandleFunc`).
4. **Uno de esta sesión.** Escriba un Dockerfile multi-stage para el
   compresor de Conway (11.04), como se hace en 11.15 con la papa caliente.
