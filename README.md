# Desafios FullCycle

---

## Desafio Multithreading

Aplicação em Go que busca informações de CEP em duas APIs simultaneamente (ViaCEP e BrasilAPI) e retorna o resultado da mais rápida, descartando a mais lenta.

### Pré-requisitos

- Go 1.21+

### Como rodar

```bash
go run main.go CEP
```

Exemplos:
```bash
go run main.go 01153-000
go run main.go 01153000
```
