module github.com/tbellembois/gochimitheque

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/GeertJohan/go.rice v1.0.2
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/Joker/jade v1.0.1-0.20200506134858-ee26e3c533bb // indirect
	github.com/Masterminds/squirrel v1.5.0
	github.com/casbin/casbin/v2 v2.23.0
	github.com/casbin/json-adapter/v2 v2.0.0
	github.com/dchest/authcookie v0.0.0-20190824115100-f900d2294c8e // indirect
	github.com/dchest/passwordreset v0.0.0-20190826080013-4518b1f41006
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.10.0
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.1
	github.com/justinas/alice v1.2.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/steambap/captcha v1.3.1
	github.com/tbellembois/gochimitheque-utils v0.0.0-20210130182605-75000d2f72a5
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5
	gopkg.in/russross/blackfriday.v2 v2.1.0
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
