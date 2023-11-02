package tunnel

import (
	"log"
	"net/url"
	"strconv"
)

func ParseUrl(rtspUrl string) (*url.URL, uint, error) {
	URL, err := url.Parse(rtspUrl)
	if err != nil {
		log.Println("Error parsing rtspUrl: ", rtspUrl)
		return nil, 0, err
	}

	queryParams, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		log.Println("Error parsing query parameters:", err)
		return nil, 0, err
	}

	// tunnel port
	var tunnelPort uint
	if tunnelPortStr := queryParams.Get("tunnelPort"); tunnelPortStr != "" {
		if tunnel, err := strconv.ParseInt(tunnelPortStr, 10, 0); err == nil {
			tunnelPort = uint(tunnel)
		}
	}

	if tunnelPort == 0 {
		tunnelPort = getPort(URL, 80)
	}

	// Modify the URL with updated query parameters
	URL.RawQuery = removeQueryKey(URL.RawQuery, "tunnelPort")

	return URL, tunnelPort, nil
}

func getPort(URL *url.URL, defaultPort uint) uint {
	if URL.Port() != "" {
		if port, err := strconv.ParseUint(URL.Port(), 10, 0); err == nil {
			return uint(port)
		}
	}
	return defaultPort
}

func removeQueryKey(rawQuery string, key string) string {
	queryParams, _ := url.ParseQuery(rawQuery)
	delete(queryParams, key)
	return queryParams.Encode()
}
