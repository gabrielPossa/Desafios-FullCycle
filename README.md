# Desafio Clean Architecture

Projeto que implementa um sistema de orders utilizando Clean Architecture em Go, expondo os serviços via REST, gRPC e GraphQL.

## Como executar

1. Clone o repositório e acesse a branch `cleanArch`:
   ```bash
   git checkout cleanArch
   ```

2. Suba a infraestrutura (MySQL + RabbitMQ) com Docker Compose:
   ```bash
   docker compose up -d
   ```

   Isso irá subir:
   - **MySQL** com o banco `orders` e a tabela criada automaticamente via migration (`mysql/startup.sql`)
   - **RabbitMQ** como message broker

3. Aguarde os containers ficarem saudáveis e execute a aplicação:
   ```bash
   cd cmd/ordersystem
   go run main.go wire_gen.go
   ```

## Portas dos serviços

| Serviço  | Porta  | Descrição                          |
|----------|--------|------------------------------------|
| REST     | `8000` | API HTTP (GET e POST em `/order`)  |
| gRPC     | `50051`| Serviço `OrderService`             |
| GraphQL  | `8080` | Playground em `/`, queries em `/query` |
| MySQL    | `3306` | Banco de dados                     |
| RabbitMQ | `5672` | Message broker (management: `15672`) |

## Endpoints

### REST

- **Criar order:** `POST http://localhost:8000/order`
  ```json
  { "id": "abc", "price": 100.5, "tax": 10.05 }
  ```
- **Listar orders:** `GET http://localhost:8000/order`

### gRPC

- **CreateOrder** e **ListOrders** no serviço `pb.OrderService` (porta `50051`)

### GraphQL

- **Mutation - Criar order:**
  ```graphql
  mutation {
    createOrder(input: { id: "abc", Price: 12.34, Tax: 1.23 }) {
      id Price Tax FinalPrice
    }
  }
  ```
- **Query - Listar orders:**
  ```graphql
  query {
    orders { id Price Tax FinalPrice }
  }
  ```

Os arquivos `api/create_order.http` e `api/list_orders.http` contêm exemplos prontos para uso.
