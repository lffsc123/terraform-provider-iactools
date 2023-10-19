package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ScaffoldingProvider{}

// ScaffoldingProvider defines the provider implementation.
type ScaffoldingProvider struct {
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Port     types.String `tfsdk:"port"`
	Address  types.String `tfsdk:"address"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dptech-demo"
	tflog.Info(ctx, "Metadata**********")
}

func (p *ScaffoldingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	tflog.Info(ctx, "Schema******** *****")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"port": schema.StringAttribute{
				MarkdownDescription: "  provider attribute",
				Optional:            true,
			}, "address": schema.StringAttribute{
				MarkdownDescription: "  provider attribute",
				Optional:            true,
			}, "username": schema.StringAttribute{
				MarkdownDescription: "  provider attribute",
				Optional:            true,
			}, "password": schema.StringAttribute{
				MarkdownDescription: "  provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *ScaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	tflog.Info(ctx, "Configure****+*******")
	if resp.Diagnostics.HasError() {
		return
	}
	address, port, name, password := "", "", "", ""
	if !data.Address.IsNull() {
		address = data.Address.ValueString()
	}
	if !data.Port.IsNull() {
		port = data.Port.ValueString()
	}
	if !data.Username.IsNull() {
		name = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if data.Port.IsNull() {
		tflog.Info(ctx, "Port is NULL")
		return
	}
	if data.Address.IsNull() {
		tflog.Info(ctx, "Address is NULL")
		return
	}
	if data.Username.IsNull() {
		tflog.Info(ctx, "Username is NULL")
		return
	}
	if data.Password.IsNull() {
		tflog.Info(ctx, "Password is NULL")
		return
	}
	tflog.Info(ctx, address+port+name+password)
	address = address + port

	autha := AuthStruct{Username: name, Password: password}

	client, err := NewClient(&address, &autha)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create HashiCups API Client",
			"An unexpected error occurred when creating the HashiCups API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"HashiCups Client Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRealServiceResource,
		NewRealServiceListResource,
		NewAddrPoolResource,
		NewVirtualServiceResource,
	}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
