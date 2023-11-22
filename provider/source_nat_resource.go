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

// 源NAT
var _ resource.Resource = &SourceNatResource{}
var _ resource.ResourceWithImportState = &SourceNatResource{}

func NewSourceNatResource() resource.Resource {
	return &SourceNatResource{}
}

type SourceNatResource struct {
	client *Client
}

type SourceNatResourceModel struct {
	AddSourceNatParameter    AddSourceNatParameter    `tfsdk:"addSourceNatParameter"`
	UpdateSourceNatParameter UpdateSourceNatParameter `tfsdk:"updateSourceNatParameter"`
	DelSourceNatParameter    DelSourceNatParameter    `tfsdk:"delSourceNatParameter"`
	ReadSourceNatParameter   ReadSourceNatParameter   `tfsdk:"readSourceNatParameter"`
}

type AddSourceNatParameter struct {
	IpVersion           types.String `tfsdk:"ipVersion"`
	VsysName            types.String `tfsdk:"vsysName"`
	Name                types.String `tfsdk:"name"`
	TargetName          types.String `tfsdk:"targetName"`
	Position            types.String `tfsdk:"position"`
	OutInterface        types.String `tfsdk:"outInterface"`
	SrcAddrObj          types.String `tfsdk:"srcAddrObj"`
	SrcAddrGroup        types.String `tfsdk:"srcAddrGroup"`
	DstAddrObj          types.String `tfsdk:"dstAddrObj"`
	PreService          types.String `tfsdk:"preService"`
	UsrService          types.String `tfsdk:"usrService"`
	ServiceGroup        types.String `tfsdk:"serviceGroup"`
	PublicIpAddressFlag types.String `tfsdk:"publicIpAddressFlag"`
	AddrpoolName        types.String `tfsdk:"addrpoolName"`
	MinPort             types.String `tfsdk:"minPort"`
	MaxPort             types.String `tfsdk:"maxPort"`
	PortHash            types.String `tfsdk:"portHash"`
	State               types.String `tfsdk:"state"`
}

type UpdateSourceNatParameter struct {
	IpVersion           types.String `tfsdk:"ipVersion"`
	VsysName            types.String `tfsdk:"vsysName"`
	OldName             types.String `tfsdk:"oldName"`
	TargetName          types.String `tfsdk:"targetName"`
	Position            types.String `tfsdk:"position"`
	OutInterface        types.String `tfsdk:"outInterface"`
	SrcAddrObj          types.String `tfsdk:"srcAddrObj"`
	SrcAddrGroup        types.String `tfsdk:"srcAddrGroup"`
	DstAddrObj          types.String `tfsdk:"dstAddrObj"`
	DstAddrGroup        types.String `tfsdk:"dstAddrGroup"`
	PreService          types.String `tfsdk:"preService"`
	UsrService          types.String `tfsdk:"usrService"`
	ServiceGroup        types.String `tfsdk:"serviceGroup"`
	PublicIpAddressFlag types.String `tfsdk:"publicIpAddressFlag"`
	AddrpoolName        types.String `tfsdk:"addrpoolName"`
	MinPort             types.String `tfsdk:"minPort"`
	MaxPort             types.String `tfsdk:"maxPort"`
	PortHash            types.String `tfsdk:"portHash"`
	State               types.String `tfsdk:"state"`
}

type DelSourceNatParameter struct {
	IpVersion    types.String `tfsdk:"ipVersion"`
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	DelallEnable types.String `tfsdk:"delallEnable"`
}

type ReadSourceNatParameter struct {
	IpVersion           types.String `tfsdk:"ipVersion"`
	VsysName            types.String `tfsdk:"vsysName"`
	Offset              types.String `tfsdk:"offset"`
	Count               types.String `tfsdk:"count"`
	Name                types.String `tfsdk:"name"`
	OutInterface        types.String `tfsdk:"outInterface"`
	SourceIp            types.String `tfsdk:"sourceIp"`
	DestinationIp       types.String `tfsdk:"destinationIp"`
	Service             types.String `tfsdk:"service"`
	State               types.String `tfsdk:"state"`
	SrcAddrObj          types.String `tfsdk:"srcAddrObj"`
	SrcAddrGroup        types.String `tfsdk:"srcAddrGroup"`
	DstAddrObj          types.String `tfsdk:"dstAddrObj"`
	DstAddrGroup        types.String `tfsdk:"dstAddrGroup"`
	PreService          types.String `tfsdk:"preService"`
	UsrService          types.String `tfsdk:"usrService"`
	ServiceGroup        types.String `tfsdk:"serviceGroup"`
	PublicIpAddressFlag types.String `tfsdk:"publicIpAddressFlag"`
	AddrpoolName        types.String `tfsdk:"addrpoolName"`
	MinPort             types.String `tfsdk:"minPort"`
	MaxPort             types.String `tfsdk:"maxPort"`
	PortHash            types.String `tfsdk:"portHash"`
	RuleId              types.String `tfsdk:"ruleId"`
	DelallEnable        types.String `tfsdk:"delallEnable"`
	TargetName          types.String `tfsdk:"targetName"`
	OldName             types.String `tfsdk:"oldName"`
	Position            types.String `tfsdk:"position"`
}

func (r *SourceNatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "firewall_SourceNat"
}

func (r *SourceNatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sourcenat": schema.SingleNestedAttribute{
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

func (r *SourceNatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *SourceNatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

func (r *SourceNatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	sendToweb_ReadSourceNatRequest(ctx, "GET", r.client, data.ReadSourceNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceNatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateSourceNatRequest(ctx, "PUT", r.client, data.UpdateSourceNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceNatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r *SourceNatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddSourceNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddSourceNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/snatlist"

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

func sendToweb_UpdateSourceNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateSourceNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/snatlist"

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

func sendToweb_DelSourceNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelSourceNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/snatlist"

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

func sendToweb_ReadSourceNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadSourceNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/snatlist"

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
