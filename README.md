# Cudy admin api

Simple library that allow to get data from Cudy routers administration panel

Tested on:
* Cudy LT700V (2.1.15 DE)

Some of the functionalities may be supported only for this version (4g) of the router.

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
api.Api.Sm.MaxRetries = 2 // Put to 0 if you need to disable automatic session management
api.Api.Sm.RetryWait = 1 * time.Second // Longer retry time is better for 500 error mitigation
api.LoginMethod = LoginMethodUnencrypted // Disable password encryption for older routers
api.Api.ExtraHeaders = map[string]string // Extra headers for each request
api.Api.SecureConnection = false // https connection to admin area
api.Timezone = "Europe/Rome" // Set custom timezone. Be sure is a value selectable from login interface
```

## Session management
The library will automatically manage the session and will login when necessary.

If retries are disabled you must manage the login manually when `ErrWrongCredentials` is received.

## Infos and quirks

When under heavy load, the router wil respond randomly with `500 Internal Server Error`. The retry mechanism should mitigate this, recommended extra failsafes.

The router does not perform any validation when sending SMS, the operation will fail silently and the message will not be sent altrough will be present in the interface. Validate your inputs.

### iface param

Here a list of interfaces available in the router. 

`4g`
`wlan00`
`wlan10`

`usb0`
`ra0`

## Ported functionalities
* `Status` - System status, fw version, uptime
* `GcomStats` - Get stats about external interface
* `ConnectedDevices` - Get list of connected devices with mac/ip/rx/tx/signal
* `SpeedStats` - Get current interface speed (usb0/ra0)
* `GetMessagesList` - Get list of messages in inbox. The messages will not have content but only a preview
    * `message.Content` - Fetch the contents of the selected message
    * `message.Respond` - Respond directly to message
    * `message.Delete` - Delete message
* `SendMessage` - Send message to phone number