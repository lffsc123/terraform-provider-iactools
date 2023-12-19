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

type UpdateUsrGroupRequest struct {
	UpdateUsrGroupRequestModel UpdateUsrGroupRequestModel `json:"grp"`
}

// 调用接口参数
type AddUsrGroupRequestModel struct {
	Name           string `json:"name"`
	VfwName        string `json:"vfwName"`
	Desc           string `json:"desc"`
	AllSerNameList string `json:"allSerNameList"`
}

type UpdateUsrGroupRequestModel struct {
	Name           string `json:"name"`
	VfwName        string `json:"vfwName"`
	OldName        string `json:"oldName"`
	Desc           string `json:"desc"`
	AllSerNameList string `json:"allSerNameList"`
}

// 查询结果结构体
type QueryUsrGroupResponseListModel struct {
	UsrGrouplist []QueryUsrGroupResponseModel `json:"grp"`
}
type QueryUsrGroupResponseModel struct {
	Id             string `json:"id"`
	VfwName        string `json:"vfwName"`
	Name           string `json:"name"`
	OldName        string `json:"oldName"`
	Desc           string `json:"desc"`
	AllSerNameList string `json:"allSerNameList"`
	PreNameList    string `json:"preNameList"`
	UsrNameList    string `json:"usrNameList"`
	ReferNum       string `json:"referNum"`
	DelAllEnable   string `json:"delAllEnable"`
	SearchValue    string `json:"searchValue"`
	Offset         string `json:"offset"`
	Count          string `json:"count"`
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

	if reqmethod == "POST" {

		// 先查询是否存在，再执行新增操作
		tflog.Info(ctx, "服务组--开始执行--查询操作")
		responseBody := sendRequest(ctx, "GET", c, nil, "/func/web_main/api/netservice/netservice/grp?vfwName=PublicSystem&searchValue="+Rsinfo.Name.ValueString()+"&offset=1&count=100", "服务组")
		var queryResList QueryUsrGroupResponseListModel
		err := json.Unmarshal([]byte(responseBody), &queryResList)
		if err != nil {
			panic("转换查询结果json出现异常")
		}
		for _, queryRes := range queryResList.UsrGrouplist {
			if queryRes.Name == Rsinfo.Name.ValueString() {
				tflog.Info(ctx, "服务组--存在重复数据，执行--修改操作")
				var sendUpdateData UpdateUsrGroupRequestModel
				sendUpdateData = UpdateUsrGroupRequestModel{
					Name:           Rsinfo.Name.ValueString(),
					VfwName:        Rsinfo.VfwName.ValueString(),
					OldName:        Rsinfo.Name.ValueString(),
					Desc:           Rsinfo.Desc.ValueString(),
					AllSerNameList: Rsinfo.AllSerNameList.ValueString(),
				}

				requstUpdateData := UpdateUsrGroupRequest{
					UpdateUsrGroupRequestModel: sendUpdateData,
				}
				body, _ := json.Marshal(requstUpdateData)

				sendRequest(ctx, "PUT", c, body, "/func/web_main/api/netservice/netservice/grp", "服务组")
				return
			}
		}

		// 新增操作
		var sendData AddUsrGroupRequestModel
		sendData = AddUsrGroupRequestModel{
			Name:           Rsinfo.Name.ValueString(),
			VfwName:        Rsinfo.VfwName.ValueString(),
			AllSerNameList: Rsinfo.AllSerNameList.ValueString(),
			Desc:           Rsinfo.Desc.ValueString(),
		}
		requstData := AddUsrGroupRequest{
			AddUsrGroupRequestModel: sendData,
		}
		body, _ := json.Marshal(requstData)
		sendRequest(ctx, reqmethod, c, body, "/func/web_main/api/netservice/netservice/grp", "服务组")
		return
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}
}
