// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package incus

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/common"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	// The name of the output artifact. Defaults to
	// name.
	OutputImage string `mapstructure:"output_image" required:"false"`
	// Reuse image alias if it already exists; defaults to false
	Reuse         bool   `mapstructure:"reuse" required:"false"`
	ContainerName string `mapstructure:"container_name"`
	// The (optional) name of the incus remote on which to publish the
	// container image.
	PublishRemoteName string `mapstructure:"publish_remote_name" required:"false"`
	// Lets you prefix all builder commands, such as
	// with ssh for a remote build host. Defaults to `{{.Command}}`; i.e. no
	// wrapper.
	CommandWrapper string `mapstructure:"command_wrapper" required:"false"`
	// The source image to use when creating the build
	// container. This can be a (local or remote) image (name or fingerprint).
	// E.G. my-base-image, ubuntu-daily:x, 08fababf6f27, ...
	Image   string   `mapstructure:"image" required:"true"`
	Profile []string `mapstructure:"profile"`
	Project string   `mapstructure:"project"`
	// The number of seconds to sleep between launching
	// the incus instance and provisioning it; defaults to 3 seconds.
	InitSleep string `mapstructure:"init_sleep" required:"false"`
	// Pass key values to the publish
	// step to be set as properties on the output image. This is most helpful to
	// set the description, but can be used to set anything needed. See
	// https://stgraber.org/2016/03/30/incus-2-0-image-management-512/
	// for more properties.
	PublishProperties map[string]string `mapstructure:"publish_properties" required:"false"`
	// List of key/value pairs you wish to
	// pass to lxc launch via --config. Defaults to empty.
	LaunchConfig map[string]string `mapstructure:"launch_config" required:"false"`
	// Create incus virtual-machine image on hosts running incus 4.0 and above; defaults to false for container image
	VirtualMachine bool `mapstructure:"virtual_machine"`
	// Skip execute `lxc publish`; defaults to false
	SkipPublish bool `mapstructure:"skip_publish" required:"false"`

	ctx interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) error {

	var md mapstructure.Metadata
	err := config.Decode(c, &config.DecodeOpts{
		Metadata:    &md,
		Interpolate: true,
	}, raws...)
	if err != nil {
		return err
	}

	// Accumulate any errors
	var errs *packersdk.MultiError

	if c.ContainerName == "" {
		c.ContainerName = fmt.Sprintf("packer-%s", c.PackerBuildName)
	}

	if c.OutputImage == "" {
		c.OutputImage = c.ContainerName
	}

	if c.CommandWrapper == "" {
		c.CommandWrapper = "{{.Command}}"
	}

	if c.Image == "" {
		errs = packersdk.MultiErrorAppend(errs, fmt.Errorf("`image` is a required parameter for incus. Please specify an image by alias or fingerprint. e.g. `ubuntu-daily:x`"))
	}

	if len(c.Profile) == 0 {
		c.Profile = []string{"default"}
	}

	if c.Project == "" {
		c.Project = "default"
	}

	// Sadly we have to wait a few seconds for /tmp to be intialized and networking
	// to finish starting. There isn't a great cross platform to check when things are ready.
	if c.InitSleep == "" {
		c.InitSleep = "3"
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}
