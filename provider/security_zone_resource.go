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

// 通用api
var _ resource.Resource = &GeneralApiResource{}
var _ resource.ResourceWithImportState = &GeneralApiResource{}

func NewGeneralApiResource() resource.Resource {
	return &GeneralApiResource{}
}

type GeneralApiResource struct {
	client *Client
}

type GeneralApiResourceModel struct {
	AddGeneralApiParameter AddGeneralApiParameter `tfsdk:"generalapiparam"`
}

type AddGeneralApiRequest struct {
	AddGeneralApiRequestModel AddGeneralApiRequestModel `json:"generalapiparam"`
}

// 调用接口参数
type AddGeneralApiRequestModel struct {
	Url         string `json:"url"`
	Method      string `json:"method"`
	RequestBody string `json:"requestBody"`
}

// 接收外部参数
type AddGeneralApiParameter struct {
	Url         types.String `tfsdk:"url"`
	Method      types.String `tfsdk:"method"`
	RequestBody types.String `tfsdk:"requestBody"`
}

func (r *GeneralApiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_GeneralApi"
}

func (r *GeneralApiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"generalapiparam": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Optional: true,
					},
					"method": schema.StringAttribute{
						Required: true,
					},
					"requestBody": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *GeneralApiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GeneralApiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *GeneralApiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource ************** ")
	sendToweb_GeneralApiRequest(ctx, r.client, data.AddGeneralApiParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GeneralApiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GeneralApiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	//sendToweb_GeneralApiRequest(ctx, "GET", r.client, data.AddGeneralApiParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GeneralApiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *GeneralApiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_GeneralApiRequest(ctx, "PUT", r.client, data.AddGeneralApiParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GeneralApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GeneralApiResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_GeneralApiRequest(ctx, "DELETE", r.client, data.AddGeneralApiParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GeneralApiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_GeneralApiRequest(ctx context.Context, c *Client, Rsinfo AddGeneralApiParameter) {
	// 新增操作
	var sendData AddGeneralApiRequestModel
	sendData = AddGeneralApiRequestModel{
		RequestBody: Rsinfo.RequestBody.ValueString(),
	}

	requstData := AddGeneralApiRequest{
		AddGeneralApiRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)
	tflog.Info(ctx, "===请求url==="+Rsinfo.Url.ValueString()+"===")
	tflog.Info(ctx, "===请求方法==="+Rsinfo.Method.ValueString()+"===")
	tflog.Info(ctx, "===请求参数转换前==="+Rsinfo.RequestBody.ValueString()+"===")
	tflog.Info(ctx, "===请求参数转换后==="+string(body)+"===")
	//sendRequest(ctx, Rsinfo.Method.ValueString(), c, body, Rsinfo.Url.ValueString(), "通用api")
	return
}
