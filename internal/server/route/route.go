package route

type Routes struct {
	ApiPath       string
	ApiHealtcheck string
	ApiListLinks  string
	ApiShowLink   string
	ApiAddLink    string
	ApiUpdateLink string
	Redirect      string
}

func Configure() Routes {
	return Routes{
		// Root
		Redirect: "/{alias}",

		// API
		ApiPath:       "/v1",
		ApiHealtcheck: "/health",
		ApiListLinks:  "/links",
		ApiShowLink:   "/link/{id}",
		ApiAddLink:    "/link",
		ApiUpdateLink: "/link/{id}",
	}
}
