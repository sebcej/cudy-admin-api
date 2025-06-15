# Cudy admin api

Simple library that allow to get data from Cudy routers administration panel

Tested on:
* Cudy LT700V (2.1.15 DE)

If you tested with success another router version/model please create a PR with updated list

## Initialization

```go
api := cudy.Init("192.168.0.1", "admin", "password") // if the interface is without username put "admin" as default

// Mandatory if maxRetries is 0
// api.Login()

data, err := api.DeviceStatus()
```

## Configuration

Extra configuration params are available after initialization if custom behavior is necessary.

```go
api.Sm.MaxRetries = 0 // Disable automatic session management
api.LoginMethod = LoginMethodUnencrypted // Disable password encryption for older routers
api.ExtraHeaders = map[string]string // Extra headers for each request
api.SecureConnection = false // https connection to admin area
api.Credentials.Timezone = "Europe/Rome" // Set custom timezone. Be sure is a value selectable from login interface
```

## Session management
The library will automatically manage the session and will login when necessary.

If retries are disabled you must manage the login manually when ErrUnknownError or ErrWrongCredentials are received.

## Available functions
`Status`
`GcomStats`

## Will be supported in next releases

`MessagesStats`
`MessagesInbox`
`SendMessage`
`ConnectedDevices`
