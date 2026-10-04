// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type BoolOption bool

func (b *BoolOption) Option() string {
	if *b {
		return "true"
	}
	return "false"
}

type IntOption int

func (i *IntOption) Option() string {
	return fmt.Sprint(*i)
}

type GeneratePasswordOptions struct {
	Length           IntOption
	IncludeNumbers   BoolOption
	IncludeUppercase BoolOption
	IncludeSymbols   BoolOption
}

type generatePasswordResponse struct {
	Password string `json:"password"`
}

func (c *Client) GeneratePassword(settings GeneratePasswordOptions) (string, error) {
	cmd := exec.Command(c.bin, "password", "generate", "random",
		"--length", settings.Length.Option(),
		"--numbers", settings.IncludeNumbers.Option(),
		"--uppercase", settings.IncludeUppercase.Option(),
		"--symbols", settings.IncludeSymbols.Option(),
		"--output", "json")

	output := bytes.Buffer{}
	cmd.Stdout = &output

	if err := cmd.Run(); err != nil {
		return "", err
	}

	var resp generatePasswordResponse

	if err := json.Unmarshal(output.Bytes(), &resp); err != nil {
		return "", err
	}

	return resp.Password, nil
}
