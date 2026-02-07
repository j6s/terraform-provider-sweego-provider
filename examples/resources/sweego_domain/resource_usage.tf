#
# Example for adding DNS records using INWX provider
#
resource sweego_domain "test_domain" {
  domain = "foo.com"
}

#
# NOTE: Depending on the DNS Provider you use, the values must be modified a bit:
# - `name`: is the prefix before the domain. If your provider requires a full domain name
#   (e.g. the INWX provider used here) the full domain name must be appended.
# - `data`: Is a CNAME record, ends with a dot. Again, depending on your DNS provider, you
#    may need to trim the dot at the end.
#
resource "inwx_nameserver_record" "test_dkim" {
  domain = "foo.com"
  name = "${resource.sweego_domain.test_domain.dkim_record.name}.foo.com"
  content = trim(resource.sweego_domain.test_domain.dkim_record.data, ".")
  type = resource.sweego_domain.test_domain.dkim_record.type
}