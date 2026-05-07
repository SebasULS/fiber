# Go + Fiber — Scaffolding con CarryCoders

> Proyecto de exposición que demuestra cómo estructurar una API REST usando **Go** y el framework **Fiber**, aplicando el patrón de **Scaffolding**.

---

## Tabla de Contenidos

1. [¿Qué es Scaffolding?](#qué-es-scaffolding)
2. [Go — El Lenguaje](#go--el-lenguaje)
3. [Fiber — El Framework](#fiber--el-framework)
4. [Fiber + Scaffolding: La Combinación Perfecta](#fiber--scaffolding-la-combinación-perfecta)
5. [Base de Datos con GORM + SQLite](#base-de-datos-con-gorm--sqlite)
6. [Instalación de Go](#instalación-de-go)
7. [Instalación de Fiber](#instalación-de-fiber)
8. [Estructura del Proyecto](#estructura-del-proyecto)
9. [Cómo correr el proyecto](#cómo-correr-el-proyecto)
10. [Endpoints disponibles](#endpoints-disponibles)

---

## ¿Qué es Scaffolding?

**Scaffolding** (andamiaje) es una técnica de desarrollo donde se genera una **estructura base y predefinida** de un proyecto antes de escribir la lógica real. Funciona como un esqueleto que incluye:

- Organización de carpetas y archivos
- Separación de responsabilidades (rutas, handlers, modelos, middleware)
- Convenciones de nomenclatura
- Archivos de configuración listos para usar

### ¿Por qué usarlo?

| Sin Scaffolding | Con Scaffolding |
|---|---|
| Todo en un solo archivo | Código organizado por capas |
| Difícil de mantener | Fácil de escalar |
| Sin convenciones claras | Equipo alineado desde el inicio |
| Reescribir configuración cada vez | Base reutilizable |

### Capas del Scaffolding en este proyecto

```
cmd/            → Punto de entrada (main.go)
internal/
  config/       → Variables de entorno y configuración
  database/     → Conexión a la BD y migraciones automáticas
  models/       → Estructuras de datos con tags GORM y JSON
  repository/   → Capa de acceso a datos (CRUD puro sobre la BD)
  handlers/     → Lógica HTTP: recibe request, llama al repo, devuelve JSON
  middleware/   → Funciones que se ejecutan entre request y response
  routes/       → Registro y agrupación de todas las rutas
```

Cada capa tiene **una única responsabilidad**. Esto sigue el principio **SRP** (Single Responsibility Principle).

---

## Go — El Lenguaje

**Go** (también llamado Golang) es un lenguaje de programación de código abierto creado por Google en 2009. Fue diseñado para ser simple, eficiente y escalable, con soporte nativo para concurrencia.

### Características principales

| Característica | Descripción |
|---|---|
| **Compilado** | Se compila a binario nativo, sin necesidad de runtime ni intérprete |
| **Tipado estático** | Los tipos se verifican en tiempo de compilación |
| **Goroutines** | Concurrencia liviana integrada al lenguaje con `go func()` |
| **Garbage Collector** | Manejo automático de memoria |
| **Sin clases** | Usa structs e interfaces, no herencia tradicional |
| **Módulos** | Sistema de dependencias propio con `go.mod` |

### ¿Por qué Go para APIs?

Go es especialmente popular para construir APIs y microservicios porque:

- **Velocidad de compilación**: proyectos grandes compilan en segundos
- **Rendimiento cercano a C**: procesa cientos de miles de requests por segundo
- **Binario único**: el resultado final es un solo ejecutable, fácil de desplegar
- **Librería estándar potente**: incluye servidor HTTP, JSON, criptografía y más sin instalar nada extra
- **Simplicidad**: pocas palabras clave, código fácil de leer para cualquier miembro del equipo

### Conceptos clave en Go que usamos aquí

```go
// Struct — define la forma de un dato (como una clase sin métodos heredados)
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
}

// Package — cada carpeta es un paquete con responsabilidad única
package handlers

// Error handling explícito — Go no tiene excepciones, devuelves el error
func GetUser(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }
    // ...
}
```

---

## Fiber — El Framework

**Fiber** es un framework web para Go, creado en 2020, inspirado directamente en **Express.js** de Node.js. Su meta es ofrecer la misma ergonomía de Express pero con el rendimiento de Go.

### ¿Qué lo hace especial?

Fiber usa por debajo **fasthttp**, la librería HTTP más rápida para Go (hasta 10x más rápida que `net/http` estándar). Esto lo convierte en uno de los frameworks web más rápidos del mundo en benchmarks independientes.

```
net/http (estándar Go)  →  ~100,000 req/s
Fiber (fasthttp)        →  ~900,000 req/s  ✓
```

### Conceptos de Fiber

#### `fiber.App` — La aplicación

```go
app := fiber.New(fiber.Config{
    AppName: "Mi API",
})
```

Es la instancia central. Todo parte de aquí: middlewares, rutas, configuración.

#### `fiber.Ctx` — El contexto de la petición

Cada handler recibe un `*fiber.Ctx` con todo lo que necesitas:

```go
func MiHandler(c *fiber.Ctx) error {
    c.Params("id")          // parámetro de ruta /users/:id
    c.Query("page")         // query string /users?page=2
    c.BodyParser(&req)      // parsear body JSON
    c.Locals("user")        // datos pasados entre middlewares
    c.Status(201).JSON(obj) // responder con código y JSON
}
```

#### Rutas y métodos HTTP

```go
app.Get("/ruta", handler)
app.Post("/ruta", handler)
app.Put("/ruta/:id", handler)
app.Delete("/ruta/:id", handler)
app.Patch("/ruta/:id", handler)
```

#### Grupos de rutas

Los grupos permiten organizar endpoints bajo un prefijo común, lo que encaja perfectamente con el scaffolding por capas:

```go
api := app.Group("/api/v1")
users := api.Group("/users")
users.Get("/", handlers.GetUsers)    // GET /api/v1/users
users.Post("/", handlers.CreateUser) // POST /api/v1/users
```

#### Middleware en Fiber

Un middleware es una función que se ejecuta **antes** o **después** de los handlers. Sirve para logging, autenticación, CORS, rate limiting, etc.

```go
// Middleware global (aplica a todas las rutas)
app.Use(middleware.Logger)

// Middleware para un grupo específico
api.Use(middleware.AuthRequired)

// Estructura típica de un middleware
func Logger(c *fiber.Ctx) error {
    start := time.Now()
    err := c.Next()  // ← llama al siguiente handler
    log.Printf("%s %s - %v", c.Method(), c.Path(), time.Since(start))
    return err
}
```

#### `fiber.Map` — Respuestas rápidas

```go
// Equivalente a map[string]interface{} pero más corto
return c.JSON(fiber.Map{
    "status": "ok",
    "data":   users,
})
```

### Fiber vs otros frameworks Go

| Framework | Basado en | Velocidad | Similitud con |
|---|---|---|---|
| **Fiber** | fasthttp | ⚡⚡⚡ Muy alta | Express.js |
| Gin | net/http | ⚡⚡ Alta | Martini |
| Echo | net/http | ⚡⚡ Alta | Express.js |
| Chi | net/http | ⚡ Media | stdlib |

Fiber es la mejor opción si vienes de JavaScript/Node.js o necesitas máximo rendimiento.

---

## Fiber + Scaffolding: La Combinación Perfecta

Fiber y Scaffolding no solo son compatibles — se **potencian mutuamente**. Aquí está por qué esta combinación tiene tanto sentido:

### 1. Fiber define el contrato, Scaffolding define el lugar

Fiber establece **qué** existe (rutas, handlers, middleware). El Scaffolding establece **dónde** vive cada cosa. Sin estructura, una app Fiber crece caótica:

```go
// ❌ Sin scaffolding — todo en main.go (anti-patrón)
func main() {
    app := fiber.New()
    app.Get("/users", func(c *fiber.Ctx) error { /* 200 líneas aquí */ })
    app.Post("/users", func(c *fiber.Ctx) error { /* otras 150 líneas */ })
    app.Listen(":3000")
}
```

```
// ✓ Con scaffolding — cada pieza en su capa
cmd/main.go           → solo arranca la app
routes/routes.go      → solo registra rutas
handlers/user.go      → solo lógica de usuarios
models/user.go        → solo define la estructura
middleware/logger.go  → solo loggea
```

### 2. El `fiber.Ctx` fluye por todas las capas

El contexto de Fiber viaja desde el middleware hasta el handler sin esfuerzo, respetando la separación de capas:

```
Request →  Logger (middleware)
              ↓  c.Next()
           Router (routes)
              ↓  despacha a...
           GetUser (handler)
              ↓  usa...
           User (model)
              ↓
         Response JSON
```

Cada capa del scaffolding corresponde exactamente a una responsabilidad de Fiber.

### 3. Los grupos de Fiber mapean 1:1 con los módulos del scaffolding

Cuando el proyecto crece y agregas, por ejemplo, un módulo de `products`:

```
// Agregas en scaffolding:
handlers/product_handler.go   → nueva capa
models/product.go             → nuevo modelo
routes/routes.go              → nuevo grupo

// Y en Fiber:
products := api.Group("/products")
products.Get("/", handlers.GetProducts)
```

No tienes que tocar nada más. El scaffolding hace que **escalar sea predecible**.

### 4. Los middlewares de Fiber son ciudadanos de primera clase en el scaffolding

La carpeta `middleware/` del scaffolding encapsula toda la lógica transversal. Fiber los consume sin acoplamiento:

```go
// routes.go — conecta scaffolding con Fiber de forma limpia
app.Use(middleware.Logger)       // de middleware/logger.go
api.Use(middleware.RateLimit)    // de middleware/rate_limit.go
admin.Use(middleware.AuthAdmin)  // de middleware/auth.go
```

### Resumen visual

```
SCAFFOLDING define la ESTRUCTURA     FIBER provee los BLOQUES
─────────────────────────────────    ──────────────────────────
cmd/            →                    fiber.New(), app.Listen()
routes/         →                    app.Group(), app.Get/Post...
handlers/       →                    func(c *fiber.Ctx) error
middleware/     →                    app.Use(fn)
models/         →                    structs con tags `json:"..."`
config/         →                    fiber.Config{}
```

---

## Base de Datos con GORM + SQLite

El proyecto usa **GORM** como ORM y **SQLite** como motor de base de datos. Los registros se guardan en un archivo `carrycoders.db` generado automáticamente al iniciar el servidor.


### ¿Qué es GORM?

**GORM** es el ORM (*Object-Relational Mapper*) más popular de Go. Convierte los structs de Go en tablas SQL y permite hacer consultas sin escribir SQL puro:

```go
// GORM lee los tags del struct y crea la tabla automáticamente
type User struct {
    ID        uint      `json:"id"    gorm:"primaryKey;autoIncrement"`
    Name      string    `json:"name"  gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### AutoMigrate — El esquema nunca se rompe

Al arrancar el servidor, GORM ejecuta `AutoMigrate` que crea o actualiza la tabla sin perder datos:

```go
// internal/database/database.go
DB.AutoMigrate(&models.User{})
// → Si la tabla no existe: la crea
// → Si le falta una columna: la agrega
// → Si ya está al día: no hace nada
```

### Patrón Repository — Separar el acceso a datos

El scaffolding agrega la capa **Repository** entre el Handler y la BD. Esto significa que el handler nunca toca la BD directamente:

```
Request HTTP
     ↓
UserHandler          ←  solo sabe de HTTP (Fiber)
     ↓  llama a
UserRepository       ←  solo sabe de SQL (GORM)
     ↓  usa
   gorm.DB           ←  solo sabe de SQLite
     ↓  escribe en
  carrycoders.db     ←  archivo en disco
```

Beneficios de esta separación:
- Cambiar de SQLite a PostgreSQL solo requiere editar `database.go` y el driver
- Los handlers se pueden testear con un repositorio simulado (*mock*)
- La lógica de consultas está centralizada en un solo lugar

### Inyección de dependencias

La BD se inyecta desde `main.go` hacia abajo, nunca se usa una variable global en los handlers:

```go
// cmd/main.go — punto de arranque
database.Connect(cfg.DBPath)   // abre carrycoders.db
routes.Setup(app, database.DB) // pasa la conexión

// internal/routes/routes.go — cableado
userRepo    := repository.NewUserRepository(db)     // BD → Repo
userHandler := handlers.NewUserHandler(userRepo)     // Repo → Handler

// internal/handlers/user_handler.go — handler con repo inyectado
type UserHandler struct {
    repo *repository.UserRepository
}
```

### Tags de GORM en el modelo

```go
gorm:"primaryKey;autoIncrement"  // clave primaria con autoincremento
gorm:"not null"                  // columna obligatoria (NOT NULL en SQL)
gorm:"uniqueIndex;not null"      // columna única con índice (como un email)
```

### Operaciones GORM utilizadas

| Operación | Método GORM | SQL equivalente |
|---|---|---|
| Obtener todos | `db.Find(&users)` | `SELECT * FROM users` |
| Obtener uno | `db.First(&user, id)` | `SELECT * FROM users WHERE id=? LIMIT 1` |
| Crear | `db.Create(&user)` | `INSERT INTO users (...)` |
| Actualizar | `db.Save(&user)` | `UPDATE users SET ... WHERE id=?` |
| Eliminar | `db.Delete(&user, id)` | `DELETE FROM users WHERE id=?` |

---

## Instalación de Go

### Windows

1. Descarga el instalador desde [https://go.dev/dl/](https://go.dev/dl/)
2. Ejecuta el `.msi` y sigue el asistente
3. Verifica la instalación:

```bash
go version
# go version go1.22.x windows/amd64
```

### macOS

```bash
brew install go
go version
```

### Linux (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install golang-go

---

## Instalación de Fiber

[Fiber](https://gofiber.io/) es un framework web para Go inspirado en Express.js. Es extremadamente rápido gracias a que usa `fasthttp` en lugar de `net/http`.

### Inicializar un módulo Go nuevo

```bash
# Crear carpeta del proyecto
mkdir mi-proyecto
cd mi-proyecto

# Inicializar el módulo (reemplaza con tu nombre de módulo)
go mod init github.com/tu-usuario/mi-proyecto
```

### Instalar Fiber

```bash
go get github.com/gofiber/fiber/v2
```

Esto actualiza `go.mod` y genera `go.sum` automáticamente.

### Ejemplo mínimo con Fiber

```go
package main

import (
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hola desde Fiber!")
    })

    app.Listen(":3000")
}
```

```bash
go run main.go
# Abre http://localhost:3000
```

---

## Estructura del Proyecto

```
Exposición/
├── cmd/
│   └── main.go                      ← Arranca el servidor e inyecta la BD
├── internal/
│   ├── config/
│   │   └── config.go                ← Lee PORT, APP_ENV, DB_PATH
│   ├── database/
│   │   └── database.go              ← Abre SQLite y ejecuta AutoMigrate
│   ├── handlers/
│   │   ├── health_handler.go        ← GET /health — estado del servidor
│   │   └── user_handler.go          ← Handlers HTTP (usan UserRepository)
│   ├── middleware/
│   │   └── logger.go                ← Middleware que loggea cada petición
│   ├── models/
│   │   └── user.go                  ← Struct User con tags GORM y JSON
│   ├── repository/
│   │   └── user_repository.go       ← CRUD directo sobre la BD con GORM
│   └── routes/
│       └── routes.go                ← Registra endpoints y cablea dependencias
├── carrycoders.db                   ← Archivo SQLite (generado al correr)
├── go.mod                           ← Módulo y dependencias
├── go.sum                           ← Hash de dependencias (auto-generado)
├── .gitignore
└── README.md
```

### Flujo de una petición HTTP

```
Cliente HTTP
    │
    ▼
fiber.App      (cmd/main.go)
    │
    ▼
Middleware     (middleware/logger.go)        ← loggea método, ruta y latencia
    │
    ▼
Router         (routes/routes.go)            ← decide qué handler usar
    │
    ▼
UserHandler    (handlers/user_handler.go)    ← parsea request, llama al repo
    │
    ▼
UserRepository (repository/user_repository.go) ← ejecuta la consulta GORM
    │
    ▼
gorm.DB        (database/database.go)        ← SQL sobre SQLite
    │
    ▼
carrycoders.db                               ← archivo en disco
    │
    ▼
Respuesta JSON al cliente
```

---

## Cómo correr el proyecto

### 1. Clonar e instalar dependencias

```bash
cd Exposición

# Descargar todas las dependencias declaradas en go.mod
go mod tidy
```

### 2. Correr en modo desarrollo

```bash
go run cmd/main.go
```

### 3. Compilar a binario

```bash
# Compilar
go build -o bin/server cmd/main.go

# Ejecutar el binario
./bin/server          # Linux/macOS
bin\server.exe        # Windows
```

### 4. Variables de entorno opcionales

```bash
# Cambiar puerto (por defecto :3000)
PORT=:8080 go run cmd/main.go

# Cambiar entorno
APP_ENV=production go run cmd/main.go

# Cambiar la ruta del archivo SQLite (por defecto carrycoders.db)
DB_PATH=./data/mydb.db go run cmd/main.go
```

---

## Endpoints disponibles

### Health Check

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/health` | Estado del servidor |

```bash
curl http://localhost:3000/health
```

```json
{
  "message": "Servidor funcionando correctamente",
  "status": "ok"
}
```

---

### Usuarios (`/api/v1/users`)

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/v1/users` | Obtener todos los usuarios |
| `GET` | `/api/v1/users/:id` | Obtener usuario por ID |
| `POST` | `/api/v1/users` | Crear nuevo usuario |
| `PUT` | `/api/v1/users/:id` | Actualizar usuario por ID |
| `DELETE` | `/api/v1/users/:id` | Eliminar usuario por ID |

#### Obtener todos los usuarios

```bash
curl http://localhost:3000/api/v1/users
```

```json
{
  "data": [
    { "id": 1, "name": "Alice García", "email": "alice@example.com", "age": 25 },
    { "id": 2, "name": "Bob Martínez", "email": "bob@example.com", "age": 30 }
  ],
  "total": 2
}
```

#### Crear un usuario

```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Carlos", "email": "carlos@example.com", "age": 22}'
```

```json
{
  "id": 3,
  "name": "Carlos",
  "email": "carlos@example.com",
  "age": 22,
  "created_at": "2026-05-07T10:00:00Z",
  "updated_at": "2026-05-07T10:00:00Z"
}
```

#### Actualizar un usuario

```bash
curl -X PUT http://localhost:3000/api/v1/users/3 \
  -H "Content-Type: application/json" \
  -d '{"name": "Carlos Updated", "email": "carlos.new@example.com", "age": 23}'
```

```json
{
  "id": 3,
  "name": "Carlos Updated",
  "email": "carlos.new@example.com",
  "age": 23,
  "created_at": "2026-05-07T10:00:00Z",
  "updated_at": "2026-05-07T10:05:00Z"
}
```

#### Eliminar un usuario

```bash
curl -X DELETE http://localhost:3000/api/v1/users/3
# 204 No Content
```

---

## Tecnologías usadas

| Tecnología | Versión | Rol |
|---|---|---|
| [Go](https://go.dev/) | 1.22+ | Lenguaje compilado, tipado estáticamente, concurrente |
| [Fiber v2](https://gofiber.io/) | v2.52+ | Framework HTTP inspirado en Express.js |
| [GORM](https://gorm.io/) | v1.31+ | ORM para Go, maneja migraciones y consultas |
| [SQLite (pure Go)](https://github.com/glebarez/sqlite) | v1.11+ | Base de datos embebida sin CGO ni servidor externo |

---

*CarryCoders — Exposición de Ingeniería Web*
