// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ssouthcity/terraform-provider-protonpass/internal/protonpass"
)

var _ provider.Provider = &ProtonPassProvider{}
var _ provider.ProviderWithFunctions = &ProtonPassProvider{}
var _ provider.ProviderWithEphemeralResources = &ProtonPassProvider{}
var _ provider.ProviderWithActions = &ProtonPassProvider{}

type ProtonPassProvider struct {
	version string
}

type ProtonPassProviderModel struct {
}

func (p *ProtonPassProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "protonpass"
	resp.Version = p.version
}

func (p *ProtonPassProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *ProtonPassProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProtonPassProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := protonpass.New()
	if err != nil {
		resp.Diagnostics.AddError("Proton Pass Error", fmt.Sprintf("Unable to configure provider, got error: %s", err))
		return
	}

	resp.EphemeralResourceData = client
}

func (p *ProtonPassProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ProtonPassProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewPasswordEphemeralResource,
	}
}

func (p *ProtonPassProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ProtonPassProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *ProtonPassProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ProtonPassProvider{
			version: version,
		}
	}
}
