module github.com/tbellembois/gochimitheque

go 1.19

// replace github.com/tbellembois/gochimitheque-utils => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-utils

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/casbin/casbin/v2 v2.57.1
	github.com/dchest/authcookie v0.0.0-20190824115100-f900d2294c8e // indirect
	github.com/dchest/passwordreset v0.0.0-20190826080013-4518b1f41006
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.18.0
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/justinas/alice v1.2.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/nicksnyder/go-i18n/v2 v2.2.1
	github.com/sirupsen/logrus v1.9.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/steambap/captcha v1.4.1
	golang.org/x/crypto v0.3.0
	golang.org/x/image v0.1.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0
)

require (
	github.com/casbin/json-adapter/v2 v2.0.0
	github.com/go-ldap/ldap/v3 v3.4.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/tbellembois/gochimitheque-utils v0.0.0-20221125140514-6e4c4a07ea3b
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20220621081337-cb9428e4ac1e // indirect
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
