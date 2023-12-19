package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// IPv4策略路由
var _ resource.Resource = &Ipv4StrategyRouterResource{}
var _ resource.ResourceWithImportState = &Ipv4StrategyRouterResource{}

func NewIpv4StrategyRouterResource() resource.Resource {
	return &Ipv4StrategyRouterResource{}
}

type Ipv4StrategyRouterResource struct {
	client *Client
}

type Ipv4StrategyRouterResourceModel struct {
	AddIpv4StrategyRouterParameter AddIpv4StrategyRouterParameter `tfsdk:"rtplist"`
}

type AddIpv4StrategyRouterRequest struct {
	AddIpv4StrategyRouterRequestModel AddIpv4StrategyRouterRequestModel `json:"rtplist"`
}

type UpdateIpv4StrategyRouterRequest struct {
	UpdateIpv4StrategyRouterRequestModel UpdateIpv4StrategyRouterRequestModel `json:"rtplist"`
}

// 调用接口参数
type AddIpv4StrategyRouterRequestModel struct {
	VsysName      string `json:"vsysName"`
	Resource      string `json:"resource"`
	ListFlag      string `json:"listFlag"`
	Sequence      string `json:"sequence"`
	RtpName       string `json:"rtpName"`
	MatchSrctype  string `json:"matchSrctype"`
	MatchDesttype string `json:"matchDesttype"`
	Src           string `json:"src"`
	Dst           string `json:"dst"`
	IIfName       string `json:"iIfName"`
	Protocol      string `json:"protocol"`
	MinSrcport    string `json:"minSrcport"`
	MaxSrcport    string `json:"maxSrcport"`
	MinDstport    string `json:"minDstport"`
	MaxDstport    string `json:"maxDstport"`
	Dscp          string `json:"dscp"`
	Act           string `json:"act"`
	Timestr       string `json:"timestr"`
	Status        string `json:"status"`
	OifName       string `json:"oifName"`
	Weight        string `json:"weight"`
	Nexthop       string `json:"nexthop"`
	HcType        string `json:"hcType"`
	HcName        string `json:"hcName"`
	BfdCheck      string `json:"bfdCheck"`
	NexthopCount  string `json:"nexthopCount"`
}

type UpdateIpv4StrategyRouterRequestModel struct {
	VsysName      string `json:"vsysName"`
	Resource      string `json:"resource"`
	ListFlag      string `json:"listFlag"`
	Sequence      string `json:"sequence"`
	RtpName       string `json:"rtpName"`
	MatchSrctype  string `json:"matchSrctype"`
	MatchDesttype string `json:"matchDesttype"`
	Src           string `json:"src"`
	Dst           string `json:"dst"`
	IIfName       string `json:"iIfName"`
	Protocol      string `json:"protocol"`
	MinSrcport    string `json:"minSrcport"`
	MaxSrcport    string `json:"maxSrcport"`
	MinDstport    string `json:"minDstport"`
	MaxDstport    string `json:"maxDstport"`
	Dscp          string `json:"dscp"`
	Act           string `json:"act"`
	Timestr       string `json:"timestr"`
	Status        string `json:"status"`
	OifName       string `json:"oifName"`
	Weight        string `json:"weight"`
	Nexthop       string `json:"nexthop"`
	HcType        string `json:"hcType"`
	HcName        string `json:"hcName"`
	BfdCheck      string `json:"bfdCheck"`
	NexthopCount  string `json:"nexthopCount"`
}

// 接收外部参数
type AddIpv4StrategyRouterParameter struct {
	VsysName      types.String `tfsdk:"vsysname"`
	Resource      types.String `tfsdk:"resource"`
	ListFlag      types.String `tfsdk:"listflag"`
	Sequence      types.String `tfsdk:"sequence"`
	RtpName       types.String `tfsdk:"rtpname"`
	MatchSrctype  types.String `tfsdk:"matchsrctype"`
	MatchDesttype types.String `tfsdk:"matchdesttype"`
	Src           types.String `tfsdk:"src"`
	Dst           types.String `tfsdk:"dst"`
	IIfName       types.String `tfsdk:"iifname"`
	Protocol      types.String `tfsdk:"protocol"`
	MinSrcport    types.String `tfsdk:"minsrcport"`
	MaxSrcport    types.String `tfsdk:"maxsrcport"`
	MinDstport    types.String `tfsdk:"mindstport"`
	MaxDstport    types.String `tfsdk:"maxdstport"`
	Dscp          types.String `tfsdk:"dscp"`
	Act           types.String `tfsdk:"act"`
	Timestr       types.String `tfsdk:"timestr"`
	Status        types.String `tfsdk:"status"`
	OifName       types.String `tfsdk:"oifname"`
	Weight        types.String `tfsdk:"weight"`
	Nexthop       types.String `tfsdk:"nexthop"`
	HcType        types.String `tfsdk:"hctype"`
	HcName        types.String `tfsdk:"hcname"`
	BfdCheck      types.String `tfsdk:"bfdcheck"`
	NexthopCount  types.String `tfsdk:"nexthopcount"`
}

