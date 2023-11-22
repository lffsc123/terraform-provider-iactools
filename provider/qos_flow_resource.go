package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Qos流模版
var _ resource.Resource = &QosFlowResource{}
var _ resource.ResourceWithImportState = &QosFlowResource{}

func NewQosFlowResource() resource.Resource {
	return &QosFlowResource{}
}

type QosFlowResource struct {
	client *Client
}

type QosFlowResourceModel struct {
	AddQosFlowParameter    AddQosFlowParameter    `tfsdk:"addQosFlowParameter"`
	UpdateQosFlowParameter UpdateQosFlowParameter `tfsdk:"updateQosFlowParameter"`
	DelQosFlowParameter    DelQosFlowParameter    `tfsdk:"delQosFlowParameter"`
	ReadQosFlowParameter   ReadQosFlowParameter   `tfsdk:"readQosFlowParameter"`
}

type AddQosFlowParameter struct {
	IpVersion  types.String `tfsdk:"ipVersion"`
	VsysName   types.String `tfsdk:"vsysName"`
	FlowName   types.String `tfsdk:"flowName"`
	SrcMac     types.String `tfsdk:"srcMac"`
	SrcMacMask types.String `tfsdk:"srcMacMask"`
	DstMac     types.String `tfsdk:"dstMac"`
	DstMacMask types.String `tfsdk:"dstMacMask"`
	SrcIp      types.String `tfsdk:"srcIp"`
	SrcIpMask  types.String `tfsdk:"srcIpMask"`
	DstIp      types.String `tfsdk:"dstIp"`
	DstIpMask  types.String `tfsdk:"dstIpMask"`
	Procotol   types.String `tfsdk:"procotol"`
	SrcMinPort types.String `tfsdk:"srcMinPort"`
	SrcMaxPort types.String `tfsdk:"srcMaxPort"`
	DstMinPort types.String `tfsdk:"dstMinPort"`
	DstMaxPort types.String `tfsdk:"dstMaxPort"`
	PriType    types.String `tfsdk:"priType"`
	Priority   types.String `tfsdk:"priority"`
	VlanId     types.String `tfsdk:"vlanId"`
	CVlanId    types.String `tfsdk:"cVlanId"`
}

type UpdateQosFlowParameter struct {
	IpVersion  types.String `tfsdk:"ipVersion"`
	VsysName   types.String `tfsdk:"vsysName"`
	FlowName   types.String `tfsdk:"flowName"`
	SrcMac     types.String `tfsdk:"srcMac"`
	SrcMacMask types.String `tfsdk:"srcMacMask"`
	DstMac     types.String `tfsdk:"dstMac"`
	DstMacMask types.String `tfsdk:"dstMacMask"`
	SrcIp      types.String `tfsdk:"srcIp"`
	SrcIpMask  types.String `tfsdk:"srcIpMask"`
	DstIp      types.String `tfsdk:"dstIp"`
	DstIpMask  types.String `tfsdk:"dstIpMask"`
	Procotol   types.String `tfsdk:"procotol"`
	SrcMinPort types.String `tfsdk:"srcMinPort"`
	SrcMaxPort types.String `tfsdk:"srcMaxPort"`
	DstMinPort types.String `tfsdk:"dstMinPort"`
	DstMaxPort types.String `tfsdk:"dstMaxPort"`
	PriType    types.String `tfsdk:"priType"`
	Priority   types.String `tfsdk:"priority"`
	VlanId     types.String `tfsdk:"vlanId"`
	CVlanId    types.String `tfsdk:"cVlanId"`
}

type DelQosFlowParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VsysName  types.String `tfsdk:"vsysName"`
	FlowName  types.String `tfsdk:"flowName"`
}

type ReadQosFlowParameter struct {
	IpVersion types.String `tfsdk:"ipVersion"`
	VsysName  types.String `tfsdk:"vsysName"`
}

func (r *QosFlowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_QosFlow"
}

func (r *QosFlowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"param": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"ip_start": schema.StringAttribute{
						Required: true,
					},
					"ip_end": schema.StringAttribute{
						Required: true,
					},
					"ip_version": schema.StringAttribute{
						Optional: true,
					},
					"vrrp_if_name": schema.StringAttribute{
						Optional: true,
					},
					"vrrp_id": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *QosFlowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *QosFlowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *QosFlowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddQosFlowRequest(ctx, "POST", r.client, data.AddQosFlowParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QosFlowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *QosFlowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadQosFlowRequest(ctx, "GET", r.client, data.ReadQosFlowParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QosFlowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *QosFlowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateQosFlowRequest(ctx, "PUT", r.client, data.UpdateQosFlowParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QosFlowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *QosFlowResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelQosFlowRequest(ctx, "DELETE", r.client, data.DelQosFlowParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *QosFlowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddQosFlowRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddQosFlowParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_flow/flowlimit"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}

func sendToweb_UpdateQosFlowRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateQosFlowParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_flow/flowlimit"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}

func sendToweb_DelQosFlowRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelQosFlowParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_flow/flowlimit"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}

func sendToweb_ReadQosFlowRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadQosFlowParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_flow/flowlimit"

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	respn, err := http.DefaultClient.Do(req)
	if err != nil {
		tflog.Info(ctx, " read Error"+err.Error())
	}
	defer respn.Body.Close()

	body, err2 := ioutil.ReadAll(respn.Body)
	if err2 == nil {
		fmt.Println(string(body))
	}
}
