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

// VPN IP地址池
var _ resource.Resource = &VpnIpPoolResource{}
var _ resource.ResourceWithImportState = &VpnIpPoolResource{}

func NewVpnIpPoolResource() resource.Resource {
	return &VpnIpPoolResource{}
}

type VpnIpPoolResource struct {
	client *Client
}

type VpnIpPoolResourceModel struct {
	AddVpnIpPoolParameter    AddVpnIpPoolParameter    `tfsdk:"addVpnIpPoolParameter"`
	UpdateVpnIpPoolParameter UpdateVpnIpPoolParameter `tfsdk:"updateVpnIpPoolParameter"`
	DelVpnIpPoolParameter    DelVpnIpPoolParameter    `tfsdk:"delVpnIpPoolParameter"`
	ReadVpnIpPoolParameter   ReadVpnIpPoolParameter   `tfsdk:"readVpnIpPoolParameter"`
}

type AddVpnIpPoolParameter struct {
	VsysName  types.String `tfsdk:"vsysName"`
	poolName  types.String `tfsdk:"poolName"`
	poolStart types.String `tfsdk:"poolStart"`
	poolEnd   types.String `tfsdk:"poolEnd"`
	poolMask  types.String `tfsdk:"poolMask"`
}

type UpdateVpnIpPoolParameter struct {
	VsysName  types.String `tfsdk:"vsysName"`
	poolName  types.String `tfsdk:"poolName"`
	poolStart types.String `tfsdk:"poolStart"`
	poolEnd   types.String `tfsdk:"poolEnd"`
	poolMask  types.String `tfsdk:"poolMask"`
}

type DelVpnIpPoolParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
	poolName types.String `tfsdk:"poolName"`
}

type ReadVpnIpPoolParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
}

func (r *VpnIpPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_VpnIpPool"
}

func (r *VpnIpPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"param": schema.SingleNestedAttribute{
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

func (r *VpnIpPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VpnIpPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VpnIpPoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddVpnIpPoolRequest(ctx, "POST", r.client, data.AddVpnIpPoolParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VpnIpPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadVpnIpPoolRequest(ctx, "GET", r.client, data.ReadVpnIpPoolParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VpnIpPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateVpnIpPoolRequest(ctx, "PUT", r.client, data.UpdateVpnIpPoolParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VpnIpPoolResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelVpnIpPoolRequest(ctx, "DELETE", r.client, data.DelVpnIpPoolParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VpnIpPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddVpnIpPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddVpnIpPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipPoolInfo"

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

func sendToweb_UpdateVpnIpPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateVpnIpPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipPoolInfo"

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

func sendToweb_DelVpnIpPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelVpnIpPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipPoolInfo"

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

func sendToweb_ReadVpnIpPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadVpnIpPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipPoolInfo?vsysName=PublicSystem"

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