// 查询结果结构体
type QueryIpv4StrategyRouterResponseModel struct {
	VsysName      string `json:"vsysName"`
	Resource      string `json:"resource"`
	ListFlag      string `json:"listFlag"`
	Sequence      string `json:"sequence"`
	RtpName       string `json:"rtpName"`
	MatchSrctype  string `json:"matchSrctype"`
	MatchDesttype string `json:"matchDesttype"`
	Src           string `json:"src"`
	Dst           string `json:"dst"`
	IIfName       string `json:"iIfName"`
	Protocol      string `json:"protocol"`
	MinSrcport    string `json:"minSrcport"`
	MaxSrcport    string `json:"maxSrcport"`
	MinDstport    string `json:"minDstport"`
	MaxDstport    string `json:"maxDstport"`
	Dscp          string `json:"dscp"`
	Act           string `json:"act"`
	Timestr       string `json:"timestr"`
	Status        string `json:"status"`
	OifName       string `json:"oifName"`
	Weight        string `json:"weight"`
	Nexthop       string `json:"nexthop"`
	HcType        string `json:"hcType"`
	HcName        string `json:"hcName"`
	BfdCheck      string `json:"bfdCheck"`
	NexthopCount  string `json:"nexthopCount"`
}

func (r *Ipv4StrategyRouterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_Ipv4Strategy"
}

