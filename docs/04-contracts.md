HTML renderizado no servidor (Gin + pongo2). Sem JSON de produto. Forms: `name`, campos de `02-data-models.md`, 1 foto de rosto e um campo de arquivo (vários scans).

POST bem-sucedido responde 303. O contrato anterior dizia 200 e redirect; navegador não redireciona em 200.

GET / → redirect 303 para GET /people

GET /files/:name → arquivo de foto ou scan guardado no diretório de arquivos. Nome fora desse diretório: 404.

POST /people → create_person
Request: form HTML.
Response 303: redirect para GET /people/{id}.
Response 422: HTML. Nome ausente ou em branco. Texto: "Informe o nome." Nada é gravado.
Idempotency: não é idempotente. Reenvio cria outra pessoa.

GET /people/{id} → show_person
Request: id na rota.
Response 200: HTML com todos os dados da pessoa, a foto de rosto e os scans. Cada campo em display tem botão de copiar para o clipboard (browser).
Response 404: HTML "Pessoa não encontrada." Pessoa inexistente. Sem texto de SQL.

POST /people/{id} → update_person
Request: form HTML.
Response 303: redirect para GET /people/{id}.
Response 404: pessoa inexistente. Nada muda.
Response 422: nome ausente ou em branco. A linha anterior permanece.
Idempotency: last write wins nos campos de texto. Arquivo vazio mantém o atual. Scans enviados substituem o conjunto. Nova foto de rosto substitui a anterior.

POST /people/{id}/delete → delete_person
Request: id na rota. Sem corpo.
Response 303: redirect para GET /people.
Response 404: HTML. Pessoa inexistente.
Idempotency: segundo apagar da mesma pessoa devolve 404. A linha some na transação; os arquivos de foto saem em seguida.

GET /people → search_people
Request: query string opcional `q`.
Response 200: HTML da lista. Sem `q`, todas as pessoas, por nome e depois por id. Com `q`, substring sem diferenciar maiúsculas no nome e nos campos de documento. `%` é texto. Lista vazia é a página com "Nenhuma pessoa.", não erro.
