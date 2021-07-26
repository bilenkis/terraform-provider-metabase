# Terraform Provider Metabase

Run the following command to build the provider

```shell
go build -o terraform-provider-metabase
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Configure Metabase URL either using provider configuration directly or using ENV variables.
Get your username and password from Metabase and export them as `TF_VARs`.

```shell
export TF_VAR_metabase_username=user@example.com TF_VAR_metabase_password=xxxxxxxxxxxxx
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
cd examples && terraform init && terraform apply
```