package sorter

import "fmt"

type Sorter struct {
	Field     string
	Direction string
	Alias     string
}

func New(field, direction string, alias ...string) *Sorter {
	Alias := ""
	if len(alias) != 0 {
		Alias = alias[0]
	}
	return &Sorter{
		Field:     field,
		Direction: direction,
		Alias:     Alias,
	}
}

func (s *Sorter) GetName() string {
	return s.Field
}

func (s *Sorter) GetDirection() string {
	return s.Direction
}

func (s *Sorter) GetSql() string {
	if len(s.Alias) == 0 {
		return fmt.Sprintf(" %s %s", s.Field, s.Direction)
	} else {
		return fmt.Sprintf(" %s.%s %s", s.Alias, s.Field, s.Direction)
	}
}
