// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protonpass

type LengthOption struct {
	length int
}

func Length(length int) *LengthOption {
	return &LengthOption{length}
}

type NumbersOption struct {
	numbers bool
}

func Numbers(includeNumbers bool) *NumbersOption {
	return &NumbersOption{includeNumbers}
}

type SymbolsOption struct {
	symbols bool
}

func Symbols(includeNumbers bool) *SymbolsOption {
	return &SymbolsOption{includeNumbers}
}

type UppercaseOption struct {
	uppercase bool
}

func Uppercase(includeNumbers bool) *UppercaseOption {
	return &UppercaseOption{includeNumbers}
}

type CapitalizeOption struct {
	capitalize bool
}

func Capitalize(includeNumbers bool) *CapitalizeOption {
	return &CapitalizeOption{includeNumbers}
}

const (
	SeparatorHyphens            = "hyphens"
	SeparatorSpaces             = "spaces"
	SeparatorPeriods            = "periods"
	SeparatorCommas             = "commas"
	SeparatorUnderscores        = "underscores"
	SeparatorNumbers            = "numbers"
	SeparatorsNumbersAndSymbols = "numbers-and-symbols"
)

type SeparatorOption struct {
	separator string
}

func Separator(separator string) *SeparatorOption {
	return &SeparatorOption{separator}
}

type CountOption struct {
	count int
}

func Count(length int) *CountOption {
	return &CountOption{length}
}
