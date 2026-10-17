// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass_test

import (
	"strings"
	"testing"

	"github.com/ssouthcity/terraform-provider-protonpass/internal/protonpass"
)

func TestGeneratePassphrase_Separator(t *testing.T) {
	client := newTestClient(t)

	validSeparators := []string{
		protonpass.SeparatorHyphens,
		protonpass.SeparatorSpaces,
		protonpass.SeparatorPeriods,
		protonpass.SeparatorCommas,
		protonpass.SeparatorUnderscores,
		protonpass.SeparatorNumbers,
		protonpass.SeparatorsNumbersAndSymbols,
	}

	for _, sep := range validSeparators {
		_, err := client.GeneratePassphrase(protonpass.Separator(sep))
		if err != nil {
			t.Fatalf("Failed to generate passphrase with separator %s, got error: %v", sep, err)
		}
	}

	invalidSeparators := []string{
		"quantum",
	}

	for _, sep := range invalidSeparators {
		passphrase, err := client.GeneratePassphrase(protonpass.Separator(sep))
		if err == nil {
			t.Fatalf("Expected error on invalid passphrase separator, got passphrase %s", passphrase)
		}
	}
}

func TestGeneratePassphrase_Count(t *testing.T) {
	client := newTestClient(t)

	counts := []int{1, 3, 5, 10, 15, 20, 25, 30, 1000}

	for _, c := range counts {
		passphrase, err := client.GeneratePassphrase(protonpass.Count(c))
		if err != nil {
			t.Fatalf("Expected passphrase with count %d to generate, got error %v", c, err)
		}

		wordC := strings.Count(passphrase, "-") + 1

		if wordC != c {
			t.Fatalf("Expected words in passphrase to equal %d, got %d", c, wordC)
		}
	}
}
