package version

var (
	Version = "0.0.0-dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func String() string {
	return Version + " (" + Commit + ", " + Date + ")"
}
