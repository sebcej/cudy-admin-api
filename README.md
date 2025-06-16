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
api.Api.Sm.MaxRetries = 0 // Disable automatic session management
api.Api.Sm.RetryWait = 1 * time.Second
api.LoginMethod = LoginMethodUnencrypted // Disable password encryption for older routers
api.Api.ExtraHeaders = map[string]string // Extra headers for each request
api.Api.SecureConnection = false // https connection to admin area
api.Timezone = "Europe/Rome" // Set custom timezone. Be sure is a value selectable from login interface
```

__FYI:__ When under heavy load, the router wil respond randomly with `500 Internal Server Error`. The retry mechanism should mitigate this but results are not guaranteed

## Session management
The library will automatically manage the session and will login when necessary.

If retries are disabled you must manage the login manually when ErrWrongCredentials is received.

## Available functions
`Status`
`GcomStats`
`ConnectedDevices`
`SpeedStats`
`GetMessagesList`

## Will be supported in next releases

`SendMessage`