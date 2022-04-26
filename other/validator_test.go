package other

import (
	"testing"
	"time"
)

type User struct {
	Name    string    `validate:"ne=admin"`
	Age     int       `validate:"gte=18"`
	Sex     string    `validate:"oneof=male female"`
	RegTime time.Time `validate:"lte"`
}

func TestValidator(t *testing.T) {
}
