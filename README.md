Unofficial Terraform Provider for Atlassian's statuspage.io
===========================================================

What does it do ?
-----------------
Given an API key and a page ID, it allows to create
- components
- component groups

It implements a backoff retry procedure to overcome Atlassian's throtteling (1 request per second max).

Why?
----
Because Atlassian does not provide an official one : https://community.atlassian.com/t5/Statuspage-questions/Will-Atlassian-StatusPage-work-on-a-Terraform-Provider/qaq-p/965153


Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-aws`

```sh

git clone git@github.com:rayene/terraform-provider-statuspageio
cd terraform-provider-statuspageio
```

Enter the provider directory and build the provider

```sh
go build
```

Using the provider
----------------------
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it. Documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/aws/index.html).

Developing the Provider
-----------------------

I need help
- adding tests
- adding a CI/CD pipeline
- adding additional resources (incidents, ...)

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is *required*).

To compile the provider, run `go build`. This will build the provider.
