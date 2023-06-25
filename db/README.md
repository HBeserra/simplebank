# Migrações

**Pacote:**  [Golang migrate](https://github.com/golang-migrate/migrate/)

Comando para criar uma migração: `migrate create -ext sql -dir db/migrate -seq init_schema`
> Cria uma migração com nome `init_schema`,
`-seq`: Numeração sequencial das migrações,
`-dir`: Diretorio para armazenar a migração,
`-ext`: Extenção do arquivo

