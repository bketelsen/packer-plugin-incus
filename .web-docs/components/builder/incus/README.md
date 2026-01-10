Type: `incus`
Artifact BuilderId: `incus`

The `incus` Packer builder builds containers for incus. The builder starts an incus
container, runs provisioners within this container, then saves the container as
an incus image.

The incus builder requires a modern linux kernel and the `incus` package. This
builder does not work with LXC.

## Basic Example

Below is a fully functioning example.

**HCL**

```hcl
source "incus" "incus-xenial" {
  image = "ubuntu-daily:xenial"
  output_image = "ubuntu-xenial"
  publish_properties {
    description = "Trivial repackage with Packer"
  }
}

build {
  sources = ["incus.incus-xenial"]
}
```

**JSON**

```json
{
  "builders": [
    {
      "type": "incus",
      "name": "incus-xenial",
      "image": "ubuntu-daily:xenial",
      "output_image": "ubuntu-xenial",
      "publish_properties": {
        "description": "Trivial repackage with Packer"
      }
    }
  ]
}
```

## Configuration Reference

### Required:

<!-- Code generated from the comments of the Config struct in builder/incus/config.go; DO NOT EDIT MANUALLY -->

- `image` (string) - The source image to use when creating the build
  container. This can be a (local or remote) image (name or fingerprint).
  E.G. my-base-image, ubuntu-daily:x, 08fababf6f27, ...

<!-- End of code generated from the comments of the Config struct in builder/incus/config.go; -->


  ~> Note: The builder may appear to pause if required to download a
  remote image, as they are usually 100-200MB. `/var/log/incus/incus.log` will
  mention starting such downloads.

### Optional:

<!-- Code generated from the comments of the Config struct in builder/incus/config.go; DO NOT EDIT MANUALLY -->

- `output_image` (string) - The name of the output artifact. Defaults to
  name.

- `reuse` (bool) - Reuse image alias if it already exists; defaults to false

- `container_name` (string) - Container Name

- `publish_remote_name` (string) - The (optional) name of the incus remote on which to publish the
  container image.

- `command_wrapper` (string) - Lets you prefix all builder commands, such as
  with ssh for a remote build host. Defaults to `{{.Command}}`; i.e. no
  wrapper.

- `profile` ([]string) - One or more Incus profiles to apply to the container or virtual machine.
  Profiles define base configuration, such as devices and limits, that is
  merged into the instance. Multiple profiles can be specified, and if none
  are provided this defaults to the "default" profile.

- `project` (string) - Project

- `init_sleep` (string) - The number of seconds to sleep between launching
  the incus instance and provisioning it; defaults to 3 seconds.

- `publish_properties` (map[string]string) - Pass key values to the publish
  step to be set as properties on the output image. This is most helpful to
  set the description, but can be used to set anything needed. See
  https://stgraber.org/2016/03/30/incus-2-0-image-management-512/
  for more properties.

- `launch_config` (map[string]string) - List of key/value pairs you wish to
  pass to lxc launch via --config. Defaults to empty.

- `virtual_machine` (bool) - Create incus virtual-machine image on hosts running incus 4.0 and above; defaults to false for container image

- `skip_publish` (bool) - Skip execute `lxc publish`; defaults to false

<!-- End of code generated from the comments of the Config struct in builder/incus/config.go; -->
