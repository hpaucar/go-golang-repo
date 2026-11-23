# Imagen del nodo de consenso (listing14.05). El contexto de construcción es la
# carpeta de la sesión (..), por eso se copia desde listing14.05/.
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY listing14.05/consenso_mesh.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /consenso_node consenso_mesh.go

FROM alpine:3.21
COPY --from=builder /consenso_node /usr/local/bin/consenso_node
EXPOSE 8000
ENTRYPOINT ["consenso_node"]
