## Bank service build in golang

This is the bank service that will provide APIs for the Frontend to do the following things:

1. Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
2. Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
3. Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both accounts’ balance are updated successfully or none of them are.

## Setup local environment

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Scoop](https://scoop.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

  ```bash
  scoop install migrate
  ```

- [DB Docs](https://dbdocs.io/docs)

  ```bash
  npm install -g dbdocs
  dbdocs login
  ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

  ```bash
  npm install -g @dbml/cli
  dbml2sql --version
  ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)
  with powershell

```powershell
    docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
```

or cmd

```CMD
    docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate
```

- [Gomock](https://github.com/golang/mock)

  ```bash
  go install github.com/golang/mock/mockgen@v1.6.0
  ```

### Setup infrastructure

- Create the bank-network:

```bash
make network
```

- Start the postgres container:

```bash
make postgres
```

- Create simple_bank database:

```bash
make createdb
```

- Run db migration up all version:

```bash
make migrateup
```

- Run db migration up 1 version:

```bash
make migrateup1
```

- Run db migration down all version:

```bash
make migratedown
```

- Run db migration down 1 version:

```bash
make migratedown1
```

### Documentation

- Generate DB documentation:

```bash
make db_docs
```

- Access database documentation at [this address](https://dbdocs.io/techschool.guru/simple_bank). Password: `secret`(https://dbdocs.io/)

### How to generate code

- Generate schema SQL file with DBML:

```bash
make db_schema
```

- Generate SQL CRUD with sqlc:

```bash
make sqlc
```

- Generate DB mock with gomock:

```bash
make mock
```

- Create a new db migration:

```bash
make new_migration name=<migration_name>
```

### How to run

- Run the server:

```bash
make server
```

- Run test:

```bash
make test
```
