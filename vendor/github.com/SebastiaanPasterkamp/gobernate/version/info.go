package version

// Info is the stuctured version data to be returned by the version
// handler.
type Info struct {
	BuildTime string `json:"build_time"`
	Commit    string `json:"commit"`
	Name      string `json:"name"`
	Release   string `json:"release"`
}
