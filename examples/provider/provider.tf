
terraform {
 required_providers {
  dptech-demo={
     source = "registry.terraform.io/xieguihua123/dptech-demo"
     version = "1.2.37"
   } 
 }
 }

provider "dptech-demo" {
  address="http://localhost:"
  port="8080"
  username="test"
  password="jsepc123!"
}

resource "dptech-demo_RealService" "cs" {
  rsinfo={
  name="rs1"
  address="192.168.1.1"
  port="8091"
}
}

resource "dptech-demo_RealServiceList" "cs2" {
 poollist={
  name="string__*"
  monitor="string"
  rs_list="string"
  schedule="string"
 }
}

resource "dptech-demo_AddrPoolList" "cs" {
addrpoollist={
  name="string__*"
  ip_start="string__*"
  	ip_end="string__*"
	ip_version="string"
	vrrp_if_name="string"//接口名称
	vrrp_id="string"    //vrid
} 
}

resource "dptech-demo_VirtualService" "cs" {
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


