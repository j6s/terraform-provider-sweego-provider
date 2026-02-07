provider "sweego" {
  api_key = "..."
  client_id = "..."
}

resource sweego_domain "test_domain" {
  domain = "foo.com"
  open_tracking_enabled = false
  click_tracking_enabled = true
}