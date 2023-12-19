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

// 目的NAT
var _ resource.Resource = &TargetNatResource{}
var _ resource.ResourceWithImportState = &TargetNatResource{}

func NewTargetNatResource() resource.Resource {
	return &TargetNatResource{}
}

type TargetNatResource struct {
	client *Client
}

type TargetNatResourceModel struct {
	AddTargetNatParameter AddTargetNatParameter `tfsdk:"dnatlist"`
}

type AddTargetNatRequest struct {
	AddTargetNatRequestModel AddTargetNatRequestModel `json:"dnatlist"`
}

type UpdateTargetNatRequest struct {
	UpdateTargetNatRequestModel UpdateTargetNatRequestModel `json:"dnatlist"`
}

// 调用接口参数
type AddTargetNatRequestModel struct {
	VsysName             string `json:"vsysName"`
	Name                 string `json:"name"`
	TargetName           string `json:"targetName"`
	Position             string `json:"position"`
	InInterface          string `json:"inInterface"`
	SrcIpObj             string `json:"srcIpObj"`
	SrcIpGroup           string `json:"srcIpGroup"`
	PublicIp             string `json:"publicIp"`
	PreService           string `json:"preService"`
	UsrService           string `json:"usrService"`
	InNetIp              string `json:"inNetIp"`
	InnetPort            string `json:"innetPort"`
	UnLimited            string `json:"unLimited"`
	SrcIpTranslate       string `json:"srcIpTranslate"`
	InterfaceAddressFlag string `json:"interfaceAddressFlag"`
	AddrpoolName         string `json:"addrpoolName"`
	VrrpIfName           string `json:"vrrpIfName"`
	VrrpId               string `json:"vrrpId"`
	State                string `json:"state"`
}

type UpdateTargetNatRequestModel struct {
	VsysName             string `json:"vsysName"`
	Name                 string `json:"name"`
	OldName              string `json:"oldName"`
	TargetName           string `json:"targetName"`
	Position             string `json:"position"`
	InInterface          string `json:"inInterface"`
	NetaddrObj           string `json:"netaddrObj"`
	NetaddrGroup         string `json:"netaddrGroup"`
	PublicIp             string `json:"publicIp"`
	PreService           string `json:"preService"`
	UsrService           string `json:"usrService"`
	InNetIp              string `json:"inNetIp"`
	InnetPort            string `json:"innetPort"`
	UnLimited            string `json:"unLimited"`
	SrcIpTranslate       string `json:"srcIpTranslate"`
	InterfaceAddressFlag string `json:"interfaceAddressFlag"`
	AddrpoolName         string `json:"addrpoolName"`
	VrrpIfName           string `json:"vrrpIfName"`
	VrrpId               string `json:"vrrpId"`
	State                string `json:"state"`
}

// 接收外部参数
type AddTargetNatParameter struct {
	VsysName             types.String `tfsdk:"vsysname"`
	Name                 types.String `tfsdk:"name"`
	TargetName           types.String `tfsdk:"targetname"`
	Position             types.String `tfsdk:"position"`
	InInterface          types.String `tfsdk:"ininterface"`
	SrcIpObj             types.String `tfsdk:"srcipobj"`
	SrcIpGroup           types.String `tfsdk:"srcipgroup"`
	PublicIp             types.String `tfsdk:"publicip"`
	PreService           types.String `tfsdk:"preservice"`
	UsrService           types.String `tfsdk:"usrservice"`
	InNetIp              types.String `tfsdk:"innetip"`
	InnetPort            types.String `tfsdk:"innetport"`
	UnLimited            types.String `tfsdk:"unlimited"`
	SrcIpTranslate       types.String `tfsdk:"srciptranslate"`
	InterfaceAddressFlag types.String `tfsdk:"interfaceaddressflag"`
	AddrpoolName         types.String `tfsdk:"addrpoolname"`
	VrrpIfName           types.String `tfsdk:"vrrpifname"`
	VrrpId               types.String `tfsdk:"vrrpid"`
	State                types.String `tfsdk:"state"`
}