func (r *Ipv4StrategyRouterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rtplist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"vsysname": schema.StringAttribute{
						Optional: true,
					},
					"resource": schema.StringAttribute{
						Optional: true,
					},
					"listflag": schema.StringAttribute{
						Optional: true,
					},
					"sequence": schema.StringAttribute{
						Optional: true,
					},
					"rtpname": schema.StringAttribute{
						Optional: true,
					},
					"matchsrctype": schema.StringAttribute{
						Optional: true,
					},
					"matchdesttype": schema.StringAttribute{
						Optional: true,
					},
					"src": schema.StringAttribute{
						Optional: true,
					},
					"dst": schema.StringAttribute{
						Optional: true,
					},
					"iifname": schema.StringAttribute{
						Optional: true,
					},
					"protocol": schema.StringAttribute{
						Optional: true,
					},
					"minsrcport": schema.StringAttribute{
						Optional: true,
					},
					"maxsrcport": schema.StringAttribute{
						Optional: true,
					},
					"mindstport": schema.StringAttribute{
						Optional: true,
					},
					"maxdstport": schema.StringAttribute{
						Optional: true,
					},
					"dscp": schema.StringAttribute{
						Optional: true,
					},
					"act": schema.StringAttribute{
						Optional: true,
					},
					"timestr": schema.StringAttribute{
						Optional: true,
					},
					"status": schema.StringAttribute{
						Optional: true,
					},
					"oifname": schema.StringAttribute{
						Optional: true,
					},
					"weight": schema.StringAttribute{
						Optional: true,
					},
					"nexthop": schema.StringAttribute{
						Optional: true,
					},
					"hctype": schema.StringAttribute{
						Optional: true,
					},
					"hcname": schema.StringAttribute{
						Optional: true,
					},
					"bfdcheck": schema.StringAttribute{
						Optional: true,
					},
					"nexthopcount": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *Ipv4StrategyRouterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*Client)

	if req.ProviderData == nil {
		return
	}
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *Ipv4StrategyRouterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Debug(ctx, "Create Start =========")

	var data *Ipv4StrategyRouterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "Create 出现异常=======")
		return
	}
	sendToweb_Ipv4StrategyRouterRequest(ctx, "POST", r.client, data.AddIpv4StrategyRouterParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4StrategyRouterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *Ipv4StrategyRouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	//sendToweb_Ipv4StrategyRouterRequest(ctx, "GET", r.client, data.AddIpv4StrategyRouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4StrategyRouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *Ipv4StrategyRouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_Ipv4StrategyRouterRequest(ctx, "PUT", r.client, data.AddIpv4StrategyRouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4StrategyRouterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *Ipv4StrategyRouterResourceModel
	tflog.Info(ctx, " Delete Start *************  ")

	//sendToweb_Ipv4StrategyRouterRequest(ctx, "DELETE", r.client, data.AddIpv4StrategyRouterParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *Ipv4StrategyRouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_Ipv4StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv4StrategyRouterParameter) {
	tflog.Info(ctx, "请求开始===========")

	if reqmethod == "POST" {
		// 先查询是否存在，再执行新增操作
		tflog.Info(ctx, "IPv4策略路由--开始执行--查询操作")
		responseBody := sendRequest(ctx, "GET", c, nil, "/func/web_main/api/rt_policy/rtpolicy/rtplist?rtpName="+Rsinfo.RtpName.ValueString()+"&listFlag="+Rsinfo.ListFlag.ValueString(), "IPv4策略路由")
		var queryRes QueryIpv4StrategyRouterResponseModel
		err := json.Unmarshal([]byte(responseBody), &queryRes)
		if err != nil {
			panic("转换查询结果json出现异常")
		}
		if queryRes.RtpName == Rsinfo.RtpName.ValueString() {
			tflog.Info(ctx, "IPv4策略路由--存在重复数据，执行--修改操作")
			var sendUpdateData UpdateIpv4StrategyRouterRequestModel
			sendUpdateData = UpdateIpv4StrategyRouterRequestModel{
				VsysName:      Rsinfo.VsysName.ValueString(),
				Resource:      Rsinfo.Resource.ValueString(),
				ListFlag:      Rsinfo.ListFlag.ValueString(),
				Sequence:      Rsinfo.Sequence.ValueString(),
				RtpName:       Rsinfo.RtpName.ValueString(),
				MatchSrctype:  Rsinfo.MatchSrctype.ValueString(),
				MatchDesttype: Rsinfo.MatchDesttype.ValueString(),
				Src:           Rsinfo.Src.ValueString(),
				Dst:           Rsinfo.Dst.ValueString(),
				IIfName:       Rsinfo.IIfName.ValueString(),
				Protocol:      Rsinfo.Protocol.ValueString(),
				MinSrcport:    Rsinfo.MinSrcport.ValueString(),
				MaxSrcport:    Rsinfo.MaxSrcport.ValueString(),
				MinDstport:    Rsinfo.MinDstport.ValueString(),
				MaxDstport:    Rsinfo.MaxDstport.ValueString(),
				Dscp:          Rsinfo.Dscp.ValueString(),
				Act:           Rsinfo.Act.ValueString(),
				Timestr:       Rsinfo.Timestr.ValueString(),
				Status:        Rsinfo.Status.ValueString(),
				OifName:       Rsinfo.OifName.ValueString(),
				Weight:        Rsinfo.Weight.ValueString(),
				Nexthop:       Rsinfo.Nexthop.ValueString(),
				HcType:        Rsinfo.HcType.ValueString(),
				HcName:        Rsinfo.HcName.ValueString(),
				BfdCheck:      Rsinfo.BfdCheck.ValueString(),
				NexthopCount:  Rsinfo.NexthopCount.ValueString(),
			}

			requstUpdateData := UpdateIpv4StrategyRouterRequest{
				UpdateIpv4StrategyRouterRequestModel: sendUpdateData,
			}
			body, _ := json.Marshal(requstUpdateData)

			sendRequest(ctx, "PUT", c, body, "/func/web_main/api/rt_policy/rtpolicy/rtplist", "IPv4策略路由")
			return
		}
		// 新增操作
		var sendData AddIpv4StrategyRouterRequestModel

		sendData = AddIpv4StrategyRouterRequestModel{
			RtpName: Rsinfo.RtpName.ValueString(),
			Act:     Rsinfo.Act.ValueString(),
		}

		requstData := AddIpv4StrategyRouterRequest{
			AddIpv4StrategyRouterRequestModel: sendData,
		}
		body, _ := json.Marshal(requstData)

		sendRequest(ctx, reqmethod, c, body, "/func/web_main/api/rt_policy/rtpolicy/rtplist", "IPv4策略路由")
		return
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}
}
