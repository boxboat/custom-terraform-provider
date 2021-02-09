provider "cmdb" {
  api_version = "v1"
  hostname = "localhost"
}

data "name_allocation" "vm_1_name" {
  provider = "cmdb"

  region = "us-east-1"
  resource_type = "COL"
}

data "name_details" "vm_1_details" {
  provider = "cmdb"

  name = data.name_allocation.vm_1_name.name
}

output "vm_1_name" {
  value = data.name_allocation.vm_1_name.name
}

output "vm_1_type" {
  value = data.name_details.vm_1_details.type
}

output "vm_1_details_raw" {
  value = data.name_details.vm_1_details.raw
}