package ocisurrogate

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreateInstance struct{}

func (s *stepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		driver = state.Get("driver").(Driver)
		ui     = state.Get("ui").(packer.Ui)
		config = state.Get("config").(*Config)
	)

	ui.Say("Creating instance...")

	instanceID, err := driver.CreateInstance(ctx, string(config.Comm.SSHPublicKey),"")
	if err != nil {
		err = fmt.Errorf("Problem creating instance: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("instance_id", instanceID)

	ui.Say(fmt.Sprintf("Created instance (%s).", instanceID))

	ui.Say("Waiting for instance to enter 'RUNNING' state...")

	if err = driver.WaitForInstanceState(ctx, instanceID, []string{"STARTING", "PROVISIONING"}, "RUNNING"); err != nil {
		err = fmt.Errorf("Error waiting for instance to start: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say("Instance 'RUNNING'.")
	ui.Say("Cloning Boot Volume to surrogate ...")

	clonedVolumeID, err := driver.CreateBootClone(ctx, instanceID)
	if err != nil {
		err = fmt.Errorf("Problem creating Boot Volume Clone: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Say("Surrogate Boot Volume Cloned.")
	state.Put("cloned_volume_id", clonedVolumeID)
	ui.Say("Waiting for Cloned Volume to enter 'AVAILABLE' state...")
	if err = driver.WaitForBootVolumeState(ctx, clonedVolumeID, []string{"PROVISIONING","RESTORING"}, "AVAILABLE"); err != nil {
		err = fmt.Errorf("Error waiting for Volume to be available: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Say("Surrogate Boot Volume in 'AVAILABLE' State.")

	ui.Say(fmt.Sprintf("Attaching Cloned Volume to instance (%s).", clonedVolumeID))
	attachedVolumeID, err := driver.AttachBootClone(ctx, instanceID, clonedVolumeID)
	if err != nil {
		err = fmt.Errorf("Problem Attaching Boot Volume Clone: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Say("Surrogate Boot Volume Attachment created.")
	ui.Say(fmt.Sprintf("Waiting for Attached Volume %s to enter 'ATTACHED' state...",attachedVolumeID))
	if err = driver.WaitForVolumeAttachmentState(ctx, attachedVolumeID, []string{"ATTACHING"}, "ATTACHED"); err != nil {
		err = fmt.Errorf("Error waiting for Volume to be attached: %s", err)
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Say(fmt.Sprintf("Cloned Volume Attached successfully to instance (%s) with id %s.", instanceID, attachedVolumeID))
	state.Put("attached_volume_id", attachedVolumeID)

	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	idRaw, ok := state.GetOk("instance_id")
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
	idSurrogate, ok := state.GetOk("instance_surrogate_id")
	if ok {
		ui.Say(fmt.Sprintf("Deleting Surrogate Volume (%s)...", idVolume))
		if err := driver.DeleteBootVolume(context.TODO(), idVolume); err != nil {
			err = fmt.Errorf("Error terminating Surrogate Boot Volume. Please terminate manually: %s", err)
			ui.Error(err.Error())
			state.Put("error", err)
			return
		}

		err = driver.WaitForBootVolumeState(context.TODO(), idVolume, []string{"TERMINATING"}, "TERMINATED")
		if err != nil {
			err = fmt.Errorf("Error terminating instance. Please terminate manually: %s", err)
			ui.Error(err.Error())
			state.Put("error", err)
			return
		}

		ui.Say("Deleted Surrogate Volume.")
	}
	
}
