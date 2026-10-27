# Feature: Cadastrar pessoa

## 4 answers (user)
1. How do you describe the feature?
> quero os dados todos de uma CNH CPF RG titulo de eleitor Passaporte espaco todos os dados crud completo e ainda um file field para guarda as fotos de scan do documento + 1 foto de rosto padrao de perfil
2. What problem does it objectively solve?
> CRUD
3. What is the expected solution? What trade-offs does it involve?
> CRUD
4. What example or context do we have of the problem and the solution?
> CRUD

## Depends on
- Nenhuma.

## Operations
| name | kind | does |
|---|---|---|
| create_person | command | grava a pessoa, os campos dos documentos e os arquivos |
| update_person | command | substitui os dados da pessoa |
| delete_person | command | apaga a pessoa e os arquivos |
| show_person | query | página com todos os dados |
| search_people | query | lista as pessoas, com busca opcional |

## Description / Problem / Solution and trade-offs
CRUD da pessoa. Cada documento guarda os dados que são dele. CPF, nome, nascimento, filiação, sexo, naturalidade e nacionalidade não se repetem. Um campo de arquivo para os scans dos documentos. Uma foto de rosto, padrão de perfil.

O usuário completou: "sim pode ser tudo vazio.. menos NOME dai NOME é obrigatorio".

Trade-off dito pelo usuário: CRUD. Nome repetido é permitido. Reenviar o cadastro cria outra pessoa. Documento repetido é permitido. Sem extração automática de arquivo.

## Flow (given/when/then) per operation
- Given nenhum nome (vazio ou só espaço), When `create_person` ou `update_person`, Then nenhuma linha muda e a resposta é HTML 422.
- Given um nome preenchido e o resto vazio, When `create_person`, Then a pessoa existe e o redirect vai para `show_person`.
- Given uma pessoa, When `show_person`, Then a página mostra todos os campos gravados, a foto de rosto e os scans. Cada campo em display tem botão de copiar para o clipboard.
- Given uma pessoa, When `update_person` com nome, Then os campos de texto enviados substituem os anteriores. Campo de arquivo vazio mantém os arquivos atuais. Arquivos de scan enviados substituem o conjunto. Nova foto de rosto substitui a anterior.
- Given uma pessoa, When `delete_person`, Then a linha some e os arquivos saem em seguida. O redirect vai para a lista.
- Given a lista, When `search_people` sem texto, Then entram todas as pessoas, ordenadas por nome e depois por id. When o texto vem preenchido, Then ficam as pessoas em que o termo aparece, sem diferenciar maiúsculas, no nome ou em qualquer campo de documento.

## Error cases
- Nome ausente ou em branco → 422, estado intacto.
- `show_person`, `update_person` ou `delete_person` com id inexistente → 404.
- Segundo `delete_person` do mesmo id → 404.
- Falha ao gravar a linha depois de escrever arquivo → apaga o arquivo órfão, a pessoa não fica pela metade.
- Falha no meio do `update_person` → a linha anterior permanece.

## Acceptance
- [x] Unit: nome em branco é inválido; nome com texto é válido; os outros campos podem ir vazios.
- [x] Integration: `create_person` só com nome persiste no SQLite real; `update_person` substitui; `delete_person` remove a linha; `search_people` lista e filtra.
- [x] Integration: scans e foto de rosto ficam no disco e os caminhos ficam na linha; falha da linha não deixa arquivo órfão.
- [x] Failure does not corrupt state.
- E2E: cadastrar só com nome, abrir a página, editar, apagar. Coberto por `internal/http/server_test.go` via HTTP.

## Example / context
CRUD

## Patterns (only if cheap to change)
- Transaction Script. SQL na operação. Arquivo fora da transação.
