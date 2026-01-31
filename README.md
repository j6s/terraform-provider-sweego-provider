# Sweego Provider

A simple terraform provider in order to manage domains at [sweego](https://www.sweego.io/).
I develop this provider in order to help developers and administrators use european based
alternatives for SMTP relaying.

## Setup

In order to start using the provider, you will need to create API Credentials. This can be done
by clicking on your account name in the top right and selecting "Credentials". You can then use
the newly created API Credentials in order to configure the provider:

```terraform
provider "sweego" {
  api_key = "YOUR_API_KEY"
  client_id = "YOUR_CLIENT_ID"
}
```

## Usage

### `sweego_domain`

The `sweego_domain` resource can be used in order to manage domains in sweego. After the domain
has been created, it will contain information about DNS records required to start using sweego as
a SMTP relay.

```terraform
resource sweego_domain "test_domain" {
  domain = "your-domain.eu"
}
```

The resulting resource will contain the following data:

| Path                     | Type                   | Description                                                              |
|--------------------------|------------------------|--------------------------------------------------------------------------|
| `domain`                 | string                 | Full domain name                                                         |
| `uuid`                   | string                 | ID of the domain in sweegos system                                       |
| `tracking_click_enabled` | bool                   | Whether or not click tracking is enabled for this domain                 |
| `tracking_open_enabled`  | bool                   | Whether or not open tracking is enabled for this domain                  |
| `is_verified`            | bool                   | Whether or not this domain is verified                                   |
| `domain_record`          | object(DnsRecord)      | CNAME DNS Record that needs to be set in order to verify the domain      |
| `dkim_record`            | object(DnsRecord)      | DKIM DNS Record that needs to be set in order to send E-Mails            |
| `dmarc_record`           | object(DnsRecord)      | DMARC DNS Record that needs to be set in order to send E-Mails           |
| `tracking_record`        | object(DnsRecord)      | CNAME DNS Record that needs to be set in order to use tracking           |
| `inbound_record_list`    | list(object(DnsRecord) | List of DNS Records that need to be set, if sweego should accept E-Mails |

With each `DnsRecord` having the following properties:

| Path   | Type   | Description                                                        |
|--------|--------|--------------------------------------------------------------------|
| `type` | string | Type of the record (e.g. `TXT`, `CNAME`, ...)                      |
| `name` | string | Name of the record without the full domain (e.g. `abc.sweego.co.`) |
| `data` | string | Value of the record                                                |

### Importing

Existing domains can be imported by their UUID. This value is not visible in sweegos user interface
directly, but can be obtained by looking at the UUID in the DNS records.

```terraform
terraform import sweego_domain.my_domain 3923bb62-f1e2-4362-ad1f-1af9f54d10f0
```

### Using records

The resulting properties of the resource can be used to create the correct DNS records with an approriate
provider. 

```terraform
resource sweego_domain "test_domain" {
  domain = "your-domain.eu"
}

resource inwx_nameserver_record "test_domain_record" {
  domain = "your-domain.eu"
  type = resource.sweego_domain.test_domain.domain_record.type
  
  # NOTE: the .name property contains the name of the DNS record without the rest of the
  #       domain name. Depending on your DNS record provider, it may need to be suffixed with
  #       that.
  name = "${resource.sweego_domain.test_domain.domain_record.name}.your-domain.eu"
  content = resource.sweego_domain.test_domain.domain_record.data
}
```