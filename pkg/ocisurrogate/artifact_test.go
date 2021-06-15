package ocisurrogate

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestArtifactImpl(t *testing.T) {
	var raw interface{}
	raw = &Artifact{}
	if _, ok := raw.(packer.Artifact); !ok {
		t.Fatalf("Artifact should be artifact")
	}
}
