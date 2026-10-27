## Cost and abstraction
- Prefer the cheap list; forbid the expensive list by default (`default-architecture` skill).
- No generic repository, no empty layers, no `shared`/`helpers` dump.
- Go: convenção em `.cursor/rules/go-convention.mdc` (MixedCaps, parâmetro `p_`, sem `init`, sem prefixo Get/Set, erro sempre checado, `panic` só em falha fatal de startup).

## Quality gate
- Lefthook orchestrates pre-commit/pre-push
- Go: golangci-lint (govet, staticcheck, errcheck, unused, ineffassign, gosec). cyclop é só relatório e não falha o gate.
- Secrets: Gitleaks.
- Coverage: `go test -cover`. Baseline em `quality-baseline.json` quando o gate existir.
- Push/CI fails if the gate is red
- CLIs do gate não entram no `go.mod`.

## Tests
- TDD phases: Red → `docs/tdd/fase{N}.md` + `fase{N}Task.md` → Green → Refactor → Verify
- Acceptance lives in `/docs/03-features` and `/docs/04-contracts` — phases do not invent behavior
- Pure → unit. Real query → DB integration (SQLite real). Critical command → transactional.
- Few E2E. Existing tests immutable (persisted-tester).

## Data
- schema-first; command with transaction; query with no side effects
- Arquivo de foto fora da transação longa
- Nenhuma dependência nova fora de `00-brief.md` (Gin, pongo2, driver SQLite)

## Security
- não há senha nem token (sem login). Não logar campos da certidão, fotos ou PII.
- no new dependency unless listed in 00-brief.md
