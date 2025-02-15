# hetzner-golang-actions

Hetzner Golang action completion - the missing part of the Hetzner Golang SDK

<!-- toc -->

* [Usage](#usage)
* [Options](#options)
  * [WithIgnoreGetByIDError](#withignoregetbyiderror)
  * [WithTimeout](#withtimeout)
* [Contributing](#contributing)
  * [Open in a container](#open-in-a-container)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

I use Hetzner for lots of things. And it's great. But the Golang SDK is missing
a function that automatically waits for [Actions](https://docs.hetzner.cloud/#actions)
to be completed.

So I wrote this.

## Usage

> This uses the [HCloud v2 SDK](https://github.com/hetznercloud/hcloud-go). V1
> is no longer supported by Hetzner

Install the package

```bash
go get github.com/mrsimonemms/hetzner-golang-actions
```

```go
package main

import (
  "context"
  "time"

  "github.com/hetznercloud/hcloud-go/v2/hcloud"
  hga "github.com/mrsimonemms/hetzner-golang-actions"
)

func main() {
  ctx := context.Background()

  // Create the hcloud client
  client := hcloud.NewClient(hcloud.WithToken("xxxxx"))

  // Create a server
  server, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{})
  if err != nil {
    panic(err)
  }

  // Create the waiter and wait - the default timeout is one minute
  if err := hga.NewWaiter(client).
    Wait(context.Background(), server.Action, server.NextActions...); err != nil {
    // Waiter has failed - most likely a timeout
    panic(err)
  }

  // Once we're here, all the actions have completed and you can move on
}
```

## Options

### WithIgnoreGetByIDError

The [Action.GetByID](https://docs.hetzner.cloud/#actions-get-an-action) method
can return an error which may be a temporary network error.

To disable this check, use the `WithIgnoreGetByIDError()` option

```go
hga.NewWaiter(client, hga.WithIgnoreGetByIDError())
```

### WithTimeout

The default timeout is 1 minute and is reset for each action/nextAction passed in.

To change the timeout, use the `WithTimeout(time.Duration)` option

```go
hga.NewWaiter(client, hga.WithTimeout(time.Minute * 5))
```

## Contributing

### Open in a container

* [Open in a container](https://code.visualstudio.com/docs/devcontainers/containers)
