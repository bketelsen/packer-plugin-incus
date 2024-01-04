// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package incus

import (
	"context"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// The unique ID for this builder
const BuilderId = "incus"

type wrappedCommandTemplate struct {
	Command string
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	errs := b.config.Prepare(raws...)
	if errs != nil {
		return nil, nil, errs
	}
	buildGeneratedData := []string{"Alias"}

	return buildGeneratedData, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packersdk.Ui, hook packersdk.Hook) (packersdk.Artifact, error) {
	wrappedCommand := func(command string) (string, error) {
		b.config.ctx.Data = &wrappedCommandTemplate{Command: command}
		return interpolate.Render(b.config.CommandWrapper, &b.config.ctx)
	}

	steps := []multistep.Step{
		&stepIncusLaunch{},
		&StepProvision{},
		&stepPublish{},
	}

	// Setup the state bag
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("wrappedCommand", CommandWrapper(wrappedCommand))

	state.Put("generated_data", map[string]interface{}{
		"Alias": b.config.OutputImage,
	})

	// Run
	b.runner = commonsteps.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(ctx, state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	id := "0"
	if !b.config.SkipPublish {
		id = state.Get("imageFingerprint").(string)
	}

	artifact := &Artifact{
		id:        id,
		StateData: map[string]interface{}{"generated_data": state.Get("generated_data")},
	}

	return artifact, nil
}
