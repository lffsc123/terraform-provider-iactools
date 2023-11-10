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

// Ensure DpProvider satisfies various provider interfaces.
var _ provider.Provider = &DpProvider{}

// DpProvider defines the provider implementation.
type DpProvider struct {
	version string
}

//type AddSourceNatProviderModel struct {
//	Port                types.String `tfsdk:"port"`
//	Ip                  types.String `tfsdk:"ip"`
//	IpVersion           types.String `tfsdk:"ipVersion"`
//	VsysName            types.String `tfsdk:"vsysName"`
//	Name                types.String `tfsdk:"name"`
//	TargetName          types.String `tfsdk:"targetName"`
//	Position            types.String `tfsdk:"position"`
//	OutInterface        types.String `tfsdk:"outInterface"`
//	SrcAddrObj          types.String `tfsdk:"srcAddrObj"`
//	SrcAddrGroup        types.String `tfsdk:"srcAddrGroup"`
//	DstAddrObj          types.String `tfsdk:"dstAddrObj"`
//	PreService          types.String `tfsdk:"preService"`
//	UsrService          types.String `tfsdk:"usrService"`
//	ServiceGroup        types.String `tfsdk:"serviceGroup"`
//	PublicIpAddressFlag types.String `tfsdk:"publicIpAddressFlag"`
//	AddrpoolName        types.String `tfsdk:"addrpoolName"`
//	MinPort             types.String `tfsdk:"minPort"`
//	MaxPort             types.String `tfsdk:"maxPort"`
//	PortHash            types.String `tfsdk:"portHash"`
//	State               types.String `tfsdk:"state"`
//}

type UpdateSourceNatProviderModel struct {
	Port                types.String `tfsdk:"port"`
	Ip                  types.String `tfsdk:"ip"`
	IpVersion           types.String `tfsdk:"ipVersion"`
	VsysName            types.String `tfsdk:"vsysName"`
	OldName             types.String `tfsdk:"oldName"`
	TargetName          types.String `tfsdk:"targetName"`
	Position            types.String `tfsdk:"position"`
	OutInterface        types.String `tfsdk:"outInterface"`
	SrcAddrObj          types.String `tfsdk:"srcAddrObj"`
	SrcAddrGroup        types.String `tfsdk:"srcAddrGroup"`
	DstAddrObj          types.String `tfsdk:"dstAddrObj"`
	DstAddrGroup        types.String `tfsdk:"dstAddrGroup"`
	PreService          types.String `tfsdk:"preService"`
	UsrService          types.String `tfsdk:"usrService"`
	ServiceGroup        types.String `tfsdk:"serviceGroup"`
	PublicIpAddressFlag types.String `tfsdk:"publicIpAddressFlag"`
	AddrpoolName        types.String `tfsdk:"addrpoolName"`
	MinPort             types.String `tfsdk:"minPort"`
	MaxPort             types.String `tfsdk:"maxPort"`
	PortHash            types.String `tfsdk:"portHash"`
	State               types.String `tfsdk:"state"`
}

type DelSourceNatProviderModel struct {
	Port         types.String `tfsdk:"port"`
	Ip           types.String `tfsdk:"ip"`
	IpVersion    types.String `tfsdk:"ipVersion"`
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	DelallEnable types.String `tfsdk:"delallEnable"`
}

type BatchGetDeviceSourceNatProviderModel struct {
	Port                types.String `tfsdk:"port"`
	Ip                  types.String `tfsdk:"ip"`
	IpVersion           types.String `tfsdk:"ipVersion"`
	VsysName            types.String `tfsdk:"vsysName"`
	Offset              types.String `tfsdk:"offset"`
	Count               types.String `tfsdk:"count"`
	Name                types.String `tfsdk:"name"`
	OutInterface        types.String `tfsdk:"outInterface"`
	SourceIp            types.String `tfsdk:"sourceIp"`
	DestinationIp       types.String `tfsdk:"destinationIp"`
	Service             types.String `tfsdk:"service"`
	State               types.String `tfsdk:"state"`
	SrcAddrObj          types.String `tfsdk:"srcAddrObj"`
	SrcAddrGroup        types.String `tfsdk:"srcAddrGroup"`
	DstAddrObj          types.String `tfsdk:"dstAddrObj"`
	DstAddrGroup        types.String `tfsdk:"dstAddrGroup"`
	PreService          types.String `tfsdk:"preService"`
	UsrService          types.String `tfsdk:"usrService"`
	ServiceGroup        types.String `tfsdk:"serviceGroup"`
	PublicIpAddressFlag types.String `tfsdk:"publicIpAddressFlag"`
	AddrpoolName        types.String `tfsdk:"addrpoolName"`
	MinPort             types.String `tfsdk:"minPort"`
	MaxPort             types.String `tfsdk:"maxPort"`
	PortHash            types.String `tfsdk:"portHash"`
	RuleId              types.String `tfsdk:"ruleId"`
	DelallEnable        types.String `tfsdk:"delallEnable"`
	TargetName          types.String `tfsdk:"targetName"`
	OldName             types.String `tfsdk:"oldName"`
	Position            types.String `tfsdk:"position"`
}

