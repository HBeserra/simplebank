# Estrutura

As migrações em `db/migrations` criam as tabelas do banco de dados,
o pacote `SQLC` utiliza essas migrações para obter os dados sobre o banco,
com base nas Querys `db/query` ele gera pacotes `db/sqlc` com base nas configurações `sqlc.yaml`.

> Não editar o código gerado pelo `SQLC`

Ao alterar o banco de dados ou querys execute o comando `make sqlc`

# Migrações

**Pacote:**  [Golang migrate](https://github.com/golang-migrate/migrate/)

Comando para criar uma migração: `migrate create -ext sql -dir db/migrate -seq init_schema`
> Cria uma migração com nome `init_schema`,
`-seq`: Numeração sequencial das migrações,
`-dir`: Diretório para armazenar a migração,
`-ext`: Extensão do arquivo

# Gerador de Querys

**Pacote:**  [SQLC](https://github.com/kyleconroy/sqlc)

Gera automaticamente os métodos de acesso ao banco de dados baseados nas querys definidas no diretorio `db/query`.
Código performático, o SQLC utiliza o pacote sql padrão do go.