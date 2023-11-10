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

// 服务对象
var _ resource.Resource = &TimeObjResource{}
var _ resource.ResourceWithImportState = &TimeObjResource{}

func NewTimeObjResource() resource.Resource {
	return &TimeObjResource{}
}

type TimeObjResource struct {
	client *Client
}

type TimeObjResourceModel struct {
	AddTimeObjParameter    AddTimeObjParameter    `tfsdk:"addTimeObjParameter"`
	UpdateTimeObjParameter UpdateTimeObjParameter `tfsdk:"updateTimeObjParameter"`
	DelTimeObjParameter    DelTimeObjParameter    `tfsdk:"delTimeObjParameter"`
	ReadTimeObjParameter   ReadTimeObjParameter   `tfsdk:"readTimeObjParameter"`
}

type AddTimeObjParameter struct {
	VsysName  types.String `tfsdk:"vsysName"`
	Name      types.String `tfsdk:"name"`
	Mode      types.String `tfsdk:"mode"`
	Week      types.String `tfsdk:"week"`
	StartDay  types.String `tfsdk:"startDay"`
	EndDay    types.String `tfsdk:"endDay"`
	StartTime types.String `tfsdk:"startTime"`
	EndTime   types.String `tfsdk:"endTime"`
}

type UpdateTimeObjParameter struct {
	VsysName  types.String `tfsdk:"vsysName"`
	Name      types.String `tfsdk:"name"`
	Mode      types.String `tfsdk:"mode"`
	Week      types.String `tfsdk:"week"`
	StartDay  types.String `tfsdk:"startDay"`
	EndDay    types.String `tfsdk:"endDay"`
	StartTime types.String `tfsdk:"startTime"`
	EndTime   types.String `tfsdk:"endTime"`
}

type DelTimeObjParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
	Name     types.String `tfsdk:"name"`
}

type ReadTimeObjParameter struct {
	VsysName  types.String `tfsdk:"vsysName"`
	Name      types.String `tfsdk:"name"`
	Mode      types.String `tfsdk:"mode"`
	Week      types.String `tfsdk:"week"`
	StartDay  types.String `tfsdk:"startDay"`
	EndDay    types.String `tfsdk:"endDay"`
	StartTime types.String `tfsdk:"startTime"`
	EndTime   types.String `tfsdk:"endTime"`
}

func (r *TimeObjResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_TimeObj"
}

func (r *TimeObjResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *TimeObjResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TimeObjResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TimeObjResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_AddTimeObjRequest(ctx, "POST", r.client, data.AddTimeObjParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TimeObjResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TimeObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadTimeObjRequest(ctx, "GET", r.client, data.ReadTimeObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TimeObjResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TimeObjResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateTimeObjRequest(ctx, "PUT", r.client, data.UpdateTimeObjParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TimeObjResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TimeObjResourceModel
	tflog.Info(ctx, " Delete Start *************")

	sendToweb_DelTimeObjRequest(ctx, "DELETE", r.client, data.DelTimeObjParameter)

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TimeObjResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddTimeObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddTimeObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj"

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

func sendToweb_UpdateTimeObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateTimeObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj"

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

func sendToweb_DelTimeObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelTimeObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj"

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

func sendToweb_ReadTimeObjRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadTimeObjParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/netservice/netservice/usrobj?vfwName=vsys&searchValue=&offset=1&count=100"

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
