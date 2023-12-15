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
	AddUsrObjParameter AddUsrObjParameter `tfsdk:"usrobj"`
}

type AddUsrObjRequest struct {
	AddUsrObjRequestModel AddUsrObjRequestModel `json:"usrobj"`
}

type UpdateUsrObjRequest struct {
	UpdateUsrObjRequestModel UpdateUsrObjRequestModel `json:"usrobj"`
}

// 调用接口参数
type AddUsrObjRequestModel struct {
	Name       string `json:"name"`
	VfwName    string `json:"vfwName"`
	Protocol   string `json:"protocol"`
	SportStart string `json:"sportStart"`
	SportEnd   string `json:"sportEnd"`
	DportStart string `json:"dportStart"`
	DportEnd   string `json:"dportEnd"`
	Services   string `json:"services"`
	Desc       string `json:"desc"`
}

type UpdateUsrObjRequestModel struct {
	Name       string `json:"name"`
	VfwName    string `json:"vfwName"`
	OldName    string `json:"oldName"`
	Protocol   string `json:"protocol"`
	SportStart string `json:"sportStart"`
	SportEnd   string `json:"sportEnd"`
	DportStart string `json:"dportStart"`
	DportEnd   string `json:"dportEnd"`
	Services   string `json:"services"`
	Desc       string `json:"desc"`
}

// 接收外部参数
type AddUsrObjParameter struct {
	Name       types.String `tfsdk:"name"`
	VfwName    types.String `tfsdk:"vfwname"`
	Protocol   types.String `tfsdk:"protocol"`
	SportStart types.String `tfsdk:"sportstart"`
	SportEnd   types.String `tfsdk:"sportend"`
	DportStart types.String `tfsdk:"dportstart"`
	DportEnd   types.String `tfsdk:"dportend"`
	Services   types.String `tfsdk:"services"`
	Desc       types.String `tfsdk:"desc"`
}

// 查询结果结构体
type QueryUsrObjResponseListModel struct {
	UsrObjlist []QueryUsrObjResponseModel `json:"usrobj"`
}
type QueryUsrObjResponseModel struct {
	Id           string `json:"id"`
	VfwName      string `json:"vfwName"`
	Name         string `json:"name"`
	OldName      string `json:"oldName"`
	Protocol     string `json:"protocol"`
	SportStart   string `json:"sportStart"`
	SportEnd     string `json:"sportEnd"`
	DportStart   string `json:"dportStart"`
	DportEnd     string `json:"dportEnd"`
	Services     string `json:"services"`
	Desc         string `json:"desc"`
	ReferNum     string `json:"referNum"`
	DelAllEnable string `json:"delAllEnable"`
	SearchValue  string `json:"searchValue"`
	Offset       string `json:"offset"`
	Count        string `json:"count"`
}

func (r *UsrObjResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_UsrObj"
}

func (r *UsrObjResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"usrobj": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"vfwname": schema.StringAttribute{
						Optional: true,
					},
					"protocol": schema.StringAttribute{
						Required: true,
					},
					"sportstart": schema.StringAttribute{
						Optional: true,
					},
					"sportend": schema.StringAttribute{
						Optional: true,
					},
					"dportstart": schema.StringAttribute{
						Optional: true,
					},
					"dportend": schema.StringAttribute{
						Optional: true,
					},
					"services": schema.StringAttribute{
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
	sendToweb_UsrObjRequest(ctx, "POST", r.client, data.AddUsrObjParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrObjResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UsrObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start *************** ")
	//sendToweb_UsrObjRequest(ctx, "GET", r.client, data.AddUsrObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrObjResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UsrObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_UsrObjRequest(ctx, "PUT", r.client, data.AddUsrObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UsrObjResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UsrObjResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_UsrObjRequest(ctx, "DELETE", r.client, data.AddUsrObjParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UsrObjResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_UsrObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddUsrObjParameter) {

	if reqmethod == "POST" {
		// 先查询是否存在，再执行新增操作
		tflog.Info(ctx, "服务对象--开始执行--查询操作")
		responseBody := sendRequest(ctx, "GET", c, nil, "/func/web_main/api/netservice/netservice/usrobj?vfwName=vsys&searchValue="+Rsinfo.Name.ValueString()+"&offset=1&count=1000", "服务对象")
		var queryResList QueryUsrObjResponseListModel
		err := json.Unmarshal([]byte(responseBody), &queryResList)
		if err != nil {
			panic("转换查询结果json出现异常")
		}
		for _, queryRes := range queryResList.UsrObjlist {
			if queryRes.Name == Rsinfo.Name.ValueString() {
				tflog.Info(ctx, "服务对象--存在重复数据，执行--修改操作")
				var sendUpdateData UpdateUsrObjRequestModel
				sendUpdateData = UpdateUsrObjRequestModel{
					Name:       Rsinfo.Name.ValueString(),
					VfwName:    Rsinfo.VfwName.ValueString(),
					Protocol:   Rsinfo.Protocol.ValueString(),
					SportStart: Rsinfo.SportStart.ValueString(),
					SportEnd:   Rsinfo.SportEnd.ValueString(),
					DportStart: Rsinfo.DportStart.ValueString(),
					DportEnd:   Rsinfo.DportEnd.ValueString(),
					Services:   Rsinfo.Services.ValueString(),
					Desc:       Rsinfo.Desc.ValueString(),
				}

				requstUpdateData := UpdateUsrObjRequest{
					UpdateUsrObjRequestModel: sendUpdateData,
				}
				body, _ := json.Marshal(requstUpdateData)

				sendRequest(ctx, "PUT", c, body, "/func/web_main/api/netservice/netservice/usrobj", "服务对象")
				return
			}
		}

		// 新增操作
		var sendData AddUsrObjRequestModel
		sendData = AddUsrObjRequestModel{
			Name:       Rsinfo.Name.ValueString(),
			VfwName:    Rsinfo.VfwName.ValueString(),
			Protocol:   Rsinfo.Protocol.ValueString(),
			SportStart: Rsinfo.SportStart.ValueString(),
			SportEnd:   Rsinfo.SportEnd.ValueString(),
			DportStart: Rsinfo.DportStart.ValueString(),
			DportEnd:   Rsinfo.DportEnd.ValueString(),
			Services:   Rsinfo.Services.ValueString(),
			Desc:       Rsinfo.Desc.ValueString(),
		}
		requstData := AddUsrObjRequest{
			AddUsrObjRequestModel: sendData,
		}
		body, _ := json.Marshal(requstData)
		sendRequest(ctx, "PUT", c, body, "/func/web_main/api/netservice/netservice/usrobj", "服务对象")
		return
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}
}
