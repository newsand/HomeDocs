# Phase 1 — Green

- [x] Tipos, `ValidName`, erros
- [x] `Open` e `Migrate`
- [x] `CreatePerson` com nome, documentos e arquivo contido no diretório
- [x] Falha do insert apaga arquivo órfão
- [x] `UpdatePerson` substitui texto, mantém ou troca arquivo, nome em branco não grava, falha de UPDATE não corrompe
- [x] `DeletePerson` remove linha e arquivo; segundo delete é `ErrNotFound`
- [x] `ShowPerson` e `SearchPeople` (ordem, caixa, `%` literal)
- [x] Phase tests pass: `go test ./internal/person/`
