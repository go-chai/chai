module github.com/go-chai/chai

go 1.18

require (
	github.com/bxcodec/faker/v3 v3.6.0
	github.com/getkin/kin-openapi v0.87.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/docgen v1.2.0
	github.com/go-openapi/spec v0.20.4
	github.com/pkg/errors v0.9.1
	github.com/swaggo/swag v1.7.6
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.7 // indirect
)

replace github.com/swaggo/swag v1.7.6 => ./third_party/swag
