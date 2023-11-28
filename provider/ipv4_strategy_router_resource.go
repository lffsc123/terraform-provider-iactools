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

// IPv4策略路由
var _ resource.Resource = &Ipv4StrategyRouterResource{}
var _ resource.ResourceWithImportState = &Ipv4StrategyRouterResource{}

func NewIpv4StrategyRouterResource() resource.Resource {
	return &Ipv4StrategyRouterResource{}
}

type Ipv4StrategyRouterResource struct {
	client *Client
}

type Ipv4StrategyRouterResourceModel struct {
	AddIpv4StrategyRouterParameter AddIpv4StrategyRouterParameter `tfsdk:"rtplist"`
}

type AddIpv4StrategyRouterRequest struct {
	AddIpv4StrategyRouterRequestModel AddIpv4StrategyRouterRequestModel `json:"rtplist"`
}

type AddIpv4StrategyRouterRequestModel struct {
	RtpName string `json:"rtpName"`
	Act     string `json:"act"`
}

type AddIpv4StrategyRouterParameter struct {
	//VsysName      types.String `tfsdk:"vsysname"`
	//Resource      types.String `tfsdk:"resource"`
	//ListFlag      types.String `tfsdk:"listflag"`
	//Sequence      types.String `tfsdk:"sequence"`
	RtpName types.String `tfsdk:"rtpname"`
	//MatchSrctype  types.String `tfsdk:"matchsrctype"`
	//MatchDesttype types.String `tfsdk:"matchdesttype"`
	//Src           types.String `tfsdk:"src"`
	//Dst           types.String `tfsdk:"dst"`
	//IIfName       types.String `tfsdk:"iifname"`
	//Protocol      types.String `tfsdk:"protocol"`
	//MinSrcport    types.String `tfsdk:"minsrcport"`
	//MaxSrcport    types.String `tfsdk:"maxsrcport"`
	//MinDstport    types.String `tfsdk:"mindstport"`
	//MaxDstport    types.String `tfsdk:"maxdstport"`
	//Dscp          types.String `tfsdk:"dscp"`
	Act types.String `tfsdk:"act"`
	//Timestr       types.String `tfsdk:"timestr"`
	//Status        types.String `tfsdk:"status"`
	//OifName       types.String `tfsdk:"oifname"`
	//Weight        types.String `tfsdk:"weight"`
	//Nexthop       types.String `tfsdk:"nexthop"`
	//HcType        types.String `tfsdk:"hctype"`
	//HcName        types.String `tfsdk:"hcname"`
	//BfdCheck      types.String `tfsdk:"bfdcheck"`
	//NexthopCount  types.String `tfsdk:"nexthopcount"`
}

//type UpdateIpv4StrategyRouterParameter struct {
//	VsysName      types.String `tfsdk:"vsysName"`
//	Resource      types.String `tfsdk:"resource"`
//	ListFlag      types.String `tfsdk:"listFlag"`
//	Sequence      types.String `tfsdk:"sequence"`
//	RtpName       types.String `tfsdk:"rtpName"`
//	MatchSrctype  types.String `tfsdk:"matchSrctype"`
//	MatchDesttype types.String `tfsdk:"matchDesttype"`
//	Src           types.String `tfsdk:"src"`
//	Dst           types.String `tfsdk:"dst"`
//	IIfName       types.String `tfsdk:"iIfName"`
//	Protocol      types.String `tfsdk:"protocol"`
//	MinSrcport    types.String `tfsdk:"minSrcport"`
//	MaxSrcport    types.String `tfsdk:"maxSrcport"`
//	MinDstport    types.String `tfsdk:"minDstport"`
//	MaxDstport    types.String `tfsdk:"maxDstport"`
//	Dscp          types.String `tfsdk:"dscp"`
//	Act           types.String `tfsdk:"act"`
//	Timestr       types.String `tfsdk:"timestr"`
//	Status        types.String `tfsdk:"status"`
//	OifName       types.String `tfsdk:"oifName"`
//	Weight        types.String `tfsdk:"weight"`
//	Nexthop       types.String `tfsdk:"nexthop"`
//	HcType        types.String `tfsdk:"hcType"`
//	HcName        types.String `tfsdk:"hcName"`
//	BfdCheck      types.String `tfsdk:"bfdCheck"`
//	NexthopCount  types.String `tfsdk:"nexthopCount"`
//}
//
//type DelIpv4StrategyRouterParameter struct {
//	VsysName types.String `tfsdk:"vsysName"`
//	Resource types.String `tfsdk:"resource"`
//	ListFlag types.String `tfsdk:"listFlag"`
//	RtpName  types.String `tfsdk:"rtpName"`
//}
//
//type ReadIpv4StrategyRouterParameter struct {
//	ListFlag types.String `tfsdk:"listFlag"`
//	RtpName  types.String `tfsdk:"rtpName"`
//}

func (r *Ipv4StrategyRouterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "dpfirewall_Ipv4Strategy"
}

