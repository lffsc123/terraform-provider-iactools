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

// IPSec保护网段
var _ resource.Resource = &TunIfResource{}
var _ resource.ResourceWithImportState = &TunIfResource{}

func NewTunIfResource() resource.Resource {
	return &TunIfResource{}
}

type TunIfResource struct {
	client *Client
}

type TunIfResourceModel struct {
	AddTunIfParameter    AddTunIfParameter    `tfsdk:"addTunIfParameter"`
	UpdateTunIfParameter UpdateTunIfParameter `tfsdk:"updateTunIfParameter"`
	DelTunIfParameter    DelTunIfParameter    `tfsdk:"delTunIfParameter"`
	ReadTunIfParameter   ReadTunIfParameter   `tfsdk:"readTunIfParameter"`
}

type AddTunIfParameter struct {
	TunId     types.String `tfsdk:"tunId"`
	TunIp     types.String `tfsdk:"tunIp"`
	IpsecDesc types.String `tfsdk:"ipsecDesc"`
}

type UpdateTunIfParameter struct {
	IfName    types.String `tfsdk:"ifName"`
	TunIp     types.String `tfsdk:"tunIp"`
	IpsecDesc types.String `tfsdk:"ipsecDesc"`
}

type DelTunIfParameter struct {
	IfName types.String `tfsdk:"ifName"`
}

type ReadTunIfParameter struct {
	VsysName  types.String `tfsdk:"vsysName"`
	VrfName   types.String `tfsdk:"vrfName"`
	IfName    types.String `tfsdk:"ifName"`
	TunId     types.String `tfsdk:"tunId"`
	TunIp     types.String `tfsdk:"tunIp"`
	AppMode   types.String `tfsdk:"appMode"`
	IpsecDesc types.String `tfsdk:"ipsecDesc"`
}

func (r *TunIfResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-TunIf"
}

func (r *TunIfResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *TunIfResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TunIfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TunIfResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddTunIfRequest(ctx, "POST", r.client, data.AddTunIfParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunIfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TunIfResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadTunIfRequest(ctx, "GET", r.client, data.ReadTunIfParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunIfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TunIfResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateTunIfRequest(ctx, "PUT", r.client, data.UpdateTunIfParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TunIfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TunIfResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelTunIfRequest(ctx, "DELETE", r.client, data.DelTunIfParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TunIfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddTunIfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddTunIfParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/tunIf"

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

func sendToweb_UpdateTunIfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateTunIfParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/tunIf"

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

func sendToweb_DelTunIfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelTunIfParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/tunIf"

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

func sendToweb_ReadTunIfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadTunIfParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/tunIf"

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
