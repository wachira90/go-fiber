# golang fiber


## run 

```
go run main.go
```
## test data

### post add data

```
POST http://localhost:3000/books
Content-Type: application/json

{
  "title": "The Great Gatsby",
  "author": "F. Scott Fitzgerald",
  "rating": 5
}
```
### result when add 

```
HTTP/1.1 200 OK
Date: Tue, 23 Apr 2024 03:12:56 GMT
Content-Type: application/json
Content-Length: 190
Connection: close

[
  {
    "ID": 1,
    "CreatedAt": "2024-04-23T10:12:53.284358+07:00",
    "UpdatedAt": "2024-04-23T10:12:53.284358+07:00",
    "DeletedAt": null,
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "rating": 5
  }
]
```

### get list data

```
HTTP/1.1 200 OK
Date: Tue, 23 Apr 2024 03:17:22 GMT
Content-Type: application/json
Content-Length: 190
Connection: close

[
  {
    "ID": 1,
    "CreatedAt": "2024-04-23T10:12:53.284358+07:00",
    "UpdatedAt": "2024-04-23T10:12:53.284358+07:00",
    "DeletedAt": null,
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "rating": 5
  }
]
```

## build 

```
go build -o main.exe main.go

```