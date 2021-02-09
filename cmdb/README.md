# CMDB Emulation

This project is a simple HTTP server written in Golang for demonstration purposes as part of a blog post on Custom Terraform Providers. 

The purpose of this microservice is to emulate a subset of the actions that could be performed by a Content Management Database (CMDB), 
with the intention of calling the microservice from Terraform via a Custom Provider. A more built-out idea of this approach would be
highly beneficial in large organizations where infrastructure is spun up and down often via Terraform, and the desire is to effectively
catalog and track which resources are deployed where. 