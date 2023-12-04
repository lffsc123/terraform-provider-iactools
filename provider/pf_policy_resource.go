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

// 包过滤
var _ resource.Resource = &PfPolicyResource{}
var _ resource.ResourceWithImportState = &PfPolicyResource{}

func NewPfPolicyResource() resource.Resource {
	return &PfPolicyResource{}
}

type PfPolicyResource struct {
	client *Client
}

type PfPolicyResourceModel struct {
	AddPfPolicyParameter AddPfPolicyParameter `tfsdk:"securitypolicylist"`
}

type AddPfPolicyRequest struct {
	AddPfPolicyRequestModel AddPfPolicyRequestModel `json:"securitypolicylist"`
}

// 调用接口参数
type AddPfPolicyRequestModel struct {
	Name                    string `json:"name"`
	Enabled                 string `json:"enabled"`
	Action                  string `json:"action"`
	DelAllEnable            string `json:"delAllEnable"`
	IpVersion               string `json:"ipVersion"`
	VsysName                string `json:"vsysName"`
	GroupName               string `json:"groupName"`
	TargetName              string `json:"targetName"`
	Position                string `json:"position"`
	SourceSecurityZone      string `json:"sourceSecurityZone"`
	DestinationSecurityZone string `json:"destinationSecurityZone"`
	SourceIpObjects         string `json:"sourceIpObjects"`
	SourceIpGroups          string `json:"sourceIpGroups"`
	SourceDomains           string `json:"sourceDomains"`
	SourceMacObjects        string `json:"sourceMacObjects"`
	SourceMacGroups         string `json:"sourceMacGroups"`
	DestinationIpObjects    string `json:"destinationIpObjects"`
	DestinationIpGroups     string `json:"destinationIpGroups"`
	DestinationDomains      string `json:"destinationDomains"`
	DestinationMacObjects   string `json:"destinationMacObjects"`
	DestinationMacGroups    string `json:"destinationMacGroups"`
	ServicePreObjects       string `json:"servicePreObjects"`
	ServiceUsrObjects       string `json:"serviceUsrObjects"`
	ServiceGroups           string `json:"serviceGroups"`
	UserObjects             string `json:"userObjects"`
	UserGroups              string `json:"userGroups"`
	Description             string `json:"describe"`
	EffectName              string `json:"effectName"`
	Matchlog                string `json:"matchlog"`
	Sessionlog              string `json:"sessionlog"`
	Longsession             string `json:"longsession"`
	Agingtime               string `json:"agingtime"`
	Fragdrop                string `json:"fragdrop"`
	Dscp                    string `json:"dscp"`
	Cos                     string `json:"cos"`
	RltGroup                string `json:"rltGroup"`
	RltUser                 string `json:"rltUser"`
	Acctl                   string `json:"acctl"`
	UrlClass                string `json:"urlClass"`
	UrlSenior               string `json:"urlSenior"`
	Cam                     string `json:"cam"`
	Ips                     string `json:"ips"`
	Av                      string `json:"av"`
}

// 接收外部参数
type AddPfPolicyParameter struct {
	Name                    types.String `tfsdk:"name"`
	Enabled                 types.String `tfsdk:"enabled"`
	Action                  types.String `tfsdk:"action"`
	DelAllEnable            types.String `tfsdk:"delAllEnable"`
	IpVersion               types.String `tfsdk:"ipVersion"`
	VsysName                types.String `tfsdk:"vsysName"`
	GroupName               types.String `tfsdk:"groupName"`
	TargetName              types.String `tfsdk:"targetName"`
	Position                types.String `tfsdk:"position"`
	SourceSecurityZone      types.String `tfsdk:"sourceSecurityZone"`
	DestinationSecurityZone types.String `tfsdk:"destinationSecurityZone"`
	SourceIpObjects         types.String `tfsdk:"sourceIpObjects"`
	SourceIpGroups          types.String `tfsdk:"sourceIpGroups"`
	SourceDomains           types.String `tfsdk:"sourceDomains"`
	SourceMacObjects        types.String `tfsdk:"sourceMacObjects"`
	SourceMacGroups         types.String `tfsdk:"sourceMacGroups"`
	DestinationIpObjects    types.String `tfsdk:"destinationIpObjects"`
	DestinationIpGroups     types.String `tfsdk:"destinationIpGroups"`
	DestinationDomains      types.String `tfsdk:"destinationDomains"`
	DestinationMacObjects   types.String `tfsdk:"destinationMacObjects"`
	DestinationMacGroups    types.String `tfsdk:"destinationMacGroups"`
	ServicePreObjects       types.String `tfsdk:"servicePreObjects"`
	ServiceUsrObjects       types.String `tfsdk:"serviceUsrObjects"`
	ServiceGroups           types.String `tfsdk:"serviceGroups"`
	UserObjects             types.String `tfsdk:"userObjects"`
	UserGroups              types.String `tfsdk:"userGroups"`
	Description             types.String `tfsdk:"describe"`
	EffectName              types.String `tfsdk:"effectName"`
	Matchlog                types.String `tfsdk:"matchlog"`
	Sessionlog              types.String `tfsdk:"sessionlog"`
	Longsession             types.String `tfsdk:"longsession"`
	Agingtime               types.String `tfsdk:"agingtime"`
	Fragdrop                types.String `tfsdk:"fragdrop"`
	Dscp                    types.String `tfsdk:"dscp"`
	Cos                     types.String `tfsdk:"cos"`
	RltGroup                types.String `tfsdk:"rltGroup"`
	RltUser                 types.String `tfsdk:"rltUser"`
	Acctl                   types.String `tfsdk:"acctl"`
	UrlClass                types.String `tfsdk:"urlClass"`
	UrlSenior               types.String `tfsdk:"urlSenior"`
	Cam                     types.String `tfsdk:"cam"`
	Ips                     types.String `tfsdk:"ips"`
	Av                      types.String `tfsdk:"av"`
}

