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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j6s/terraform-provider-sweego-provider/internal/sweego"
)

var _ provider.Provider = &SweegoProvider{}
var _ provider.ProviderWithFunctions = &SweegoProvider{}
var _ provider.ProviderWithEphemeralResources = &SweegoProvider{}
var _ provider.ProviderWithActions = &SweegoProvider{}

// SweegoProvider defines the provider implementation.
type SweegoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SweegoProviderModel describes the provider data model.
type SweegoProviderModel struct {
	BaseUrl  types.String `tfsdk:"base_url"`
	ApiKey   types.String `tfsdk:"api_key"`
	ClientId types.String `tfsdk:"client_id"`
}

func (p *SweegoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sweego"
	resp.Version = p.version
}

func (p *SweegoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Base URL of the sweego API. Defaults to %s", sweego.DefaultBaseUrl),
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key used to authenticate the sweego API",
				Optional:            false,
				Required:            true,
				Sensitive:           true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID used to authenticate the sweego API",
				Optional:            false,
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SweegoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SweegoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	baseUrl := data.BaseUrl.ValueString()
	if baseUrl == "" {
		baseUrl = sweego.DefaultBaseUrl
	}
	client := sweego.NewSweegoApiWithBaseUrl(baseUrl, data.ApiKey.ValueString(), data.ClientId.ValueString())

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SweegoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSweegoDomainResource,
	}
}

func (p *SweegoProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *SweegoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *SweegoProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *SweegoProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SweegoProvider{
			version: version,
		}
	}
}
