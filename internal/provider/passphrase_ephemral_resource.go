// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ssouthcity/terraform-provider-protonpass/internal/protonpass"
)

var _ ephemeral.EphemeralResource = &PassphraseEphemeralResource{}

func NewPassphraseEphemeralResource() ephemeral.EphemeralResource {
	return &PassphraseEphemeralResource{}
}

type PassphraseEphemeralResource struct {
	client *protonpass.Client
}

type PassphraseEphemeralResourceModel struct {
	Separator  types.String `tfsdk:"separator"`
	Capitalize types.Bool   `tfsdk:"capitalize"`
	Numbers    types.Bool   `tfsdk:"numbers"`
	WordCount  types.Int32  `tfsdk:"word_count"`

	Value types.String `tfsdk:"value"`
}

func (r *PassphraseEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passphrase"
}

func (r *PassphraseEphemeralResource) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Passphrase ephemeral resource",

		Attributes: map[string]schema.Attribute{
			"separator": schema.StringAttribute{
				MarkdownDescription: "Which separator to use",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						protonpass.SeparatorHyphens,
						protonpass.SeparatorSpaces,
						protonpass.SeparatorPeriods,
						protonpass.SeparatorCommas,
						protonpass.SeparatorUnderscores,
						protonpass.SeparatorNumbers,
						protonpass.SeparatorsNumbersAndSymbols,
					),
				},
			},
			"capitalize": schema.BoolAttribute{
				MarkdownDescription: "Whether to capitalize each word in the passphrase",
				Optional:            true,
			},
			"numbers": schema.BoolAttribute{
				MarkdownDescription: "Whether to include numbers in the generated passphrase",
				Optional:            true,
			},
			"word_count": schema.Int32Attribute{
				MarkdownDescription: "Amount of words in the passphrase",
				Optional:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Generated passphrase",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *PassphraseEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*protonpass.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Ephemeral Resource Configure Type",
			fmt.Sprintf("Expected *protonpass.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PassphraseEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data PassphraseEphemeralResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	options := []protonpass.GeneratePassphraseOption{}

	if separator := data.Separator.ValueStringPointer(); separator != nil {
		options = append(options, protonpass.Separator(*separator))
	}

	if capitalize := data.Capitalize.ValueBoolPointer(); capitalize != nil {
		options = append(options, protonpass.Capitalize(*capitalize))
	}

	if numbers := data.Numbers.ValueBoolPointer(); numbers != nil {
		options = append(options, protonpass.Numbers(*numbers))
	}

	if wordCount := data.WordCount.ValueInt32Pointer(); wordCount != nil {
		options = append(options, protonpass.Count(int(*wordCount)))
	}

	password, err := r.client.GeneratePassphrase(options...)
	if err != nil {
		resp.Diagnostics.AddError("Proton Pass Error", fmt.Sprintf("Unable to generate password, got error: %s", err))
		return
	}

	data.Value = types.StringValue(password)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
