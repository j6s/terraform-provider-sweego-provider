
resource sweego_domain "test_domain" {
  domain = "foo.com"

  # Optional
  open_tracking_enabled = false
  click_tracking_enabled = false
}