package abstraction

type Tracker interface {
	WriteToCSV(string, []string) error
}
