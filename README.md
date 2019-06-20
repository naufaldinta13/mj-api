# MJ API

[![build status](https://git.qasico.com/mj/api/badges/develop/build.svg)](https://git.qasico.com/mj/api/commits/develop)
[![coverage report](https://git.qasico.com/mj/api/badges/develop/coverage.svg)](https://git.qasico.com/mj/api/commits/develop)

API that will consumed by mj.

## Installation

Requirement:
```bash
    1. github.com/Masterminds/glide
    2. github.com/mattes/migrate
```

```bash
    1. go get git.qasico.com/mj/api
    2. run glide install
    3. run migration
```

## Database

How to use migration please see https://github.com/mattes/migrate.

### Migrations

Database migrations are on migrations folder, see https://github.com/mattes/migrate.

Example migrating:
```bash
migrate -url="mysql://username:password@tcp(127.0.0.1:3306)/project_mj" -path="./migrations" up
```

```bash
migrate -url="mysql://username:password@tcp(127.0.0.1:3306)/project_mj" -path="./migrations" down
```

### Running Test

Test can be executed by makefile
```bash
	make test -s
```

For fixing formating and lint on development mode
```bash
	make format -s
```

Completed command just read the help
```bash
	make
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a merge request :D