func (r *Ipv4StrategyRouterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rtplist": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"rtpname": schema.StringAttribute{
						Required: true,
					},
					"act": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *Ipv4StrategyRouterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Ipv4StrategyRouterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *Ipv4StrategyRouterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "created a resource **************")
	sendToweb_Ipv4StrategyRouterRequest(ctx, "POST", r.client, data.AddIpv4StrategyRouterParameter)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4StrategyRouterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *Ipv4StrategyRouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " read Start")
	sendToweb_Ipv4StrategyRouterRequest(ctx, "GET", r.client, data.AddIpv4StrategyRouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4StrategyRouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *Ipv4StrategyRouterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, " Update Start ************")
	sendToweb_Ipv4StrategyRouterRequest(ctx, "PUT", r.client, data.AddIpv4StrategyRouterParameter)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Ipv4StrategyRouterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *Ipv4StrategyRouterResourceModel
	tflog.Info(ctx, " Delete Start *************  ")

	//sendToweb_Ipv4StrategyRouterRequest(ctx, "DELETE", r.client, data.AddIpv4StrategyRouterParameter)

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

func (r *Ipv4StrategyRouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func sendToweb_Ipv4StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv4StrategyRouterParameter) {

	var sendData AddIpv4StrategyRouterRequestModel
	if reqmethod == "POST" {
		sendData = AddIpv4StrategyRouterRequestModel{
			RtpName: Rsinfo.RtpName.ValueString(),
			Act:     Rsinfo.Act.ValueString(),
		}
	} else if reqmethod == "GET" {
		sendData = AddIpv4StrategyRouterRequestModel{
			RtpName: Rsinfo.RtpName.ValueString(),
			Act:     Rsinfo.Act.ValueString(),
		}
	} else if reqmethod == "PUT" {
		sendData = AddIpv4StrategyRouterRequestModel{
			RtpName: Rsinfo.RtpName.ValueString(),
			Act:     Rsinfo.Act.ValueString(),
		}
	} else if reqmethod == "DELETE" {
		sendData = AddIpv4StrategyRouterRequestModel{
			RtpName: Rsinfo.RtpName.ValueString(),
		}
	}

	requstData := AddIpv4StrategyRouterRequest{
		AddIpv4StrategyRouterRequestModel: sendData,
	}
	body, _ := json.Marshal(requstData)
	//targetUrl := c.HostURL + "/func/web_main/api/rt_policy/rtpolicy/rtplist"
	targetUrl := "http://192.168.131.115:1888/api/ems-data-maintenance/equipment/detail?id=1716752407464554497"

	req, _ := http.NewRequest("GET", targetUrl, bytes.NewBuffer(body))
	//req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Blade-Auth", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZW5hbnRfaWQiOiIwMDAwMDAiLCJ1c2VyX25hbWUiOiJhZG1pbiIsInJlYWxfbmFtZSI6IueuoeeQhuWRmCIsImF2YXRhciI6Imh0dHBzOi8vZ3cuYWxpcGF5b2JqZWN0cy5jb20vem9zL3Jtc3BvcnRhbC9CaWF6ZmFueG1hbU5Sb3h4VnhrYS5wbmciLCJhdXRob3JpdGllcyI6WyJhZG1pbmlzdHJhdG9yIl0sImNsaWVudF9pZCI6InNhYmVyIiwicm9sZV9uYW1lIjoiYWRtaW5pc3RyYXRvciIsImxpY2Vuc2UiOiJwb3dlcmVkIGJ5IGJsYWRleCIsInBvc3RfaWQiOiIxMTIzNTk4ODE3NzM4Njc1MjAxIiwidXNlcl9pZCI6IjExMjM1OTg4MjE3Mzg2NzUyMDEiLCJyb2xlX2lkIjoiMTEyMzU5ODgxNjczODY3NTIwMSIsInNjb3BlIjpbImFsbCJdLCJuaWNrX25hbWUiOiLnrqHnkIblkZgiLCJvYXV0aF9pZCI6IiIsImRldGFpbCI6eyJ0eXBlIjoid2ViIn0sImV4cCI6MTcwMTE0NTIxMywiZGVwdF9pZCI6IjExMjM1OTg4MTM3Mzg2NzUyMDEiLCJqdGkiOiI1M2YxNWQxMi1lZjRlLTQyMDEtODAzNC0xMzRmMWQ3NjA0MTUiLCJhY2NvdW50IjoiYWRtaW4ifQ.84k6JLWWC6v60-ii6xlQcuV8YzShd6hquZgqdA3Y2wQ")
	//req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	// 创建一个HTTP客户端并发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	respn, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败：", err)
		return
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		fmt.Println("读取响应失败：", err2)
		return
	}
	// 打印响应结果
	fmt.Println("响应状态码:", respn.Status)
	fmt.Println("响应体:", string(body))
}

//func sendToweb_UpdateIpv4StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv4StrategyRouterParameter) {
//	requstData := Rsinfo
//
//	body, _ := json.Marshal(requstData)
//	targetUrl := c.HostURL + "/func/web_main/api/rt_policy/rtpolicy/rtplist"
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
//
//func sendToweb_DelIpv4StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv4StrategyRouterParameter) {
//	requstData := Rsinfo
//
//	body, _ := json.Marshal(requstData)
//	targetUrl := c.HostURL + "/func/web_main/api/rt_policy/rtpolicy/rtplist"
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
//
//func sendToweb_ReadIpv4StrategyRouterRequest(ctx context.Context, reqmethod string, c *Client, Rsinfo AddIpv4StrategyRouterParameter) {
//	requstData := Rsinfo
//
//	body, _ := json.Marshal(requstData)
//	targetUrl := c.HostURL + "/func/web_main/api/rt_policy/rtpolicy/rtplist?rtpName=ss&listFlag=0"
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
