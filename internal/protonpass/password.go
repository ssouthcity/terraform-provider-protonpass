// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
)

type generatePasswordConfig struct {
	length           int
	includeNumbers   bool
	includeUppercase bool
	includeSymbols   bool
}

var defaultGeneratePasswordConfig = generatePasswordConfig{
	length:           32,
	includeNumbers:   true,
	includeSymbols:   true,
	includeUppercase: true,
}

type GeneratePasswordOption interface {
	configurePasswordConfig(*generatePasswordConfig)
}

func (opt *LengthOption) configurePasswordConfig(config *generatePasswordConfig) {
	config.length = opt.length
}

func (opt *NumbersOption) configurePasswordConfig(config *generatePasswordConfig) {
	config.includeNumbers = opt.numbers
}

func (opt *SymbolsOption) configurePasswordConfig(config *generatePasswordConfig) {
	config.includeSymbols = opt.symbols
}

func (opt *UppercaseOption) configurePasswordConfig(config *generatePasswordConfig) {
	config.includeUppercase = opt.uppercase
}

type generatePasswordResponse struct {
	Password string `json:"password"`
}

func (c *Client) GeneratePassword(opts ...GeneratePasswordOption) (string, error) {
	config := defaultGeneratePasswordConfig

	for _, opt := range opts {
		opt.configurePasswordConfig(&config)
	}

	cmd := exec.Command(c.bin, "password", "generate", "random",
		"--length", strconv.FormatInt(int64(config.length), 10),
		"--numbers", strconv.FormatBool(config.includeNumbers),
		"--uppercase", strconv.FormatBool(config.includeUppercase),
		"--symbols", strconv.FormatBool(config.includeSymbols),
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
