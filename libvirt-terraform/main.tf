terraform {
  required_version = ">= 1.0"

  required_providers {
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = "~> 0.9.0"
    }
  }

  # Local backend configuration
  backend "local" {
    path = "terraform.tfstate"
  }
}

provider "libvirt" {
  uri = var.libvirt_uri
}

# Create a storage pool
resource "libvirt_pool" "vm_pool" {
  name = var.storage_pool_name
  type = "dir"
  target = {
    path = var.storage_pool_path
  }
}

# Download base image
resource "libvirt_volume" "base_image" {
  name   = "${var.vm_name}-base.qcow2"
  pool   = libvirt_pool.vm_pool.name
  format = "qcow2"

  create = {
    content = {
      url = var.base_image_url
    }
  }
}

# Create volume for VM from base image
resource "libvirt_volume" "vm_disk" {
  name     = "${var.vm_name}-disk.qcow2"
  pool     = libvirt_pool.vm_pool.name
  format   = "qcow2"
  capacity = var.disk_size

  backing_store = {
    path   = libvirt_volume.base_image.path
    format = "qcow2"
  }
}

# Cloud-init configuration
data "template_file" "user_data" {
  template = file("${path.module}/cloud_init.cfg")

  vars = {
    hostname = var.vm_name
    fqdn     = "${var.vm_name}.${var.domain}"
    ssh_key  = var.ssh_public_key
    vm_user  = var.vm_user
  }
}

data "template_file" "meta_data" {
  template = file("${path.module}/meta_data.cfg")

  vars = {
    hostname = var.vm_name
  }
}

data "template_file" "network_config" {
  template = <<-EOF
    version: 2
    ethernets:
      eth0:
        dhcp4: true
  EOF
}

# Cloud-init disk for initial configuration
resource "libvirt_cloudinit_disk" "cloudinit" {
  name           = "${var.vm_name}-cloudinit.iso"
  user_data      = data.template_file.user_data.rendered
  meta_data      = data.template_file.meta_data.rendered
  network_config = data.template_file.network_config.rendered
}

# Upload the cloud-init ISO into the pool.
resource "libvirt_volume" "cloudinit_volume" {
  name = "${var.vm_name}-cloudinit.iso"
  pool = libvirt_pool.vm_pool.name

  create = {
    content = {
      url = libvirt_cloudinit_disk.cloudinit.path
    }
  }
}

# Define the VM domain
resource "libvirt_domain" "vm" {
  name   = var.vm_name
  memory = var.memory
  vcpu   = var.vcpu
  unit   = "MiB"
  type   = "kvm"

  os = {
    type         = "hvm"
    machine      = "q35"
    type_arch    = "x86_64"
    boot_devices = ["hd"]
    firmware     = "efi"
    loader       = "/usr/share/edk2/x64/OVMF_CODE.secboot.4m.fd"
    nvram = {
      undefine_on_destroy = true
    }
  }

  # Configure NVRAM to be removed on destroy


  # Enable ACPI for cloud images
  features = {
    acpi = true
    apic = true
  }


  devices = {
    # Disks
    disks = [
      {
        source = {
          pool   = libvirt_pool.vm_pool.name
          volume = libvirt_volume.vm_disk.name
        }
        target = {
          dev = "vda"
          bus = "virtio"
        }
      },
      {
        source = {
          pool   = libvirt_volume.cloudinit_volume.pool
          volume = libvirt_volume.cloudinit_volume.name
        }
        target = {
          dev = "sda"
          bus = "sata"
        }
        device = "cdrom"
      }
    ]

    # Network interfaces
    interfaces = [
      {
        type  = "network"
        model = "virtio"
        source = {
          network = var.network_name
        }
        # wait_for_ip = {
        #   source  = "lease"
        #   timeout = 300
        # }
      }
    ]

    # Console for debugging
    consoles = [
      {
        type        = "pty"
        target_type = "serial"
        target_port = 0
      }
    ]

    # Graphics for virt-manager access
    graphics = {
      vnc = {
        autoport = "yes"
        listen   = "127.0.0.1"
      }
    }
  }
}
