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

// IPv6静态路由
var _ resource.Resource = &Ipv6RouterResource{}
var _ resource.ResourceWithImportState = &Ipv6RouterResource{}

func NewIpv6RouterResource() resource.Resource {
	return &Ipv6RouterResource{}
}

type Ipv6RouterResource struct {
	client *Client
}

type Ipv6RouterResourceModel struct {
	AddIpv6RouterParameter    AddIpv6RouterParameter    `tfsdk:"addIpv6RouterParameter"`
	UpdateIpv6RouterParameter UpdateIpv6RouterParameter `tfsdk:"updateIpv6RouterParameter"`
	DelIpv6RouterParameter    DelIpv6RouterParameter    `tfsdk:"delIpv6RouterParameter"`
	ReadIpv6RouterParameter   ReadIpv6RouterParameter   `tfsdk:"readIpv6RouterParameter"`
}

type AddIpv6RouterParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VpnName   types.String `tfsdk:"vpnName"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	Gateway   types.String `tfsdk:"gateway"`
	Interface types.String `tfsdk:"interface"`
	RouteType types.String `tfsdk:"routeType"`
	Distance  types.String `tfsdk:"distance"`
	BfdCheck  types.String `tfsdk:"bfdCheck"`
}

type UpdateIpv6RouterParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VpnName   types.String `tfsdk:"vpnName"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	Gateway   types.String `tfsdk:"gateway"`
	Interface types.String `tfsdk:"interface"`
	RouteType types.String `tfsdk:"routeType"`
	Distance  types.String `tfsdk:"distance"`
	BfdCheck  types.String `tfsdk:"bfdCheck"`
}

type DelIpv6RouterParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VpnName   types.String `tfsdk:"vpnName"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	Gateway   types.String `tfsdk:"gateway"`
	Interface types.String `tfsdk:"interface"`
	RouteType types.String `tfsdk:"routeType"`
}

type ReadIpv6RouterParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VpnName   types.String `tfsdk:"vpnName"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
}

func (r *Ipv6RouterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-Ipv6Router"
}

func (r *Ipv6RouterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *Ipv6RouterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Ipv6RouterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *Ipv6RouterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddIpv6RouterRequest(ctx, "POST", r.client, data.AddIpv6RouterParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv6RouterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *Ipv6RouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	sendToweb_ReadIpv6RouterRequest(ctx, "GET", r.client, data.ReadIpv6RouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv6RouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *Ipv6RouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateIpv6RouterRequest(ctx, "PUT", r.client, data.UpdateIpv6RouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv6RouterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *Ipv6RouterResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelIpv6RouterRequest(ctx, "DELETE", r.client, data.DelIpv6RouterParameter)

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

func (r *Ipv6RouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddIpv6RouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv6RouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rtm_api/restfulstaticroute/routeEntries"

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

func sendToweb_UpdateIpv6RouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateIpv6RouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rtm_api/restfulstaticroute/routeEntries"

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

func sendToweb_DelIpv6RouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelIpv6RouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rtm_api/restfulstaticroute/routeEntries"

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

func sendToweb_ReadIpv6RouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadIpv6RouterParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/rtm_api/restfulstaticroute/routeEntries?ip=123:100::123&mask=128&ipVersion=6"

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
