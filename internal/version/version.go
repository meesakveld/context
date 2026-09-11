package version

var Current = "dev"

func String() string {
	return "context v" + Current
}