package request

import (
	"net/http"
	"net/url"
	"regexp"
)

// ChimithequeContextKey is the Go request context
// used in each request.
type ChimithequeContextKey string

// Container is a struct passed to the view.
type Container struct {
	PersonEmail    string `json:"PersonEmail"`
	PersonLanguage string `json:"PersonLanguage"`
	PersonID       int64  `json:"PersonID"`
	AppURL         string `json:"AppURL"`
	AppPath        string `json:"AppPath"`
	BuildID        string `json:"BuildID"`
	DisableCache   bool   `json:"DisableCache"`
}

// ContainerFromRequestContext returns a ViewContainer from the request context
// initialized in the AuthenticateMiddleware and AuthorizeMiddleware middlewares.
func ContainerFromRequestContext(r *http.Request) Container {
	// getting the request context
	var (
		container Container
	)

	ctx := r.Context()
	ctxcontainer := ctx.Value(ChimithequeContextKey("container"))

	if ctxcontainer != nil {
		container = ctxcontainer.(Container)
	}

	return container
}

// EndsPathWithDigits checks if the path of the given URL ends with digits (ignoring query parameters)
func EndsPathWithDigits(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Regex to match if the path ends with digits
	re := regexp.MustCompile(`\d+$`)
	return re.MatchString(parsedURL.Path)
}

// HasIDParam checks if the request URL has a query parameter named "id"
func HasIDParam(r *http.Request) bool {
	// Get the "id" query parameter
	_, ok := r.URL.Query()["id"]
	return ok
}
