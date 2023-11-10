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

// 服务对象
var _ resource.Resource = &UsrObjResource{}
var _ resource.ResourceWithImportState = &UsrObjResource{}

func NewUsrObjResource() resource.Resource {
	return &UsrObjResource{}
}

type UsrObjResource struct {
	client *Client
}

type UsrObjResourceModel struct {
	AddUsrObjParameter    AddUsrObjParameter    `tfsdk:"addUsrObjParameter"`
	UpdateUsrObjParameter UpdateUsrObjParameter `tfsdk:"updateUsrObjParameter"`
	DelUsrObjParameter    DelUsrObjParameter    `tfsdk:"delUsrObjParameter"`
	ReadUsrObjParameter   ReadUsrObjParameter   `tfsdk:"readUsrObjParameter"`
}

type AddUsrObjParameter struct {
	Name       types.String `tfsdk:"name"`
	VfwName    types.String `tfsdk:"vfwName"`
	Protocol   types.String `tfsdk:"protocol"`
	SportStart types.String `tfsdk:"sportStart"`
	SportEnd   types.String `tfsdk:"sportEnd"`
	DportStart types.String `tfsdk:"dportStart"`
	DportEnd   types.String `tfsdk:"dportEnd"`
	Services   types.String `tfsdk:"services"`
	Desc       types.String `tfsdk:"desc"`
}

type UpdateUsrObjParameter struct {
	Name       types.String `tfsdk:"name"`
	VfwName    types.String `tfsdk:"vfwName"`
	OldName    types.String `tfsdk:"oldName"`
	Protocol   types.String `tfsdk:"protocol"`
	SportStart types.String `tfsdk:"sportStart"`
	SportEnd   types.String `tfsdk:"sportEnd"`
	DportStart types.String `tfsdk:"dportStart"`
	DportEnd   types.String `tfsdk:"dportEnd"`
	Services   types.String `tfsdk:"services"`
	Desc       types.String `tfsdk:"desc"`
}

type DelUsrObjParameter struct {
	Name         types.String `tfsdk:"name"`
	VfwName      types.String `tfsdk:"vfwName"`
	DelAllEnable types.String `tfsdk:"delAllEnable"`
}

type ReadUsrObjParameter struct {
	Id           types.String `tfsdk:"id"`
	VfwName      types.String `tfsdk:"vfwName"`
	Name         types.String `tfsdk:"name"`
	OldName      types.String `tfsdk:"oldName"`
	Protocol     types.String `tfsdk:"protocol"`
	SportStart   types.String `tfsdk:"sportStart"`
	SportEnd     types.String `tfsdk:"sportEnd"`
	DportStart   types.String `tfsdk:"dportStart"`
	DportEnd     types.String `tfsdk:"dportEnd"`
	Services     types.String `tfsdk:"services"`
	ReferNum     types.String `tfsdk:"referNum"`
	DelAllEnable types.String `tfsdk:"delAllEnable"`
	SearchValue  types.String `tfsdk:"searchValue"`
	Offset       types.String `tfsdk:"offset"`
	Count        types.String `tfsdk:"count"`
}

func (r *UsrObjResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_UsrObj"
}

func (r *UsrObjResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *UsrObjResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UsrObjResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *UsrObjResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddUsrObjRequest(ctx, "POST", r.client, data.AddUsrObjParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrObjResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UsrObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadUsrObjRequest(ctx, "GET", r.client, data.ReadUsrObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrObjResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UsrObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateUsrObjRequest(ctx, "PUT", r.client, data.UpdateUsrObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrObjResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UsrObjResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelUsrObjRequest(ctx, "DELETE", r.client, data.DelUsrObjParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UsrObjResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddUsrObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddUsrObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj"

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

func sendToweb_UpdateUsrObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateUsrObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj"

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

func sendToweb_DelUsrObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelUsrObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj"

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

func sendToweb_ReadUsrObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadUsrObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj?vfwName=vsys&searchValue=&offset=1&count=100"

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
