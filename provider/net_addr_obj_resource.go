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

// IP 地址对象
var _ resource.Resource = &NetAddrObjResource{}
var _ resource.ResourceWithImportState = &NetAddrObjResource{}

func NewNetAddrObjResource() resource.Resource {
	return &NetAddrObjResource{}
}

type NetAddrObjResource struct {
	client *Client
}

type NetAddrObjResourceModel struct {
	AddNetAddrObjParameter AddNetAddrObjParameter `tfsdk:"netaddrobjlist"`
}

type AddNetAddrObjRequest struct {
	AddNetAddrObjRequestModel AddNetAddrObjRequestModel `json:"netaddrobjlist"`
}

// 调用接口参数
type AddNetAddrObjRequestModel struct {
	IpVersion string `json:"ipVersion"`
	VsysName  string `json:"vsysName"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Ip        string `json:"ip"`
	ExpIp     string `json:"expIp"`
}

// 接收外部参数
type AddNetAddrObjParameter struct {
	IpVersion types.String `tfsdk:"ipversion"`
	VsysName  types.String `tfsdk:"vsysname"`
	Name      types.String `tfsdk:"name"`
	Desc      types.String `tfsdk:"desc"`
	Ip        types.String `tfsdk:"ip"`
	ExpIp     types.String `tfsdk:"expip"`
}

func (r *NetAddrObjResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_NetAddrObj"
}

func (r *NetAddrObjResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"netaddrobjlist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"ipversion": schema.StringAttribute{
						Optional: true,
					},
					"vsysname": schema.StringAttribute{
						Optional: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"desc": schema.StringAttribute{
						Optional: true,
					},
					"ip": schema.StringAttribute{
						Optional: true,
					},
					"expip": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *NetAddrObjResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetAddrObjResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetAddrObjResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource ****** ********")
	sendToweb_NetAddrObjRequest(ctx, "POST", r.client, data.AddNetAddrObjParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrObjResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetAddrObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	//sendToweb_NetAddrObjRequest(ctx, "GET", r.client, data.AddNetAddrObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrObjResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetAddrObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_NetAddrObjRequest(ctx, "PUT", r.client, data.AddNetAddrObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrObjResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetAddrObjResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_NetAddrObjRequest(ctx, "DELETE", r.client, data.AddNetAddrObjParameter)

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

func (r *NetAddrObjResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_NetAddrObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddNetAddrObjParameter) {

	var sendData AddNetAddrObjRequestModel
	if reqmethod == "POST" {
		sendData = AddNetAddrObjRequestModel{
			IpVersion: Rsinfo.IpVersion.ValueString(),
			Name:      Rsinfo.Name.ValueString(),
			VsysName:  Rsinfo.VsysName.ValueString(),
			Ip:        Rsinfo.Ip.ValueString(),
			Desc:      Rsinfo.Desc.ValueString(),
			ExpIp:     Rsinfo.ExpIp.ValueString(),
		}
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}

	requstData := AddNetAddrObjRequest{
		AddNetAddrObjRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)

	tflog.Info(ctx, "IP地址对象--请求体============:"+string(body)+"======")

	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist"

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
		tflog.Error(ctx, "IP地址对象--发送请求失败======="+err.Error())
		panic("IP地址对象--发送请求失败=======")
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		tflog.Error(ctx, "IP地址对象--发送请求失败======="+err2.Error())
		panic("IP地址对象--发送请求失败=======")
	}

	if strings.HasSuffix(respn.Status, "200") && strings.HasSuffix(respn.Status, "201") && strings.HasSuffix(respn.Status, "204") {
		tflog.Info(ctx, "IP地址对象--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "IP地址对象--响应体======="+string(body)+"======")
		panic("IP地址对象--请求响应失败=======")
	} else {
		// 打印响应结果
		tflog.Info(ctx, "IP地址对象--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "IP地址对象--响应体======="+string(body)+"======")
	}
}
