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

package hetznergolangactions

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ErrTimeout = fmt.Errorf("action timeout")

type WaitForAction struct {
	client             *hcloud.Client
	ignoreGetByIDError bool
	timeout            time.Duration
}

func (w *WaitForAction) waitForAction(ctx context.Context, action *hcloud.Action) error {
	startTime := time.Now()
	timeoutTime := startTime.Add(w.timeout)

	for {
		if time.Now().After(timeoutTime) {
			return ErrTimeout
		}

		time.Sleep(time.Second)

		status, _, err := w.client.Action.GetByID(ctx, action.ID)
		if err != nil {
			if w.ignoreGetByIDError {
				continue
			}
			return fmt.Errorf("error getting action: %w", err)
		}

		if status.Status == hcloud.ActionStatusError {
			return fmt.Errorf("error completing action - code: %s, message: %s", status.ErrorCode, status.ErrorMessage)
		}

		if status.Status == hcloud.ActionStatusSuccess {
			break
		}
	}

	return nil
}

func (w *WaitForAction) Wait(ctx context.Context, action *hcloud.Action, nextActions ...*hcloud.Action) error {
	actions := []*hcloud.Action{action}
	actions = append(actions, nextActions...)

	for _, action := range actions {
		if err := w.waitForAction(ctx, action); err != nil {
			return err
		}
	}

	return nil
}

func NewWaiter(client *hcloud.Client, opts ...WaitOption) *WaitForAction {
	wfa := &WaitForAction{
		client:             client,
		timeout:            time.Minute,
		ignoreGetByIDError: false,
	}

	for _, option := range opts {
		option(wfa)
	}

	return wfa
}

type WaitOption func(*WaitForAction)

// Ignore the HCloud Action.GetByID error - there may be a temporary network error
func WithIgnoreGetByIDError() WaitOption {
	return func(wfa *WaitForAction) {
		wfa.ignoreGetByIDError = true
	}
}

// Set timeout duration - default is 1 minutes
func WithTimeout(timeout time.Duration) WaitOption {
	return func(wfa *WaitForAction) {
		wfa.timeout = timeout
	}
}
