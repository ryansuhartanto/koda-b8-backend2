# GoREST

Go and React exercise implementing a full stack REST server.

![Screenshot 1](docs/1.png)
![Screenshot 2](docs/2.png)

## ERD

```mermaid
---
title: GoREST
---
erDiagram

users ||--o| profiles : "has"

users {
    bigint id PK

    timestamptz created_at
    timestamptz updated_at

    string email    UK
    string password
}

profiles {
    bigint id PK,FK

    timestamptz created_at
    timestamptz updated_at

    string  name
    string? picture_url
}
```

## License

[MIT](LICENSE)
