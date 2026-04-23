# Listagem de Orders (Clean Architecture)

Projeto alinhado ao módulo [20-CleanArch do curso Go Expert](https://github.com/devfullcycle/goexpert/tree/main/20-CleanArch): entidades em `internal/entity`, casos de uso em `internal/usecase`, adaptadores em `internal/infra` (banco, web com **Chi**, GraphQL com **gqlgen**, gRPC com serviço em `internal/infra/grpc/service`), configuração com **Viper** em `configs` e injeção de dependência com **Wire** em `cmd/ordersystem` (`wire.go` + `wire_gen.go`).

## Execução (comando único)

```bash
docker compose up
```

Sobem o PostgreSQL, as migrações, e a aplicação (REST, gRPC, GraphQL).

## Portas

| Serviço   | Porta    | Detalhe |
|-----------|----------|---------|
| REST      | **8080** | Chi — `GET /order`, `POST /order` |
| gRPC      | **50051** | `order.v1.OrderService` / `ListOrders` (reflection) |
| GraphQL   | **8081** | Playground em `/`, API em `/query` (gqlgen) |
| PostgreSQL| **5432** | (mapeada para inspeção local) |

## Variáveis de ambiente

| Variável | Exemplo | Descrição |
|----------|---------|-----------|
| `DATABASE_URL` | `postgres://…` | DSN do Postgres |
| `MIGRATIONS_DIR` | `/app/migrations` | Diretório das migrações SQL |
| `WEB_SERVER_PORT` | `:8080` | Porta HTTP (inclui `:`) |
| `GRPC_SERVER_PORT` | `50051` | Porta gRPC (sem `:`) |
| `GRAPHQL_SERVER_PORT` | `8081` | Porta GraphQL (sem `:`) — playground + `/query` |

O carregamento segue o padrão do curso (`configs` + Viper), com leitura opcional de `configs/.env` e **sobrescrita por variáveis de ambiente** (recomendado no Docker).

## Código gRPC e GraphQL

- Protobuf: `api/proto/order/v1/order.proto`  
- Código Go gerado: `internal/infra/grpc/pb/…`  
- GraphQL: `gqlgen.yml` + `internal/infra/graph/*.graphqls`  
- Regenerar (após alterar schema):  
  - `go run github.com/99designs/gqlgen generate`  
- Regenerar `wire_gen.go` após editar `wire.go`:  
  - `go run github.com/google/wire/cmd/wire ./cmd/ordersystem`

## Exemplos

Ver `api.http` na raiz.

## Arquitetura (resumo)

- **Uso único de regra de negócio**: `ListOrdersUseCase` (e `CreateOrderUseCase` só para popular / REST / mutation GraphQL).  
- **Persistência**: `*sql.DB` (driver `pgx`) + `internal/infra/database`.  
- **Não** inclui RabbitMQ ou MySQL do repositório de referência: o desafio pede Postgres e listagem; o *estilo* de pastas e bibliotecas segue o curso.
