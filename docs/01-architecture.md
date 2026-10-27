## Unit
Named operation (command / query / job). Not a layer, not an entity.

Operações nomeadas a partir do brief (adicionar, editar, buscar, página da pessoa). Contratos finos estão em `04-contracts.md`. Arquivos de feature ainda não existem.

| name | kind | does |
|---|---|---|
| create_person | command | cadastra uma pessoa |
| update_person | command | edita uma pessoa |
| delete_person | command | apaga uma pessoa |
| show_person | query | página com todos os dados da pessoa |
| search_people | query | busca por pessoa ou por campo |

Copiar para o clipboard é comportamento do browser na página de display, não é operação de servidor.

## Prefer (cheap to change)
Modular monolith. Transaction Script or short pipeline. Plain functions. Schema-first.
Local query object / SQL. Facade only if it hides real complexity.
Active Record or table-driven when data and behavior stay together.

## Avoid by default (expensive)
Microservices, event sourcing, full CQRS, heavy Clean/Hexagonal, DI container,
generic repository, duplicated DTO/Entity/ViewModel, needless factory, MVC as skeleton.

## Before abstracting
- real duplication or only elegance?
- shorter or longer main path?
- better locality or spread behavior?
- less risk or more coupling?
- stable enough to become structure?

## Folder tree
```
cmd/idcenter/main.go
internal/person/create_person.go
internal/person/update_person.go
internal/person/delete_person.go
internal/person/show_person.go
internal/person/search_people.go
internal/http/          # Gin: schema HTTP, render pongo2
web/templates/          # templates pongo2
```
SQLite ao lado da operação que lê ou grava. Sem `domain/`, `infra/`, `usecase/`.

## One operation (example)
`POST` HTML form → schema → `create_person` → grava na transação SQLite → redirect para `show_person`.

Command: transação abre e fecha dentro de `create_person`, `update_person` e `delete_person`. Query: sem efeito colateral. Não há I/O externo além do SQLite e do disco das fotos. Arquivo (scans e foto de rosto) fica fora da transação longa: grava o arquivo, depois a linha; se a linha falhar, apaga o arquivo órfão. `delete_person` apaga a linha e, em seguida, os arquivos.

## Boundaries
- HTTP entra por schema explícito (form HTML).
- command/query devolvem resultado explícito.
- erro interno não vaza cru para o HTML.
- efeito colateral aparece na operação (SQL e arquivo de foto).

## Consistency
- `create_person`, `update_person` e `delete_person`: uma transação SQLite por comando. Atômico = a linha da pessoa (e os campos dela) naquela transação. Arquivos de foto não entram na transação.
- Idempotência: `create_person` não é idempotente; reenvio cria outra pessoa. `delete_person` repetido devolve 404.
- Timeout/cancel: o request HTTP delimita o comando. Sem retry automático.
- Concorrência: duas edições da mesma pessoa — last write wins.

## Observability
correlation id por request; log estruturado por operação; duração e ponto de falha. Nunca logar campos de documento, fotos, nem PII.

## Evolution (when abstraction is no longer premature)
recurring latency, hotspot, concurrency, change conflicts, operational isolation,
unstable external dependency, domain larger than acceptable locality
