package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// IPv6Strategy策略路由
var _ resource.Resource = &Ipv6StrategyRouterResource{}
var _ resource.ResourceWithImportState = &Ipv6StrategyRouterResource{}

func NewIpv6StrategyRouterResource() resource.Resource {
	return &Ipv6StrategyRouterResource{}
}

type Ipv6StrategyRouterResource struct {
	client *Client
}

type Ipv6StrategyRouterResourceModel struct {
	AddIpv6StrategyRouterParameter    AddIpv6StrategyRouterParameter    `tfsdk:"addIpv6StrategyRouterParameter"`
	UpdateIpv6StrategyRouterParameter UpdateIpv6StrategyRouterParameter `tfsdk:"updateIpv6StrategyRouterParameter"`
	DelIpv6StrategyRouterParameter    DelIpv6StrategyRouterParameter    `tfsdk:"delIpv6StrategyRouterParameter"`
	ReadIpv6StrategyRouterParameter   ReadIpv6StrategyRouterParameter   `tfsdk:"readIpv6StrategyRouterParameter"`
}

type AddIpv6StrategyRouterParameter struct {
	VsysName     types.String `tfsdk:"vsysName"`
	Resource     types.String `tfsdk:"resource"`
	ListFlag     types.String `tfsdk:"listFlag"`
	Sequence     types.String `tfsdk:"sequence"`
	RtpName      types.String `tfsdk:"rtpName"`
	Src          types.String `tfsdk:"src"`
	Dst          types.String `tfsdk:"dst"`
	IIfName      types.String `tfsdk:"iIfName"`
	Protocol     types.String `tfsdk:"protocol"`
	MinSrcport   types.String `tfsdk:"minSrcport"`
	MaxSrcport   types.String `tfsdk:"maxSrcport"`
	MinDstport   types.String `tfsdk:"minDstport"`
	MaxDstport   types.String `tfsdk:"maxDstport"`
	Dscp         types.String `tfsdk:"dscp"`
	Act          types.String `tfsdk:"act"`
	OifName      types.String `tfsdk:"oifName"`
	Weight       types.String `tfsdk:"weight"`
	Nexthop      types.String `tfsdk:"nexthop"`
	HcType       types.String `tfsdk:"hcType"`
	HcName       types.String `tfsdk:"hcName"`
	BfdCheck     types.String `tfsdk:"bfdCheck"`
	NexthopCount types.String `tfsdk:"nexthopCount"`
}

type UpdateIpv6StrategyRouterParameter struct {
	VsysName     types.String `tfsdk:"vsysName"`
	Resource     types.String `tfsdk:"resource"`
	ListFlag     types.String `tfsdk:"listFlag"`
	Sequence     types.String `tfsdk:"sequence"`
	RtpName      types.String `tfsdk:"rtpName"`
	Src          types.String `tfsdk:"src"`
	Dst          types.String `tfsdk:"dst"`
	IIfName      types.String `tfsdk:"iIfName"`
	Protocol     types.String `tfsdk:"protocol"`
	MinSrcport   types.String `tfsdk:"minSrcport"`
	MaxSrcport   types.String `tfsdk:"maxSrcport"`
	MinDstport   types.String `tfsdk:"minDstport"`
	MaxDstport   types.String `tfsdk:"maxDstport"`
	Dscp         types.String `tfsdk:"dscp"`
	Act          types.String `tfsdk:"act"`
	OifName      types.String `tfsdk:"oifName"`
	Weight       types.String `tfsdk:"weight"`
	Nexthop      types.String `tfsdk:"nexthop"`
	HcType       types.String `tfsdk:"hcType"`
	HcName       types.String `tfsdk:"hcName"`
	BfdCheck     types.String `tfsdk:"bfdCheck"`
	NexthopCount types.String `tfsdk:"nexthopCount"`
}

type DelIpv6StrategyRouterParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
	Resource types.String `tfsdk:"resource"`
	ListFlag types.String `tfsdk:"listFlag"`
	RtpName  types.String `tfsdk:"rtpName"`
}

type ReadIpv6StrategyRouterParameter struct {
	Resource types.String `tfsdk:"resource"`
	ListFlag types.String `tfsdk:"listFlag"`
	RtpName  types.String `tfsdk:"rtpName"`
	Offset   types.String `tfsdk:"offset"`
	Count    types.String `tfsdk:"count"`
}

func (r *Ipv6StrategyRouterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_Ipv6Strategy"
}

func (r *Ipv6StrategyRouterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"targetnat": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"ip_start": schema.StringAttribute{
						Required: true,
					},
					"ip_end": schema.StringAttribute{
						Required: true,
					},
					"ip_version": schema.StringAttribute{
						Optional: true,
					},
					"vrrp_if_name": schema.StringAttribute{
						Optional: true,
					},
					"vrrp_id": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *Ipv6StrategyRouterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Ipv6StrategyRouterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *Ipv6StrategyRouterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddIpv6StrategyRouterRequest(ctx, "POST", r.client, data.AddIpv6StrategyRouterParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv6StrategyRouterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *Ipv6StrategyRouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	sendToweb_ReadIpv6StrategyRouterRequest(ctx, "GET", r.client, data.ReadIpv6StrategyRouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv6StrategyRouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *Ipv6StrategyRouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateIpv6StrategyRouterRequest(ctx, "PUT", r.client, data.UpdateIpv6StrategyRouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv6StrategyRouterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *Ipv6StrategyRouterResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelIpv6StrategyRouterRequest(ctx, "DELETE", r.client, data.DelIpv6StrategyRouterParameter)

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

func (r *Ipv6StrategyRouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddIpv6StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv6StrategyRouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rt_policy_ipv6/rtpolicyipv6/rtp6list"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}

func sendToweb_UpdateIpv6StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateIpv6StrategyRouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rt_policy_ipv6/rtpolicyipv6/rtp6list"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}

func sendToweb_DelIpv6StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelIpv6StrategyRouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rt_policy_ipv6/rtpolicyipv6/rtp6list"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}

func sendToweb_ReadIpv6StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadIpv6StrategyRouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rt_policy_ipv6/rtpolicyipv6/rtp6list?rtpName=aaa&listFlag=0"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}
