# Imagen del nodo Hot Potato (listing11.06). El contexto de construcción es la
# carpeta de la sesión (..), por eso se copia desde listing11.06/.
FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY listing11.06/hot_potato_node.go .
RUN CGO_ENABLED=0 go build -o /hot_potato_node hot_potato_node.go

FROM alpine:3.20
COPY --from=builder /hot_potato_node /usr/local/bin/hot_potato_node
ENTRYPOINT ["hot_potato_node"]
