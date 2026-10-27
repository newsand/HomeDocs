## [2026-10-26] Sync

### Built
- Campos repetidos saíram do cadastro. CPF fica só no bloco CPF.

### Decisions outside the spec
- Nome, nascimento, filiação, sexo, naturalidade e nacionalidade também ficam uma vez.

### Spec deviations
- None

### Known follow-ups
- Banco já criado conserva as colunas antigas. O formulário não lê mais essas cópias.

## [2026-10-26] Build

### Built
- `create_person`, `update_person`, `delete_person`, `show_person`, `search_people`
- Páginas HTML Gin + pongo2, foto de rosto e scans

### Decisions outside the spec
- POST de sucesso responde 303, não 200. Redirect de navegador não acontece em 200. Contratos atualizados.
- `GET /` redireciona para `/people`.
- `GET /files/:name` entrega o arquivo guardado.

### Spec deviations
- None

### Known follow-ups
- Copy para o clipboard está na página, sem feature própria.
- Quality gate (Lefthook) ainda não instalado neste repositório.

## [2026-10-26] Build / Sync

### Built
- Nenhum. DOC apenas: brief, arquitetura, modelo, contratos, non-negotiables, discretion.

### Decisions outside the spec
- Nenhuma.

### Spec deviations
- None

### Known follow-ups
- Feature escrita: `03-features/feature-cadastrar-pessoa.md`. Respostas 2, 3 e 4 verbatim: "CRUD".
- Nome é o único obrigatório. O resto pode ir vazio.
- Campos impressos de CNH, CPF, RG, título e passaporte fechados em `02-data-models.md`.
- Copy para o clipboard segue no contrato de `show_person`, sem feature própria.