// 查询结果结构体
type QueryTargetNatResponseListModel struct {
	Dnatlist QueryTargetNatResponseModel `json:"dnatlist"`
}
type QueryTargetNatResponseModel struct {
	Vsysname             string `json:"vsysName"`
	Count                string `json:"count"`
	Offset               string `json:"offset"`
	SearchValue          string `json:"searchValue"`
	Name                 string `json:"name"`
	InInterface          string `json:"inInterface"`
	SourceIp             string `json:"sourceIp"`
	PublicIp             string `json:"publicIp"`
	Protocol             string `json:"protocol"`
	Port                 string `json:"port"`
	InNetIp              string `json:"inNetIp"`
	State                string `json:"state"`
	SrcIpObj             string `json:"srcIpObj"`
	SrcIpGroup           string `json:"srcIpGroup"`
	PreService           string `json:"preService"`
	UsrService           string `json:"usrService"`
	InnetPort            string `json:"innetPort"`
	UnLimited            string `json:"unLimited"`
	SrcIpTranslate       string `json:"srcIpTranslate"`
	InterfaceAddressFlag string `json:"interfaceAddressFlag"`
	AddrpoolName         string `json:"addrpoolName"`
	VrrpIfName           string `json:"vrrpIfName"`
	VrrpId               string `json:"vrrpId"`
	RuleId               string `json:"ruleId"`
	DelallEnable         string `json:"delallEnable"`
	TargetName           string `json:"targetName"`
	OldName              string `json:"oldName"`
	Position             string `json:"position"`
}

func (r *TargetNatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_TargetNat"
}

