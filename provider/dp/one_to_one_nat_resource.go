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

// 一对一NAT
var _ resource.Resource = &OneToOneNatResource{}
var _ resource.ResourceWithImportState = &OneToOneNatResource{}

func NewOneToOneNatResource() resource.Resource {
	return &OneToOneNatResource{}
}

type OneToOneNatResource struct {
	client *Client
}

type OneToOneNatResourceModel struct {
	AddOneToOneNatParameter    AddOneToOneNatParameter    `tfsdk:"addOneToOneNatParameter"`
	UpdateOneToOneNatParameter UpdateOneToOneNatParameter `tfsdk:"updateOneToOneNatParameter"`
	DelOneToOneNatParameter    DelOneToOneNatParameter    `tfsdk:"delOneToOneNatParameter"`
	ReadOneToOneNatParameter   ReadOneToOneNatParameter   `tfsdk:"readOneToOneNatParameter"`
}

type AddOneToOneNatParameter struct {
	VsysName            types.String `tfsdk:"vsysName"`
	Name                types.String `tfsdk:"name"`
	TargetName          types.String `tfsdk:"targetName"`
	Position            types.String `tfsdk:"position"`
	GlobalInterfaceName types.String `tfsdk:"globalInterfaceName"`
	GlobalAddress       types.String `tfsdk:"globalAddress"`
	LocalAddress        types.String `tfsdk:"localAddress"`
	VrrpIfName          types.String `tfsdk:"vrrpIfName"`
	VrrpId              types.String `tfsdk:"vrrpId"`
}

type UpdateOneToOneNatParameter struct {
	VsysName            types.String `tfsdk:"vsysName"`
	Name                types.String `tfsdk:"name"`
	OldName             types.String `tfsdk:"oldName"`
	TargetName          types.String `tfsdk:"targetName"`
	Position            types.String `tfsdk:"position"`
	GlobalInterfaceName types.String `tfsdk:"globalInterfaceName"`
	GlobalAddress       types.String `tfsdk:"globalAddress"`
	LocalAddress        types.String `tfsdk:"localAddress"`
	VrrpIfName          types.String `tfsdk:"vrrpIfName"`
	VrrpId              types.String `tfsdk:"vrrpId"`
}

type DelOneToOneNatParameter struct {
	VsysName     types.String `tfsdk:"vsysName"`
	Name         types.String `tfsdk:"name"`
	DelallEnable types.String `tfsdk:"delallEnable"`
}

type ReadOneToOneNatParameter struct {
	VsysName            types.String `tfsdk:"vsysName"`
	Offset              types.String `tfsdk:"offset"`
	Count               types.String `tfsdk:"count"`
	Name                types.String `tfsdk:"name"`
	GlobalInterfaceName types.String `tfsdk:"globalInterfaceName"`
	GlobalAddress       types.String `tfsdk:"globalAddress"`
	LocalAddress        types.String `tfsdk:"localAddress"`
	VrrpIfName          types.String `tfsdk:"vrrpIfName"`
	VrrpId              types.String `tfsdk:"vrrpId"`
	RuleId              types.String `tfsdk:"ruleId"`
	DelallEnable        types.String `tfsdk:"delallEnable"`
	Position            types.String `tfsdk:"position"`
	OldName             types.String `tfsdk:"oldName"`
	TargetName          types.String `tfsdk:"targetName"`
}

func (r *OneToOneNatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_TargetNat"
}

func (r *OneToOneNatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *OneToOneNatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OneToOneNatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OneToOneNatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddOneToOneNatRequest(ctx, "POST", r.client, data.AddOneToOneNatParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OneToOneNatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OneToOneNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	sendToweb_ReadOneToOneNatRequest(ctx, "GET", r.client, data.ReadOneToOneNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OneToOneNatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OneToOneNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateOneToOneNatRequest(ctx, "PUT", r.client, data.UpdateOneToOneNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OneToOneNatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OneToOneNatResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelOneToOneNatRequest(ctx, "DELETE", r.client, data.DelOneToOneNatParameter)

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

func (r *OneToOneNatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddOneToOneNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddOneToOneNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/nat1to1list"

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

func sendToweb_UpdateOneToOneNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateOneToOneNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/nat1to1list"

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

func sendToweb_DelOneToOneNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelOneToOneNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/nat1to1list"

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

func sendToweb_ReadOneToOneNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadOneToOneNatParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/nat1to1list"

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
