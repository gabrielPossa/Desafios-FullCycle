# Desafios FullCycle

## Desafio Client-Server-API

Sistema em Go que consulta a cotação do Dólar (USD-BRL) e persiste no banco de dados SQLite e responde à cotação bid ao cliente. O cliente por sua vez recebe o valor e salva no arquivo cotacao.txt.

### Pré-requisitos

- Go 1.21+
- GCC (necessário para o driver SQLite)

### Como rodar

1. Instale as dependências:
```bash
go mod tidy
```

2. Inicie o servidor:
```bash
go run server.go
```

3. Em outro terminal, execute o client:
```bash
go run client.go
```
