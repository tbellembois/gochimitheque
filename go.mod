module github.com/tbellembois/gochimitheque

go 1.21

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/casbin/casbin/v2 v2.85.0
	github.com/dchest/authcookie v0.0.0-20190824115100-f900d2294c8e // indirect
	github.com/dchest/passwordreset v0.0.0-20190826080013-4518b1f41006
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.19.0
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/justinas/alice v1.2.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/nicksnyder/go-i18n/v2 v2.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/steambap/captcha v1.4.1
	golang.org/x/crypto v0.21.0
	golang.org/x/image v0.13.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0
)

require (
	github.com/barweiss/go-tuple v1.1.2
	github.com/casbin/json-adapter/v2 v2.1.1
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/coreos/go-oidc/v3 v3.10.0
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/pebbe/zmq4 v1.2.11
	github.com/russross/blackfriday/v2 v2.1.0
	golang.org/x/net v0.22.0
	golang.org/x/oauth2 v0.18.0
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/casbin/govaluate v1.1.1 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	github.com/tidwall/gjson v1.17.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
