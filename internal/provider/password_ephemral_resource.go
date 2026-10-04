// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ssouthcity/terraform-provider-protonpass/internal/protonpass"
)

var _ ephemeral.EphemeralResource = &PasswordEphemeralResource{}

func NewPasswordEphemeralResource() ephemeral.EphemeralResource {
	return &PasswordEphemeralResource{}
}

type PasswordEphemeralResource struct {
	client *protonpass.Client
}

type PasswordEphemeralResourceModel struct {
	Length           types.Int32 `tfsdk:"length"`
	IncludeNumbers   types.Bool  `tfsdk:"include_numbers"`
	IncludeUppercase types.Bool  `tfsdk:"include_uppercase"`
	IncludeSymbols   types.Bool  `tfsdk:"include_symbols"`

	Value types.String `tfsdk:"value"`
}

func (r *PasswordEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

func (r *PasswordEphemeralResource) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Password ephemeral resource",

		Attributes: map[string]schema.Attribute{
			"length": schema.Int32Attribute{
				MarkdownDescription: "Amount of characters in the password",
				Optional:            true,
			},
			"include_numbers": schema.BoolAttribute{
				MarkdownDescription: "Whether to include numbers in the generated password",
				Optional:            true,
			},
			"include_uppercase": schema.BoolAttribute{
				MarkdownDescription: "Whether to include uppercase characters in the generated password",
				Optional:            true,
			},
			"include_symbols": schema.BoolAttribute{
				MarkdownDescription: "Whether to include symbols in the generated password",
				Optional:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Generated password",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *PasswordEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (r *PasswordEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data PasswordEphemeralResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Login(); err != nil {
		resp.Diagnostics.AddError("Proton Pass Error", fmt.Sprintf("Unable to login, got error: %s", err))
		return
	}

	password, err := r.client.GeneratePassword(protonpass.GeneratePasswordOptions{
		Length:           protonpass.IntOption(int(data.Length.ValueInt32())),
		IncludeNumbers:   protonpass.BoolOption(data.IncludeNumbers.ValueBool()),
		IncludeUppercase: protonpass.BoolOption(data.IncludeUppercase.ValueBool()),
		IncludeSymbols:   protonpass.BoolOption(data.IncludeSymbols.ValueBool()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Proton Pass Error", fmt.Sprintf("Unable to generate password, got error: %s", err))
		return
	}

	data.Value = types.StringValue(password)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
