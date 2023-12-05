package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
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
	AddUsrGroupParameter AddUsrGroupParameter `tfsdk:"grp"`
}

type AddUsrGroupRequest struct {
	AddUsrGroupRequestModel AddUsrGroupRequestModel `json:"grp"`
}

// 调用接口参数
type AddUsrGroupRequestModel struct {
	Name           string `json:"name"`
	VfwName        string `json:"vfwName"`
	Desc           string `json:"desc"`
	AllSerNameList string `json:"allSerNameList"`
}

// 接收外部参数
type AddUsrGroupParameter struct {
	Name           types.String `tfsdk:"name"`
	VfwName        types.String `tfsdk:"vfwname"`
	Desc           types.String `tfsdk:"desc"`
	AllSerNameList types.String `tfsdk:"allsernamelist"`
}

func (r *UsrGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_UsrGroup"
}

func (r *UsrGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"grp": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"vfwname": schema.StringAttribute{
						Optional: true,
					},
					"desc": schema.StringAttribute{
						Optional: true,
					},
					"allsernamelist": schema.StringAttribute{
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
	sendToweb_UsrGroupRequest(ctx, "POST", r.client, data.AddUsrGroupParameter)
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
	//sendToweb_UsrGroupRequest(ctx, "GET", r.client, data.AddUsrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UsrGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_UsrGroupRequest(ctx, "PUT", r.client, data.AddUsrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UsrGroupResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_UsrGroupRequest(ctx, "DELETE", r.client, data.AddUsrGroupParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UsrGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_UsrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddUsrGroupParameter) {

	var sendData AddUsrGroupRequestModel
	if reqmethod == "POST" {
		sendData = AddUsrGroupRequestModel{
			Name:           Rsinfo.Name.ValueString(),
			VfwName:        Rsinfo.VfwName.ValueString(),
			AllSerNameList: Rsinfo.AllSerNameList.ValueString(),
			Desc:           Rsinfo.Desc.ValueString(),
		}
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}

	requstData := AddUsrGroupRequest{
		AddUsrGroupRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)

	tflog.Info(ctx, "请求体============:"+string(body))

	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/grp"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	// 创建一个HTTP客户端并发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	respn, err := client.Do(req)
	if err != nil {
		tflog.Error(ctx, "发送请求失败======="+err.Error())
		panic("发送请求失败=======")
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		tflog.Error(ctx, "发送请求失败======="+err2.Error())
		panic("发送请求失败=======")
	}
	// 打印响应结果
	tflog.Info(ctx, "响应状态码======="+string(respn.Status))
	tflog.Info(ctx, "响应体======="+string(body))

	if respn.Status != "200" || respn.Status != "201" || respn.Status != "204" {
		panic("请求响应失败=======")
	}

}
