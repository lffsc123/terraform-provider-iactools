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

// 黑名单
var _ resource.Resource = &BlackNameCfgResource{}
var _ resource.ResourceWithImportState = &BlackNameCfgResource{}

func NewBlackNameCfgResource() resource.Resource {
	return &BlackNameCfgResource{}
}

type BlackNameCfgResource struct {
	client *Client
}

type BlackNameCfgResourceModel struct {
	AddBlackNameCfgParameter    AddBlackNameCfgParameter    `tfsdk:"addBlackNameCfgParameter"`
	UpdateBlackNameCfgParameter UpdateBlackNameCfgParameter `tfsdk:"updateBlackNameCfgParameter"`
	DelBlackNameCfgParameter    DelBlackNameCfgParameter    `tfsdk:"delBlackNameCfgParameter"`
	ReadBlackNameCfgParameter   ReadBlackNameCfgParameter   `tfsdk:"readBlackNameCfgParameter"`
}

type AddBlackNameCfgParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	BoardType types.String `tfsdk:"boardType"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	Age       types.String `tfsdk:"age"`
	State     types.String `tfsdk:"state"`
}

type UpdateBlackNameCfgParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	BoardType types.String `tfsdk:"boardType"`
	OldIp     types.String `tfsdk:"oldIp"`
	OldMask   types.String `tfsdk:"oldMask"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	Age       types.String `tfsdk:"age"`
	State     types.String `tfsdk:"state"`
}

type DelBlackNameCfgParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	BoardType types.String `tfsdk:"boardType"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
}

type ReadBlackNameCfgParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	BoardType types.String `tfsdk:"boardType"`
	Id        types.String `tfsdk:"id"`
	Ip        types.String `tfsdk:"ip"`
	Mask      types.String `tfsdk:"mask"`
	CreatTime types.String `tfsdk:"creatTime"`
	LeftAge   types.String `tfsdk:"leftAge"`
	Age       types.String `tfsdk:"age"`
	State     types.String `tfsdk:"state"`
}

func (r *BlackNameCfgResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_BlackNameCfg"
}

func (r *BlackNameCfgResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"param": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *BlackNameCfgResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BlackNameCfgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *BlackNameCfgResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddBlackNameCfgRequest(ctx, "POST", r.client, data.AddBlackNameCfgParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlackNameCfgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *BlackNameCfgResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadBlackNameCfgRequest(ctx, "GET", r.client, data.ReadBlackNameCfgParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlackNameCfgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *BlackNameCfgResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateBlackNameCfgRequest(ctx, "PUT", r.client, data.UpdateBlackNameCfgParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlackNameCfgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *BlackNameCfgResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelBlackNameCfgRequest(ctx, "DELETE", r.client, data.DelBlackNameCfgParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BlackNameCfgResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddBlackNameCfgRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddBlackNameCfgParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/blackname/blackname_cfg/blackname_cfg/blackNameList"

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

func sendToweb_UpdateBlackNameCfgRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateBlackNameCfgParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/blackname/blackname_cfg/blackname_cfg/blackNameList"

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

func sendToweb_DelBlackNameCfgRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelBlackNameCfgParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/blackname/blackname_cfg/blackname_cfg/blackNameList"

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

func sendToweb_ReadBlackNameCfgRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadBlackNameCfgParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/blackname/blackname_cfg/blackname_cfg/blackNameList"

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
