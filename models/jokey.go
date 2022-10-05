package jokey

type Jokey struct {
	Id          int
	Name        string
	SuspendDate string
	AllowDate   string // TODO: change type time.Time?
	Description string // TODO: change type time.Time?
}
