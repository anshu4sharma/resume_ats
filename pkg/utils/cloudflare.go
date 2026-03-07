package utils

import (
	"bufio"
	"net/http"
)

func FetchCloudflareIPs() ([]string, error) {
	urls := []string{
		"https://www.cloudflare.com/ips-v4",
		"https://www.cloudflare.com/ips-v6",
	}

	var ips []string

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			ips = append(ips, scanner.Text())
		}
	}

	return ips, nil
}