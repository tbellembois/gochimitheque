# technical notes

## permissions

possible permissions:
`r`, `w`, `all` on each item

possible items:
- all
- products
- rproducts
- storages
- storelocations
- entities

- `entity_id` = -1 : all items
- `entity_id` = -2 : any/one of item -> -2 is never stored in the db but when used the item_id field is ommited requesting the permission table

possible combinations for a given person:

| item_name (ex: storage) | perm_name (ex: r) | entity_id (ex: 2) | notes |
| :-- | :--: | --: | :-- |
| `all`       |     `all`     |  ? | manager  ex: `all` permission on all items of entity with id `3` |
| `all`       |     `all`     | -1 | super admin |
| `all`       |     `all`     | -2 | manager of at least one entity |
| `all`     |   ?    |      ? | ex: `r` on `all` items of entity with `id` 3 |
| `all`     |   ?    |      -1| ex: `r` permission on all items of all entities |
| `all`     |   ?    |      -2| *not used* ex: `r` permission on all items of at least one entity|
| ?     |   `all`    |      ? | ex: `all` permission on `storage` of entity with `id` 3 |
| ?     |   `all`    |      -1| ex: `all` permission on `storage` of all entities |
| ?     |   `all`    |      -2| *not used* ex: `all` permission on `storage` of at least one entity |
| ?     |   ?    |   -1 | ex: `r` permission on all `storage` on all entities |
| ?     |   ?    |   ?  | ex: `r` permission on `storage` on entity with `id` 3 |
| ?     |   ?    |   -2  | *not used* ex: `r` permission on `storage` on at least one entity |

permissions on other items (suppliers, classes of compounds...) are linked to the items above

entities / people management:
- only super admins can create new entities and modify entities
- entity managers can create/update/delete people in their entities

implemented in `globals/global.go` with `PermMatrix`

## authorization

A custom middleware. Look at `func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler ` in `handlers\auth.go`.

## static content

JS install/upgrade:

*these 3 modules are version dependent*
npm install jquery
npm install @popperjs/core
npm install bootstrap

npm install bootstrap-colorpicker
npm install bootstrap-table
npm install jquery-validation
npm install select2
npm install --save-dev @fortawesome/fontawesome-free
npm install @mdi/font
npm install animate.css --save

rsync -av ./node_modules/bootstrap/dist/js/bootstrap.min.js ./static/js/
rsync -av ./node_modules/bootstrap/dist/js/bootstrap.min.js.map ./static/js/
rsync -av ./node_modules/bootstrap-colorpicker/dist/js/bootstrap-colorpicker.min.js ./static/js/
rsync -av ./node_modules/bootstrap-colorpicker/dist/js/bootstrap-colorpicker.min.js.map ./static/js/
rsync -av ./node_modules/bootstrap-table/dist/bootstrap-table.min.js  ./static/js/
rsync -av ./node_modules/jquery/dist/jquery.min.js ./static/js/
rsync -av ./node_modules/jquery-validation/dist/jquery.validate.min.js ./static/js/
rsync -av ./node_modules/jquery-validation/dist/additional-methods.min.js  ./static/js/jquery.validate.additional-methods.min.js 
rsync -av ./node_modules/@popperjs/core/dist/umd/popper.min.js ./static/js/
rsync -av ./node_modules/@popperjs/core/dist/umd/popper.min.js.map ./static/js/
rsync -av ./node_modules/select2/dist/js/select2.full.min.js ./static/js/

rsync -av ./node_modules/bootstrap/dist/css/bootstrap.min.css ./static/css/
rsync -av ./node_modules/bootstrap/dist/css/bootstrap.min.css.map ./static/css/
rsync -av ./node_modules/@fortawesome/fontawesome-free/css/all.min.css ./static/css/fontawesome.all.min.css
rsync -av ./node_modules/@fortawesome/fontawesome-free/webfonts/* ./static/webfonts/
rsync -av ./node_modules/bootstrap-colorpicker/dist/css/bootstrap-colorpicker.min.css ./static/css/
rsync -av ./node_modules/bootstrap-colorpicker/dist/css/bootstrap-colorpicker.min.css.map ./static/css/
rsync -av ./node_modules/bootstrap-table/dist/bootstrap-table.min.css ./static/css/
rsync -av ./node_modules/@mdi/font/css/materialdesignicons.min.css ./static/css/
rsync -av ./node_modules/@mdi/font/css/materialdesignicons.min.css.map ./static/css/
rsync -av ./node_modules/@mdi/font/fonts/* ./static/fonts/
rsync -av ./node_modules/select2/dist/css/select2.min.css ./static/css/
rsync -av ./node_modules/animate.css/animate.min.css ./static/css/

## windows cross compilation (officially not supported)

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