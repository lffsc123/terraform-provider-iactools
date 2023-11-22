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

// IPSec策略方式
var _ resource.Resource = &NetConnResource{}
var _ resource.ResourceWithImportState = &NetConnResource{}

func NewNetConnResource() resource.Resource {
	return &NetConnResource{}
}

type NetConnResource struct {
	client *Client
}

type NetConnResourceModel struct {
	AddNetConnParameter    AddNetConnParameter    `tfsdk:"addNetConnParameter"`
	UpdateNetConnParameter UpdateNetConnParameter `tfsdk:"updateNetConnParameter"`
	DelNetConnParameter    DelNetConnParameter    `tfsdk:"delNetConnParameter"`
	ReadNetConnParameter   ReadNetConnParameter   `tfsdk:"readNetConnParameter"`
}

type AddNetConnParameter struct {
	VsysName      types.String `tfsdk:"vsysName"`
	Name          types.String `tfsdk:"name"`
	State         types.String `tfsdk:"state"`
	Ifname        types.String `tfsdk:"ifname"`
	SaName        types.String `tfsdk:"saName"`
	LocalIp       types.String `tfsdk:"localIp"`
	RemoteIp      types.String `tfsdk:"remoteIp"`
	LocalIdType   types.String `tfsdk:"localIdType"`
	LocalId       types.String `tfsdk:"localId"`
	RemoteIdType  types.String `tfsdk:"remoteIdType"`
	RemoteId      types.String `tfsdk:"remoteId"`
	Subnet        types.String `tfsdk:"subnet"`
	Psk           types.String `tfsdk:"psk"`
	PolicyMode    types.String `tfsdk:"policyMode"`
	PolicyAct     types.String `tfsdk:"policyAct"`
	Antiaction    types.String `tfsdk:"antiaction"`
	Antireplay    types.String `tfsdk:"antireplay"`
	Accesslocal   types.String `tfsdk:"accesslocal"`
	Aggrmode      types.String `tfsdk:"aggrmode"`
	Routeenable   types.String `tfsdk:"routeenable"`
	Routemode     types.String `tfsdk:"routemode"`
	Routenexthop  types.String `tfsdk:"routenexthop"`
	Routepriority types.String `tfsdk:"routepriority"`
	Dpdstatus     types.String `tfsdk:"dpdstatus"`
	Dpddelay      types.String `tfsdk:"dpddelay"`
	Dpdtimeout    types.String `tfsdk:"dpdtimeout"`
}

type UpdateNetConnParameter struct {
	VsysName      types.String `tfsdk:"vsysName"`
	Name          types.String `tfsdk:"name"`
	State         types.String `tfsdk:"state"`
	Ifname        types.String `tfsdk:"ifname"`
	SaName        types.String `tfsdk:"saName"`
	LocalIp       types.String `tfsdk:"localIp"`
	RemoteIp      types.String `tfsdk:"remoteIp"`
	LocalIdType   types.String `tfsdk:"localIdType"`
	LocalId       types.String `tfsdk:"localId"`
	RemoteIdType  types.String `tfsdk:"remoteIdType"`
	RemoteId      types.String `tfsdk:"remoteId"`
	Subnet        types.String `tfsdk:"subnet"`
	Psk           types.String `tfsdk:"psk"`
	PolicyMode    types.String `tfsdk:"policyMode"`
	PolicyAct     types.String `tfsdk:"policyAct"`
	Antiaction    types.String `tfsdk:"antiaction"`
	Antireplay    types.String `tfsdk:"antireplay"`
	Accesslocal   types.String `tfsdk:"accesslocal"`
	Aggrmode      types.String `tfsdk:"aggrmode"`
	Routeenable   types.String `tfsdk:"routeenable"`
	Routemode     types.String `tfsdk:"routemode"`
	Routenexthop  types.String `tfsdk:"routenexthop"`
	Routepriority types.String `tfsdk:"routepriority"`
	Dpdstatus     types.String `tfsdk:"dpdstatus"`
	Dpddelay      types.String `tfsdk:"dpddelay"`
	Dpdtimeout    types.String `tfsdk:"dpdtimeout"`
}

type DelNetConnParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
	Name     types.String `tfsdk:"name"`
}

type ReadNetConnParameter struct {
	VsysName      types.String `tfsdk:"vsysName"`
	Name          types.String `tfsdk:"name"`
	State         types.String `tfsdk:"state"`
	Ifname        types.String `tfsdk:"ifname"`
	SaName        types.String `tfsdk:"saName"`
	LocalIp       types.String `tfsdk:"localIp"`
	RemoteIp      types.String `tfsdk:"remoteIp"`
	LocalIdType   types.String `tfsdk:"localIdType"`
	LocalId       types.String `tfsdk:"localId"`
	RemoteIdType  types.String `tfsdk:"remoteIdType"`
	RemoteId      types.String `tfsdk:"remoteId"`
	Subnet        types.String `tfsdk:"subnet"`
	Psk           types.String `tfsdk:"psk"`
	PolicyMode    types.String `tfsdk:"policyMode"`
	PolicyAct     types.String `tfsdk:"policyAct"`
	Antiaction    types.String `tfsdk:"antiaction"`
	Antireplay    types.String `tfsdk:"antireplay"`
	Accesslocal   types.String `tfsdk:"accesslocal"`
	Aggrmode      types.String `tfsdk:"aggrmode"`
	Routeenable   types.String `tfsdk:"routeenable"`
	Routemode     types.String `tfsdk:"routemode"`
	Routenexthop  types.String `tfsdk:"routenexthop"`
	Routepriority types.String `tfsdk:"routepriority"`
	Dpdstatus     types.String `tfsdk:"dpdstatus"`
	Dpddelay      types.String `tfsdk:"dpddelay"`
	Dpdtimeout    types.String `tfsdk:"dpdtimeout"`
}

func (r *NetConnResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_NetConn"
}

func (r *NetConnResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *NetConnResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetConnResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetConnResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddNetConnRequest(ctx, "POST", r.client, data.AddNetConnParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetConnResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetConnResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadNetConnRequest(ctx, "GET", r.client, data.ReadNetConnParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetConnResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetConnResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateNetConnRequest(ctx, "PUT", r.client, data.UpdateNetConnParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetConnResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetConnResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelNetConnRequest(ctx, "DELETE", r.client, data.DelNetConnParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NetConnResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddNetConnRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddNetConnParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/netconnlist"

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

func sendToweb_UpdateNetConnRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateNetConnParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/netconnlist"

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

func sendToweb_DelNetConnRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelNetConnParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/netconnlist"

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

func sendToweb_ReadNetConnRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadNetConnParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ipsec_vpn/ipsec/netconnlist"

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
