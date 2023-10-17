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

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RealServiceListResource{}
var _ resource.ResourceWithImportState = &RealServiceListResource{}

func NewRealServiceListResource() resource.Resource {
	return &RealServiceListResource{}
}

// ExampleResource defines the resource implementation.
type RealServiceListResource struct {
	client *Client
}

// ExampleResourceModel describes the resource data model.
type RealServiceListResourceModel struct {
	// Uuid_count types.String `tfsdk:"uuid_count"`
	Poollist RealServiceListParameter `tfsdk:"poollist"`
}

type RealServiceListParameter struct {
	Name     types.String `tfsdk:"name"`
	Monitor  types.String `tfsdk:"monitor"`
	RsList   types.String `tfsdk:"rs_list"`
	Schedule types.String `tfsdk:"schedule"`
}

func (r *RealServiceListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo_RealServiceList"
}

func (r *RealServiceListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"poollist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"monitor": schema.StringAttribute{
						Optional: true,
					},
					"rs_list": schema.StringAttribute{
						Optional: true,
					},
					"schedule": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

func (r *RealServiceListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RealServiceListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RealServiceListResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource")
	sendToweb_RealServiceListRequest(ctx, "POST", r.client, data.Poollist)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RealServiceListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RealServiceListResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	// sendToweb_RealServiceListRequest(ctx,"POST", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RealServiceListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *RealServiceListResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_RealServiceListRequest(ctx, "PUT", r.client, data.Poollist)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RealServiceListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RealServiceListResourceModel
	tflog.Info(ctx, " Delete Start")

	sendToweb_RealServiceListRequest(ctx, "DELETE", r.client, data.Poollist)

	// Read Terraform prior state data into the model
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

func (r *RealServiceListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_RealServiceListRequest(ctx context.Context, reqmethod string, c *Client, datamodel RealServiceListParameter) {
	sendData := RealServiceListRequestModel{
		Name:     datamodel.Name.ValueString(),
		Monitor:  datamodel.Monitor.ValueString(),
		RsList:   datamodel.RsList.ValueString(),
		Schedule: datamodel.Schedule.ValueString(),
	}

	requstData := RealServiceListRequest{
		Poollist: sendData,
	}
	body, _ := json.Marshal(requstData)

	targetUrl := c.HostURL + "/func/web_main/api/slb/pool/adx_slb_pool/poollist"

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
