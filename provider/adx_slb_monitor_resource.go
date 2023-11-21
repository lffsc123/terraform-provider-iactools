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

// 健康监测
var _ resource.Resource = &AdxSlbMonitorResource{}
var _ resource.ResourceWithImportState = &AdxSlbMonitorResource{}

func NewAdxSlbMonitorResource() resource.Resource {
	return &AdxSlbMonitorResource{}
}

// ExampleResource defines the resource implementation.
type AdxSlbMonitorResource struct {
	client *Client
}

// ExampleResourceModel describes the resource data model.
type AdxSlbMonitorResourceModel struct {
	Rsinfo AdxSlbMonitorParameter `tfsdk:"monitorinfo"`
}

type AdxSlbMonitorParameter struct {
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	OverTime types.String `tfsdk:"overtime"`
	Interval types.String `tfsdk:"interval"`
}

func (r *AdxSlbMonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_AdxSlbMonitor"
}

func (r *AdxSlbMonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"monitorinfo": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
					"overtime": schema.StringAttribute{
						Required: true,
					},
					"interval": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *AdxSlbMonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AdxSlbMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AdxSlbMonitorResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource")
	sendToweb_AdxSlbMonitorRequest(ctx, "POST", r.client, data.Rsinfo)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AdxSlbMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AdxSlbMonitorResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	// sendToweb_AdxSlbMonitorRequest(ctx,"POST", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AdxSlbMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AdxSlbMonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_AdxSlbMonitorRequest(ctx, "PUT", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AdxSlbMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AdxSlbMonitorResourceModel
	tflog.Info(ctx, " Delete Start")

	//sendToweb_AdxSlbMonitorRequest(ctx, "DELETE", r.client, data.Rsinfo)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AdxSlbMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AdxSlbMonitorRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AdxSlbMonitorParameter) {
	sendData := AdxSlbMonitorRequestModel{
		Name:     Rsinfo.Name.ValueString(),
		Type:     Rsinfo.Type.ValueString(),
		OverTime: Rsinfo.OverTime.ValueString(),
		Interval: Rsinfo.Interval.ValueString(),
	}

	requstData := AdxSlbMonitorRequest{
		adxSlbMonitor: sendData,
	}

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/slb/vs/virtual/virtualservice"

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
