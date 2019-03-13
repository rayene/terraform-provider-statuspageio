variable "statuspageio_api_key" {
  type = "string"
}

variable "statuspageio_api_url" {
  type = "string"
}

variable "statuspageio_page" {
  type = "string"
}

provider "statuspageio" {
  api_key = "${var.statuspageio_api_key}"
  api_url = "${var.statuspageio_api_url}"
}

resource "statuspageio_component_group" "group0" {
  name = "my group 0" 
  page = "${var.statuspageio_page}"
  components = [
    "${statuspageio_component.component1.id}",
    "${statuspageio_component.component0.id}"
  ]
}

resource "statuspageio_component" "component0" {
  name = "my component 0" 
  page = "${var.statuspageio_page}"
}

resource "statuspageio_component" "component1" {
  name = "my component 1" 
  page = "${var.statuspageio_page}"
}
