package provider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id`
	Username string `json:"username`
	Token    string `json:"token"`
}
type RealServiceRequest struct {
	Rsinfo RealServiceRequestModel `json:"rsinfo"`
}

type RealServiceRequestModel struct {
	Name                string `json:"name"`
	Address             string `json:"address"`
	Port                string `json:"port"`
	Weight              string `json:"weight,omitempty"`
	ConnectionLimit     string `json:"connectionLimit,omitempty"`
	ConnectionRateLimit string `json:"connectionRateLimit,omitempty"`
	RecoveryTime        string `json:"recoveryTime,omitempty"`
	WarmTime            string `json:"warmTime,omitempty"`
	Monitor             string `json:"monitor,omitempty"`
	MonitorList         string `json:"monitorList,omitempty"`
	LeastNumber         string `json:"leastNumber,omitempty"`
	Priority            string `json:"priority,omitempty"`
	MonitorLog          string `json:"monitorLog,omitempty"`
	SimulTunnelsLimit   string `json:"simulTunnelsLimit,omitempty"`
	CpuWeight           string `json:"cpuWeight,omitempty"`
	MemoryWeight        string `json:"memoryWeight,omitempty"`
	State               string `json:"state,omitempty"`
	VsysName            string `json:"vsysName,omitempty"`
}
type RealServiceListRequest struct {
	Poollist RealServiceListRequestModel `json:"poollist"`
}

type RealServiceListRequestModel struct {
	Name     string `json:"name"`
	Monitor  string `json:"monitor,omitempty"`
	RsList   string `json:"rsList,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

type AddrPoolRequest struct {
	Addrpoollist AddrPoolRequestModel `json:"addrpoollist"`
}

type AddrPoolRequestModel struct {
	Name       string `json:"name"`
	IpVersion  string `json:"ipVersion,omitempty"`
	IpStart    string `json:"ipStart"`
	IpEnd      string `json:"ipEnd"`
	VrrpIfName string `json:"vrrpIfName,omitempty"` //接口名称
	VrrpId     string `json:"vrrpId,omitempty"`     //vrid
}

type VirtualServiceRequest struct {
	Virtualservice VirtualServiceRequestModel `json:"virtualservice"`
}

type VirtualServiceRequestModel struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	Mode        string `json:"mode"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
	Protocol    string `json:"protocol"`
	SessionKeep string `json:"sessionKeep"`
	DefaultPool string `json:"defaultPool"`
	TcpPolicy   string `json:"tcpPolicy"` //引用tcp超时时间，不引用默认600s
	Snat        string `json:"snat"`
	SessionBkp  string `json:"sessionBkp"` //必须配置集群模式
	Vrrp        string `json:"vrrp"`       //涉及普通双机热备场景，需要关联具体的vrrp组
}

type SessionKeepRequest struct {
	sessionkeep SessionKeepRequestModel `json:"sessionkeep"`
}

type SessionKeepRequestModel struct {
	Name string `json:"name"`
	Type string `json:"type"`
	//ActiveTime        string `json:"activeTime"`
	//OverrideLimit     string `json:"overrideLimit"`
	//MatchAcrossVs     string `json:"matchAcrossVs"`
	//MatchFailAction   string `json:"matchFailAction"`
	//PrefixType        string `json:"prefixType"`
	//PrefixLen         string `json:"prefixLen"`
	//Sessioncookie     string `json:"sessioncookie"`
	//Cookiename        string `json:"cookiename"`
	//Cookiemode        string `json:"cookiemode"`
	//CookieEncryMode   string `json:"cookieEncryMode"`
	//CookieEncryPasswd string `json:"cookieEncryPasswd"`
	//RadiusAttribute   string `json:"radiusAttribute"`
	//RelatePolicy      string `json:"relatePolicy"`
	//RadiusStopProc    string `json:"radiusStopProc"`
	//HashCondition     string `json:"hashCondition"`
	//SipHeadType       string `json:"sipHeadType"`
	//SipHeadName       string `json:"sipHeadName"`
	//RequestContent    string `json:"requestContent"`
	//ReplyContent      string `json:"replyContent"`
	//VirtualSystem     string `json:"virtualSystem"`
	//HttponlyAttribute string `json:"httponlyAttribute"`
	//SecureAttribute   string `json:"secureAttribute"`
}

type AdxSlbMonitorRequest struct {
	adxSlbMonitor AdxSlbMonitorRequestModel `json:"monitorinfo"`
}

type AdxSlbMonitorRequestModel struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	OverTime string `json:"overTime"`
	Interval string `json:"interval"`
}

func NewClient(host *string, auth *AuthStruct) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
		HostURL: *host,
		Auth:    *auth,
	}

	// req, err := http.NewRequest("POST", c.HostURL, nil)
	// req.Header.Add("Content-type", "application/json")
	// req.Header.Set("Accept", "application/json")
	// req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	// if err != nil {
	// 	return nil, err
	// }
	// body, err := c.doRequest(req)
	// if err != nil {
	// 	return nil, err
	// }
	// ar := AuthResponse{}
	// err = json.Unmarshal(body, &ar)
	// if err != nil {
	// 	return nil, err
	// }

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
