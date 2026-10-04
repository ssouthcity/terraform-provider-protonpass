// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass

import (
	"os/exec"
)

func (c *Client) Login() error {
	cmd := exec.Command(c.bin, "login")

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Logout() error {
	cmd := exec.Command(c.bin, "logout")

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
