package names

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestNamesComponent(t *testing.T) {
	n := NewNameComponent("test")
	testutils.Equal(t, n.GetName(), "test")
	n.SetName("abc")
	testutils.Equal(t, n.GetName(), "abc")
}
