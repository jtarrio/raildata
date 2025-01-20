Package `raildata` is a library to access NJ Transit's RailData API.

It provides automatic API token management and a higher-level, enriched API.

# Library usage

```go
import "github.com/jtarrio/raildata"
```

# CLI tool

This repository also includes a command-line utility that uses the `raildata` library to download
and display NJ Transit information. For more information, execute this:

```shell
go run github.com/jtarrio/raildata/raildata-cli
```

# API access

In order to use this library, you need to visit https://developer.njtransit.com/registration/login
to request your NJ Transit developer API credentials for the RailData API.

# Example

```go
// Read the token from a file
token, err := io.ReadFile("/path/to/token-file")
if err != nil { return err }
client, err := raildata.NewClient(
    // Provide the token to the client so it doesn't need to get a new one
    raildata.WithToken(string(token)),
    // Provide the username and password to the client so it can get a new token if the old one expires
    raildata.WithCredentials(username, password),
    // If the token changes, save it to the file so we can use it in the future
    raildata.WithTokenUpdateListener(func(newToken string, oldToken string) {
        _ = io.WriteFile("/path/to/token-file", []byte(newToken))
    }),
)
if err != nil { return err }

vehicles, err := client.GetVehicleData(context.Background())
if err != nil { return err }
for _, vehicle := vehicles.Vehicles {
    println(vehicle.TrainId)
}
```

# Caveats

This client is not provided by NJ Transit. It may fail to parse some messages as we are still discovering
which fields are optional and which aren't. Use at your own risk.
