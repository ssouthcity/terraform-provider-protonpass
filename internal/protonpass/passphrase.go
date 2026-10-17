// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
)

type generatePassphraseConfig struct {
	separator  string
	capitalize bool
	numbers    bool
	count      int
}

var defaultGeneratePassphraseConfig = generatePassphraseConfig{
	separator:  "hyphens",
	capitalize: true,
	numbers:    true,
	count:      5,
}

type GeneratePassphraseOption interface {
	configurePassphraseConfig(*generatePassphraseConfig)
}

func (opt *SeparatorOption) configurePassphraseConfig(config *generatePassphraseConfig) {
	config.separator = opt.separator
}

func (opt *CapitalizeOption) configurePassphraseConfig(config *generatePassphraseConfig) {
	config.capitalize = opt.capitalize
}

func (opt *NumbersOption) configurePassphraseConfig(config *generatePassphraseConfig) {
	config.numbers = opt.numbers
}

func (opt *CountOption) configurePassphraseConfig(config *generatePassphraseConfig) {
	config.count = opt.count
}

type generatePassphraseResponse struct {
	Password string `json:"password"`
}

func (c *Client) GeneratePassphrase(opts ...GeneratePassphraseOption) (string, error) {
	config := defaultGeneratePassphraseConfig

	for _, opt := range opts {
		opt.configurePassphraseConfig(&config)
	}

	cmd := exec.Command(c.bin, "password", "generate", "passphrase",
		"--separator", config.separator,
		"--capitalize", strconv.FormatBool(config.capitalize),
		"--numbers", strconv.FormatBool(config.numbers),
		"--count", strconv.FormatInt(int64(config.count), 10),
		"--output", "json")

	output := bytes.Buffer{}
	cmd.Stdout = &output

	if err := cmd.Run(); err != nil {
		return "", err
	}

	var resp generatePassphraseResponse

	if err := json.Unmarshal(output.Bytes(), &resp); err != nil {
		return "", err
	}

	return resp.Password, nil
}
