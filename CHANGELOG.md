## 0.2.1 - 2026-02-07
### Changed
* Improved documentation for provider registry, no functional changes

## 0.2.0 - 2026-02-06
### Added
* `open_tracking_enabled` and `click_tracking_enabled` are now configurable properties
* When creating, updating or refreshing a domain, a check request is sent in order to
  get up-to-date verification information.

### Fixed
* Non-20x HTTP Status codes now reliably lead to errors

## 0.1.0 - 2026-01-31
### Added
* The whole provider