type AddTargetNatProviderModel struct {
	Port                 types.String `tfsdk:"port"`
	Ip                   types.String `tfsdk:"ip"`
	VsysName             types.String `tfsdk:"vsysName"`
	Name                 types.String `tfsdk:"name"`
	TargetName           types.String `tfsdk:"targetName"`
	Position             types.String `tfsdk:"position"`
	InInterface          types.String `tfsdk:"inInterface"`
	SrcIpObj             types.String `tfsdk:"srcIpObj"`
	SrcIpGroup           types.String `tfsdk:"srcIpGroup"`
	PublicIp             types.String `tfsdk:"publicIp"`
	PreService           types.String `tfsdk:"preService"`
	UsrService           types.String `tfsdk:"usrService"`
	InNetIp              types.String `tfsdk:"inNetIp"`
	InnetPort            types.String `tfsdk:"innetPort"`
	UnLimited            types.String `tfsdk:"unLimited"`
	SrcIpTranslate       types.String `tfsdk:"srcIpTranslate"`
	InterfaceAddressFlag types.String `tfsdk:"interfaceAddressFlag"`
	AddrpoolName         types.String `tfsdk:"addrpoolName"`
	VrrpIfName           types.String `tfsdk:"vrrpIfName"`
	VrrpId               types.String `tfsdk:"vrrpId"`
	State                types.String `tfsdk:"state"`
}

type UpdateTargetNatProviderModel struct {
	Port                 types.String `tfsdk:"port"`
	Ip                   types.String `tfsdk:"ip"`
	VsysName             types.String `tfsdk:"vsysName"`
	OldName              types.String `tfsdk:"oldName"`
	TargetName           types.String `tfsdk:"targetName"`
	Position             types.String `tfsdk:"position"`
	InInterface          types.String `tfsdk:"inInterface"`
	NetaddrObj           types.String `tfsdk:"netaddrObj"`
	NetaddrGroup         types.String `tfsdk:"netaddrGroup"`
	PublicIp             types.String `tfsdk:"publicIp"`
	PreService           types.String `tfsdk:"preService"`
	UsrService           types.String `tfsdk:"usrService"`
	InNetIp              types.String `tfsdk:"inNetIp"`
	InnetPort            types.String `tfsdk:"innetPort"`
	UnLimited            types.String `tfsdk:"unLimited"`
	SrcIpTranslate       types.String `tfsdk:"srcIpTranslate"`
	InterfaceAddressFlag types.String `tfsdk:"interfaceAddressFlag"`
	AddrpoolName         types.String `tfsdk:"addrpoolName"`
	VrrpIfName           types.String `tfsdk:"vrrpIfName"`
	VrrpId               types.String `tfsdk:"vrrpId"`
	State                types.String `tfsdk:"state"`
}

type DelTargetNatProviderModel struct {
	Port         types.String `tfsdk:"port"`
	Ip           types.String `tfsdk:"ip"`
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	DelallEnable types.String `tfsdk:"delallEnable"`
}

type BatchGetDeviceTargetNatProviderModel struct {
	Port                 types.String `tfsdk:"port"`
	Ip                   types.String `tfsdk:"ip"`
	VsysName             types.String `tfsdk:"vsysName"`
	Offset               types.String `tfsdk:"offset"`
	Count                types.String `tfsdk:"count"`
	Name                 types.String `tfsdk:"name"`
	SearchValue          types.String `tfsdk:"searchValue"`
	InInterface          types.String `tfsdk:"inInterface"`
	SourceIp             types.String `tfsdk:"sourceIp"`
	PublicIp             types.String `tfsdk:"publicIp"`
	Protocol             types.String `tfsdk:"protocol"`
	InNetIp              types.String `tfsdk:"inNetIp"`
	State                types.String `tfsdk:"state"`
	SrcIpObj             types.String `tfsdk:"srcIpObj"`
	SrcIpGroup           types.String `tfsdk:"srcIpGroup"`
	PreService           types.String `tfsdk:"preService"`
	UsrService           types.String `tfsdk:"usrService"`
	InnetPort            types.String `tfsdk:"innetPort"`
	UnLimited            types.String `tfsdk:"unLimited"`
	SrcIpTranslate       types.String `tfsdk:"srcIpTranslate"`
	InterfaceAddressFlag types.String `tfsdk:"interfaceAddressFlag"`
	AddrpoolName         types.String `tfsdk:"addrpoolName"`
	VrrpIfName           types.String `tfsdk:"vrrpIfName"`
	VrrpId               types.String `tfsdk:"vrrpId"`
	RuleId               types.String `tfsdk:"ruleId"`
	DelallEnable         types.String `tfsdk:"delallEnable"`
	TargetName           types.String `tfsdk:"targetName"`
	OldName              types.String `tfsdk:"oldName"`
	Position             types.String `tfsdk:"position"`
}

func (p *DpProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dpsc_lff"
	tflog.Info(ctx, "Metadata**********")
}

// 配置选项验证: Schema 允许你明确定义支持的配置选项以及它们的类型、默认值、是否必填等信息。
//
//	当 Terraform 用户在配置文件中定义资源或数据源时，Schema 会验证这些配置选项是否符合插件的预期格式
func (p *DpProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
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

// 该方法用于设置插件的配置，包括认证信息、连接参数以及其他插件特定的配置选项
func (p *DpProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configure****+*******")
}

func (p *DpProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSourceNatResource,
		NewTargetNatResource,
		NewIpv4RouterResource,
		NewIpv4StrategyRouterResource,
	}
}

func (p *DpProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewExampleDataSource,
	}
}

func NewDp(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DpProvider{
			version: version,
		}
	}
}
