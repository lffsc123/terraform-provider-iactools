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

// 源NAT
var _ resource.Resource = &SourceNatResource{}
var _ resource.ResourceWithImportState = &SourceNatResource{}

func NewSourceNatResource() resource.Resource {
	return &SourceNatResource{}
}

type SourceNatResource struct {
	client *Client
}

type SourceNatResourceModel struct {
	AddSourceNatParameter AddSourceNatParameter `tfsdk:"snatlist"`
}

type AddSourceNatRequest struct {
	AddSourceNatRequestModel AddSourceNatRequestModel `json:"snatlist"`
}

type UpdateSourceNatRequest struct {
	UpdateSourceNatRequestModel UpdateSourceNatRequestModel `json:"snatlist"`
}

// 调用接口参数
type AddSourceNatRequestModel struct {
	IpVersion           string `json:"ipVersion"`
	VsysName            string `json:"vsysName"`
	Name                string `json:"name"`
	TargetName          string `json:"targetName"`
	Position            string `json:"position"`
	OutInterface        string `json:"outInterface"`
	SrcAddrObj          string `json:"srcAddrObj"`
	SrcAddrGroup        string `json:"srcAddrGroup"`
	DstAddrObj          string `json:"dstAddrObj"`
	DstAddrGroup        string `json:"dstaddrgroup"`
	PreService          string `json:"preService"`
	UsrService          string `json:"usrService"`
	ServiceGroup        string `json:"serviceGroup"`
	PublicIpAddressFlag string `json:"publicIpAddressFlag"`
	AddrpoolName        string `json:"addrpoolName"`
	MinPort             string `json:"minPort"`
	MaxPort             string `json:"maxPort"`
	PortHash            string `json:"portHash"`
	State               string `json:"state"`
}

type UpdateSourceNatRequestModel struct {
	IpVersion           string `json:"ipVersion"`
	VsysName            string `json:"vsysName"`
	Name                string `json:"name"`
	OldName             string `json:"oldName"`
	TargetName          string `json:"targetName"`
	Position            string `json:"position"`
	OutInterface        string `json:"outInterface"`
	SrcAddrObj          string `json:"srcAddrObj"`
	SrcAddrGroup        string `json:"srcAddrGroup"`
	DstAddrObj          string `json:"dstAddrObj"`
	DstAddrGroup        string `json:"dstaddrgroup"`
	PreService          string `json:"preService"`
	UsrService          string `json:"usrService"`
	ServiceGroup        string `json:"serviceGroup"`
	PublicIpAddressFlag string `json:"publicIpAddressFlag"`
	AddrpoolName        string `json:"addrpoolName"`
	MinPort             string `json:"minPort"`
	MaxPort             string `json:"maxPort"`
	PortHash            string `json:"portHash"`
	State               string `json:"state"`
}

// 接收外部参数
type AddSourceNatParameter struct {
	IpVersion           types.String `tfsdk:"ipversion"`
	VsysName            types.String `tfsdk:"vsysname"`
	Name                types.String `tfsdk:"name"`
	TargetName          types.String `tfsdk:"targetname"`
	Position            types.String `tfsdk:"position"`
	OutInterface        types.String `tfsdk:"outinterface"`
	SrcAddrObj          types.String `tfsdk:"srcaddrobj"`
	SrcAddrGroup        types.String `tfsdk:"srcaddrgroup"`
	DstAddrObj          types.String `tfsdk:"dstaddrobj"`
	DstAddrGroup        types.String `tfsdk:"dstaddrgroup"`
	PreService          types.String `tfsdk:"preservice"`
	UsrService          types.String `tfsdk:"usrservice"`
	ServiceGroup        types.String `tfsdk:"servicegroup"`
	PublicIpAddressFlag types.String `tfsdk:"publicipaddressflag"`
	AddrpoolName        types.String `tfsdk:"addrpoolname"`
	MinPort             types.String `tfsdk:"minport"`
	MaxPort             types.String `tfsdk:"maxport"`
	PortHash            types.String `tfsdk:"porthash"`
	State               types.String `tfsdk:"state"`
}

