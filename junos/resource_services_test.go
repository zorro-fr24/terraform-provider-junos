package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosServices_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosServicesConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services.testacc",
							"security_intelligence.#", "1"),
					),
				},
				{
					Config: testAccJunosServicesConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_services.testacc",
							"security_intelligence.#", "1"),
						resource.TestCheckResourceAttr("junos_services.testacc",
							"security_intelligence.0.default_policy.#", "1"),
					),
				},
				{
					Config: testAccJunosServicesConfigUpdate2(),
				},
				{
					ResourceName:      "junos_services_proxy_profile.testacc_services",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_services.testacc",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosServicesConfigPostCheck(),
				},
			},
		})
	}
}

func testAccJunosServicesConfigCreate() string {
	return `
resource "junos_services_proxy_profile" "testacc_services" {
  name               = "testacc_services"
  protocol_http_host = "192.0.2.1"
  protocol_http_port = 3128
}
resource "junos_services" "testacc" {
  application_identification {
    application_system_cache {}
    download {
      automatic_start_time = "12-24.22:00"
    }
    enable_performance_mode {}
    max_transactions = 10
  }
  security_intelligence {
    authentication_token = "abcdefghijklmnopqrstuvwxyz123456"
    category_disable     = ["all"]
    proxy_profile        = junos_services_proxy_profile.testacc_services.name
    url                  = "https://example.com/api/manifest.xml"
    url_parameter        = "test_param"
  }
}
`
}

func testAccJunosServicesConfigUpdate() string {
	return `
resource "junos_services_proxy_profile" "testacc_services" {
  name               = "testacc_services"
  protocol_http_host = "192.0.2.2"
  protocol_http_port = 3129
}
resource "junos_services_security_intelligence_profile" "testacc_services" {
  name     = "testacc_services"
  category = "IPFilter"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services" "testacc" {
  application_identification {
    application_system_cache {
      security_services = true
    }
    application_system_cache_timeout = 120
    download {
      automatic_interval       = 120
      automatic_start_time     = "12-24.22:00"
      ignore_server_validation = true
      proxy_profile            = junos_services_proxy_profile.testacc_services.name
      url                      = "https://example.com/"
    }
    enable_performance_mode {
      max_packet_threshold = 50
    }
    imap_cache_size     = 120
    imap_cache_timeout  = 120
    max_transactions    = 10
    statistics_interval = 120
  }
  security_intelligence {
    authentication_token = "abcdefghijklmnopqrstuvwxyz123400"
    category_disable     = ["CC"]
    default_policy {
      category_name = "IPFilter"
      profile_name  = junos_services_security_intelligence_profile.testacc_services.name
    }
    proxy_profile = junos_services_proxy_profile.testacc_services.name
    url           = "https://example.com/api/manifest.xml"
    url_parameter = "test_param_update"
  }
}
`
}

func testAccJunosServicesConfigUpdate2() string {
	return `
resource "junos_services_proxy_profile" "testacc_services" {
  name               = "testacc_services"
  protocol_http_host = "192.0.2.2"
  protocol_http_port = 3129
}
resource "junos_services_security_intelligence_profile" "testacc_services" {
  name     = "testacc_services"
  category = "IPFilter"
  rule {
    name = "rule_1"
    match {
      threat_level = [1]
    }
    then_action = "permit"
  }
}
resource "junos_services" "testacc" {
  application_identification {
    no_application_system_cache = true
  }
}
  `
}

func testAccJunosServicesConfigPostCheck() string {
	return `
resource "junos_services_proxy_profile" "testacc_services" {
  name               = "testacc_services"
  protocol_http_host = "192.0.2.2"
  protocol_http_port = 3129
}
resource "junos_services" "testacc" {
}
`
}
