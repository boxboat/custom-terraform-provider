# Terraform Provider CMDB

This folder encapsulates the Terraform Provider that issues API calls to the CMDB microservice.

## Disclaimer

This project is meant to be for demonstration purposes only. There are many things that should be improved upon within
the code before this project would be considered production ready. 

## Running the example

To run the Terraform Provider locally there are a few steps to complete:

Step 1: Build the source code locally

```
go build -o terraform-provider-cmdb_v1.0.0
```

Step 2: Move the executable into the local terraform plugin folder:

```
mv terraform-provider-cmdb_v1.0.0 ~/terraform.d/plugins/linux_amd64/
```

> Note: The plugin folder may need to be created.

Step 3: From within this directory, initialize Terraform:

```
terraform init
```

Step 4: Run an apply via Terraform:

```
terraform apply
```

The output generated should look similar to the following:

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

vm_1_details_raw = {"Name":"M2540TCOLRus","Type":"COL","Region":"us-east-1"}
vm_1_name = M2540TCOLRus
vm_1_type = COL
```