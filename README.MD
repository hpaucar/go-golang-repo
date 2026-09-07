# Crear, compilar y ejecutar proyectos Go en Ubuntu y Windows 11

Esta guía explica cómo instalar Go en Ubuntu y Windows 11, crear un proyecto, ejecutarlo, compilarlo y realizar las operaciones habituales desde la terminal.

## 1. Instalar Go

### 1.1. Instalación en Ubuntu

Actualiza la lista de paquetes e instala Go:

```bash
sudo apt update
sudo apt install golang-go
```

Comprueba la instalación:

```bash
go version
```

También puedes consultar la configuración activa:

```bash
go env
```

> Para instalar una versión más reciente que la disponible en Ubuntu, consulta la página oficial: <https://go.dev/doc/install>.

### 1.2. Instalación en Windows 11

1. Accede a la página oficial de descargas de Go: <https://go.dev/dl/>.
2. Descarga el instalador de Windows con extensión `.msi` correspondiente a tu arquitectura. En la mayoría de las computadoras con Windows 11 se utiliza la versión `x86-64` (`amd64`). Los equipos con procesador ARM deben utilizar `ARM64`.
3. Ejecuta el archivo descargado y sigue los pasos del asistente de instalación.
4. Mantén la carpeta de instalación predeterminada:

```text
C:\Program Files\Go
```

5. Cuando finalice la instalación, cierra las ventanas de PowerShell o Símbolo del sistema que estuvieran abiertas y abre una nueva terminal.

Comprueba desde PowerShell que Go se instaló correctamente:

```powershell
go version
```

Consulta la configuración activa:

```powershell
go env
```

Para comprobar dónde se encuentra el ejecutable de Go:

```powershell
Get-Command go
```

El instalador normalmente agrega automáticamente esta carpeta a la variable de entorno `Path`:

```text
C:\Program Files\Go\bin
```

Si PowerShell muestra que `go` no se reconoce como un comando, reinicia la terminal. Si el problema continúa, abre **Configuración avanzada del sistema → Variables de entorno**, edita la variable `Path` y verifica que la carpeta anterior esté incluida.

> Instala Go desde el sitio oficial. No es necesario configurar manualmente `GOROOT` para una instalación estándar.

## 2. Crear un proyecto nuevo

Crea una carpeta y entra en ella:

```bash
mkdir hola-go
cd hola-go
```

Inicializa un módulo. En un proyecto real conviene usar la dirección de su repositorio:

```bash
go mod init github.com/usuario/hola-go
```

Para un ejercicio local también puedes usar:

```bash
go mod init hola-go
```

Este comando crea `go.mod`, que identifica el módulo y registra sus dependencias.

## 3. Crear el programa principal

Crea el archivo `main.go`:

```bash
nano main.go
```

Agrega el siguiente código:

```go
package main

import "fmt"

func main() {
	message := "Hola desde Go en Ubuntu"
	fmt.Println(message)
}
```

En Nano, guarda con `Ctrl+O`, presiona `Enter` y sal con `Ctrl+X`.

La estructura inicial será:

```text
hola-go/
├── go.mod
└── main.go
```

## 4. Ejecutar sin conservar un binario

Para un archivo individual:

```bash
go run main.go
```

Para ejecutar el paquete completo ubicado en la carpeta actual:

```bash
go run .
```

La segunda opción es preferible cuando el paquete contiene varios archivos `.go`.

Para pasar argumentos al programa:

```bash
go run . argumento1 argumento2
```

## 5. Compilar el proyecto

Compila el paquete actual:

```bash
go build .
```

Se generará un ejecutable llamado `hola-go`. Ejecútalo con:

```bash
./hola-go
```

Para elegir otro nombre:

```bash
go build -o programa .
./programa
```

Para compilar solamente un archivo:

```bash
go build main.go
./main
```

## 6. Compilar y guardar el ejecutable en una carpeta

```bash
mkdir -p bin
go build -o bin/hola-go .
./bin/hola-go
```

Conviene agregar `bin/` al archivo `.gitignore`:

```gitignore
bin/
```

## 7. Gestionar dependencias

Para agregar una biblioteca:

```bash
go get github.com/google/uuid
```

Después puede importarse en el código:

```go
import "github.com/google/uuid"
```

Descarga las dependencias y elimina del módulo las que ya no se usan:

```bash
go mod tidy
```

Descarga las dependencias sin ejecutar el programa:

```bash
go mod download
```

Lista los módulos usados:

```bash
go list -m all
```

## 8. Formatear y revisar el código

Formatea todos los archivos del proyecto:

```bash
go fmt ./...
```

