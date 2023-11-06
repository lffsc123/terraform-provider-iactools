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

// 地址池
var _ resource.Resource = &AddrPoolResource{}
var _ resource.ResourceWithImportState = &AddrPoolResource{}

func NewAddrPoolResource() resource.Resource {
	return &AddrPoolResource{}
}

type AddrPoolResource struct {
	client *Client
}

type AddrPoolResourceModel struct {
	AddAddrPoolParameter    AddAddrPoolParameter    `tfsdk:"addAddrPoolParameter"`
	UpdateAddrPoolParameter UpdateAddrPoolParameter `tfsdk:"updateAddrPoolParameter"`
	DelAddrPoolParameter    DelAddrPoolParameter    `tfsdk:"delAddrPoolParameter"`
	ReadAddrPoolParameter   ReadAddrPoolParameter   `tfsdk:"readAddrPoolParameter"`
}

type AddAddrPoolParameter struct {
	IpVersion     types.String `tfsdk:"ipVersion"`
	VsysName      types.String `tfsdk:"vsysName"`
	Name          types.String `tfsdk:"name"`
	IpStart       types.String `tfsdk:"ipStart"`
	IpEnd         types.String `tfsdk:"ipEnd"`
	LoopRoute     types.String `tfsdk:"loopRoute"`
	Slot          types.String `tfsdk:"slot"`
	GratuitousArp types.String `tfsdk:"gratuitousArp"`
	VrrpIfName    types.String `tfsdk:"vrrpIfName"`
	VrrpId        types.String `tfsdk:"vrrpId"`
}

type UpdateAddrPoolParameter struct {
	IpVersion     types.String `tfsdk:"ipVersion"`
	VsysName      types.String `tfsdk:"vsysName"`
	Name          types.String `tfsdk:"name"`
	OldName       types.String `tfsdk:"oldName"`
	IpStart       types.String `tfsdk:"ipStart"`
	IpEnd         types.String `tfsdk:"ipEnd"`
	LoopRoute     types.String `tfsdk:"loopRoute"`
	Slot          types.String `tfsdk:"slot"`
	GratuitousArp types.String `tfsdk:"gratuitousArp"`
	VrrpIfName    types.String `tfsdk:"vrrpIfName"`
	VrrpId        types.String `tfsdk:"vrrpId"`
}

type DelAddrPoolParameter struct {
	IpVersion    types.String `tfsdk:"ipVersion"`
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	DelAllEnable types.String `tfsdk:"delAllEnable"`
}

type ReadAddrPoolParameter struct {
	IpVersion     types.String `tfsdk:"ipVersion"`
	VsysName      types.String `tfsdk:"vsysName"`
	Count         types.String `tfsdk:"count"`
	Offset        types.String `tfsdk:"offset"`
	Name          types.String `tfsdk:"name"`
	IpStart       types.String `tfsdk:"ipStart"`
	IpEnd         types.String `tfsdk:"ipEnd"`
	LoopRoute     types.String `tfsdk:"loopRoute"`
	Slot          types.String `tfsdk:"slot"`
	GratuitousArp types.String `tfsdk:"gratuitousArp"`
	VrrpIfName    types.String `tfsdk:"vrrpIfName"`
	VrrpId        types.String `tfsdk:"vrrpId"`
	ReferNum      types.String `tfsdk:"referNum"`
	Id            types.String `tfsdk:"id"`
	DelallEnable  types.String `tfsdk:"delallEnable"`
	OldName       types.String `tfsdk:"oldName"`
}

func (r *AddrPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-AddrPool"
}

func (r *AddrPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *AddrPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AddrPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AddrPoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddAddrPoolRequest(ctx, "POST", r.client, data.AddAddrPoolParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddrPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AddrPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadAddrPoolRequest(ctx, "GET", r.client, data.ReadAddrPoolParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddrPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AddrPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateAddrPoolRequest(ctx, "PUT", r.client, data.UpdateAddrPoolParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddrPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AddrPoolResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelAddrPoolRequest(ctx, "DELETE", r.client, data.DelAddrPoolParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AddrPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddAddrPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddAddrPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/addrpool/addrpool/addrpoollist"

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

func sendToweb_UpdateAddrPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateAddrPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/addrpool/addrpool/addrpoollist"

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

func sendToweb_DelAddrPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelAddrPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/addrpool/addrpool/addrpoollist"

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

func sendToweb_ReadAddrPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadAddrPoolParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/addrpool/addrpool/addrpoollist"

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