// 查询结果结构体
type QuerySourceNatResponseListModel struct {
	Snatlist QuerySourceNatResponseModel `json:"snatlist"`
}
type QuerySourceNatResponseModel struct {
	IpVersion           string `json:"ipVersion"`
	VsysName            string `json:"vsysName"`
	Offset              string `json:"offset"`
	Count               string `json:"count"`
	Name                string `json:"name"`
	OutInterface        string `json:"outInterface"`
	SourceIp            string `json:"sourceIp"`
	DestinationIp       string `json:"destinationIp"`
	Service             string `json:"service"`
	State               string `json:"state"`
	SrcAddrObj          string `json:"srcAddrObj"`
	SrcAddrGroup        string `json:"srcAddrGroup"`
	DstAddrObj          string `json:"dstAddrObj"`
	DstAddrGroup        string `json:"dstaddrgroup"`
	PreService          string `json:"preService"`
	UsrService          string `json:"usrService"`
	ServiceGroup        string `json:"serviceGroup"`
	PublicIpAddressFlag string `json:"publicIpAddressFlag"`
	AddrpoolName        string `json:"addrpoolName"`
	MinPort             string `json:"minPort"`
	MaxPort             string `json:"maxPort"`
	PortHash            string `json:"portHash"`
	RuleId              string `json:"ruleId"`
	DelallEnable        string `json:"delallEnable"`
	TargetName          string `json:"targetName"`
	OldName             string `json:"oldName"`
}

func (r *SourceNatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_SourceNat"
}