Ejecuta el análisis estático incluido con Go:

```bash
go vet ./...
```

Comprueba que todos los paquetes compilan:

```bash
go build ./...
```

## 9. Crear y ejecutar pruebas

Crea `main_test.go`:

```go
package main

import "testing"

func sumar(firstNumber, secondNumber int) int {
	return firstNumber + secondNumber
}

func TestSumar(t *testing.T) {
	result := sumar(2, 3)
	wanted := 5

	if result != wanted {
		t.Fatalf("sumar(2, 3) = %d; se esperaba %d", result, wanted)
	}
}
```

Ejecuta todas las pruebas:

```bash
go test ./...
```

Muestra información detallada:

```bash
go test -v ./...
```

Ejecuta una prueba específica:

```bash
go test -run TestSumar
```

Detecta posibles condiciones de carrera en programas concurrentes:

```bash
go test -race ./...
```

Calcula la cobertura:

```bash
go test -cover ./...
```

## 10. Instalar el ejecutable

Compila e instala el programa en la carpeta de binarios configurada por Go:

```bash
go install .
```

Consulta dónde se guardará:

```bash
go env GOBIN GOPATH
```

Normalmente se guarda en `$GOPATH/bin`, que suele corresponder a `$HOME/go/bin`.

## 11. Compilar para otro sistema operativo

Crear un ejecutable para Windows de 64 bits desde Ubuntu:

```bash
GOOS=windows GOARCH=amd64 go build -o programa.exe .
```

Para Linux ARM64:

```bash
GOOS=linux GOARCH=arm64 go build -o programa-linux-arm64 .
```

Consulta los destinos compatibles:

```bash
go tool dist list
```

## 12. Compilar en modo optimizado para distribución

Para reducir el tamaño del binario eliminando información de depuración:

```bash
go build -ldflags="-s -w" -o bin/hola-go .
```

## 13. Trabajar con un proyecto descargado de GitHub

Clona el repositorio y entra en él:

```bash
git clone https://github.com/usuario/repositorio.git
cd repositorio
```

Descarga y organiza dependencias:

```bash
go mod download
go mod tidy
```

Revisa, prueba, compila y ejecuta:

```bash
go fmt ./...
go vet ./...
go test ./...
go build .
go run .
```

Si el repositorio contiene varios comandos bajo `cmd/`, revisa su estructura:

```bash
find cmd -maxdepth 2 -type f -name '*.go'
```

Por ejemplo, para ejecutar un comando llamado `servidor`:

```bash
go run ./cmd/servidor
```

Y para compilarlo:

```bash
go build -o bin/servidor ./cmd/servidor
```

## 14. Comandos de uso frecuente

| Objetivo | Comando |
|---|---|
| Ver la versión | `go version` |
| Crear el módulo | `go mod init nombre-del-modulo` |
| Ejecutar el proyecto | `go run .` |
| Compilar el proyecto | `go build .` |
| Compilar con otro nombre | `go build -o programa .` |
| Ejecutar el binario | `./programa` |
| Formatear | `go fmt ./...` |
| Revisar errores frecuentes | `go vet ./...` |
| Ejecutar pruebas | `go test ./...` |
| Probar condiciones de carrera | `go test -race ./...` |
| Organizar dependencias | `go mod tidy` |
| Descargar dependencias | `go mod download` |
| Ver documentación de un paquete | `go doc fmt` |
| Limpiar la caché de compilación | `go clean -cache` |

## 15. Flujo de trabajo recomendado

Durante el desarrollo:

```bash
go fmt ./...
go test ./...
go run .
```

Antes de entregar o publicar:

```bash
go mod tidy
go fmt ./...
go vet ./...
go test -race ./...
go build -o bin/hola-go .
./bin/hola-go
```

## 16. Errores frecuentes

### `go: command not found`

Go no está instalado o su carpeta `bin` no está incluida en `PATH`. Comprueba:

```bash
which go
go version
```

### `go: go.mod file not found`

Estás fuera del proyecto o todavía no inicializaste el módulo:

```bash
go mod init nombre-del-modulo
```

### `package ... is not in std`

Comprueba que el nombre del módulo y la ruta de importación sean correctos. Luego ejecuta:

```bash
go mod tidy
```

### `permission denied` al ejecutar

Los binarios generados por Go normalmente ya son ejecutables. Si el permiso se perdió:

```bash
chmod +x programa
./programa
```

### El programa compilado no refleja los cambios

Vuelve a compilarlo antes de ejecutarlo:

```bash
go build -o programa .
./programa
```

## Licencia

Puedes adaptar libremente esta guía para tus prácticas, clases o proyectos.
