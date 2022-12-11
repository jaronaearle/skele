package data

type AvyCenterUrlPaths struct {
	BaseUrl      string
	Forecast     string
	Avalanches   string
	Observations string
}

var AvyUrlPaths = AvyCenterUrlPaths{
	BaseUrl:      "https://utahavalanchecenter.org",
	Forecast:     "/forecast/salt-lake",
	Avalanches:   "/avalanches",
	Observations: "/observations",
}
