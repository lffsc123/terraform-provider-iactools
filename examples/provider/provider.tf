terraform {
 required_providers {
  dpfirewall={
     source = "registry.terraform.io/lffsc123/dpfirewall"
     version = "1.2.37"
   } 
 }
 }

provider "dpfirewall" {
  address="http://localhost:"
  port="8080"
  username="test"
  password="jsepc123!"
}

resource "dpfirewall_RealServiceList" "cs2" {
 poollist={
  name="string__*"
  monitor="string" // 健康监测
  rs_list="string"
  schedule="string"
 }
}

resource "dpfirewall_AddrPoolList" "cs" {
addrpoollist={
    name="string__*"
    ip_start="string__*"
  	ip_end="string__*"
	ip_version="string"
	vrrp_if_name="string"//接口名称
	vrrp_id="string"    //vrid
} 
}

resource "dpfirewall_VirtualService" "cs" {
    virtualservice={
      name ="string__*"
      mode ="string__*"
      ip ="string__*"
      port ="string__*"
      protocol ="string"
      state ="string"
      session_keep ="string" // 会话保持
      default_pool ="string"
      tcp_policy ="string"
      snat ="string"
      session_bkp ="string"
      vrrp ="string"      //涉及普通双机热备场景，需要关联具体的vrrp组
  }
}

resource "dpfirewall-TargetNat" "dpcs" {
  targetnat={
    name="string__*"
    ip_start="string__*"
    ip_end="string__*"
    ip_version="string"
    vrrp_if_name="string"//接口名称
    vrrp_id="string"    //vrid
  }
}

resource "dpfirewall-SourceNat" "dpcs" {
  sourcenat={
    name="string__*"
    ip_start="string__*"
    ip_end="string__*"
    ip_version="string"
    vrrp_if_name="string"//接口名称
    vrrp_id="string"    //vrid
  }
}

// 会话保持
resource "dpfirewall-SessionKeep" "dpcs" {
  sessionkeep={
    // 1）必须配置。
    // 2）会话保持策略名称。
    // 3）用户名称不可以包涵”[~`!@#$%^&,*\\\"'<>?]\\+”等特殊字符,用户名称最大支持63个字符。
    name="string__*"

    // 1）可选配置，默认配置为源IP会话保持source-ip。
    // 2）会话保持策略类型 。
    // 3）各类型表示方法： 源IP会话保持：source-ip、目的IP会话保持：destination-ip、Cookie会话保持：cookie、RADIUS会话保持：radius、
    // DHCP会话保持：dhcp、L2TP会话保持：l2tp、HASH会话保持：hash、SIP会话保持：sip、HTTP被动会话保持：http-passive、SSL-ID会话保持：ssl-id
    type="string__*"
  }
}

// 健康监测
resource "dpfirewall-AdxSlbMonitor" "dpcs" {
  monitorinfo={
    // 1）必须配置
    // 2）新建健康监测的名称
    // 3）最大支持63个字符，不可以包涵
    // ”[~`!@#$%^&,*\\\"'<>?]\\+”特殊字符
    name="string__*"

    // 1）可选配置，默认为” ICMP”
    // 2）健康监测类型
    // 3）有如下健康监测类型可选：TCP、HTTP、ICMP、DNS、Radius、SNMP、SMTP、FTP、POP3、UDP、SSL、SIP、External；
    type="string__*"

    // 1）可选配置，默认值为16
    // 2）健康监测超时时间配置
    // 3）取值范围1-2147483647
    overtime="string__*"

    // 1）可选配置，默认取值为5
    // 2）健康监测检测间隔配置
    // 3）取值范围1-2147483647
    interval="string__*"
  }
}


