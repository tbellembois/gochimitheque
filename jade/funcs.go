package jade

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tbellembois/gochimitheque/global"
	"github.com/tbellembois/gochimitheque/helpers"
)

// T returns the translated messageID string
func T(messageID string, pluralCount int) string {
	return global.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID, PluralCount: pluralCount})
}

// ViewContainer is a struct passed to the view
type ViewContainer = helpers.ViewContainer
