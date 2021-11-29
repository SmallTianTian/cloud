package pkg

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"tianxu.xin/phone/cloud/internal/http"
	"tianxu.xin/phone/cloud/internal/model"
	"tianxu.xin/phone/cloud/pkg/photo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ Client = new(ICloud)

type ICloud struct {
	log      *zap.Logger
	appleID  string
	password string

	clientID    string
	commonParam string
	client      *http.Client

	webServiceRoute *model.ICloudWebservices
}

func NewICloud(account, password, cookieDir string, log *zap.Logger) (*ICloud, error) {
	clientID, err := parseClientID(cookieDir)
	if err != nil {
		return nil, err
	}

	client, err := http.New(path.Join(cookieDir, "icloud_cookie"), model.ICloudCommonHeader(), log)
	if err != nil {
		return nil, err
	}

	return &ICloud{
		appleID:     account,
		password:    password,
		commonParam: strings.Join(model.ICloudCommonParam(clientID), "&"),
		clientID:    clientID,
		client:      client,
	}, nil
}

func (ic *ICloud) Login(ctx context.Context) error {
	loginURL := fmt.Sprintf("%s/login?%s", model.ICloudSetupEndpoint, ic.commonParam)

	req := &model.ICloudLoginRequest{AppleID: ic.appleID, Password: ic.password}
	var resp model.ICloudLoginResponse
	if err := ic.client.Do(ctx, http.POST, loginURL, req, nil, &resp); err != nil {
		return err
	}
	defer ic.client.FlushCookie()

	if resp.Requires2sa() {
		return ic.Requires2sa(ctx)
	}
	ic.commonParam += "&dsid=" + resp.DsInfo.Dsid
	ic.webServiceRoute = &resp.Webservices
	return nil
}

func (ic *ICloud) Requires2sa(ctx context.Context) error {
	listDevicesURL := fmt.Sprintf("%s/listDevices?%s", model.ICloudSetupEndpoint, ic.commonParam)
	var resp model.ICloudDevicesResponse
	if err := ic.client.Do(ctx, http.GET, listDevicesURL, nil, nil, &resp); err != nil {
		return err
	}

	var index int
	for _, d := range resp.Devices {
		fmt.Printf("%d: %s to %s", index, d.DeviceType, d.DeviceType)
	}

	for {
		var code string
		fmt.Println("请选择将设备，或直接输入验证码")
		if _, err := fmt.Scanln(&code); err != nil {
			return err
		}
		codeInt, err := strconv.Atoi(code)
		if err != nil {
			fmt.Println(code, "不是数字")
			continue
		}
		if len(code) != 6 {
			if codeInt >= len(resp.Devices) {
				fmt.Println("没有该序列的设备")
				continue
			}
			if err = ic.sendVerificationCode(ctx, resp.Devices[codeInt]); err != nil {
				return err
			}
		}
		ok, err := ic.validateVerificationCode(ctx, code)
		if err != nil {
			return err
		}
		if ok {
			break
		}
	}
	// relogin
	return ic.Login(ctx)
}

func (ic *ICloud) PhotoService(ctx context.Context) (photo.Service, error) {
	return photo.NewICloudService(ctx, ic.client, ic.webServiceRoute.Ckdatabasews.URL, ic.commonParam)
}

func (ic *ICloud) sendVerificationCode(ctx context.Context, d *model.ICloudDevice) error {
	sendCodeURL := fmt.Sprintf("%s/sendVerificationCode?%s", model.ICloudSetupEndpoint, ic.commonParam)
	return ic.client.Do(ctx, http.POST, sendCodeURL, d, nil, nil)
}

func (ic *ICloud) validateVerificationCode(ctx context.Context, code string) (bool, error) {
	cerificationCodeURL := fmt.Sprintf("%s/validateVerificationCode?%s",
		model.ICloudSetupEndpoint, ic.commonParam)
	var rp model.ICloudVerificationCodeResponse
	if err := ic.client.Do(ctx, http.GET, cerificationCodeURL, &model.ICloudVerificationCodeReq{
		VerificationCode: code,
		TrustBrowser:     true,
	}, nil, &rp); err != nil {
		return false, err
	}
	if rp.Success {
		return true, nil
	}
	ic.log.Error("icloud validate verification code failed", zap.String("code", code),
		zap.Int("error_code", rp.ErrorCode), zap.String("error_msg", rp.ErrorMessage),
		zap.String("error_title", rp.ErrorTitle))
	return false, nil
}

func parseClientID(dir string) (string, error) {
	var clientID string
	file := path.Join(dir, "icloud_client_id")
	bs, err := os.ReadFile(file)
	if err == nil {
		clientID = string(bs)
	} else {
		clientID = uuid.NewString()
		if err = os.WriteFile(file, []byte(clientID), 0644); err != nil {
			return "", err
		}
	}
	return clientID, nil
}
