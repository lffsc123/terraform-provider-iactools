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

// VPN IP资源组
var _ resource.Resource = &VpnIpResGroupResource{}
var _ resource.ResourceWithImportState = &VpnIpResGroupResource{}

func NewVpnIpResGroupResource() resource.Resource {
	return &VpnIpResGroupResource{}
}

type VpnIpResGroupResource struct {
	client *Client
}

type VpnIpResGroupResourceModel struct {
	AddVpnIpResGroupParameter    AddVpnIpResGroupParameter    `tfsdk:"addVpnIpResGroupParameter"`
	UpdateVpnIpResGroupParameter UpdateVpnIpResGroupParameter `tfsdk:"updateVpnIpResGroupParameter"`
	DelVpnIpResGroupParameter    DelVpnIpResGroupParameter    `tfsdk:"delVpnIpResGroupParameter"`
	ReadVpnIpResGroupParameter   ReadVpnIpResGroupParameter   `tfsdk:"readVpnIpResGroupParameter"`
}

type AddVpnIpResGroupParameter struct {
	VsysName   types.String `tfsdk:"vsysName"`
	resGrpName types.String `tfsdk:"resGrpName"`
	ipResNames types.String `tfsdk:"ipResNames"`
}

type UpdateVpnIpResGroupParameter struct {
	VsysName   types.String `tfsdk:"vsysName"`
	resGrpName types.String `tfsdk:"resGrpName"`
	ipResNames types.String `tfsdk:"ipResNames"`
}

type DelVpnIpResGroupParameter struct {
	VsysName   types.String `tfsdk:"vsysName"`
	resGrpName types.String `tfsdk:"resGrpName"`
}

type ReadVpnIpResGroupParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
}

func (r *VpnIpResGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-VpnIpResGroup"
}

func (r *VpnIpResGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *VpnIpResGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VpnIpResGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VpnIpResGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddVpnIpResGroupRequest(ctx, "POST", r.client, data.AddVpnIpResGroupParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpResGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VpnIpResGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadVpnIpResGroupRequest(ctx, "GET", r.client, data.ReadVpnIpResGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpResGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VpnIpResGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateVpnIpResGroupRequest(ctx, "PUT", r.client, data.UpdateVpnIpResGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpResGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VpnIpResGroupResourceModel
	tflog.Info(ctx, " Delete Start  *************")

	sendToweb_DelVpnIpResGroupRequest(ctx, "DELETE", r.client, data.DelVpnIpResGroupParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VpnIpResGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddVpnIpResGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddVpnIpResGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/resGrpInfo"

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

func sendToweb_UpdateVpnIpResGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateVpnIpResGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/resGrpInfo"

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

func sendToweb_DelVpnIpResGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelVpnIpResGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/resGrpInfo"

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

func sendToweb_ReadVpnIpResGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadVpnIpResGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/resGrpInfo?vsysName=PublicSystem"

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
