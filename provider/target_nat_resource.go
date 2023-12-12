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

	var sendData AddTargetNatRequestModel
	if reqmethod == "POST" {
		sendData = AddTargetNatRequestModel{
			VsysName:             Rsinfo.VsysName.ValueString(),
			Name:                 Rsinfo.Name.ValueString(),
			TargetName:           Rsinfo.TargetName.ValueString(),
			Position:             Rsinfo.Position.ValueString(),
			InInterface:          Rsinfo.InInterface.ValueString(),
			SrcIpObj:             Rsinfo.SrcIpObj.ValueString(),
			SrcIpGroup:           Rsinfo.SrcIpGroup.ValueString(),
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
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}

	requstData := AddTargetNatRequest{
		AddTargetNatRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)

	tflog.Info(ctx, "目的NAT--请求体============:"+string(body))

	targetUrl := c.HostURL + "/func/web_main/api/nat/nat/dnatlist"

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
		tflog.Error(ctx, "目的NAT--发送请求失败======="+err.Error())
		panic("目的NAT--发送请求失败=======")
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		tflog.Error(ctx, "目的NAT--发送请求失败======="+err2.Error())
		panic("目的NAT--发送请求失败=======")
	}

	if respn.Status != "200" && respn.Status != "201" && respn.Status != "204" {
		tflog.Info(ctx, "目的NAT--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "目的NAT--响应体======="+string(body))
		panic("目的NAT--请求响应失败=======")
	} else {
		// 打印响应结果
		tflog.Info(ctx, "目的NAT--响应状态码======="+string(respn.Status)+"======")
		tflog.Info(ctx, "目的NAT--响应体======="+string(body))
	}
}
