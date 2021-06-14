package ocisurrogate

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{}
	raw = &Builder{}
	if _, ok := raw.(packer.Builder); !ok {
		t.Fatalf("Builder should be a builder")
	}
}
