package locales

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	Bundle    *i18n.Bundle
	Localizer *i18n.Localizer
)

func init() {

	// load translations
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	Bundle.MustParseMessageFileBytes(LOCALES_EN, "en.toml")
	Bundle.MustParseMessageFileBytes(LOCALES_FR, "fr.toml")

	Localizer = i18n.NewLocalizer(Bundle)
}
