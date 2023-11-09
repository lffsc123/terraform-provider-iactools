terraform {
 required_providers {
  dptech-demo={
     source = "registry.terraform.io/lffsc123/dpsc"
     version = "1.2.37"
   } 
 }
 }

provider "dpsc" {
  address="http://localhost:"
  port="8080"
  username="test"
  password="jsepc123!"
}

resource "dpsc_SourceNat" "cs" {
  sourcenat={
    ipVersion="string"
    vsysName="string"
    name="string"
    targetName="string"
    position="string"
    outInterface="string"
    srcAddrObj="string"
    srcAddrGroup="string"
    dstAddrObj="string"
    dstAddrGroup="string"
    preService="string"
    usrService="string"
    serviceGroup="string"
    publicIpAddressFlag="string"
    addrpoolName="string"
    minPort="string"
    maxPort="string"
    portHash="string"
    state="string"
  }
}

resource "dpsc_TargetNat" "cs" {
  targetnat={
    vsysName="string"
    name="string"
    targetName="string"
    position="string"
    inInterface="string"
    srcIpObj="string"
    srcIpGroup="string"
    publicIp="string"
    preService="string"
    usrService="string"
    serviceGroup="string"
    inNetIp="string"
    innetPort="string"
    unLimited="string"
    srcIpTranslate="string"
    interfaceAddressFlag="string"
    addrpoolName="string"
    vrrpIfName="string"
    vrrpId="string"
    state="string"
  }
}

resource "dpsc_RealServiceList" "cs2" {
 poollist={
  name="string__*"
  monitor="string"
  rs_list="string"
  schedule="string"
 }
}

resource "dpsc_AddrPoolList" "cs" {
addrpoollist={
  name="string__*"
  ip_start="string__*"
  	ip_end="string__*"
	ip_version="string"
	vrrp_if_name="string"//接口名称
	vrrp_id="string"    //vrid
} 
}

resource "dpsc_VirtualService" "cs" {
    virtualservice={
    name ="string__*"
      mode ="string__*"
      ip ="string__*"
      port ="string__*"
      protocol ="string"
    state ="string"
      session_keep ="string"
      default_pool ="string"
      tcp_policy ="string"
      snat ="string"
      session_bkp ="string"
      vrrp ="string"      //涉及普通双机热备场景，需要关联具体的vrrp组
  }
}

resource "dpsc-TargetNat" "dpcs" {
  targetnat={
    name="string__*"
    ip_start="string__*"
    ip_end="string__*"
    ip_version="string"
    vrrp_if_name="string"//接口名称
    vrrp_id="string"    //vrid
  }
}

resource "dpsc-SourceNat" "dpcs" {
  sourcenat={
    name="string__*"
    ip_start="string__*"
    ip_end="string__*"
    ip_version="string"
    vrrp_if_name="string"//接口名称
    vrrp_id="string"    //vrid
  }
}


