# Instrucciones
¡Bienvenido/a a este coding challenge de Go! Recomendamos leer este archivo completo antes de empezar.

En el proyecto proporcionado encontrarás una base de código con algunos contratos definidos y funcionalidades por implementar. A continuación, te detallamos los objetivos a cumplir.

1. Tu primera tarea es crear e integrar una implementación de GetBooksProvider para obtener información de libros desde un servicio externo. Esto incluye realizar la solicitud HTTP, procesar la respuesta y asegurarte de manejar los posibles errores. Además, deberás integrar esta implementación en el flujo principal del programa. La información se encuentra en https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books

2. El siguiente paso consiste en separar la lógica de negocio y la de presentación, siguiendo los principios de separación de capas. Actualmente toda esta lógica se encuentra "mezclada" en una sola capa, en el archivo handlers/handlers.go.

3. Queremos garantizar un uso correcto del contexto en el proyecto. Revisa y ajusta las funciones para asegurarte de que el contexto sea utilizado únicamente donde sea necesario.

4. Es importante que el código esté bien cubierto por tests. Actualiza las pruebas existentes para reflejar los cambios realizados en la lógica de negocio y añade nuevos casos que validen el manejo de errores. Podés usar el comando `go test ./handlers` en la shell.

## Instalación y Ejecución

### Prerrequisitos
- Go 1.21 o superior
- Git

### Pasos para ejecutar el proyecto

1. **Clonar el repositorio** (si aún no lo has hecho)
   ```bash
   git clone <repository-url>
   cd resolution-code-challenge-go
   ```

2. **Instalar las dependencias**
   ```bash
   go mod download
   ```

3. **Configurar las variables de entorno**
   
   Crear un archivo `.env` basándote en el archivo `.env-example`:
   ```bash
   cp .env-example .env
   ```
   
   Editar el archivo `.env` y agregar la URL de la API de libros:
   ```
   BOOKS_API_URL=
   ```

4. **Ejecutar el proyecto**
   ```bash
   go run main.go
   ```

5. **Acceder a la aplicación**
   
   Una vez que el servidor esté ejecutándose, podrás acceder a:
   
   - **API Endpoints:**
     - `GET http://localhost:3000/books` - Obtener todos los libros
     - `GET http://localhost:3000/books/metrics?author=<nombre>` - Obtener métricas de libros
   
   - **Documentación Swagger:**
     - `http://localhost:3000/swagger/index.html` - Interfaz interactiva de la API
   
   Desde la documentación de Swagger podrás probar todos los endpoints de la API de forma interactiva.

### Ejecutar Tests

Para ejecutar todos los tests del proyecto:
```bash
go test ./...
```

Para ejecutar tests de un paquete específico:
```bash
go test ./handlers
go test ./providers
go test ./repositories
```
