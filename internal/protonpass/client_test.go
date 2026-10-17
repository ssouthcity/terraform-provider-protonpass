// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass_test

import (
	"testing"

	"github.com/ssouthcity/terraform-provider-protonpass/internal/protonpass"
)

func newTestClient(t *testing.T) *protonpass.Client {
	t.Helper()
	client, err := protonpass.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client, got error: %v", err)
	}
	return client
}
