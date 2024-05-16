packer {
  required_plugins {
    incus = {
      version = ">= 1.0.0"
      source  = "github.com/bketelsen/incus"
    }
  }
}

source "incus" "jammy" {
  image        = "images:ubuntu/jammy"
  output_image = "ubuntu-jammy"
  reuse        = true
}

build {
  sources = ["incus.jammy"]
  provisioner "file" {
    source = "examples/common/90-incus"
    destination = "/tmp/90-incus"
  }

  provisioner "shell" {
    scripts = [
      "examples/common/packages.sh",
      "examples/common/user.sh",
      "examples/common/sudoers.sh",
    ]
  }

}

