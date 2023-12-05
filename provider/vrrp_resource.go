package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// vrrp
var _ resource.Resource = &VrrpResource{}
var _ resource.ResourceWithImportState = &VrrpResource{}

func NewVrrpResource() resource.Resource {
	return &VrrpResource{}
}

type VrrpResource struct {
	client *Client
}

type VrrpResourceModel struct {
	AddVrrpParameter AddVrrpParameter `tfsdk:"vrrpv3list"`
}

type AddVrrpRequest struct {
	AddVrrpRequestModel AddVrrpRequestModel `json:"vrrpv3list"`
}

// 调用接口参数
type AddVrrpRequestModel struct {
	Ipversion    string `json:"ipVersion"`
	Vsysname     string `json:"vsysName"`
	Ifname       string `json:"ifname"`
	Vrid         string `json:"vrid"`
	Virtualip    string `json:"virtualip"`
	Version      string `json:"version"`
	Priority     string `json:"priority"`
	Timerlearn   string `json:"timerLearn"`
	Adverint     string `json:"adverInt"`
	Preemptmode  string `json:"preemptMode"`
	Preemptdelay string `json:"preemptDelay"`
	Authmode     string `json:"authMode"`
	Authpass     string `json:"authPass"`
	Chksumflag   string `json:"chksumFlag"`
	Trackifs     string `json:"trackifs"`
	Trackips     string `json:"trackips"`
	Bfdsip       string `json:"bfdSip"`
	Bfddip       string `json:"bfdDip"`
	Bfdpro       string `json:"bfdPro"`
	Cfgstate     string `json:"cfgState"`
}

// 接收外部参数
type AddVrrpParameter struct {
	Ipversion    types.String `tfsdk:"ipversion"`
	Vsysname     types.String `tfsdk:"vsysname"`
	Ifname       types.String `tfsdk:"ifname"`
	Vrid         types.String `tfsdk:"vrid"`
	Virtualip    types.String `tfsdk:"virtualip"`
	Version      types.String `tfsdk:"version"`
	Priority     types.String `tfsdk:"priority"`
	Timerlearn   types.String `tfsdk:"timerlearn"`
	Adverint     types.String `tfsdk:"adverint"`
	Preemptmode  types.String `tfsdk:"preemptmode"`
	Preemptdelay types.String `tfsdk:"preemptdelay"`
	Authmode     types.String `tfsdk:"authmode"`
	Authpass     types.String `tfsdk:"authpass"`
	Chksumflag   types.String `tfsdk:"chksumflag"`
	Trackifs     types.String `tfsdk:"trackifs"`
	Trackips     types.String `tfsdk:"trackips"`
	Bfdsip       types.String `tfsdk:"bfdsip"`
	Bfddip       types.String `tfsdk:"bfddip"`
	Bfdpro       types.String `tfsdk:"bfdpro"`
	Cfgstate     types.String `tfsdk:"cfgstate"`
}

func (r *VrrpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_Vrrp"
}

func (r *VrrpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vrrpv3list": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"ipversion": schema.StringAttribute{
						Required: true,
					},
					"vsysname": schema.StringAttribute{
						Required: true,
					},
					"ifname": schema.StringAttribute{
						Required: true,
					},
					"vrid": schema.StringAttribute{
						Required: true,
					},
					"virtualip": schema.StringAttribute{
						Required: true,
					},
					"version": schema.StringAttribute{
						Required: true,
					},
					"priority": schema.StringAttribute{
						Required: true,
					},
					"timerlearn": schema.StringAttribute{
						Required: true,
					},
					"adverint": schema.StringAttribute{
						Required: true,
					},
					"preemptmode": schema.StringAttribute{
						Required: true,
					},
					"preemptdelay": schema.StringAttribute{
						Required: true,
					},
					"authmode": schema.StringAttribute{
						Required: true,
					},
					"authpass": schema.StringAttribute{
						Required: true,
					},
					"chksumflag": schema.StringAttribute{
						Required: true,
					},
					"trackifs": schema.StringAttribute{
						Required: true,
					},
					"trackips": schema.StringAttribute{
						Required: true,
					},
					"bfdsip": schema.StringAttribute{
						Required: true,
					},
					"bfddip": schema.StringAttribute{
						Required: true,
					},
					"bfdpro": schema.StringAttribute{
						Required: true,
					},
					"cfgstate": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *VrrpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VrrpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VrrpResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_VrrpRequest(ctx, "POST", r.client, data.AddVrrpParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VrrpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VrrpResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	//sendToweb_VrrpRequest(ctx, "GET", r.client, data.AddVrrpParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VrrpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VrrpResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_VrrpRequest(ctx, "PUT", r.client, data.AddVrrpParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VrrpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VrrpResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_VrrpRequest(ctx, "DELETE", r.client, data.AddVrrpParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VrrpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_VrrpRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddVrrpParameter) {

	var sendData AddVrrpRequestModel
	if reqmethod == "POST" {
		sendData = AddVrrpRequestModel{
			Ipversion:    Rsinfo.Ipversion.ValueString(),
			Vsysname:     Rsinfo.Vsysname.ValueString(),
			Ifname:       Rsinfo.Ifname.ValueString(),
			Vrid:         Rsinfo.Vrid.ValueString(),
			Virtualip:    Rsinfo.Virtualip.ValueString(),
			Version:      Rsinfo.Version.ValueString(),
			Priority:     Rsinfo.Priority.ValueString(),
			Timerlearn:   Rsinfo.Timerlearn.ValueString(),
			Adverint:     Rsinfo.Adverint.ValueString(),
			Preemptmode:  Rsinfo.Preemptmode.ValueString(),
			Preemptdelay: Rsinfo.Preemptdelay.ValueString(),
			Authmode:     Rsinfo.Authmode.ValueString(),
			Authpass:     Rsinfo.Authpass.ValueString(),
			Chksumflag:   Rsinfo.Chksumflag.ValueString(),
			Trackifs:     Rsinfo.Trackifs.ValueString(),
			Trackips:     Rsinfo.Trackips.ValueString(),
			Bfddip:       Rsinfo.Bfddip.ValueString(),
			Bfdsip:       Rsinfo.Bfdsip.ValueString(),
			Bfdpro:       Rsinfo.Bfdpro.ValueString(),
			Cfgstate:     Rsinfo.Cfgstate.ValueString(),
		}
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {

	}

	requstData := AddVrrpRequest{
		AddVrrpRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)

	tflog.Info(ctx, "请求体============:"+string(body))

	targetUrl := c.HostURL + "/func/web_main/api/vrrpv3/vrrpv3/vrrpv3list"

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
		tflog.Error(ctx, "发送请求失败======="+err.Error())
		panic("发送请求失败=======")
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		tflog.Error(ctx, "发送请求失败======="+err2.Error())
		panic("发送请求失败=======")
	}
	// 打印响应结果
	tflog.Info(ctx, "响应状态码======="+string(respn.Status))
	tflog.Info(ctx, "响应体======="+string(body))

	if respn.Status != "200" || respn.Status != "201" || respn.Status != "204" {
		panic("请求响应失败=======")
	}
}
