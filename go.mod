module github.com/tbellembois/gochimitheque

go 1.17

replace github.com/tbellembois/gochimitheque-utils v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-utils

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/casbin/casbin/v2 v2.34.0
	github.com/casbin/json-adapter/v2 v2.0.0
	github.com/dchest/authcookie v0.0.0-20190824115100-f900d2294c8e // indirect
	github.com/dchest/passwordreset v0.0.0-20190826080013-4518b1f41006
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.13.0
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/justinas/alice v1.2.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/steambap/captcha v1.3.1
	github.com/tbellembois/gochimitheque-utils v0.0.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.6
	gopkg.in/russross/blackfriday.v2 v2.1.0
)

require (
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
