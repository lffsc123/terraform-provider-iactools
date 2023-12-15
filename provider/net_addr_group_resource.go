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

type UpdateAddrGroupRequest struct {
	UpdateAddrGroupRequestModel UpdateAddrGroupRequestModel `json:"netaddrgrplist"`
}

// 调用接口参数
type AddAddrGroupRequestModel struct {
	VsysName    string `json:"vsysName"`
	Name        string `json:"name"`
	ObjNameList string `json:"objNameList"`
	Desc        string `json:"desc"`
}

type UpdateAddrGroupRequestModel struct {
	VsysName    string `json:"vsysName"`
	Name        string `json:"name"`
	OldName     string `json:"oldName"`
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

// 查询结果结构体
type QueryNetAddrGroupResponseListModel struct {
	Netaddrgrplist []QueryNetAddrGroupResponseModel `json:"netaddrgrplist"`
}
type QueryNetAddrGroupResponseModel struct {
	Id           string `json:"id"`
	VsysName     string `json:"vsysName"`
	Name         string `json:"name"`
	Oldname      string `json:"oldname"`
	ObjNameList  string `json:"objNameList"`
	Desc         string `json:"desc"`
	ReferNum     string `json:"referNum"`
	DelAllEnable string `json:"delAllEnable"`
	SearchValue  string `json:"searchValue"`
	Offset       string `json:"offset"`
	Count        string `json:"count"`
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

	if reqmethod == "POST" {
		// 先查询是否存在，再执行新增操作
		tflog.Info(ctx, "ip地址组--开始执行--查询操作")
		responseBody := sendRequest(ctx, "GET", c, nil, "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist?vsysName=PublicSystem&offset=0&count=1000", "ip地址组")
		var queryResList QueryNetAddrGroupResponseListModel
		err := json.Unmarshal([]byte(responseBody), &queryResList)
		if err != nil {
			panic("转换查询结果json出现异常")
		}
		for _, queryRes := range queryResList.Netaddrgrplist {
			if queryRes.Name == Rsinfo.Name.ValueString() {
				tflog.Info(ctx, "ip地址组--存在重复数据，执行--修改操作")
				var sendUpdateData UpdateAddrGroupRequestModel
				sendUpdateData = UpdateAddrGroupRequestModel{
					VsysName:    Rsinfo.VsysName.ValueString(),
					Name:        Rsinfo.Name.ValueString(),
					OldName:     Rsinfo.Name.ValueString(),
					ObjNameList: Rsinfo.ObjNameList.ValueString(),
					Desc:        Rsinfo.Desc.ValueString(),
				}

				requstUpdateData := UpdateAddrGroupRequest{
					UpdateAddrGroupRequestModel: sendUpdateData,
				}
				body, _ := json.Marshal(requstUpdateData)

				sendRequest(ctx, "PUT", c, body, "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist", "ip地址组")
				return
			}
		}

		// 新增操作
		var sendData AddAddrGroupRequestModel
		sendData = AddAddrGroupRequestModel{
			VsysName:    Rsinfo.VsysName.ValueString(),
			Name:        Rsinfo.Name.ValueString(),
			ObjNameList: Rsinfo.ObjNameList.ValueString(),
			Desc:        Rsinfo.Desc.ValueString(),
		}

		requstData := AddAddrGroupRequest{
			AddAddrGroupRequestModel: sendData,
		}
		body, _ := json.Marshal(requstData)
		sendRequest(ctx, reqmethod, c, body, "/func/web_main/api/netaddr/netaddr_group/netaddrgrplist", "ip地址组")
		return
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}

}
