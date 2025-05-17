package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type CurlRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Body    string
}

func ParseCurl(curlCmd string) (*CurlRequest, error) {
	parts := splitCurlArgs(curlCmd)

	var method string
	url := ""
	headers := make(map[string]string)
	var body string

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		switch part {
		case "--location":
			if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "--") && url == "" {
				i++
				url = parts[i]
			}
		case "-X", "--request":
			i++
			if i < len(parts) {
				method = parts[i]
			}
		case "-H", "--header":
			i++
			if i < len(parts) {
				key, value := ParseHeader(parts[i])
				headers[key] = value
			}
		case "-d", "--data", "--data-raw":
			i++
			if i < len(parts) {
				body = parts[i]
			}
			if method == "" {
				method = "POST"
			}
		default:
			if !strings.HasPrefix(part, "-") && url == "" && !strings.HasPrefix(part, "curl") {
				url = part
			}
		}
	}

	if method == "" {
		method = "GET"
	}

	if url == "" {
		return nil, fmt.Errorf("解析失败")
	}

	baseURL, params := ParseURLWithParams(url)

	res := &CurlRequest{
		Method:  method,
		URL:     baseURL,
		Headers: headers,
		Params:  params,
		Body:    body,
	}
	return res, nil

}

func splitCurlArgs(curlCmd string) []string {
	var parts []string
	var currentPart strings.Builder
	var inSingleQuote bool

	for i := 0; i < len(curlCmd); i++ {
		char := curlCmd[i]

		if char == '\'' {
			inSingleQuote = !inSingleQuote
		}

		if char == ' ' && !inSingleQuote {
			if currentPart.Len() > 0 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
			}
		} else if char != '\'' {
			currentPart.WriteByte(char)
		}

		if i == len(curlCmd)-1 && currentPart.Len() > 0 {
			parts = append(parts, currentPart.String())
		}
	}

	return parts
}

func ParseHeader(header string) (string, string) {
	parts := strings.SplitN(header, ": ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func ParseURLWithParams(inputURL string) (string, map[string]string) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return inputURL, nil
	}

	params := make(map[string]string)
	for key, values := range u.Query() {
		params[key] = values[0]
	}

	return u.Scheme + "://" + u.Host + u.Path, params
}
