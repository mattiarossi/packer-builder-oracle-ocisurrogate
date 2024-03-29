package ocisurrogate

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v36/core"
)

// Artifact is an artifact implementation that contains a built Custom Image.
type Artifact struct {
	Image  core.Image
	Region string
	driver Driver

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]interface{}
}

// BuilderId uniquely identifies the builder.
func (a *Artifact) BuilderId() string {
	return BuilderId
}

// Files lists the files associated with an artifact. We don't have any files
// as the custom image is stored server side.
func (a *Artifact) Files() []string {
	return nil
}

// Id returns the OCID of the associated Image.
func (a *Artifact) Id() string {
	return *a.Image.Id
}

func (a *Artifact) String() string {
	var displayName string
	if a.Image.DisplayName != nil {
		displayName = *a.Image.DisplayName
	}

	return fmt.Sprintf(
		"An image was created: '%v' (OCID: %v) in region '%v'",
		displayName, *a.Image.Id, a.Region,
	)
}

// State ...
func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

// Destroy deletes the custom image associated with the artifact.
func (a *Artifact) Destroy() error {
	return a.driver.DeleteImage(context.TODO(), *a.Image.Id)
}
