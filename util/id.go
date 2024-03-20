package util

import (
	"database/sql/driver"
	"time"

	"github.com/segmentio/ksuid"
)

type ID struct {
	id ksuid.KSUID
}

func NewID(t time.Time) (ID, error) {
	id, err := ksuid.NewRandomWithTime(t)
	if err != nil {
		return NilIdentifier, err
	}
	return ID{id: id}, nil
}

var NilIdentifier ID

func (i ID) Time() time.Time {
	return i.id.Time()
}

func (i ID) String() string {
	return i.id.String()
}

func (i ID) MarshalText() ([]byte, error) {
	return i.id.MarshalText()
}

func (i ID) MarshalBinary() ([]byte, error) {
	return i.id.MarshalBinary()
}

func (i *ID) UnmarshalText(b []byte) error {
	return i.id.UnmarshalText(b)
}

func (i *ID) UnmarshalBinary(b []byte) error {
	return i.id.UnmarshalBinary(b)
}

func (i ID) Value() (driver.Value, error) {
	return i.id.Value()
}

// Scan implements the sql.Scanner interface. It supports converting from
// string, []byte, or nil into a Identifier value. Attempting to convert from
// another type will return an error.
func (i *ID) Scan(src interface{}) error {
	return i.id.Scan(src)
}

func Parse(s string) (ID, error) {
	id, err := ksuid.Parse(s)
	if err != nil {
		return NilIdentifier, err
	}
	return ID{id: id}, nil
}

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
