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

// 目的NAT
var _ resource.Resource = &TargetNatResource{}
var _ resource.ResourceWithImportState = &TargetNatResource{}

func NewTargetNatResource() resource.Resource {
	return &TargetNatResource{}
}

type TargetNatResource struct {
	client *Client
}

type TargetNatResourceModel struct {
	AddSourceNatParameter    AddSourceNatParameter    `tfsdk:"addSourceNatParameter"`
	UpdateSourceNatParameter UpdateSourceNatParameter `tfsdk:"updateSourceNatParameter"`
	DelSourceNatParameter    DelSourceNatParameter    `tfsdk:"delSourceNatParameter"`
}

type AddTargetNatParameter struct {
	Port                 types.String `tfsdk:"port"`
	Ip                   types.String `tfsdk:"ip"`
	VsysName             types.String `tfsdk:"vsysName"`
	Name                 types.String `tfsdk:"name"`
	TargetName           types.String `tfsdk:"targetName"`
	Position             types.String `tfsdk:"position"`
	InInterface          types.String `tfsdk:"inInterface"`
	SrcIpObj             types.String `tfsdk:"srcIpObj"`
	SrcIpGroup           types.String `tfsdk:"srcIpGroup"`
	PublicIp             types.String `tfsdk:"publicIp"`
	PreService           types.String `tfsdk:"preService"`
	UsrService           types.String `tfsdk:"usrService"`
	InNetIp              types.String `tfsdk:"inNetIp"`
	InnetPort            types.String `tfsdk:"innetPort"`
	UnLimited            types.String `tfsdk:"unLimited"`
	SrcIpTranslate       types.String `tfsdk:"srcIpTranslate"`
	InterfaceAddressFlag types.String `tfsdk:"interfaceAddressFlag"`
	AddrpoolName         types.String `tfsdk:"addrpoolName"`
	VrrpIfName           types.String `tfsdk:"vrrpIfName"`
	VrrpId               types.String `tfsdk:"vrrpId"`
	State                types.String `tfsdk:"state"`
}

type UpdateTargetNatParameter struct {
	Port                 types.String `tfsdk:"port"`
	Ip                   types.String `tfsdk:"ip"`
	VsysName             types.String `tfsdk:"vsysName"`
	OldName              types.String `tfsdk:"oldName"`
	TargetName           types.String `tfsdk:"targetName"`
	Position             types.String `tfsdk:"position"`
	InInterface          types.String `tfsdk:"inInterface"`
	NetaddrObj           types.String `tfsdk:"netaddrObj"`
	NetaddrGroup         types.String `tfsdk:"netaddrGroup"`
	PublicIp             types.String `tfsdk:"publicIp"`
	PreService           types.String `tfsdk:"preService"`
	UsrService           types.String `tfsdk:"usrService"`
	InNetIp              types.String `tfsdk:"inNetIp"`
	InnetPort            types.String `tfsdk:"innetPort"`
	UnLimited            types.String `tfsdk:"unLimited"`
	SrcIpTranslate       types.String `tfsdk:"srcIpTranslate"`
	InterfaceAddressFlag types.String `tfsdk:"interfaceAddressFlag"`
	AddrpoolName         types.String `tfsdk:"addrpoolName"`
	VrrpIfName           types.String `tfsdk:"vrrpIfName"`
	VrrpId               types.String `tfsdk:"vrrpId"`
	State                types.String `tfsdk:"state"`
}

type DelTargetNatParameter struct {
	Port         types.String `tfsdk:"port"`
	Ip           types.String `tfsdk:"ip"`
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	DelallEnable types.String `tfsdk:"delallEnable"`
}

func (r *TargetNatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-TargetNat"
}

func (r *TargetNatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *TargetNatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TargetNatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddSourceNatRequest(ctx, "POST", r.client, data.AddSourceNatParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TargetNatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	// sendToweb_AddSourceNatRequest(ctx,"POST", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TargetNatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateSourceNatRequest(ctx, "PUT", r.client, data.UpdateSourceNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TargetNatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SourceNatResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelSourceNatRequest(ctx, "DELETE", r.client, data.DelSourceNatParameter)

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

func (r *TargetNatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddTargetNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddSourceNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/dnatlist"

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

func sendToweb_UpdateTargetNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateSourceNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/dnatlist"

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

func sendToweb_DelTargetNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelSourceNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/dnatlist"

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
