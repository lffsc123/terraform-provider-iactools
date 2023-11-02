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

// IPv4策略路由
var _ resource.Resource = &Ipv4RouterResource{}
var _ resource.ResourceWithImportState = &Ipv4RouterResource{}

func NewIpv4RouterResource() resource.Resource {
	return &Ipv4RouterResource{}
}

type Ipv4RouterResource struct {
	client *Client
}

type Ipv4RouterResourceModel struct {
	AddIpv4RouterParameter    AddIpv4RouterParameter    `tfsdk:"addIpv4RouterParameter"`
	UpdateIpv4RouterParameter UpdateIpv4RouterParameter `tfsdk:"updateIpv4RouterParameter"`
	DelIpv4RouterParameter    DelIpv4RouterParameter    `tfsdk:"delIpv4RouterParameter"`
}

type AddIpv4RouterParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VpnName   types.String `tfsdk:"vpnName"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	Gateway   types.String `tfsdk:"gateway"`
	Interface types.String `tfsdk:"interface"`
	RouteType types.String `tfsdk:"routeType"`
	Distance  types.String `tfsdk:"distance"`
	Weight    types.String `tfsdk:"weight"`
	BfdCheck  types.String `tfsdk:"bfdCheck"`
	HcName    types.String `tfsdk:"hcName"`
	Describe  types.String `tfsdk:"describe"`
	Tag       types.String `tfsdk:"tag"`
}

type UpdateIpv4RouterParameter struct {
	IpVersion    types.String `tfsdk:"ipVersion"`
	VpnName      types.String `tfsdk:"vpnName"`
	Ip           types.String `tfsdk:"ip"`
	Mask         types.String `tfsdk:"mask"`
	Gateway      types.String `tfsdk:"gateway"`
	Interface    types.String `tfsdk:"interface"`
	IpMod        types.String `tfsdk:"ipMod"`
	MaskMod      types.String `tfsdk:"MaskMod"`
	GatewayMod   types.String `tfsdk:"gatewayMod"`
	InterfaceMod types.String `tfsdk:"interfaceMod"`
	RouteTypeMod types.String `tfsdk:"routeTypeMod"`
	Distance     types.String `tfsdk:"distance"`
	Weight       types.String `tfsdk:"weight"`
	BfdCheck     types.String `tfsdk:"bfdCheck"`
	HcName       types.String `tfsdk:"hcName"`
	Describe     types.String `tfsdk:"describe"`
	Tag          types.String `tfsdk:"tag"`
}

type DelIpv4RouterParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VpnName   types.String `tfsdk:"vpnName"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	Gateway   types.String `tfsdk:"gateway"`
	Interface types.String `tfsdk:"interface"`
	RouteType types.String `tfsdk:"routeType"`
}

func (r *Ipv4RouterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-Ipv4Router"
}

func (r *Ipv4RouterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"updatesourcenat": schema.SingleNestedAttribute{
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

func (r *Ipv4RouterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Ipv4RouterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *Ipv4RouterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddIpv4RouterRequest(ctx, "POST", r.client, data.AddIpv4RouterParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4RouterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *Ipv4RouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	// sendToweb_AddSourceNatRequest(ctx,"POST", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4RouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *Ipv4RouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateIpv4RouterRequest(ctx, "PUT", r.client, data.UpdateIpv4RouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4RouterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *Ipv4RouterResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelIpv4RouterRequest(ctx, "DELETE", r.client, data.DelIpv4RouterParameter)

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

func (r *Ipv4RouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddIpv4RouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv4RouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist"

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

func sendToweb_UpdateIpv4RouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateIpv4RouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist"

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

func sendToweb_DelIpv4RouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelIpv4RouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist"

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
