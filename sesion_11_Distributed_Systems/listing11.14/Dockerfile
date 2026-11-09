# Dockerfile simple (versión de la guía): útil en laboratorio, NO en producción.
# Problemas: incluye todo el SDK de Go (imagen grande) y compila con `go run`
# cada vez que arranca el contenedor (arranque lento).
FROM golang:1.24-alpine
WORKDIR /app
# Copia y renombra el archivo fuente
COPY ./codeApi.go ./api-svc001.go
# Documenta el puerto (no lo publica; eso lo hace `-p` o `ports:`)
EXPOSE 9015
CMD ["go", "run", "api-svc001.go"]
