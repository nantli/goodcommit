package commit

type Config struct {
	Type         string
	Scope        string
	Description  string
	Body         string
	Footer       string
	Breaking     bool
	CoAuthoredBy []string
	Extras       map[string]*string
}
