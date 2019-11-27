# technical notes

## databases

### sqlite

- [https://github.com/mattn/go-sqlite3]

### helpers

- [https://github.com/jmoiron/sqlx]
- [https://github.com/Masterminds/squirrel]

### naming convention
 
 - lowercase names
 - separate words with underscore
 - singular names
 - pk: tablename_id
 - fk: tablename_target_tablename_id
 - columns: tablename_fieldname

## permissions

`r`, `rw` or `` on earch item

possible combinations for a given person:

| item_name (ex: storage) | item_permname (ex: r) | entity_id (ex: 2) | notes |
| :-- | :--: | --: | :-- |
| `all`       |     `all`     |  ? | ex: `all` permission on all items of entity with id `3` : manager |
| `all`       |     `all`     | -1 | super admin |
| `all`     |   ?    |      ? | *not used* - non sense (ex: `r` on `all` items of entity with `id` 3) |
| `all`     |   ?    |      -1| ex: `r` permission on all items of all entities |
| ?     |   `all`    |      ? | ex: `all` permission on `storage` of entity with `id` 3 |
| ?     |   `all`    |      -1| ex: `all` permission on all `storage` of all entities |
| ?     |   ?    |   -1 | ex: `r` permission on all `storage` of all entities |
| ?     |   ?    |   ?  | ex: `r` permission on `storage` of entity with `id` 3 |

- `item_id` = -1 : all items
- `item_id` = -2 : any item -> -2 is never stored in the db but when used the item_id field is ommited requesting the permission table
- `item_name` can be `all`
- `perm_name` can be `all`

possible items:
- product card
- restricted product card
- storage card
- store location
- entity

permissions on other items (suppliers, classes of compounds...) are linked to the items above

entities / people management:
- only super admins can create new entities and modify entities
- entity managers can create/update/delete people in their entities

## toolkits used (plus db helpers)

### middlewares

- [https://github.com/justinas/alice]

### routing

- [https://github.com/gorilla/mux]

### authentication

- [https://github.com/dgrijalva/jwt-go]

### web form to struct

- [https://github.com/gorilla/schema]

### UI

- [http://bootstrap-confirmation.js.org/]
- [http://bootstrap-table.wenzhixin.net.cn/]
- [https://github.com/creative-area/jQuery-form-autofill]
- [https://github.com/Joker/jade]
- [https://jquery.com/]
- [https://jqueryvalidation.org/]
- [https://select2.org/]
- [https://v4-alpha.getbootstrap.com/]

### authorization

A custom middleware. Look at `func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler ` in `handlers\auth.go`.

## windows cross compilation

### windows 10

```bash
    go generate
    CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build
```

### windows 7

```bash
    go generate
    CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build
```

### installed arch linux packages

> local/mingw-w64-binutils-bin 2.31.1-1 (mingw-w64-toolchain mingw-w64)
>     Cross binutils for the MinGW-w64 cross-compiler (pre-compiled)
> local/mingw-w64-crt-bin 6.0.0-1 (mingw-w64-toolchain mingw-w64)
>     MinGW-w64 CRT for Windows (pre-compiled)
> local/mingw-w64-gcc-bin 8.2.0-1 (mingw-w64-toolchain mingw-w64)
>     Cross GCC for the MinGW-w64 cross-compiler (pre-compiled)
> local/mingw-w64-headers-bin 6.0.0-1 (mingw-w64-toolchain mingw-w64)
>     MinGW-w64 headers for Windows (pre-compiled)
> local/mingw-w64-winpthreads-bin 6.0.0-1 (mingw-w64-toolchain mingw-w64)
>     MinGW-w64 winpthreads library (pre-compiled)

## chimithèque V1 databases export

### postgresql to csv

```bash
SCHEMA="public"; DB="chimitheque"; psql -U [user] -h [host] -p [port] -Atc "select tablename from pg_tables where schemaname='$SCHEMA'" $DB | while read TBL; do psql -U [user] -h [host] -p [port] -c "COPY $SCHEMA.$TBL TO STDOUT WITH CSV HEADER" $DB > $TBL.csv; done;
```

## testers

delphine.pitrat@ens-lyon.fr, laurelise.chapellet@ens-lyon.fr, guillaume.george@ens-lyon.fr, sylvain.david@ens-lyon.fr, laure.guy@ens-lyon.fr, yann.bretonniere@ens-lyon.fr, loic.richard@irstea.fr, christophe.le-bourlot@insa-lyon.fr, julien.devemy@uca.fr, alix.tordo@ens-lyon.fr, clement.courtin@ens-lyon.fr

## REST API samples

### auth

- get an authentication token
```
curl -X POST http://localhost:8081/get-token -d "person_email=admin@chimitheque.fr&person_password=test"
```

### people

- list people
```
curl -X GET --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAc3VwZXIuY29tIiwiZXhwIjoxNTUwODU5MTQxfQ.Az49BEqLmxmsS5OxSe49K9Cbli3yhaWMJe_wDsp8A4w" http://localhost:8081/people
```

- fake deleting a person
```
curl -X DELETE --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAc3VwZXIuY29tIiwiZXhwIjoxNTUwODU5MTQxfQ.Az49BEqLmxmsS5OxSe49K9Cbli3yhaWMJe_wDsp8A4w" http://localhost:8081/f/people/1
```

### storage

-- fake accessing any storage
```
curl -X GET --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAc3VwZXIuY29tIiwiZXhwIjoxNTUwODU5MTQxfQ.Az49BEqLmxmsS5OxSe49K9Cbli3yhaWMJe_wDsp8A4w" http://localhost:8081/f/storages/-2
```