package local_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/driver"
	. "github.com/mlafeldt/chef-runner/driver/local"
	"github.com/stretchr/testify/assert"
)

func TestDriverInterface(t *testing.T) {
	assert.Implements(t, (*driver.Driver)(nil), new(Driver))
}

func TestString(t *testing.T) {
	expect := "Local driver (hostname: some-host)"
	actual := Driver{Hostname: "some-host"}.String()
	assert.Equal(t, expect, actual)
}
