# Chimithèque 2.0 (Go version)

This is a work in progress !

# introduction

This project is a complete rewrite of the [Chimithèque](https://github.com/tbellembois/chimitheque) project in [Go](https://golang.org/).

![screenshot](screenshot.png)

# technical

No framework, just some toolkits.  
Restfull API.

## databases

- [https://github.com/jmoiron/sqlx]
- [https://github.com/Masterminds/squirrel]

## permissions

`r`, `rw` or `` on earch item

possible combinations for a given person:

| item_name (ex: entity) | item_permname (ex: r) | entity_id (ex: 2) | notes |
| :-- | :--: | --: | :-- |
| `all`       |     `all`     |  ? | *not used* - we will use -1 |
| `all`       |     `all`     | -1 | super admin |
| `all`     |   ?    |      ? | *not used* - non sense (ex: `r` on `all` items of entity with `id` 3) |
| `all`     |   ?    |      -1| ex: `r` permission on all items of all entities |
| ?     |   `all`    |      ? | ex: `all` permission on `storage` of entity with `id` 3 |
| ?     |   `all`    |      -1| ex: `all` permission on all `storage` of all entities |
| ?     |   ?    |   -1 | ex: `r` permission on all `storage` of all entities |
| ?     |   ?    |   ?  | ex: `r` permission on `storage` of entity with `id` 3 |

=> 6

final clean table:

| item_name (ex: entity) | item_permname (ex: r) | entity_id (ex: 2) | notes |
| :-- | :--: | --: | :-- |
| `all`       |     `all`     | -1 | super admin |
| `all`     |   ?    |      -1| ex: `r` permission on all items of all entities |
| ?     |   `all`    |      ? | ex: `all` permission on `storage` of entity with `id` 3 |
| ?     |   `all`    |      -1| ex: `all` permission on all `storage` of all entities |
| ?     |   ?    |   -1 | ex: `r` permission on all `storage` of all entities |
| ?     |   ?    |   ?  | ex: `r` permission on `storage` of entity with `id` 3 |

- `item_id` = -1 : all items
- `item_name` can be `all`
- `perm_name` can be `all`

possible items:
- product card
- restricted product card
- storage card
- archived storage card
- store location
- entity
- person
- class of compounds
- supplier

entities / people management:
- only super admins can create new entities and modify entities
- entity managers can create/update/delete people in their entities

### database name convention:
 
 - lowercase names
 - separate words with underscore
 - singular names
 - pk: tablename_id
 - fk: tablename_target_tablename_id
 - columns: tablename_fieldname

## middlewares

- [https://github.com/justinas/alice]

## routing

- [https://github.com/gorilla/mux]

## authentication

- [https://github.com/dgrijalva/jwt-go]

## authorization

A custom middleware. Look at `func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler ` in `handlers\auth.go`.

## UI

- [http://bootstrap-confirmation.js.org/]
- [http://bootstrap-table.wenzhixin.net.cn/]
- [https://github.com/creative-area/jQuery-form-autofill]
- [https://github.com/Joker/jade]
- [https://jquery.com/]
- [https://jqueryvalidation.org/]
- [https://select2.org/]
- [https://v4-alpha.getbootstrap.com/]

## web form to struct

- [https://github.com/gorilla/schema]
