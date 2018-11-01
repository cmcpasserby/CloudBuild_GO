package responses

type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
)

type Link struct {
	Method string `json:"method"`
	Href   string `json:"href"`
}
