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

// ip地址组
var _ resource.Resource = &NetAddrGroupResource{}
var _ resource.ResourceWithImportState = &NetAddrGroupResource{}

func NewNetAddrGroupResource() resource.Resource {
	return &NetAddrGroupResource{}
}

type NetAddrGroupResource struct {
	client *Client
}

type NetAddrGroupResourceModel struct {
	AddNetAddrGroupParameter AddNetAddrGroupParameter `tfsdk:"netaddrgrplist"`
}

type AddAddrGroupRequest struct {
	AddAddrGroupRequestModel AddAddrGroupRequestModel `json:"netaddrgrplist"`
}

// 调用接口参数
type AddAddrGroupRequestModel struct {
	VsysName    string `json:"vsysName"`
	Name        string `json:"name"`
	ObjNameList string `json:"objNameList"`
	Desc        string `json:"desc"`
}

// 接收外部参数
type AddNetAddrGroupParameter struct {
	VsysName    types.String `tfsdk:"vsysname"`
	Name        types.String `tfsdk:"name"`
	ObjNameList types.String `tfsdk:"objnamelist"`
	Desc        types.String `tfsdk:"desc"`
}

func (r *NetAddrGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_NetAddrGroup"
}

func (r *NetAddrGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"netaddrgrplist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"vsysname": schema.StringAttribute{
						Optional: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"objnamelist": schema.StringAttribute{
						Optional: true,
					},
					"desc": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *NetAddrGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetAddrGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetAddrGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_NetAddrGroupRequest(ctx, "POST", r.client, data.AddNetAddrGroupParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetAddrGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	//sendToweb_NetAddrGroupRequest(ctx, "GET", r.client, data.AddNetAddrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetAddrGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_NetAddrGroupRequest(ctx, "PUT", r.client, data.AddNetAddrGroupParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetAddrGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetAddrGroupResourceModel
	tflog.Info(ctx, " Delete Start ***** *******")

	//sendToweb_NetAddrGroupRequest(ctx, "DELETE", r.client, data.AddNetAddrGroupParameter)

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

func (r *NetAddrGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_NetAddrGroupRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddNetAddrGroupParameter) {

	var sendData AddAddrGroupRequestModel
	if reqmethod == "POST" {
		sendData = AddAddrGroupRequestModel{
			VsysName:    Rsinfo.VsysName.ValueString(),
			Name:        Rsinfo.Name.ValueString(),
			ObjNameList: Rsinfo.ObjNameList.ValueString(),
			Desc:        Rsinfo.Desc.ValueString(),
		}
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}

	requstData := AddAddrGroupRequest{
		AddAddrGroupRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)

	tflog.Info(ctx, "ip地址组--请求体============:"+string(body))

	targetUrl := c.HostURL + "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist"

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
		tflog.Error(ctx, "ip地址组--发送请求失败======="+err.Error())
		panic("ip地址组--发送请求失败=======")
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		tflog.Error(ctx, "ip地址组--发送请求失败======="+err2.Error())
		panic("ip地址组--发送请求失败=======")
	}

	if respn.Status != "200" && respn.Status != "201" && respn.Status != "204" {
		tflog.Info(ctx, "ip地址组--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "ip地址组--响应体======="+string(body))
		panic("ip地址组--请求响应失败=======")
	} else {
		// 打印响应结果
		tflog.Info(ctx, "ip地址组--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "ip地址组--响应体======="+string(body))
	}

}
