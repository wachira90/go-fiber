# golang fiber


## run 

```
go run main.go
```
## Test data

- Create a new book:
  ```
  POST http://localhost:3000/books
  {
    "title": "The Catcher in the Rye",
    "author": "J.D. Salinger",
    "rating": 4
  }
  ```

- Retrieve a book by ID:
  ```
  GET http://localhost:3000/books/1
  ```

- Update a book:
  ```
  PUT http://localhost:3000/books/1
  {
    "title": "The Catcher in the Rye (Updated)",
    "author": "J.D. Salinger",
    "rating": 5
  }
  ```

- Delete a book:
  ```
  DELETE http://localhost:3000/books/1
  ```


## build 

```
go build -o main.exe main.go

```