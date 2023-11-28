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

// Qos策略
var _ resource.Resource = &QosPolicyResource{}
var _ resource.ResourceWithImportState = &QosPolicyResource{}

func NewQosPolicyResource() resource.Resource {
	return &QosPolicyResource{}
}

type QosPolicyResource struct {
	client *Client
}

type QosPolicyResourceModel struct {
	AddQosPolicyParameter    AddQosPolicyParameter    `tfsdk:"addQosPolicyParameter"`
	UpdateQosPolicyParameter UpdateQosPolicyParameter `tfsdk:"updateQosPolicyParameter"`
	DelQosPolicyParameter    DelQosPolicyParameter    `tfsdk:"delQosPolicyParameter"`
	ReadQosPolicyParameter   ReadQosPolicyParameter   `tfsdk:"readQosPolicyParameter"`
}

type AddQosPolicyParameter struct {
	VsysName        types.String `tfsdk:"vsysName"`
	QosName         types.String `tfsdk:"qosName"`
	Dir             types.String `tfsdk:"Dir"`
	Interface       types.String `tfsdk:"Interface"`
	PolicyType      types.String `tfsdk:"policyType"`
	FlowName        types.String `tfsdk:"flowName"`
	Ip              types.String `tfsdk:"Ip"`
	UpMinlinkRate   types.String `tfsdk:"UpMinlinkRate"`
	UplinkRate      types.String `tfsdk:"UplinkRate"`
	DownMinlinkRate types.String `tfsdk:"DownMinlinkRate"`
	DownlinkRate    types.String `tfsdk:"DownlinkRate"`
	MappingType     types.String `tfsdk:"MappingType"`
	Pri             types.String `tfsdk:"Pri"`
}

type UpdateQosPolicyParameter struct {
	VsysName        types.String `tfsdk:"vsysName"`
	QosName         types.String `tfsdk:"qosName"`
	Dir             types.String `tfsdk:"Dir"`
	Interface       types.String `tfsdk:"Interface"`
	PolicyType      types.String `tfsdk:"policyType"`
	FlowName        types.String `tfsdk:"flowName"`
	Ip              types.String `tfsdk:"Ip"`
	UpMinlinkRate   types.String `tfsdk:"UpMinlinkRate"`
	UplinkRate      types.String `tfsdk:"UplinkRate"`
	DownMinlinkRate types.String `tfsdk:"DownMinlinkRate"`
	DownlinkRate    types.String `tfsdk:"DownlinkRate"`
	MappingType     types.String `tfsdk:"MappingType"`
	Pri             types.String `tfsdk:"Pri"`
}

type DelQosPolicyParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
	QosName  types.String `tfsdk:"qosName"`
	Dir      types.String `tfsdk:"Dir"`
}

type ReadQosPolicyParameter struct {
	VsysName        types.String `tfsdk:"vsysName"`
	QosName         types.String `tfsdk:"qosName"`
	Dir             types.String `tfsdk:"Dir"`
	Interface       types.String `tfsdk:"Interface"`
	PolicyType      types.String `tfsdk:"policyType"`
	FlowName        types.String `tfsdk:"flowName"`
	Ip              types.String `tfsdk:"Ip"`
	UpMinlinkRate   types.String `tfsdk:"UpMinlinkRate"`
	UplinkRate      types.String `tfsdk:"UplinkRate"`
	DownMinlinkRate types.String `tfsdk:"DownMinlinkRate"`
	DownlinkRate    types.String `tfsdk:"DownlinkRate"`
	UplinkCounter   types.String `tfsdk:"UplinkCounter"`
	UplinkByte      types.String `tfsdk:"UplinkByte"`
	UplinkDrop      types.String `tfsdk:"UplinkDrop"`
	DownlinkCounter types.String `tfsdk:"DownlinkCounter"`
	DownlinkByte    types.String `tfsdk:"DownlinkByte"`
	DownlinkDrop    types.String `tfsdk:"DownlinkDrop"`
	MappingType     types.String `tfsdk:"MappingType"`
	Pri             types.String `tfsdk:"Pri"`
}

func (r *QosPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_QosPolicy"
}

func (r *QosPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *QosPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *QosPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *QosPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddQosPolicyRequest(ctx, "POST", r.client, data.AddQosPolicyParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QosPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *QosPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadQosPolicyRequest(ctx, "GET", r.client, data.ReadQosPolicyParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QosPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *QosPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateQosPolicyRequest(ctx, "PUT", r.client, data.UpdateQosPolicyParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QosPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *QosPolicyResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelQosPolicyRequest(ctx, "DELETE", r.client, data.DelQosPolicyParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *QosPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddQosPolicyRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddQosPolicyParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_policy/policylimit"

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

func sendToweb_UpdateQosPolicyRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateQosPolicyParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_policy/policylimit"

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

func sendToweb_DelQosPolicyRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelQosPolicyParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_policy/policylimit"

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

func sendToweb_ReadQosPolicyRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadQosPolicyParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/sw_qos/qos_policy/policylimit?qosName=”150”"

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
