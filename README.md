# Curso de Go (Golang) - Platzi 🚀🐹

Este repositorio contiene los ejercicios y proyectos desarrollados durante el curso de Go de Platzi, abarcando desde conceptos básicos hasta temas avanzados como concurrencia, channels, y manejo de JSON.

## 📋 1. Requisitos

- Go versión 1.19 o superior (recomendado 1.21+)
- Editor de texto o IDE (como GoLand, Visual Studio Code con la extensión Go)
- Git para clonar el repositorio

## 🚀 2. Ejecutar

### 🖥️ 2.1. Terminal

```sh
# Ejemplos básicos
go run 01-HolaMundo/main.go
go run 02-Paquetes/main.go

# Ejemplos con estructuras de control
go run 05-Condiciones/main.go
go run 06-For/main.go
go run 07-Switch/main.go

# Ejemplos de concurrencia
go run 20-Goroutines/main.go
go run 21-Channels/main.go
go run 22-Timeouts/main.go

# Ejercicios prácticos
go run 21.2-Ejercicio_Carrera_Channels/main.go
go run 23.1-Ejercicio_Emails/main.go
```

### 2.2. Instalar dependencias para módulos específicos

Esta sección explica cómo instalar las dependencias necesarias para módulos específicos del proyecto.

```sh
# Para el módulo de variables de entorno
cd 28-VariablesEntorno

# Descarga e instala las dependencias listadas en go.mod
go mod tidy

# Instala la librería godotenv para cargar variables desde archivos .env
go get github.com/joho/godotenv
```

### ▶️ 2.3. Code Runner (Extensión de VS Code)

- Se selecciona todo el código
- Se da clic derecho
- `Run Code`

## 🗂️ 3. Estructura del Proyecto

### 📚 Conceptos Básicos

