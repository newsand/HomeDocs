Schema-first: tipo e validação no mesmo lugar. Texto é guardado como foi digitado. Data não é interpretada. Vazio significa string vazia. O único obrigatório é `name`.

O mesmo dado da pessoa aparece uma vez. CPF só no bloco CPF. Nome só em `Person.name`. Nascimento só no CPF. Filiação, sexo e naturalidade só no RG. Nacionalidade só no passaporte. Cada documento guarda o que é dele: número, datas, órgão.

## Person
| field | type | required | rule |
|---|---|---|---|
| id | string | yes | gerado no `create_person`; único |
| name | string | yes | único obrigatório; trim; só espaço conta como vazio e rejeita o comando; nome repetido é permitido |
| profile_photo | file path | no | 1 foto de rosto, padrão de perfil. Ausente → vazio. Disco + caminho no SQLite. |
| document_scans | file paths | no | Um campo de arquivo, vários scans. Ausente → lista vazia. |

## CNH
| field | type | required | rule |
|---|---|---|---|
| register_number | string | no | nº de registro; repetido entre pessoas é permitido |
| category | string | no | categoria |
| expires_on | string | no | validade |
| first_license_on | string | no | 1ª habilitação |
| issued_on | string | no | emissão |
| notes | string | no | observações |

## CPF
| field | type | required | rule |
|---|---|---|---|
| number | string | no | único CPF da pessoa; repetido entre pessoas é permitido; sem validação de dígito |
| birth_date | string | no | único nascimento |
| status | string | no | situação cadastral |
| registered_on | string | no | data da inscrição |

## RG
| field | type | required | rule |
|---|---|---|---|
| number | string | no | registro geral; repetido entre pessoas é permitido |
| social_name | string | no | nome social |
| birth_place | string | no | única naturalidade |
| father_name | string | no | filiação |
| mother_name | string | no | filiação |
| issued_on | string | no | expedição |
| issuer | string | no | órgão expedidor |
| issuer_state | string | no | UF |
| expires_on | string | no | validade |
| sex | string | no | único sexo |

## VoterID
| field | type | required | rule |
|---|---|---|---|
| number | string | no | inscrição; repetido entre pessoas é permitido |
| zone | string | no | zona |
| section | string | no | seção |
| city | string | no | município |
| state | string | no | UF |
| issued_on | string | no | emissão |

## Passport
| field | type | required | rule |
|---|---|---|---|
| kind | string | no | tipo |
| issuing_country | string | no | país emissor |
| number | string | no | número; repetido entre pessoas é permitido |
| nationality | string | no | única nacionalidade |
| issued_on | string | no | |
| authority | string | no | autoridade |
| expires_on | string | no | validade |

## Edges
- Duplicata de pessoa: permitida. A chave única é `id`.
- Campo ausente no formulário: gravado vazio. Não é erro, exceto `name`.
- Foto de rosto ausente: vazio.
- Scan ausente: lista vazia.
- Na edição, arquivo não enviado mantém o arquivo atual. Scan enviado substitui o conjunto inteiro. Foto de rosto enviada substitui a anterior.
