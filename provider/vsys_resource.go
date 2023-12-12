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

// 虚拟系统
var _ resource.Resource = &VsysResource{}
var _ resource.ResourceWithImportState = &VsysResource{}

func NewVsysResource() resource.Resource {
	return &VsysResource{}
}

type VsysResource struct {
	client *Client
}

type VsysResourceModel struct {
	AddVsysParameter AddVsysParameter `tfsdk:"vsyslist"`
}

type AddVsysRequest struct {
	AddVsysRequestModel AddVsysRequestModel `json:"vsyslist"`
}

// 调用接口参数
type AddVsysRequestModel struct {
	VsysName string `json:"vsysName"`
	VsysType string `json:"vsysType"`
	VsysId   string `json:"vsysId"`
	VsysInfo string `json:"vsysInfo"`
}

// 接收外部参数
type AddVsysParameter struct {
	VsysName types.String `tfsdk:"vsysname"`
	VsysType types.String `tfsdk:"vsystype"`
	VsysId   types.String `tfsdk:"vsysid"`
	VsysInfo types.String `tfsdk:"vsysinfo"`
}

func (r *VsysResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_Vsys"
}

func (r *VsysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vsyslist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"vsysname": schema.StringAttribute{
						Required: true,
					},
					"vsystype": schema.StringAttribute{
						Required: true,
					},
					"vsysid": schema.StringAttribute{
						Optional: true,
					},
					"vsysinfo": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *VsysResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VsysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VsysResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_VsysRequest(ctx, "POST", r.client, data.AddVsysParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VsysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VsysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	//sendToweb_VsysRequest(ctx, "GET", r.client, data.AddVsysParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VsysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VsysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_VsysRequest(ctx, "PUT", r.client, data.AddVsysParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VsysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VsysResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_VsysRequest(ctx, "DELETE", r.client, data.AddVsysParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VsysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_VsysRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddVsysParameter) {

	var sendData AddVsysRequestModel
	if reqmethod == "POST" {
		sendData = AddVsysRequestModel{
			VsysName: Rsinfo.VsysName.ValueString(),
			VsysType: Rsinfo.VsysType.ValueString(),
			VsysId:   Rsinfo.VsysId.ValueString(),
			VsysInfo: Rsinfo.VsysInfo.ValueString(),
		}
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}

	requstData := AddVsysRequest{
		AddVsysRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)

	tflog.Info(ctx, "虚拟系统--请求体============:"+string(body))

	targetUrl := c.HostURL + "/func/web_main/api/vfw/vsyslist/vsyslist"

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
		tflog.Error(ctx, "虚拟系统--发送请求失败======="+err.Error())
		panic("虚拟系统--发送请求失败=======")
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		tflog.Error(ctx, "虚拟系统--发送请求失败======="+err2.Error())
		panic("虚拟系统--发送请求失败=======")
	}

	if strings.HasSuffix(respn.Status, "200") && strings.HasSuffix(respn.Status, "201") && strings.HasSuffix(respn.Status, "204") {
		tflog.Info(ctx, "虚拟系统--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "虚拟系统--响应体======="+string(body))
		panic("虚拟系统--请求响应失败=======")
	} else {
		// 打印响应结果
		tflog.Info(ctx, "虚拟系统--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "虚拟系统--响应体======="+string(body))
	}
}
