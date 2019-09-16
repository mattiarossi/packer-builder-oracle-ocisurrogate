package ocisurrogate

import (
	"context"

	"github.com/oracle/oci-go-sdk/core"
)

// Driver interfaces between the builder steps and the OCI SDK.
type Driver interface {
	CreateInstance(ctx context.Context, publicKey string, surrogateVolumeId string) (string, error)
	CreateBootClone(ctx context.Context, InstanceId string) (string, error)
	AttachBootClone(ctx context.Context, InstanceId string, VolumeId string) (string, error)
	DetachBootClone(ctx context.Context, VolumeId string) (string, error)
	CreateImage(ctx context.Context, id string) (core.Image, error)
	DeleteImage(ctx context.Context, id string) error
	GetInstanceIP(ctx context.Context, id string) (string, error)
	TerminateInstance(ctx context.Context, id string) error
	DeleteBootVolume(ctx context.Context, id string) error
	WaitForImageCreation(ctx context.Context, id string) error
	WaitForInstanceState(ctx context.Context, id string, waitStates []string, terminalState string) error
	WaitForBootVolumeState(ctx context.Context, id string, waitStates []string, terminalState string) error
	WaitForVolumeAttachmentState(ctx context.Context, id string, waitStates []string, terminalState string) error
}
