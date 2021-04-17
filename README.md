[![Go Build & Deploy](https://github.com/IceToast/zoom_schedule_backend_go/actions/workflows/go_deploy.yml/badge.svg?branch=master)](https://github.com/IceToast/zoom_schedule_backend_go/actions/workflows/go_deploy.yml)
[![code style: prettier](https://img.shields.io/badge/code_style-prettier-ff69b4.svg)](https://github.com/prettier/prettier)


# Zoom Schedule Backend

This is the backend for the [Zoom Schedule App](https://github.com/IceToast/zoom_schedule)

To host this App yourself (without editing much) you have to change the "Host" constant, to your desired OAuth Callback Domain, in the main.go and provide some environment variables (for Example in a .env-file) like the following:
- MongoDB Connection string ("CONNECTION_STRING=")
- OAuth Provider IDs ("PROVIDER_CLIENT_ID=")
- OAuth Provider Secrets ("PROVIDER_SECRET=")


## Resolve missing dependencies
```
go get
```

## Build the backend

```
go build
```

## Run the backend

```
go run
```
