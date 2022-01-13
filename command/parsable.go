package command

// Parsable can be parsed from string
type Parsable interface {
	ParseFromString(string) error
}

type any = interface{}
