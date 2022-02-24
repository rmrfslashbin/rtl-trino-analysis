package rtl

import (
	"time"

	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/geoip"
	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/useragent"
)

// Record represents a properly formatted and typed Trino entry
type Record struct {
	Timestamp                time.Time         `json:"timestamp"`
	ClientIPAddr             string            `json:"client_ip_addr"`
	Status                   int               `json:"status"`
	Bytes                    int64             `json:"bytes"`
	Method                   string            `json:"method"`
	Protocol                 string            `json:"protocol"`
	Host                     string            `json:"host"`
	UriStem                  string            `json:"uri_stem"`
	EdgeLocation             string            `json:"edge_location"`
	EdgeRequestID            string            `json:"edge_request_id"`
	HostHeader               string            `json:"host_header"`
	TimeTaken                float64           `json:"time_taken"`
	ProtoVersion             string            `json:"proto_version"`
	IPVersion                string            `json:"ip_version"`
	Referer                  string            `json:"referer"`
	Cookie                   string            `json:"cookie"`
	UriQuery                 string            `json:"uri_query"`
	EdgeResponseResultType   string            `json:"edge_response_result_type"`
	SslProtocol              string            `json:"ssl_protocol"`
	SslCipher                string            `json:"ssl_cipher"`
	EdgeResultType           string            `json:"edge_result_type"`
	ContentType              string            `json:"content_type"`
	ContentLength            int64             `json:"content_length"`
	EdgeDetailedResultType   string            `json:"edge_detailed_result_type"`
	Country                  string            `json:"country"`
	CacheBehaviorPathPattern string            `json:"cache_behavior_path_pattern"`
	Year                     int               `json:"year"`
	Month                    int               `json:"month"`
	Day                      int               `json:"day"`
	ClientIP                 *geoip.GeoIPData  `json:"geoip_data" gorm:"embedded"`
	UserAgent                *useragent.Record `json:"useragent_data" gorm:"embedded"`
}