- **01-HolaMundo/** 📝: Ejemplo básico de un programa "Hola Mundo" en Go.
  - `main.go` 📄: Archivo fuente con el código del programa.
  - `main.exe` ⚙️: Archivo ejecutable compilado.
- **02-Paquetes/** 📦: Ejemplos de manejo de paquetes y funciones del sistema.
  - `main.go` 📄: Acceso a variables de entorno multiplataforma.
  - `ejemplo.txt` 🗒️: Archivo de ejemplo para operaciones de lectura.
- **03-Variables/** 🔤: Declaración y uso de variables en Go.
  - `main.go` 📄: Tipos de variables, declaraciones explícitas e implícitas.
- **04-Valores/** 🔢: Tipos de datos primitivos y valores en Go.
  - `main.go` 📄: Strings, números, booleanos y constantes.

### 🔀 Estructuras de Control

- **05-Condiciones/** 🔀: Estructuras condicionales (if, else if, else).
  - `main.go` 📄: Ejemplos de condicionales simples y anidadas.
- **06-For/** 🔁: Bucles for en todas sus variantes.
  - `main.go` 📄: For tradicional, while, range y bucles infinitos.
- **07-Switch/** 🎯: Estructura switch y sus diferentes usos.
  - `main.go` 📄: Switch básico, múltiples casos y switch sin expresión.

### 📊 Estructuras de Datos

- **08-Arreglos/** �: Arrays de tamaño fijo en Go.
  - `main.go` 📄: Declaración, inicialización y manipulación de arrays.
- **09-Slices/** 🍰: Slices (arrays dinámicos) y sus operaciones.
  - `main.go` 📄: Creación, append, copy y manipulación de slices.
- **10-Maps/** 🗺️: Mapas (hash tables) en Go.
  - `main.go` 📄: Creación, acceso, modificación y eliminación en maps.

### ⚙️ Funciones

- **11-Funciones/** 🔧: Funciones básicas en Go.
  - `main.go` 📄: Declaración, parámetros y valores de retorno.
- **12-FuncionesMultiples/** ↩️: Funciones con múltiples valores de retorno.
  - `main.go` 📄: Retorno múltiple y manejo de errores.
- **13-FuncionesVariaticas/** 📏: Funciones con número variable de argumentos.
  - `main.go` 📄: Funciones variádicas usando ... (spread operator).
- **14-FuncionesRecursivas/** 🔄: Funciones recursivas y casos base.
  - `main.go` 📄: Ejemplos de recursión (factorial, fibonacci).

### 🏗️ Conceptos Avanzados

- **15-Punteros/** 👉: Punteros y referencias en Go.
  - `main.go` 📄: Declaración, desreferenciación y uso de punteros.
- **16-Runas/** �: Runas y manejo de caracteres Unicode.
  - `main.go` 📄: Diferencias entre byte y rune, caracteres especiales.
- **17-Structs/** 🏢: Estructuras y tipos personalizados.
  - `main.go` 📄: Definición, inicialización y métodos de structs.
- **18-Enums/** 📝: Enumeraciones usando iota.
  - `main.go` 📄: Creación de constantes enumeradas.
- **19-Errores/** ❌: Manejo de errores en Go.
  - `main.go` 📄: Creación y manejo de errores personalizados.

### 🚀 Concurrencia

- **20-Goroutines/** ⚡: Programación concurrente con goroutines.
  - `main.go` 📄: Creación y ejecución de goroutines.
- **21-Channels/** 📡: Canales para comunicación entre goroutines.
  - `main.go` 📄: Channels buffered y unbuffered, select statements.
- **21.2-Ejercicio_Carrera_Channels/** 🏁: Ejercicio práctico de simulación de carrera.
  - `main.go` 📄: Implementación de carrera usando channels y goroutines.
- **21.3-Ejercicio_Hacer_Pedido_Channels/** �: Ejercicio de sistema de pedidos.
  - `main.go` 📄: Simulación de procesamiento de pedidos con concurrencia.
- **22-Timeouts/** ⏱️: Manejo de timeouts en operaciones concurrentes.
  - `main.go` 📄: Implementación de timeouts con select y time.After.
- **23-Workers/** 👷: Patrón Worker Pool para procesamiento concurrente.
  - `main.go` 📄: Pool de workers para procesamiento distribuido.
- **23.1-Ejercicio_Emails/** 📧: Sistema de procesamiento de emails.
  - `main.go` 📄: Ejercicio completo con workers, channels y sincronización.

### 📄 Serialización y Utilidades

- **25-json-decode/** 📥: Decodificación de JSON a structs.
  - `main.go` 📄: Unmarshaling JSON y manejo de tags.
- **26-json-encode/** 📤: Codificación de structs a JSON.
  - `main.go` 📄: Marshaling de estructuras a JSON.
- **27-hashes/** 🔐: Generación de hashes criptográficos.
  - `main.go` 📄: SHA256 y otras funciones hash.
- **28-VariablesEntorno/** 🌍: Manejo de variables de entorno.
  - `main.go` 📄: Lectura de variables usando godotenv.
  - `go.mod` 📋: Módulo Go con dependencias.
  - `.env` 🔒: Archivo de configuración de entorno.
  - `.env.example` 📋: Plantilla de variables de entorno.

## � 4. Cómo Empezar

### Para Principiantes

1. Comienza con los conceptos básicos: `01-HolaMundo`, `02-Paquetes`, `03-Variables`
2. Practica estructuras de control: `05-Condiciones`, `06-For`, `07-Switch`
3. Explora estructuras de datos: `08-Arreglos`, `09-Slices`, `10-Maps`

### Para Nivel Intermedio

1. Domina las funciones: `11-Funciones` hasta `14-FuncionesRecursivas`
2. Conceptos avanzados: `15-Punteros`, `17-Structs`, `19-Errores`

### Para Nivel Avanzado

1. Concurrencia: `20-Goroutines`, `21-Channels`, `22-Timeouts`
2. Ejercicios prácticos: `21.2-Ejercicio_Carrera_Channels`, `23.1-Ejercicio_Emails`
3. Utilidades: `25-json-decode`, `27-hashes`, `28-VariablesEntorno`

## 💡 5. Conceptos Clave Cubiertos

- ✅ **Sintaxis Básica**: Variables, tipos de datos, operadores
- ✅ **Control de Flujo**: Condicionales, bucles, switch
- ✅ **Estructuras de Datos**: Arrays, slices, maps, structs
- ✅ **Funciones**: Declaración, parámetros, retorno múltiple, recursión
- ✅ **Programación Orientada a Objetos**: Structs, métodos, interfaces implícitas
- ✅ **Manejo de Errores**: Error handling idiomático de Go
- ✅ **Concurrencia**: Goroutines, channels, select, patterns
- ✅ **Serialización**: JSON encoding/decoding
- ✅ **Seguridad**: Hashing, variables de entorno
- ✅ **Herramientas**: Módulos, dependencias, testing

## �📚 6. Recursos Adicionales

- [Documentación oficial de Go](https://golang.org/doc/) 🌐
- [Tour de Go](https://tour.golang.org/welcome/1) 🧭
- [Go by Example](https://gobyexample.com/) 📖
- [Effective Go](https://golang.org/doc/effective_go.html) ⚡
- [Curso de Go en Platzi](https://platzi.com/cursos/go-golang/) 🎓
- [Go Blog](https://blog.golang.org/) 📝
- [Awesome Go](https://github.com/avelino/awesome-go) ⭐

## 🤝 7. Contribución

Si deseas contribuir a este repositorio:

1. Haz fork del proyecto
2. Crea una nueva rama (`git checkout -b feature/nueva-caracteristica`)
3. Confirma tus cambios (`git commit -am 'Añade nueva característica'`)
4. Sube tu rama (`git push origin feature/nueva-caracteristica`)
5. Abre un Pull Request

---

**Última actualización**: 13 de julio de 2025  
**Autor**: [afperdomo2](https://github.com/afperdomo2)  
**Repositorio**: [curso_go_platzi](https://github.com/afperdomo2/curso_go_platzi)
