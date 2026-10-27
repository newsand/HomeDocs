# Fase 2

## Context
Fase 2. Feature `docs/03-features/feature-cadastrar-pessoa.md`. Contratos HTTP em `docs/04-contracts.md`. A fase 1 já cobre SQLite e arquivo.

Red: `internal/http/server_test.go`.
Comando: `go test ./internal/http/`
Resultado: falha de compilação. `New` ainda não existe.

## Red command
`go test ./internal/http/`

## Green steps
1. `New` sobe o Gin com templates pongo2 em `web/templates`.
2. `GET /` redireciona para `/people` com 303.
3. `POST /people` cria e redireciona com 303 para a página. Nome em branco devolve HTML 422 com "Informe o nome." e não grava.
4. `GET /people/:id` mostra todos os campos, com botão de copiar. Nome com `<script>` sai escapado. Id ausente é HTML 404 "Pessoa não encontrada." sem texto de SQL.
5. `POST /people/:id` substitui os campos. `POST /people/:id/delete` apaga e volta para a lista.
6. `GET /people` lista e filtra com `q`. `%` é texto.
7. Foto de rosto e scans entram por multipart e saem em `GET /files/:name`, sem sair do diretório.

Fora desta fase: login, JSON, extração de PDF.

## Verify
`go test ./internal/http/`
Feito quando esse comando passa.
