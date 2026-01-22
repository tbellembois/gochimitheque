module github.com/tbellembois/gochimitheque

go 1.24

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/gorilla/mux v1.8.1
	github.com/justinas/alice v1.2.0
	github.com/nicksnyder/go-i18n/v2 v2.6.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0
)

require github.com/stretchr/testify v1.8.4 // indirect

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
