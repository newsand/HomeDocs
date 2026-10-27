# Fase 1

## Context
Fase 1. Feature `docs/03-features/feature-cadastrar-pessoa.md`. Operações `create_person`, `update_person`, `delete_person`, `show_person`, `search_people`. Contrato em `docs/04-contracts.md`. Modelo em `docs/02-data-models.md`.

Red: `internal/person/person_test.go`.
Comando: `go test ./internal/person/`
Resultado: falha de compilação. Símbolos ainda não existem (`ValidName`, `CreatePerson`, `Input`, e os demais).

## Red command
`go test ./internal/person/`

## Green steps
1. Tipos `Person`, `Input`, documentos, `Upload`, `ErrNameRequired`, `ErrNotFound`, `ValidName`.
2. `Open` + `Migrate` em SQLite real (`modernc.org/sqlite`). Tabela `people` com os campos de `02-data-models.md`. Tabela `person_scan`.
3. `CreatePerson`: nome em branco não grava; nome com espaço é trim; o resto pode ir vazio; nome repetido gera outro id; texto é dado, não SQL.
4. Arquivo: nome gerado dentro do diretório; `../` não escapa. Falha do insert apaga o arquivo novo.
5. `UpdatePerson`: substitui texto; nome em branco não muda a linha nem grava arquivo; upload vazio mantém arquivos; upload novo substitui e apaga os antigos; falha no UPDATE mantém a linha e não deixa arquivo novo; id ausente é `ErrNotFound`.
6. `DeletePerson`: apaga linha e arquivos; segundo delete e id no formato de injeção devolvem `ErrNotFound` sem apagar as outras linhas.
7. `ShowPerson`: devolve a pessoa ou `ErrNotFound`.
8. `SearchPeople`: vazio lista todos por nome e depois id; substring sem caixa no nome e nos campos de documento; `%` é texto.

Fora desta fase: HTTP, Gin, pongo2, botão de copiar.

## Verify
`go test ./internal/person/`
Feito quando esse comando passa.
