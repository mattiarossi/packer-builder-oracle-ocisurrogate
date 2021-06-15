package ocisurrogate

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepImage struct{}

func (s *stepImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		driver           = state.Get("driver").(Driver)
		ui               = state.Get("ui").(packer.Ui)
		idVolume         = state.Get("cloned_volume_id").(string)
		attachedVolumeID = state.Get("attached_volume_id").(string)
		config           = state.Get("config").(*Config)
	)

	ui.Say("Detaching Boot Volume from main instance...")
	detachedVolumeID, err := driver.DetachBootClone(ctx, attachedVolumeID)
	if err != nil {
		err = fmt.Errorf("Problem Detaching Boot Volume Clone: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Say(fmt.Sprintf("Surrogate Boot Volume Detachment request created for %s.", detachedVolumeID))
	ui.Say(fmt.Sprintf("Waiting for Attached Volume %s to enter 'DETACHED' state...", attachedVolumeID))
	if err = driver.WaitForVolumeAttachmentState(ctx, attachedVolumeID, []string{"DETACHING"}, "DETACHED"); err != nil {
		err = fmt.Errorf("Error waiting for Volume to be detached: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say("Cloned Volume detached...")
	ui.Say("Creating Surrogate instance...")

	instanceSurrogateID, err := driver.CreateInstance(ctx, string(config.Comm.SSHPublicKey), idVolume)
	if err != nil {
		err = fmt.Errorf("Problem creating surrogate instance: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("instance_surrogate_id", instanceSurrogateID)

	ui.Say(fmt.Sprintf("Created Surrogate instance (%s).", instanceSurrogateID))

	ui.Say("Waiting for Surrogate instance to enter 'RUNNING' state...")

	if err = driver.WaitForInstanceState(ctx, instanceSurrogateID, []string{"STARTING", "PROVISIONING"}, "RUNNING"); err != nil {
		err = fmt.Errorf("Error waiting for instance to start: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say("Surrogate Instance 'RUNNING'.")

	ui.Say("Creating image from Surrogate instance...")

	image, err := driver.CreateImage(ctx, instanceSurrogateID)
	if err != nil {
		err = fmt.Errorf("Error creating image from instance: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	err = driver.WaitForImageCreation(ctx, *image.Id)
	if err != nil {
		err = fmt.Errorf("Error waiting for image creation to finish: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	// TODO(apryde): This is stale as .LifecycleState has changed to
	// AVAILABLE at this point. Does it matter?
	state.Put("image", image)

	ui.Say("Image created.")

	return multistep.ActionContinue
}

func (s *stepImage) Cleanup(state multistep.StateBag) {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	idRaw, ok := state.GetOk("instance_surrogate_id")
	if !ok {
		return
	}
	id := idRaw.(string)

	ui.Say(fmt.Sprintf("Terminating instance (%s)...", id))

	if err := driver.TerminateInstance(context.TODO(), id); err != nil {
		err = fmt.Errorf("Error terminating instance. Please terminate manually: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return
	}

	err := driver.WaitForInstanceState(context.TODO(), id, []string{"TERMINATING"}, "TERMINATED")
	if err != nil {
		err = fmt.Errorf("Error terminating instance. Please terminate manually: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return
	}

	ui.Say("Terminated instance.")
}
