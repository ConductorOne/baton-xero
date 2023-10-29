![Baton Logo](./docs/images/baton-logo.png)

# `baton-xero` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-xero.svg)](https://pkg.go.dev/github.com/conductorone/baton-xero) ![main ci](https://github.com/conductorone/baton-xero/actions/workflows/main.yaml/badge.svg)

`baton-xero` is a connector for Xero built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Xero Accounting API to sync data about organizations and their members. 

Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Prerequisites

To run the connector, you will need to create a Xero App that ensures the connection to the API, and obtain a client ID and client secret from it.

There are three types of apps that you can create and their authentication flow differs. The connector supports the [Client Credentials Flow](https://developer.xero.com/documentation/guides/oauth2/client-credentials) and [Refresh Token Flow](https://developer.xero.com/documentation/guides/oauth2/auth-flow/#refreshing-access-and-refresh-tokens). 

To use the Client Credentials Flow, you will need to create an app of type ["Custom connection"](https://developer.xero.com/documentation/guides/oauth2/custom-connections) and use the connector with client ID and client secret. There are multiple other prerequisities for this flow, so please read the documentation carefully.

To use the Refresh Token Flow, you will need to create an app of type "Web app" and use the connector with client ID, client secret and refresh token. This flow is part of the [OAuth 2.0 Authorization Code Flow](https://developer.xero.com/documentation/guides/oauth2/auth-flow) and it requires user interaction to obtain the refresh token. This refresh token, based on documentation, is valid for 60 days unless it is refreshed. 

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-xero

BATON_XERO_CLIENT_ID=xeroClientId BATON_XERO_CLIENT_SECRET=xeroClientSecret BATON_REFRESH_TOKEN=refreshToken baton-xero
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_XERO_CLIENT_ID=xeroClientId BATON_XERO_CLIENT_SECRET=xeroClientSecret ghcr.io/conductorone/baton-xero:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-xero/cmd/baton-xero@main

BATON_TOKEN=token baton-xero
baton resources
```

# Data Model

`baton-xero` will pull down information about the following resources from Accounting API:

- Organizations
- Users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-xero` Command Line Usage

```
baton-xero

Usage:
  baton-xero [flags]
  baton-xero [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string            The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string        The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                 The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                        help for baton-xero
      --log-format string           The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string            The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning                This must be set in order for provisioning actions to be enabled. ($BATON_PROVISIONING)
      --refresh-token string        The Xero refresh token used to exchange for a new access token. ($BATON_REFRESH_TOKEN)
      --token string                The Xero access token used to connect to the Xero API. ($BATON_TOKEN)
  -v, --version                     version for baton-xero
      --xero-client-id string       The Xero client ID used to connect to the Xero API. ($BATON_XERO_CLIENT_ID)
      --xero-client-secret string   The Xero client secret used to connect to the Xero API. ($BATON_XERO_CLIENT_SECRET)

Use "baton-xero [command] --help" for more information about a command.
```
