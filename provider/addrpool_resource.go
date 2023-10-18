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
var _ resource.Resource = &AddrPoolResource{}
var _ resource.ResourceWithImportState = &AddrPoolResource{}

func NewAddrPoolResource() resource.Resource {
	return &AddrPoolResource{}
}

// ExampleResource defines the resource implementation.
type AddrPoolResource struct {
	client *Client
}

// ExampleResourceModel describes the resource data model.
type AddrPoolResourceModel struct {
	// Uuid_count types.String `tfsdk:"uuid_count"`
	Addrpoollist AddrPoolParameter `tfsdk:"addrpoollist"`
}

type AddrPoolParameter struct {
	Name       types.String `tfsdk:"name"`
	IpVersion  types.String `tfsdk:"ip_version"`
	IpStart    types.String `tfsdk:"ip_start"`
	IpEnd      types.String `tfsdk:"ip_end"`
	VrrpIfName types.String `tfsdk:"vrrp_if_name"` //接口名称
	VrrpId     types.String `tfsdk:"vrrp_id"`      //vrid
}

func (r *AddrPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dptech-demo_AddrPoolList"
}

func (r *AddrPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"addrpoollist": schema.SingleNestedAttribute{
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

func (r *AddrPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AddrPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AddrPoolResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource")
	sendToweb_AddrPoolRequest(ctx, "POST", r.client, data.Addrpoollist)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddrPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AddrPoolResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	// sendToweb_AddrPoolRequest(ctx,"POST", r.client, data.Rsinfo)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddrPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AddrPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_AddrPoolRequest(ctx, "PUT", r.client, data.Addrpoollist)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddrPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AddrPoolResourceModel
	tflog.Info(ctx, " Delete Start")

	sendToweb_AddrPoolRequest(ctx, "DELETE", r.client, data.Addrpoollist)

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

func (r *AddrPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_AddrPoolRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddrPoolParameter) {
	sendData := AddrPoolRequestModel{
		Name:       Rsinfo.Name.ValueString(),
		IpStart:    Rsinfo.IpStart.ValueString(),
		IpEnd:      Rsinfo.IpEnd.ValueString(),
		IpVersion:  Rsinfo.IpVersion.ValueString(),
		VrrpIfName: Rsinfo.VrrpIfName.ValueString(),
		VrrpId:     Rsinfo.VrrpId.ValueString(),
	}
	requstData := AddrPoolRequest{
		Addrpoollist: sendData,
	}

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/addrpool/addrpool/addrpoollist"

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
