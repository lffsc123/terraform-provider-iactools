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

// VPN IP资源
var _ resource.Resource = &VpnIpResResource{}
var _ resource.ResourceWithImportState = &VpnIpResResource{}

func NewVpnIpResResource() resource.Resource {
	return &VpnIpResResource{}
}

type VpnIpResResource struct {
	client *Client
}

type VpnIpResResourceModel struct {
	AddVpnIpResParameter    AddVpnIpResParameter    `tfsdk:"addVpnIpResParameter"`
	UpdateVpnIpResParameter UpdateVpnIpResParameter `tfsdk:"updateVpnIpResParameter"`
	DelVpnIpResParameter    DelVpnIpResParameter    `tfsdk:"delVpnIpResParameter"`
	ReadVpnIpResParameter   ReadVpnIpResParameter   `tfsdk:"readVpnIpResParameter"`
}

type AddVpnIpResParameter struct {
	VsysName     types.String `tfsdk:"vsysName"`
	IpResName    types.String `tfsdk:"ipResName"`
	IpResNetAddr types.String `tfsdk:"ipResNetAddr"`
	IpResMaskLen types.String `tfsdk:"ipResMaskLen"`
}

type UpdateVpnIpResParameter struct {
	VsysName     types.String `tfsdk:"vsysName"`
	IpResName    types.String `tfsdk:"ipResName"`
	IpResNetAddr types.String `tfsdk:"ipResNetAddr"`
	IpResMaskLen types.String `tfsdk:"ipResMaskLen"`
}

type DelVpnIpResParameter struct {
	VsysName  types.String `tfsdk:"vsysName"`
	IpResName types.String `tfsdk:"ipResName"`
}

type ReadVpnIpResParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
}

func (r *VpnIpResResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-VpnIpRes"
}

func (r *VpnIpResResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *VpnIpResResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VpnIpResResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VpnIpResResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddVpnIpResRequest(ctx, "POST", r.client, data.AddVpnIpResParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpResResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VpnIpResResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadVpnIpResRequest(ctx, "GET", r.client, data.ReadVpnIpResParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpResResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VpnIpResResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateVpnIpResRequest(ctx, "PUT", r.client, data.UpdateVpnIpResParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnIpResResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VpnIpResResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelVpnIpResRequest(ctx, "DELETE", r.client, data.DelVpnIpResParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VpnIpResResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddVpnIpResRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddVpnIpResParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipResInfo"

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

func sendToweb_UpdateVpnIpResRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateVpnIpResParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipResInfo"

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

func sendToweb_DelVpnIpResRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelVpnIpResParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipResInfo"

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

func sendToweb_ReadVpnIpResRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadVpnIpResParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipResInfo?vsysName=PublicSystem"

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
