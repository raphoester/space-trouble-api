package date_test

import (
	"testing"

	"github.com/raphoester/space-trouble-api/internal/pkg/date"
	"github.com/stretchr/testify/assert"
)

func TestNew_NominalCase(t *testing.T) {
	d, err := date.Parse("13/10/2024")
	assert.NoError(t, err)
	assert.Equal(t, "13/10/2024", d.String())
}
