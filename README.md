# Terraform Provider Metabase

This is the Metabase provider for Terraform.

The provider manages the databases and (not yet implemented) other resources in Metabase through Terraform.

## Test sample configuration


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Configure Metabase URL either using provider configuration directly or using ENV variables.
Get your username and password from Metabase and export them as `TF_VARs`.

```shell
export TF_VAR_metabase_username=user@example.com TF_VAR_metabase_password=xxxxxxxxxxxxx
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd examples && terraform init && terraform apply
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
go install
```

To generate or update documentation, run `go generate`.

To test your code run `make testacc`.

```sh
$ make testacc
```