func (r *TargetNatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dnatlist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"vsysname": schema.StringAttribute{
						Optional: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"targetname": schema.StringAttribute{
						Optional: true,
					},
					"position": schema.StringAttribute{
						Optional: true,
					},
					"ininterface": schema.StringAttribute{
						Optional: true,
					},
					"srcipobj": schema.StringAttribute{
						Optional: true,
					},
					"srcipgroup": schema.StringAttribute{
						Optional: true,
					},
					"publicip": schema.StringAttribute{
						Required: true,
					},
					"preservice": schema.StringAttribute{
						Optional: true,
					},
					"usrservice": schema.StringAttribute{
						Optional: true,
					},
					"innetip": schema.StringAttribute{
						Required: true,
					},
					"innetport": schema.StringAttribute{
						Required: true,
					},
					"unlimited": schema.StringAttribute{
						Optional: true,
					},
					"srciptranslate": schema.StringAttribute{
						Optional: true,
					},
					"interfaceaddressflag": schema.StringAttribute{
						Optional: true,
					},
					"addrpoolname": schema.StringAttribute{
						Optional: true,
					},
					"vrrpifname": schema.StringAttribute{
						Optional: true,
					},
					"vrrpid": schema.StringAttribute{
						Optional: true,
					},
					"state": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *TargetNatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TargetNatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TargetNatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_TargetNatRequest(ctx, "POST", r.client, data.AddTargetNatParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TargetNatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TargetNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	//sendToweb_TargetNatRequest(ctx, "GET", r.client, data.AddTargetNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TargetNatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TargetNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_TargetNatRequest(ctx, "PUT", r.client, data.AddTargetNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TargetNatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TargetNatResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_TargetNatRequest(ctx, "DELETE", r.client, data.AddTargetNatParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TargetNatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_TargetNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddTargetNatParameter) {

	if reqmethod == "POST" {

		// 先查询是否存在，再执行新增操作
		tflog.Info(ctx, "目的NAT--开始执行--查询操作")
		responseBody := sendRequest(ctx, "GET", c, nil, "/func/web_main/api/nat/nat/dnatlist", "目的NAT")
		var queryResList QueryTargetNatResponseListModel
		err := json.Unmarshal([]byte(responseBody), &queryResList)
		if err != nil {
			panic("转换查询结果json出现异常")
		}
		queryRes := queryResList.Dnatlist
		if queryRes.Name == Rsinfo.Name.ValueString() {
			tflog.Info(ctx, "目的NAT--存在重复数据，执行--修改操作")
			var sendUpdateData UpdateTargetNatRequestModel
			sendUpdateData = UpdateTargetNatRequestModel{
				VsysName:             Rsinfo.VsysName.ValueString(),
				Name:                 Rsinfo.Name.ValueString(),
				OldName:              Rsinfo.Name.ValueString(),
				TargetName:           Rsinfo.TargetName.ValueString(),
				Position:             Rsinfo.Position.ValueString(),
				InInterface:          Rsinfo.InInterface.ValueString(),
				NetaddrObj:           Rsinfo.SrcIpObj.ValueString(),
				NetaddrGroup:         Rsinfo.SrcIpGroup.ValueString(),
				PublicIp:             Rsinfo.PublicIp.ValueString(),
				PreService:           Rsinfo.PreService.ValueString(),
				UsrService:           Rsinfo.UsrService.ValueString(),
				InNetIp:              Rsinfo.InNetIp.ValueString(),
				InnetPort:            Rsinfo.InnetPort.ValueString(),
				UnLimited:            Rsinfo.UnLimited.ValueString(),
				SrcIpTranslate:       Rsinfo.SrcIpTranslate.ValueString(),
				InterfaceAddressFlag: Rsinfo.InterfaceAddressFlag.ValueString(),
				AddrpoolName:         Rsinfo.AddrpoolName.ValueString(),
				VrrpIfName:           Rsinfo.VrrpIfName.ValueString(),
				VrrpId:               Rsinfo.VrrpId.ValueString(),
				State:                Rsinfo.State.ValueString(),
			}

			requstUpdateData := UpdateTargetNatRequest{
				UpdateTargetNatRequestModel: sendUpdateData,
			}
			body, _ := json.Marshal(requstUpdateData)

			sendRequest(ctx, "PUT", c, body, "/func/web_main/api/nat/nat/dnatlist", "目的NAT")
			return
		}
		// 新增操作
		var sendData AddTargetNatRequestModel
		sendData = AddTargetNatRequestModel{
			VsysName:             Rsinfo.VsysName.ValueString(),
			Name:                 Rsinfo.Name.ValueString(),
			TargetName:           Rsinfo.TargetName.ValueString(),
			Position:             Rsinfo.Position.ValueString(),
			InInterface:          Rsinfo.InInterface.ValueString(),
			SrcIpObj:             Rsinfo.SrcIpObj.ValueString(),
			SrcIpGroup:           Rsinfo.SrcIpGroup.ValueString(),
			PublicIp:             Rsinfo.PublicIp.ValueString(),
			PreService:           Rsinfo.PreService.ValueString(),
			UsrService:           Rsinfo.UsrService.ValueString(),
			InNetIp:              Rsinfo.InNetIp.ValueString(),
			InnetPort:            Rsinfo.InnetPort.ValueString(),
			UnLimited:            Rsinfo.UnLimited.ValueString(),
			SrcIpTranslate:       Rsinfo.SrcIpTranslate.ValueString(),
			InterfaceAddressFlag: Rsinfo.InterfaceAddressFlag.ValueString(),
			AddrpoolName:         Rsinfo.AddrpoolName.ValueString(),
			VrrpIfName:           Rsinfo.VrrpIfName.ValueString(),
			VrrpId:               Rsinfo.VrrpId.ValueString(),
			State:                Rsinfo.State.ValueString(),
		}
		requstData := AddTargetNatRequest{
			AddTargetNatRequestModel: sendData,
		}
		body, _ := json.Marshal(requstData)
		sendRequest(ctx, reqmethod, c, body, "/func/web_main/api/nat/nat/dnatlist", "目的NAT")
		return
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}
}
