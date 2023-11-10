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
var _ resource.Resource = &ProtectNetResource{}
var _ resource.ResourceWithImportState = &ProtectNetResource{}

func NewProtectNetResource() resource.Resource {
	return &ProtectNetResource{}
}

type ProtectNetResource struct {
	client *Client
}

type ProtectNetResourceModel struct {
	AddProtectNetParameter    AddProtectNetParameter    `tfsdk:"addProtectNetParameter"`
	UpdateProtectNetParameter UpdateProtectNetParameter `tfsdk:"updateProtectNetParameter"`
	DelProtectNetParameter    DelProtectNetParameter    `tfsdk:"delProtectNetParameter"`
	ReadProtectNetParameter   ReadProtectNetParameter   `tfsdk:"readProtectNetParameter"`
}

type AddProtectNetParameter struct {
	IpVersion  types.String `tfsdk:"ipVersion"`
	VsysName   types.String `tfsdk:"vsysName"`
	Group      types.String `tfsdk:"group"`
	Srcnetaddr types.String `tfsdk:"srcnetaddr"`
	Srcnetmask types.String `tfsdk:"srcnetmask"`
	Dstnetaddr types.String `tfsdk:"dstnetaddr"`
	Dstnetmask types.String `tfsdk:"dstnetmask"`
}

type UpdateProtectNetParameter struct {
	IpVersion  types.String `tfsdk:"ipVersion"`
	VsysName   types.String `tfsdk:"vsysName"`
	ResName    types.String `tfsdk:"resName"`
	Group      types.String `tfsdk:"group"`
	Srcnetaddr types.String `tfsdk:"srcnetaddr"`
	Srcnetmask types.String `tfsdk:"srcnetmask"`
	Dstnetaddr types.String `tfsdk:"dstnetaddr"`
	Dstnetmask types.String `tfsdk:"dstnetmask"`
}

type DelProtectNetParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
	Group    types.String `tfsdk:"group"`
}

type ReadProtectNetParameter struct {
	IpVersion  types.String `tfsdk:"ipVersion"`
	VsysName   types.String `tfsdk:"vsysName"`
	ResName    types.String `tfsdk:"resName"`
	Group      types.String `tfsdk:"group"`
	Srcnetaddr types.String `tfsdk:"srcnetaddr"`
	Srcnetmask types.String `tfsdk:"srcnetmask"`
	Dstnetaddr types.String `tfsdk:"dstnetaddr"`
	Dstnetmask types.String `tfsdk:"dstnetmask"`
}

func (r *ProtectNetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_ProtectNet"
}

func (r *ProtectNetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *ProtectNetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ProtectNetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ProtectNetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddProtectNetRequest(ctx, "POST", r.client, data.AddProtectNetParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProtectNetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ProtectNetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadProtectNetRequest(ctx, "GET", r.client, data.ReadProtectNetParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProtectNetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ProtectNetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateProtectNetRequest(ctx, "PUT", r.client, data.UpdateProtectNetParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProtectNetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ProtectNetResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelProtectNetRequest(ctx, "DELETE", r.client, data.DelProtectNetParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ProtectNetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddProtectNetRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddProtectNetParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/protectnetlist"

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

func sendToweb_UpdateProtectNetRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateProtectNetParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/protectnetlist"

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

func sendToweb_DelProtectNetRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelProtectNetParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/protectnetlist"

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

func sendToweb_ReadProtectNetRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadProtectNetParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/protectnetlist"

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
