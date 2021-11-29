package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
)

const (
	POST = "POST"
	GET  = "GET"
)

type Client struct {
	log        *zap.Logger
	cookieFile string

	commonHeader http.Header
	cookieCache  map[*url.URL][]*http.Cookie
	client       *http.Client
}

func New(cookieFile string, commonHeader http.Header, log *zap.Logger) (*Client, error) {
	if _, err := os.Stat(path.Dir(cookieFile)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(cookieFile), os.ModeDir); err != nil {
			return nil, err
		}
	}

	cookieMap, err := parseCookie(cookieFile)
	if err != nil {
		return nil, err
	}

	sc, _ := cookiejar.New(nil)
	for k, v := range cookieMap {
		sc.SetCookies(k, v)
	}

	c := &http.Client{}
	c.Jar = sc
	return &Client{
		log:          log,
		client:       c,
		cookieFile:   cookieFile,
		cookieCache:  cookieMap,
		commonHeader: commonHeader,
	}, nil
}

func (client *Client) Do(ctx context.Context, method, url string,
	body interface{}, header http.Header, respStruct interface{}) error {
	stream, err := client.Stream(ctx, method, url, body, header)
	if err != nil {
		return err
	}
	defer stream.Close()

	bs, err := io.ReadAll(stream)
	if err != nil {
		return err
	}
	client.log.Sugar().Debug("<--", string(bs))
	if respStruct == nil {
		return nil
	}
	return json.Unmarshal(bs, respStruct)
}

func (client *Client) Stream(ctx context.Context, method, url string,
	body interface{}, header http.Header) (io.ReadCloser, error) {
	var bodyBytes io.Reader
	var bs []byte
	if body != nil {
		var err error
		switch v := body.(type) {
		case string:
			bs = []byte(v)
		default:
			bs, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		}
		bodyBytes = bytes.NewBuffer(bs)
	}

	client.log.Sugar().Debug("-->", method, url, string(bs))
	req, err := http.NewRequest(method, url, bodyBytes)
	if err != nil {
		return nil, err
	}

	req.Header = client.commonHeader.Clone()
	if len(header) != 0 {
		for k, v := range header {
			var value string
			if len(v) != 0 {
				value = v[0]
			}
			req.Header.Set(k, value)
		}
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	client.cookieCache[req.URL] = client.client.Jar.Cookies(req.URL)
	return resp.Body, nil
}

func parseCookie(file string) (map[*url.URL][]*http.Cookie, error) {
	uc := make(map[*url.URL][]*http.Cookie)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return uc, nil
	}

	bs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	tmp := make(map[string][]*http.Cookie)
	if err = json.Unmarshal(bs, &tmp); err != nil {
		return nil, err
	}
	for k, v := range tmp {
		u, err := url.Parse(k)
		if err != nil {
			continue
		}
		uc[u] = v
	}

	return uc, nil
}

func base64Enc(s string) string {
	codes := base64.StdEncoding.EncodeToString([]byte(s))
	codes = strings.ReplaceAll(codes, "/", "_")
	return codes
}

func base64Dec(s string) string {
	s = strings.ReplaceAll(s, "_", "/")
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s
	}
	return string(bs)
}

func (client *Client) FlushCookie() error {
	if len(client.cookieCache) == 0 {
		return nil
	}
	tmp := make(map[string][]*http.Cookie)
	for k, v := range client.cookieCache {
		tmp[k.String()] = v
	}
	bs, err := json.Marshal(tmp)
	if err != nil {
		return err
	}

	return os.WriteFile(client.cookieFile, bs, 0644)
}
