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

### database name convention:
 
 - lowercase names
 - separate words with underscore
 - singular names
 - pk: tablename_id
 - fk: tablename_targettablename_id
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
