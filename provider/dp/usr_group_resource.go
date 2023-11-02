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

// 服务组
var _ resource.Resource = &UsrGroupResource{}
var _ resource.ResourceWithImportState = &UsrGroupResource{}

func NewUsrGroupResource() resource.Resource {
	return &UsrGroupResource{}
}

type UsrGroupResource struct {
	client *Client
}

type UsrGroupResourceModel struct {
	AddUsrGroupParameter    AddUsrGroupParameter    `tfsdk:"addUsrGroupParameter"`
	UpdateUsrGroupParameter UpdateUsrGroupParameter `tfsdk:"updateUsrGroupParameter"`
	DelUsrGroupParameter    DelUsrGroupParameter    `tfsdk:"delUsrGroupParameter"`
	ReadUsrGroupParameter   ReadUsrGroupParameter   `tfsdk:"readUsrGroupParameter"`
}

type AddUsrGroupParameter struct {
	Name           types.String `tfsdk:"name"`
	VfwName        types.String `tfsdk:"vfwName"`
	Desc           types.String `tfsdk:"desc"`
	AllSerNameList types.String `tfsdk:"allSerNameList"`
}

type UpdateUsrGroupParameter struct {
	Name           types.String `tfsdk:"name"`
	VfwName        types.String `tfsdk:"vfwName"`
	OldName        types.String `tfsdk:"oldName"`
	Desc           types.String `tfsdk:"desc"`
	AllSerNameList types.String `tfsdk:"allSerNameList"`
}

type DelUsrGroupParameter struct {
	Name         types.String `tfsdk:"name"`
	VfwName      types.String `tfsdk:"vfwName"`
	DelAllEnable types.String `tfsdk:"delAllEnable"`
}

type ReadUsrGroupParameter struct {
	Id             types.String `tfsdk:"id"`
	VfwName        types.String `tfsdk:"vfwName"`
	Name           types.String `tfsdk:"name"`
	OldName        types.String `tfsdk:"oldName"`
	Desc           types.String `tfsdk:"desc"`
	AllSerNameList types.String `tfsdk:"allSerNameList"`
	PreNameList    types.String `tfsdk:"preNameList"`
	UsrNameList    types.String `tfsdk:"usrNameList"`
	ReferNum       types.String `tfsdk:"referNum"`
	DelAllEnable   types.String `tfsdk:"delAllEnable"`
	SearchValue    types.String `tfsdk:"searchValue"`
	Offset         types.String `tfsdk:"offset"`
	Count          types.String `tfsdk:"count"`
}

func (r *UsrGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo-UsrGroup"
}

func (r *UsrGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *UsrGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UsrGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *UsrGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddUsrGroupRequest(ctx, "POST", r.client, data.AddUsrGroupParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UsrGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadUsrGroupRequest(ctx, "GET", r.client, data.ReadUsrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UsrGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateUsrGroupRequest(ctx, "PUT", r.client, data.UpdateUsrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UsrGroupResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelUsrGroupRequest(ctx, "DELETE", r.client, data.DelUsrGroupParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UsrGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddUsrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddUsrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/grp"

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

func sendToweb_UpdateUsrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateUsrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/grp"

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

func sendToweb_DelUsrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelUsrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/grp"

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

func sendToweb_ReadUsrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadUsrGroupParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/grp?vfwName=vsys&searchValue=&offset=1&count=100"

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
