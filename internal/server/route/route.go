package route

type Routes struct {
	Root root
	API  api
}

type root struct {
	Path     string
	Redirect string
}

type api struct {
	Path       string
	Healtcheck string
	ListLinks  string
	ShowLink   string
	AddLink    string
	UpdateLink string
	DeleteLink string
}

func Configure() Routes {
	return Routes{
		Root: root{
			Path:     "/",
			Redirect: "/{alias}",
		},
		API: api{
			Path:       "/v1",
			Healtcheck: "/health",
			ListLinks:  "/links",
			ShowLink:   "/link/{id}",
			AddLink:    "/link",
			UpdateLink: "/link/{id}",
			DeleteLink: "/link/{id}",
		},
	}
}
