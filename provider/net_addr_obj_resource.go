package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
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

type UpdateNetAddrObjRequest struct {
	UpdateNetAddrObjRequestModel UpdateNetAddrObjRequestModel `json:"netaddrobjlist"`
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

type UpdateNetAddrObjRequestModel struct {
	IpVersion string `json:"ipVersion"`
	VsysName  string `json:"vsysName"`
	Name      string `json:"name"`
	OldName   string `json:"oldName"`
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

// 查询结果结构体
type QueryNetAddrObjResponseListModel struct {
	Netaddrobjlist []QueryNetAddrObjResponseModel `json:"netaddrobjlist"`
}
type QueryNetAddrObjResponseModel struct {
	IpVersion    string `json:"ipVersion"`
	Vsysname     string `json:"vsysName"`
	Offset       string `json:"offset"`
	Count        string `json:"count"`
	SearchValue  string `json:"searchValue"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	OldName      string `json:"oldName"`
	Desc         string `json:"desc"`
	Ip           string `json:"ip"`
	ExpIp        string `json:"expIp"`
	ReferNum     string `json:"referNum"`
	DelAllEnable string `json:"delAllEnable"`
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

	if reqmethod == "POST" {
		// 先查询是否存在，再执行新增操作
		//tflog.Info(ctx, "IP地址对象--开始执行--查询操作")
		//responseBody := sendRequest(ctx, "GET", c, nil, "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist?vsysName=PublicSystem&offset=0&count=1000", "IP地址对象")
		//var queryResList QueryNetAddrObjResponseListModel
		//err := json.Unmarshal([]byte(responseBody), &queryResList)
		//if err != nil {
		//	panic("转换查询结果json出现异常")
		//}
		//for _, queryRes := range queryResList.Netaddrobjlist {
		//	if queryRes.Name == Rsinfo.Name.ValueString() {
		//		tflog.Info(ctx, "IP地址对象--存在重复数据，执行--修改操作")
		//		var sendUpdateData UpdateNetAddrObjRequestModel
		//		sendUpdateData = UpdateNetAddrObjRequestModel{
		//			IpVersion: Rsinfo.IpVersion.ValueString(),
		//			VsysName:  Rsinfo.VsysName.ValueString(),
		//			Name:      Rsinfo.Name.ValueString(),
		//			OldName:   Rsinfo.Name.ValueString(),
		//			Desc:      Rsinfo.Desc.ValueString(),
		//			Ip:        Rsinfo.Ip.ValueString(),
		//			ExpIp:     Rsinfo.ExpIp.ValueString(),
		//		}
		//
		//		requstUpdateData := UpdateNetAddrObjRequest{
		//			UpdateNetAddrObjRequestModel: sendUpdateData,
		//		}
		//		body, _ := json.Marshal(requstUpdateData)
		//
		//		sendRequest(ctx, "PUT", c, body, "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist", "IP地址对象")
		//		return
		//	}
		//}

		// 新增操作
		var sendData AddNetAddrObjRequestModel
		sendData = AddNetAddrObjRequestModel{
			IpVersion: Rsinfo.IpVersion.ValueString(),
			Name:      Rsinfo.Name.ValueString(),
			VsysName:  Rsinfo.VsysName.ValueString(),
			Ip:        Rsinfo.Ip.ValueString(),
			Desc:      Rsinfo.Desc.ValueString(),
			ExpIp:     Rsinfo.ExpIp.ValueString(),
		}
		requstData := AddNetAddrObjRequest{
			AddNetAddrObjRequestModel: sendData,
		}
		body, _ := json.Marshal(requstData)

		responseBody := sendRequest(ctx, reqmethod, c, body, "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist", "IP地址对象")
		if strings.Contains(responseBody, "already exists") {
			tflog.Info(ctx, "IP地址对象--存在重复数据，执行--修改操作")
			// 更新操作
			var sendUpdateData UpdateNetAddrObjRequestModel
			sendUpdateData = UpdateNetAddrObjRequestModel{
				IpVersion: Rsinfo.IpVersion.ValueString(),
				VsysName:  Rsinfo.VsysName.ValueString(),
				Name:      Rsinfo.Name.ValueString(),
				OldName:   Rsinfo.Name.ValueString(),
				Desc:      Rsinfo.Desc.ValueString(),
				Ip:        Rsinfo.Ip.ValueString(),
				ExpIp:     Rsinfo.ExpIp.ValueString(),
			}

			requstUpdateData := UpdateNetAddrObjRequest{
				UpdateNetAddrObjRequestModel: sendUpdateData,
			}
			body, _ := json.Marshal(requstUpdateData)

			sendRequest(ctx, "PUT", c, body, "/func/web_main/api/netaddr/netaddr_obj/netaddrobjlist", "IP地址对象")
			return
		}
		return
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}
}
