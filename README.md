## EndPoint Monitor - Terraform Provider

### Build

Run the below in the project root directory:
`go build -o terraform-provider-endpointmonitor`


### Documentation

We use a Terraform Docs plugin to generate the provider documentation.

https://github.com/hashicorp/terraform-plugin-docs

Once downloaded and compiled, the below command can be used to generate the documentation.

`tfplugindocs generate --provider-name endpointmonitor  --rendered-provider-name EndPointMonitor`

We have the main provider page also set up in templates/index.md.tmpl too as not all of the 
information in there can be generated. Should bear this in mind when modifying the provider 
as some information in that template may become outdated.