func (r *SourceNatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"snatlist": schema.SingleNestedAttribute{
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
					"targetname": schema.StringAttribute{
						Optional: true,
					},
					"position": schema.StringAttribute{
						Optional: true,
					},
					"outinterface": schema.StringAttribute{
						Required: true,
					},
					"srcaddrobj": schema.StringAttribute{
						Optional: true,
					},
					"srcaddrgroup": schema.StringAttribute{
						Optional: true,
					},
					"dstaddrobj": schema.StringAttribute{
						Optional: true,
					},
					"dstaddrgroup": schema.StringAttribute{
						Optional: true,
					},
					"preservice": schema.StringAttribute{
						Optional: true,
					},
					"usrservice": schema.StringAttribute{
						Optional: true,
					},
					"servicegroup": schema.StringAttribute{
						Optional: true,
					},
					"publicipaddressflag": schema.StringAttribute{
						Optional: true,
					},
					"addrpoolname": schema.StringAttribute{
						Optional: true,
					},
					"minport": schema.StringAttribute{
						Optional: true,
					},
					"maxport": schema.StringAttribute{
						Optional: true,
					},
					"porthash": schema.StringAttribute{
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

func (r *SourceNatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *SourceNatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_SourceNatRequest(ctx, "POST", r.client, data.AddSourceNatParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceNatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	//sendToweb_SourceNatRequest(ctx, "GET", r.client, data.AddSourceNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceNatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SourceNatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_SourceNatRequest(ctx, "PUT", r.client, data.AddSourceNatParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceNatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SourceNatResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_SourceNatRequest(ctx, "DELETE", r.client, data.AddSourceNatParameter)

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

func (r *SourceNatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_SourceNatRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddSourceNatParameter) {

	if reqmethod == "POST" {
		// 先查询是否存在，再执行新增操作
		tflog.Info(ctx, "源NAT--开始执行--查询操作")
		responseBody := sendRequest(ctx, "GET", c, nil, "/func/web_main/api/nat/nat/snatlist", "源NAT")
		var queryResList QuerySourceNatResponseListModel
		err := json.Unmarshal([]byte(responseBody), &queryResList)
		if err != nil {
			panic("转换查询结果json出现异常")
		}
		queryRes := queryResList.Snatlist
		if queryRes.Name == Rsinfo.Name.ValueString() {
			tflog.Info(ctx, "源NAT--存在重复数据，执行--修改操作")
			var sendUpdateData UpdateSourceNatRequestModel
			sendUpdateData = UpdateSourceNatRequestModel{
				IpVersion:           Rsinfo.IpVersion.ValueString(),
				VsysName:            Rsinfo.VsysName.ValueString(),
				Name:                Rsinfo.Name.ValueString(),
				OldName:             Rsinfo.Name.ValueString(),
				TargetName:          Rsinfo.TargetName.ValueString(),
				Position:            Rsinfo.Position.ValueString(),
				OutInterface:        Rsinfo.OutInterface.ValueString(),
				SrcAddrObj:          Rsinfo.SrcAddrObj.ValueString(),
				SrcAddrGroup:        Rsinfo.SrcAddrGroup.ValueString(),
				DstAddrObj:          Rsinfo.DstAddrObj.ValueString(),
				DstAddrGroup:        Rsinfo.DstAddrGroup.ValueString(),
				PreService:          Rsinfo.PreService.ValueString(),
				UsrService:          Rsinfo.UsrService.ValueString(),
				ServiceGroup:        Rsinfo.ServiceGroup.ValueString(),
				PublicIpAddressFlag: Rsinfo.PublicIpAddressFlag.ValueString(),
				AddrpoolName:        Rsinfo.AddrpoolName.ValueString(),
				MinPort:             Rsinfo.MinPort.ValueString(),
				MaxPort:             Rsinfo.MaxPort.ValueString(),
				PortHash:            Rsinfo.PortHash.ValueString(),
				State:               Rsinfo.State.ValueString(),
			}

			requstUpdateData := UpdateSourceNatRequest{
				UpdateSourceNatRequestModel: sendUpdateData,
			}
			body, _ := json.Marshal(requstUpdateData)

			sendRequest(ctx, "PUT", c, body, "/func/web_main/api/nat/nat/snatlist", "源NAT")
			return
		}
		// 新增操作
		var sendData AddSourceNatRequestModel
		sendData = AddSourceNatRequestModel{
			IpVersion:           Rsinfo.IpVersion.ValueString(),
			VsysName:            Rsinfo.VsysName.ValueString(),
			Name:                Rsinfo.Name.ValueString(),
			TargetName:          Rsinfo.TargetName.ValueString(),
			Position:            Rsinfo.Position.ValueString(),
			OutInterface:        Rsinfo.OutInterface.ValueString(),
			SrcAddrObj:          Rsinfo.SrcAddrObj.ValueString(),
			SrcAddrGroup:        Rsinfo.SrcAddrGroup.ValueString(),
			DstAddrObj:          Rsinfo.DstAddrObj.ValueString(),
			DstAddrGroup:        Rsinfo.DstAddrGroup.ValueString(),
			PreService:          Rsinfo.PreService.ValueString(),
			UsrService:          Rsinfo.UsrService.ValueString(),
			ServiceGroup:        Rsinfo.ServiceGroup.ValueString(),
			PublicIpAddressFlag: Rsinfo.PublicIpAddressFlag.ValueString(),
			AddrpoolName:        Rsinfo.AddrpoolName.ValueString(),
			MinPort:             Rsinfo.MinPort.ValueString(),
			MaxPort:             Rsinfo.MaxPort.ValueString(),
			PortHash:            Rsinfo.PortHash.ValueString(),
			State:               Rsinfo.State.ValueString(),
		}
		requstData := AddSourceNatRequest{
			AddSourceNatRequestModel: sendData,
		}
		body, _ := json.Marshal(requstData)
		sendRequest(ctx, reqmethod, c, body, "/func/web_main/api/nat/nat/snatlist", "源NAT")
		return
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}
}
