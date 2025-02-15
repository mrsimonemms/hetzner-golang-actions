/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	hga "github.com/mrsimonemms/hetzner-golang-actions"
)

func main() {
	ctx := context.Background()

	// Create the hcloud client
	client := hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))

	// Create a server
	slog.Info("Create new server")
	server, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name: "example",
		ServerType: &hcloud.ServerType{
			Name: "cx22",
		},
		Image: &hcloud.Image{
			Name: "ubuntu-24.04",
		},
	})
	if err != nil {
		slog.Any("Error creating server", err)
		panic(err)
	}

	// Create the waiter and wait - the default timeout is one minute
	slog.Info("Waiting....")
	if err := hga.NewWaiter(client, hga.WithTimeout(time.Minute)).
		Wait(context.Background(), server.Action, server.NextActions...); err != nil {
		// Waiter has failed - most likely a timeout
		slog.Any("Error waiting", err)
		panic(err)
	}

	// Once we're here, all the actions have completed and you can move on
	slog.Info("All done successfully")
}
