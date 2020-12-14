package jade

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tbellembois/gochimitheque/globals"
	"github.com/tbellembois/gochimitheque/models"
)

// T returns the translated messageID string
func T(messageID string, pluralCount int) string {
	return globals.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID, PluralCount: pluralCount})
}

// ViewContainer is a struct passed to the view
type ViewContainer = models.ViewContainer
