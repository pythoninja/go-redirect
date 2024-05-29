package route

type Routes struct {
	ApiPath       string
	ApiHealtcheck string
	ApiListLinks  string
	ApiShowLink   string
	Redirect      string
}

func New() Routes {
	return Routes{
		// Root
		Redirect: "/{alias}",

		// API
		ApiPath:       "/v1",
		ApiHealtcheck: "/health",
		ApiListLinks:  "/links",
		ApiShowLink:   "/link/{id}",
	}
}
