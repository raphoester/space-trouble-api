package birthday_test

import (
	"testing"

	"github.com/raphoester/space-trouble-api/internal/pkg/birthday"
	"github.com/stretchr/testify/assert"
)

func TestNew_NominalCase(t *testing.T) {
	b, err := birthday.Parse("13/02")
	assert.NoError(t, err)
	assert.Equal(t, "13/02", b.String())
}

func TestNew_MonthTooBig(t *testing.T) {
	_, err := birthday.Parse("13/13")
	assert.Error(t, err)
}

func TestNew_DayTooBig(t *testing.T) {
	_, err := birthday.Parse("32/12")
	assert.Error(t, err)
}
