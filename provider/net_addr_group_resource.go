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

// ip地址组
var _ resource.Resource = &NetAddrGroupResource{}
var _ resource.ResourceWithImportState = &NetAddrGroupResource{}

func NewNetAddrGroupResource() resource.Resource {
	return &NetAddrGroupResource{}
}

type NetAddrGroupResource struct {
	client *Client
}

type NetAddrGroupResourceModel struct {
	AddNetAddrGroupParameter    AddNetAddrGroupParameter    `tfsdk:"addNetAddrGroupParameter"`
	UpdateNetAddrGroupParameter UpdateNetAddrGroupParameter `tfsdk:"updateNetAddrGroupParameter"`
	DelNetAddrGroupParameter    DelNetAddrGroupParameter    `tfsdk:"delNetAddrGroupParameter"`
	ReadNetAddrGroupParameter   ReadNetAddrGroupParameter   `tfsdk:"readNetAddrGroupParameter"`
}

type AddNetAddrGroupParameter struct {
	VsysName    types.String `tfsdk:"vsysName"`
	Name        types.String `tfsdk:"name"`
	ObjNameList types.String `tfsdk:"objNameList"`
	Desc        types.String `tfsdk:"desc"`
}

type UpdateNetAddrGroupParameter struct {
	VsysName    types.String `tfsdk:"vsysName"`
	Name        types.String `tfsdk:"name"`
	OldName     types.String `tfsdk:"oldName"`
	ObjNameList types.String `tfsdk:"objNameList"`
	Desc        types.String `tfsdk:"desc"`
}

type DelNetAddrGroupParameter struct {
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	DelAllEnable types.String `tfsdk:"delAllEnable"`
}

type ReadNetAddrGroupParameter struct {
	Id           types.String `tfsdk:"id"`
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	Oldname      types.String `tfsdk:"oldname"`
	ObjNameList  types.String `tfsdk:"objNameList"`
	Desc         types.String `tfsdk:"desc"`
	ReferNum     types.String `tfsdk:"referNum"`
	DelAllEnable types.String `tfsdk:"delAllEnable"`
	SearchValue  types.String `tfsdk:"searchValue"`
	Offset       types.String `tfsdk:"offset"`
	Count        types.String `tfsdk:"count"`
}

func (r *NetAddrGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_NetAddrGroup"
}

func (r *NetAddrGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *NetAddrGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetAddrGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetAddrGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddNetAddrGroupRequest(ctx, "POST", r.client, data.AddNetAddrGroupParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetAddrGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	sendToweb_ReadNetAddrGroupRequest(ctx, "GET", r.client, data.ReadNetAddrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetAddrGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateNetAddrGroupRequest(ctx, "PUT", r.client, data.UpdateNetAddrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetAddrGroupResourceModel
	tflog.Info(ctx, " Delete Start ***** *******")

	sendToweb_DelNetAddrGroupRequest(ctx, "DELETE", r.client, data.DelNetAddrGroupParameter)

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

func (r *NetAddrGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddNetAddrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddNetAddrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist"

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

func sendToweb_UpdateNetAddrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateNetAddrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist"

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

func sendToweb_DelNetAddrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelNetAddrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist"

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

func sendToweb_ReadNetAddrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadNetAddrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist?vsysName=PublicSystem&offset=0&count=25"

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
