// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass

import (
	"errors"
	"os/exec"
)

var (
	ErrBinaryNotFound = errors.New("protonpass: pass-cli binary not found on $PATH")
)

type Client struct {
	bin string
}

type Option func(*Client)

func WithBinary(path string) Option {
	return func(c *Client) {
		c.bin = path
	}
}

func New(opts ...Option) (*Client, error) {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.bin == "" {
		bin, err := exec.LookPath("pass-cli")
		if err != nil {
			return nil, errors.Join(ErrBinaryNotFound, err)
		}
		c.bin = bin
	}

	return c, nil
}
