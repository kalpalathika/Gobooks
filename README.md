# GoBooks
A lightweight Go REST API for managing books with full CRUD operations. Built using Gorilla Mux, JSON helpers, and clean modular handlers.

## Installation & Setup

**Prerequisites:** Go 1.16+

```bash
git clone <repository-url>
cd GoBooks
go mod download
go run .
```

Server runs on `http://localhost:8080`

## API Endpoints

- `GET /books` - Get all books
- `GET /books/{id}` - Get book by ID
- `POST /book` - Create a new book
- `PUT /book/{id}` - Update a book
- `DELETE /book/{id}` - Delete a book

## Usage Examples

**Create a book:**
```bash
curl -X POST http://localhost:8080/book \
  -H "Content-Type: application/json" \
  -d '{"title": "New Book", "author": "Author Name", "year": 2024}'
```

**Get all books:**
```bash
curl http://localhost:8080/books
```

## Data Model

```go
type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}
```

## Running Tests

```bash
go test -v
```

## Project Structure

- `main.go` - Server setup and routes
- `book.go` - Book model, handlers, and helpers
- `main_test.go` - Test suite

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
