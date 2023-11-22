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

// SSL VPN 用户组
var _ resource.Resource = &VpnUserGroupResource{}
var _ resource.ResourceWithImportState = &VpnUserGroupResource{}

func NewVpnUserGroupResource() resource.Resource {
	return &VpnUserGroupResource{}
}

type VpnUserGroupResource struct {
	client *Client
}

type VpnUserGroupResourceModel struct {
	AddVpnUserGroupParameter    AddVpnUserGroupParameter    `tfsdk:"addVpnUserGroupParameter"`
	UpdateVpnUserGroupParameter UpdateVpnUserGroupParameter `tfsdk:"updateVpnUserGroupParameter"`
	DelVpnUserGroupParameter    DelVpnUserGroupParameter    `tfsdk:"delVpnUserGroupParameter"`
	ReadVpnUserGroupParameter   ReadVpnUserGroupParameter   `tfsdk:"readVpnUserGroupParameter"`
}

type AddVpnUserGroupParameter struct {
	VsysName    types.String `tfsdk:"vsysName"`
	UserGrpName types.String `tfsdk:"userGrpName"`
	ResGrpNames types.String `tfsdk:"resGrpNames"`
	IpPoolName  types.String `tfsdk:"ipPoolName"`
}

type UpdateVpnUserGroupParameter struct {
	VsysName    types.String `tfsdk:"vsysName"`
	UserGrpName types.String `tfsdk:"userGrpName"`
	ResGrpNames types.String `tfsdk:"resGrpNames"`
	IpPoolName  types.String `tfsdk:"ipPoolName"`
}

type DelVpnUserGroupParameter struct {
	VsysName    types.String `tfsdk:"vsysName"`
	UserGrpName types.String `tfsdk:"userGrpName"`
}

type ReadVpnUserGroupParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
}

func (r *VpnUserGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "firewall_VpnUserGroup"
}

func (r *VpnUserGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *VpnUserGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VpnUserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VpnUserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddVpnUserGroupRequest(ctx, "POST", r.client, data.AddVpnUserGroupParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnUserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VpnUserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadVpnUserGroupRequest(ctx, "GET", r.client, data.ReadVpnUserGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnUserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VpnUserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateVpnUserGroupRequest(ctx, "PUT", r.client, data.UpdateVpnUserGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnUserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VpnUserGroupResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelVpnUserGroupRequest(ctx, "DELETE", r.client, data.DelVpnUserGroupParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VpnUserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddVpnUserGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddVpnUserGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/userGrpInfo"

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

func sendToweb_UpdateVpnUserGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateVpnUserGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/userGrpInfo"

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

func sendToweb_DelVpnUserGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelVpnUserGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/userGrpInfo"

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

func sendToweb_ReadVpnUserGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadVpnUserGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/userGrpInfo?vsysName=PublicSystem"

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
