package useragent

import (
	"github.com/ua-parser/uap-go/uaparser"
)

// Record is the data returned from the useragent package.
type Record struct {
	Raw            string `json:"raw"`
	UAFamily       string `json:"ua_family"`
	UAMajor        string `json:"ua_major"`
	UAMinor        string `json:"ua_minor"`
	UAPatch        string `json:"ua_patch"`
	UAOSFamily     string `json:"ua_os_family"`
	UAOSMajor      string `json:"ua_os_major"`
	UAOSMinor      string `json:"ua_os_minor"`
	UAOSPatch      string `json:"ua_os_patch"`
	UAOSPatchMinor string `json:"ua_os_patch_minor"`
	UADeviceFamily string `json:"ua_device_family"`
	UADeviceBrand  string `json:"ua_device_brand"`
	UADeviceModel  string `json:"ua_device_model"`
}

// ParseFile parses the given file.
func Parse(line string) (*Record, error) {
	parser, err := uaparser.NewFromBytes(uaparser.DefinitionYaml)
	if err != nil {
		return nil, err
	}

	client := parser.Parse(line)
	return &Record{
		Raw:            line,
		UAFamily:       client.UserAgent.Family,
		UAMajor:        client.UserAgent.Major,
		UAMinor:        client.UserAgent.Minor,
		UAPatch:        client.UserAgent.Patch,
		UAOSFamily:     client.Os.Family,
		UAOSMajor:      client.Os.Major,
		UAOSMinor:      client.Os.Minor,
		UAOSPatch:      client.Os.Patch,
		UAOSPatchMinor: client.Os.PatchMinor,
		UADeviceFamily: client.Device.Family,
		UADeviceBrand:  client.Device.Brand,
		UADeviceModel:  client.Device.Model,
	}, nil
}