func (r *PfPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_PfPolicy"
}

func (r *PfPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"securitypolicylist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"enabled": schema.StringAttribute{
						Required: true,
					},
					"action": schema.StringAttribute{
						Required: true,
					},
					"ipVersion": schema.StringAttribute{
						Required: false,
					},
					"vsysName": schema.StringAttribute{
						Required: true,
					},
					"groupName": schema.StringAttribute{
						Required: true,
					},
					"targetName": schema.StringAttribute{
						Required: true,
					},
					"position": schema.StringAttribute{
						Required: true,
					},
					"effectName": schema.StringAttribute{
						Required: true,
					},
					"matchlog": schema.StringAttribute{
						Required: true,
					},
					"sessionlog": schema.StringAttribute{
						Required: true,
					},
					"longsession": schema.StringAttribute{
						Required: true,
					},
					"agingtime": schema.StringAttribute{
						Required: true,
					},
					"fragdrop": schema.StringAttribute{
						Required: true,
					},
					"sourceSecurityZone": schema.StringAttribute{
						Required: true,
					},
					"destinationSecurityZone": schema.StringAttribute{
						Required: true,
					},
					"sourceIpObjects": schema.StringAttribute{
						Required: true,
					},
					"sourceIpGroups": schema.StringAttribute{
						Required: true,
					},
					"sourceDomains": schema.StringAttribute{
						Required: true,
					},
					"sourceMacObjects": schema.StringAttribute{
						Required: true,
					},
					"sourceMacGroups": schema.StringAttribute{
						Required: true,
					},
					"destinationIpObjects": schema.StringAttribute{
						Required: true,
					},
					"destinationIpGroups": schema.StringAttribute{
						Required: true,
					},
					"destinationDomains": schema.StringAttribute{
						Required: true,
					},
					"destinationMacObjects": schema.StringAttribute{
						Required: true,
					},
					"destinationMacGroups": schema.StringAttribute{
						Required: true,
					},
					"servicePreObjects": schema.StringAttribute{
						Required: true,
					},
					"serviceUsrObjects": schema.StringAttribute{
						Required: true,
					},
					"serviceGroups": schema.StringAttribute{
						Required: true,
					},
					"userObjects": schema.StringAttribute{
						Required: true,
					},
					"userGroups": schema.StringAttribute{
						Required: true,
					},
					"describe": schema.StringAttribute{
						Required: true,
					},
					"dscp": schema.StringAttribute{
						Required: true,
					},
					"cos": schema.StringAttribute{
						Required: true,
					},
					"rltGroup": schema.StringAttribute{
						Required: true,
					},
					"rltUser": schema.StringAttribute{
						Required: true,
					},
					"acctl": schema.StringAttribute{
						Required: true,
					},
					"urlClass": schema.StringAttribute{
						Required: true,
					},
					"urlSenior": schema.StringAttribute{
						Required: true,
					},
					"cam": schema.StringAttribute{
						Required: true,
					},
					"ips": schema.StringAttribute{
						Required: true,
					},
					"av": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *PfPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PfPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PfPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	//sendToweb_Request(ctx, "POST", r.client, data.AddPfPolicyParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PfPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PfPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	//sendToweb_Request(ctx, "GET", r.client, data.AddPfPolicyParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PfPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *PfPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	//sendToweb_Request(ctx, "PUT", r.client, data.AddPfPolicyParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PfPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PfPolicyResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_Request(ctx, "DELETE", r.client, data.AddPfPolicyParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PfPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_Request(ctx context.Context, reqmethod string, c *Client, Rsinfo AddPfPolicyParameter) {

	var sendData AddPfPolicyRequestModel
	if reqmethod == "POST" {
		sendData = AddPfPolicyRequestModel{
			Name:    Rsinfo.Name.ValueString(),
			Enabled: Rsinfo.Enabled.ValueString(),
			Action:  Rsinfo.Action.ValueString(),
		}
	} else if reqmethod == "GET" {

	} else if reqmethod == "PUT" {

	} else if reqmethod == "DELETE" {
		sendData = AddPfPolicyRequestModel{
			DelAllEnable: "1",
		}
	}

	requstData := AddPfPolicyRequest{
		AddPfPolicyRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)

	tflog.Info(ctx, "请求体============:"+string(body))

	targetUrl := c.HostURL + "/func/web_main/api/pf_policy/pf_policy/pf_policy/securitypolicylist"

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
