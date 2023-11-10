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

// VPN全局配置
var _ resource.Resource = &VpnSystemConfResource{}
var _ resource.ResourceWithImportState = &VpnSystemConfResource{}

func NewVpnSystemConfResource() resource.Resource {
	return &VpnSystemConfResource{}
}

type VpnSystemConfResource struct {
	client *Client
}

type VpnSystemConfResourceModel struct {
	//AddVpnSystemConfParameter    AddVpnSystemConfParameter    `tfsdk:"addVpnSystemConfParameter"`
	UpdateVpnSystemConfParameter UpdateVpnSystemConfParameter `tfsdk:"updateVpnSystemConfParameter"`
	//DelVpnSystemConfParameter    DelVpnSystemConfParameter    `tfsdk:"delVpnSystemConfParameter"`
	ReadVpnSystemConfParameter ReadVpnSystemConfParameter `tfsdk:"readVpnSystemConfParameter"`
}

type AddVpnSystemConfParameter struct {
}

type UpdateVpnSystemConfParameter struct {
	VsysName   types.String `tfsdk:"vsysName"`
	status     types.String `tfsdk:"status"`
	port       types.String `tfsdk:"port"`
	tunnelType types.String `tfsdk:"tunnelType"`
}

type DelVpnSystemConfParameter struct {
}

type ReadVpnSystemConfParameter struct {
	VsysName types.String `tfsdk:"vsysName"`
}

func (r *VpnSystemConfResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpsc_VpnSystemConf"
}

func (r *VpnSystemConfResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *VpnSystemConfResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VpnSystemConfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//var data *VpnSystemConfResourceModel
	//resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	tflog.Trace(ctx, "created a resource **************")
	//sendToweb_AddVpnSystemConfRequest(ctx, "POST", r.client, data.AddVpnSystemConfParameter)
	//// Save data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnSystemConfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VpnSystemConfResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start ***************")
	sendToweb_ReadVpnSystemConfRequest(ctx, "GET", r.client, data.ReadVpnSystemConfParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnSystemConfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VpnSystemConfResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_UpdateVpnSystemConfRequest(ctx, "PUT", r.client, data.UpdateVpnSystemConfParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VpnSystemConfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data *VpnSystemConfResourceModel
	tflog.Info(ctx, " Delete Start *************")

	//sendToweb_DelVpnSystemConfRequest(ctx, "DELETE", r.client, data.DelVpnSystemConfParameter)
	//
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}
}

func (r *VpnSystemConfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//func sendToweb_AddVpnSystemConfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddVpnSystemConfParameter) {
//	requstData := Rsinfo
//
//	body, _ := json.Marshal(requstData)
//	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipPoolInfo"
//
//	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Accept", "application/json")
//	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
//	respn, err := http.DefaultClient.Do(req)
//	if err != nil {
//		tflog.Info(ctx, " read Error"+err.Error())
//	}
//	defer respn.Body.Close()
//
//	body, err2 := ioutil.ReadAll(respn.Body)
//	if err2 == nil {
//		fmt.Println(string(body))
//	}
//}

func sendToweb_UpdateVpnSystemConfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo UpdateVpnSystemConfParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/systemConf"

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

//func sendToweb_DelVpnSystemConfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo DelVpnSystemConfParameter) {
//	requstData := Rsinfo
//
//	body, _ := json.Marshal(requstData)
//	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/ipPoolInfo"
//
//	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Accept", "application/json")
//	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
//	respn, err := http.DefaultClient.Do(req)
//	if err != nil {
//		tflog.Info(ctx, " read Error"+err.Error())
//	}
//	defer respn.Body.Close()
//
//	body, err2 := ioutil.ReadAll(respn.Body)
//	if err2 == nil {
//		fmt.Println(string(body))
//	}
//}

func sendToweb_ReadVpnSystemConfRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo ReadVpnSystemConfParameter) {
	requstData := Rsinfo

	body, _ := json.Marshal(requstData)
	targetUrl := c.HostURL + "/func/web_main/api/vpn/ssl_vpn/sslvpn/systemConf?vsysName=PublicSystem"

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
