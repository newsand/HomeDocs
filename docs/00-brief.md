# Project: ID-Center

## Goal in one paragraph
ID Center é uma central para adicionar, editar e buscar informações de pessoas. Não tem login. Não tem nada disso. São campos de dados, fotos e pronto.

Feature fechada: `03-features/feature-cadastrar-pessoa.md` (CRUD, documentos, scans, foto de rosto). O único campo obrigatório é o nome.

Ainda sem arquivo de feature:

- botão de copy para o clipboard em todos os campos em display (já está no contrato de `show_person`)

## Mandatory stack
- Language/framework: Go + Gin
- Database: SQLite
- Libraries already decided (do not pick others): pongo2 (templates no estilo Jinja2). Páginas HTML renderizadas no servidor. Sem rotas JSON como superfície do produto.

## Architecture
Modular monolith + named operations (command/query/job). Not MVC / Clean / Hexagonal / DDD / CQRS by default.

## Non-goals
- Não tem login nem autenticação.
- Não expõe API JSON como produto.
- Não é SPA.
