provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_vpc" "test_vpc" {
  name       = "hello"
  cidr_block = "10.1.0.0/16"
}