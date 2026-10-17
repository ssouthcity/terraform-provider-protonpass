// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass_test

import "testing"

func TestGeneratePassword(t *testing.T) {
	client := newTestClient(t)

	password, err := client.GeneratePassword()
	if err != nil {
		t.Fatalf("Expected plain generate password to work, got error: %v", err)
	}

	if len(password) == 0 {
		t.Fatalf("Expected generated password to not be blank")
	}
}
