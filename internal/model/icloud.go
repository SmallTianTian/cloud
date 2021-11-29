package model

import "net/http"

const (
	ICloudSetupEndpoint = "https://setup.icloud.com/setup/ws/1"
)

var (
	ICloudCommonParam = func(clientID string) []string {
		return []string{
			"clientBuildNumber=17DHotfix5",
			"clientMasteringNumber=17DHotfix5",
			"ckjsBuildVersion=17DProjectDev77",
			"ckjsVersion=2.0.5",
			"clientId=" + clientID,
		}
	}
	ICloudCommonHeader = func() http.Header {
		return http.Header{
			"Origin":     []string{"https://www.icloud.com"},
			"Referer":    []string{"https://www.icloud.com/"},
			"User-Agent": []string{"Opera/9.52 (X11; Linux i686; U; en)"},
		}
	}
)

type ICloudLoginRequest struct {
	AppleID  string `json:"apple_id,omitempty"`
	Password string `json:"password,omitempty"`
}

type ICloudLoginResponse struct {
	DsInfo struct {
		LastName                        string        `json:"lastName"`
		ICDPEnabled                     bool          `json:"iCDPEnabled"`
		TantorMigrated                  bool          `json:"tantorMigrated"`
		Dsid                            string        `json:"dsid"`
		HsaEnabled                      bool          `json:"hsaEnabled"`
		IroncadeMigrated                bool          `json:"ironcadeMigrated"`
		Locale                          string        `json:"locale"`
		BrZoneConsolidated              bool          `json:"brZoneConsolidated"`
		IsManagedAppleID                bool          `json:"isManagedAppleID"`
		IsCustomDomainsFeatureAvailable bool          `json:"isCustomDomainsFeatureAvailable"`
		IsHideMyEmailFeatureAvailable   bool          `json:"isHideMyEmailFeatureAvailable"`
		GilliganInvited                 bool          `json:"gilligan-invited"`
		AppleIDAliases                  []interface{} `json:"appleIdAliases"`
		HsaVersion                      int           `json:"hsaVersion"`
		IsPaidDeveloper                 bool          `json:"isPaidDeveloper"`
		CountryCode                     string        `json:"countryCode"`
		NotificationID                  string        `json:"notificationId"`
		PrimaryEmailVerified            bool          `json:"primaryEmailVerified"`
		ADsID                           string        `json:"aDsID"`
		Locked                          bool          `json:"locked"`
		HasICloudQualifyingDevice       bool          `json:"hasICloudQualifyingDevice"`
		PrimaryEmail                    string        `json:"primaryEmail"`
		AppleIDEntries                  []struct {
			IsPrimary bool   `json:"isPrimary"`
			Type      string `json:"type"`
			Value     string `json:"value"`
		} `json:"appleIdEntries"`
		GilliganEnabled    bool   `json:"gilligan-enabled"`
		FullName           string `json:"fullName"`
		LanguageCode       string `json:"languageCode"`
		AppleID            string `json:"appleId"`
		HasUnreleasedOS    bool   `json:"hasUnreleasedOS"`
		FirstName          string `json:"firstName"`
		ICloudAppleIDAlias string `json:"iCloudAppleIdAlias"`
		NotesMigrated      bool   `json:"notesMigrated"`
		BeneficiaryInfo    struct {
			IsBeneficiary bool `json:"isBeneficiary"`
		} `json:"beneficiaryInfo"`
		HasPaymentInfo bool   `json:"hasPaymentInfo"`
		PcsDeleted     bool   `json:"pcsDeleted"`
		AppleIDAlias   string `json:"appleIdAlias"`
		BrMigrated     bool   `json:"brMigrated"`
		StatusCode     int    `json:"statusCode"`
		FamilyEligible bool   `json:"familyEligible"`
	} `json:"dsInfo"`
	HasMinimumDeviceForPhotosWeb bool              `json:"hasMinimumDeviceForPhotosWeb"`
	ICDPEnabled                  bool              `json:"iCDPEnabled"`
	Webservices                  ICloudWebservices `json:"webservices"`
	PcsEnabled                   bool              `json:"pcsEnabled"`
	TermsUpdateNeeded            bool              `json:"termsUpdateNeeded"`
	ConfigBag                    struct {
		Urls struct {
			AccountCreateUI     string `json:"accountCreateUI"`
			AccountLoginUI      string `json:"accountLoginUI"`
			AccountLogin        string `json:"accountLogin"`
			AccountRepairUI     string `json:"accountRepairUI"`
			DownloadICloudTerms string `json:"downloadICloudTerms"`
			RepairDone          string `json:"repairDone"`
			AccountAuthorizeUI  string `json:"accountAuthorizeUI"`
			VettingURLForEmail  string `json:"vettingUrlForEmail"`
			AccountCreate       string `json:"accountCreate"`
			GetICloudTerms      string `json:"getICloudTerms"`
			VettingURLForPhone  string `json:"vettingUrlForPhone"`
		} `json:"urls"`
		AccountCreateEnabled string `json:"accountCreateEnabled"`
	} `json:"configBag"`
	HsaTrustedBrowser            bool     `json:"hsaTrustedBrowser"`
	ICloudAppsOrder              []string `json:"appsOrder"`
	Version                      int      `json:"version"`
	IsExtendedLogin              bool     `json:"isExtendedLogin"`
	PcsServiceIdentitiesIncluded bool     `json:"pcsServiceIdentitiesIncluded"`
	IsRepairNeeded               bool     `json:"isRepairNeeded"`
	HsaChallengeRequired         bool     `json:"hsaChallengeRequired"`
	RequestInfo                  struct {
		Country  string `json:"country"`
		TimeZone string `json:"timeZone"`
	} `json:"requestInfo"`
	PcsDeleted bool `json:"pcsDeleted"`
	ICloudInfo struct {
		SafariBookmarksHasMigratedToCloudKit bool `json:"SafariBookmarksHasMigratedToCloudKit"`
	} `json:"iCloudInfo"`
	ICloudApps struct {
		Calendar struct {
		} `json:"calendar"`
		Reminders struct {
		} `json:"reminders"`
		Keynote struct {
			IsQualifiedForBeta bool `json:"isQualifiedForBeta"`
		} `json:"keynote"`
		Settings struct {
			CanLaunchWithOneFactor bool `json:"canLaunchWithOneFactor"`
		} `json:"settings"`
		Pages struct {
			IsQualifiedForBeta bool `json:"isQualifiedForBeta"`
		} `json:"pages"`
		Mail struct {
			IsCKMail bool `json:"isCKMail"`
		} `json:"mail"`
		Notes3 struct {
		} `json:"notes3"`
		Find struct {
			CanLaunchWithOneFactor bool `json:"canLaunchWithOneFactor"`
		} `json:"find"`
		Iclouddrive struct {
		} `json:"iclouddrive"`
		Numbers struct {
			IsQualifiedForBeta bool `json:"isQualifiedForBeta"`
		} `json:"numbers"`
		Photos struct {
		} `json:"photos"`
		Contacts struct {
		} `json:"contacts"`
	} `json:"apps"`
}

// Requires2sa 是否需要两步验证
func (lrp *ICloudLoginResponse) Requires2sa() bool {
	return lrp.HsaChallengeRequired && lrp.DsInfo.HsaVersion >= 1
}

type ICloudWebservices struct {
	Reminders    ICloudApps `json:"reminders"`
	Ckdatabasews ICloudApps `json:"ckdatabasews"`
	Photosupload ICloudApps `json:"photosupload"`
	Photos       struct {
		ICloudApps
		UploadURL string `json:"uploadUrl"`
	} `json:"photos"`
	Drivews             ICloudApps `json:"drivews"`
	Uploadimagews       ICloudApps `json:"uploadimagews"`
	Schoolwork          ICloudApps `json:"schoolwork"`
	Cksharews           ICloudApps `json:"cksharews"`
	Findme              ICloudApps `json:"findme"`
	Ckdeviceservice     ICloudApps `json:"ckdeviceservice"`
	Iworkthumbnailws    ICloudApps `json:"iworkthumbnailws"`
	Calendar            ICloudApps `json:"calendar"`
	Docws               ICloudApps `json:"docws"`
	Settings            ICloudApps `json:"settings"`
	Premiummailsettings ICloudApps `json:"premiummailsettings"`
	Ubiquity            ICloudApps `json:"ubiquity"`
	Streams             ICloudApps `json:"streams"`
	Keyvalue            ICloudApps `json:"keyvalue"`
	Archivews           ICloudApps `json:"archivews"`
	Push                ICloudApps `json:"push"`
	Iwmb                ICloudApps `json:"iwmb"`
	Iworkexportws       ICloudApps `json:"iworkexportws"`
	Geows               ICloudApps `json:"geows"`
	Account             struct {
		ICloudApps
		ICloudEnv struct {
			ShortID   string `json:"shortId"`
			VipSuffix string `json:"vipSuffix"`
		} `json:"iCloudEnv"`
	} `json:"account"`
	Contacts ICloudApps `json:"contacts"`
}

type ICloudApps struct {
	URL         string `json:"url"`
	Status      string `json:"status"`
	PcsRequired bool   `json:"pcsRequired"`
}

type ICloudDevicesResponse struct {
	Devices []*ICloudDevice `json:"devices"`
}

type ICloudDevice struct {
	AreaCode    string `json:"areaCode"`
	DeviceID    string `json:"deviceId"`
	DeviceType  string `json:"deviceType"`
	PhoneNumber string `json:"phoneNumber"`
}

type ICloudVerificationCodeReq struct {
	VerificationCode string `json:"verificationCode"`
	TrustBrowser     bool   `json:"trustBrowser"`
}

type ICloudVerificationCodeResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorTitle   string `json:"errorTitle"`
	Success      bool   `json:"success"`
}

type ICloudPhotoInitRecords struct {
	Records []struct {
		Created struct {
			DeviceID       string `json:"deviceID"`
			Timestamp      int64  `json:"timestamp"`
			UserRecordName string `json:"userRecordName"`
		} `json:"created"`
		Deleted bool `json:"deleted"`
		Fields  struct {
			Progress struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"progress"`
			State struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"state"`
		} `json:"fields"`
		Modified struct {
			DeviceID       string `json:"deviceID"`
			Timestamp      int64  `json:"timestamp"`
			UserRecordName string `json:"userRecordName"`
		} `json:"modified"`
		PluginFields struct {
		} `json:"pluginFields"`
		RecordChangeTag string `json:"recordChangeTag"`
		RecordName      string `json:"recordName"`
		RecordType      string `json:"recordType"`
		ZoneID          struct {
			OwnerRecordName string `json:"ownerRecordName"`
			ZoneName        string `json:"zoneName"`
			ZoneType        string `json:"zoneType"`
		} `json:"zoneID"`
	} `json:"records"`
	SyncToken string `json:"syncToken"`
}

type ICloudPhotoFolders struct {
	Records []struct {
		Created struct {
			DeviceID       string `json:"deviceID"`
			Timestamp      int64  `json:"timestamp"`
			UserRecordName string `json:"userRecordName"`
		} `json:"created"`
		Deleted bool `json:"deleted"`
		Fields  struct {
			AlbumNameEnc struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"albumNameEnc"`
			AlbumType struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"albumType"`
			Position struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"position"`
			RecordModificationDate struct {
				Type  string `json:"type"`
				Value int64  `json:"value"`
			} `json:"recordModificationDate"`
			SortAscending struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"sortAscending"`
			SortType struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"sortType"`
			SortTypeExt struct {
				Type  string `json:"type"`
				Value int    `json:"value"`
			} `json:"sortTypeExt"`
		} `json:"fields,omitempty"`
		Modified struct {
			DeviceID       string `json:"deviceID"`
			Timestamp      int64  `json:"timestamp"`
			UserRecordName string `json:"userRecordName"`
		} `json:"modified"`
		PluginFields    struct{} `json:"pluginFields"`
		RecordChangeTag string   `json:"recordChangeTag"`
		RecordName      string   `json:"recordName"`
		RecordType      string   `json:"recordType"`
		ZoneID          struct {
			OwnerRecordName string `json:"ownerRecordName"`
			ZoneName        string `json:"zoneName"`
			ZoneType        string `json:"zoneType"`
		} `json:"zoneID"`
	} `json:"records"`
	SyncToken string `json:"syncToken"`
}

type ICloudPhotoInfos struct {
	Records            []ICloudPhotoInfo `json:"records"`
	ContinuationMarker string            `json:"continuationMarker"`
	SyncToken          string            `json:"syncToken"`
}

type ICloudPhotoInfo struct {
	RecordName string `json:"recordName"`
	RecordType string `json:"recordType"`
	Fields     struct {
		MasterRef struct {
			Value struct {
				Action     string `json:"action"`
				RecordName string `json:"recordName"`
				ZoneID     struct {
					OwnerRecordName string `json:"ownerRecordName"`
					ZoneName        string `json:"zoneName"`
					ZoneType        string `json:"zoneType"`
				} `json:"zoneID"`
			} `json:"value"`
			Type string `json:"type"`
		} `json:"masterRef"`
		ResVidSmallHeight struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidSmallHeight"`
		ItemType struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"itemType"`
		ResJPEGThumbFingerprint struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGThumbFingerprint"`
		FilenameEnc struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"filenameEnc"`
		OriginalOrientation struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"originalOrientation"`
		ResOriginalVidComplFileType struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalVidComplFileType"`
		ResJPEGMedFileType struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGMedFileType"`
		ResOriginalVidComplWidth struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalVidComplWidth"`
		ResJPEGThumbWidth struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGThumbWidth"`
		ResOriginalWidth struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalWidth"`
		DataClassType struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"dataClassType"`
		ResVidMedFileType struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidMedFileType"`
		ResOriginalFingerprint struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalFingerprint"`
		ResJPEGMedWidth struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGMedWidth"`
		ResVidMedRes struct {
			Value struct {
				FileChecksum      string `json:"fileChecksum"`
				Size              int    `json:"size"`
				WrappingKey       string `json:"wrappingKey"`
				ReferenceChecksum string `json:"referenceChecksum"`
				DownloadURL       string `json:"downloadURL"`
			} `json:"value"`
			Type string `json:"type"`
		} `json:"resVidMedRes"`
		ResOriginalHeight struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalHeight"`
		ResVidSmallFingerprint struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidSmallFingerprint"`
		ResVidMedWidth struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidMedWidth"`
		AssetDate struct {
			Value int64  `json:"value"`
			Type  string `json:"type"`
		} `json:"assetDate"`
		AddDate struct {
			Value int64  `json:"value"`
			Type  string `json:"type"`
		} `json:"addDate"`
		ResJPEGMedRes struct {
			Value struct {
				FileChecksum      string `json:"fileChecksum"`
				Size              int    `json:"size"`
				WrappingKey       string `json:"wrappingKey"`
				ReferenceChecksum string `json:"referenceChecksum"`
				DownloadURL       string `json:"downloadURL"`
			} `json:"value"`
			Type string `json:"type"`
		} `json:"resJPEGMedRes"`
		ResOriginalVidComplFingerprint struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalVidComplFingerprint"`
		ResJPEGMedHeight struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGMedHeight"`
		ResOriginalRes struct {
			Value struct {
				FileChecksum      string `json:"fileChecksum"`
				Size              int    `json:"size"`
				WrappingKey       string `json:"wrappingKey"`
				ReferenceChecksum string `json:"referenceChecksum"`
				DownloadURL       string `json:"downloadURL"`
			} `json:"value"`
			Type string `json:"type"`
		} `json:"resOriginalRes"`
		ResVidSmallFileType struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidSmallFileType"`
		ResVidSmallRes struct {
			Value struct {
				FileChecksum      string `json:"fileChecksum"`
				Size              int    `json:"size"`
				WrappingKey       string `json:"wrappingKey"`
				ReferenceChecksum string `json:"referenceChecksum"`
				DownloadURL       string `json:"downloadURL"`
			} `json:"value"`
			Type string `json:"type"`
		} `json:"resVidSmallRes"`
		ResJPEGThumbHeight struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGThumbHeight"`
		ResOriginalVidComplRes struct {
			Value struct {
				FileChecksum      string `json:"fileChecksum"`
				Size              int    `json:"size"`
				WrappingKey       string `json:"wrappingKey"`
				ReferenceChecksum string `json:"referenceChecksum"`
				DownloadURL       string `json:"downloadURL"`
			} `json:"value"`
			Type string `json:"type"`
		} `json:"resOriginalVidComplRes"`
		ResJPEGThumbFileType struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGThumbFileType"`
		ResVidMedFingerprint struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidMedFingerprint"`
		ResVidMedHeight struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidMedHeight"`
		ResOriginalVidComplHeight struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalVidComplHeight"`
		ResVidSmallWidth struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"resVidSmallWidth"`
		ResJPEGThumbRes struct {
			Value struct {
				FileChecksum      string `json:"fileChecksum"`
				Size              int    `json:"size"`
				WrappingKey       string `json:"wrappingKey"`
				ReferenceChecksum string `json:"referenceChecksum"`
				DownloadURL       string `json:"downloadURL"`
			} `json:"value"`
			Type string `json:"type"`
		} `json:"resJPEGThumbRes"`
		ResOriginalFileType struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resOriginalFileType"`
		ResJPEGMedFingerprint struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resJPEGMedFingerprint"`
	} `json:"fields"`
	PluginFields struct {
	} `json:"pluginFields"`
	RecordChangeTag string `json:"recordChangeTag"`
	Created         struct {
		Timestamp      int64  `json:"timestamp"`
		UserRecordName string `json:"userRecordName"`
		DeviceID       string `json:"deviceID"`
	} `json:"created"`
	Modified struct {
		Timestamp      int64  `json:"timestamp"`
		UserRecordName string `json:"userRecordName"`
		DeviceID       string `json:"deviceID"`
	} `json:"modified"`
	Deleted bool `json:"deleted"`
	ZoneID  struct {
		ZoneName        string `json:"zoneName"`
		OwnerRecordName string `json:"ownerRecordName"`
		ZoneType        string `json:"zoneType"`
	} `json:"zoneID"`
}

type ICloudPhotoFolderInfo struct {
	Batch []struct {
		Records []struct {
			Created struct {
				DeviceID       string `json:"deviceID"`
				Timestamp      int64  `json:"timestamp"`
				UserRecordName string `json:"userRecordName"`
			} `json:"created"`
			Deleted bool `json:"deleted"`
			Fields  struct {
				ItemCount struct {
					Type  string `json:"type"`
					Value int    `json:"value"`
				} `json:"itemCount"`
			} `json:"fields"`
			Modified struct {
				DeviceID       string `json:"deviceID"`
				Timestamp      int64  `json:"timestamp"`
				UserRecordName string `json:"userRecordName"`
			} `json:"modified"`
			PluginFields struct {
			} `json:"pluginFields"`
			RecordChangeTag string `json:"recordChangeTag"`
			RecordName      string `json:"recordName"`
			RecordType      string `json:"recordType"`
			ZoneID          struct {
				OwnerRecordName string `json:"ownerRecordName"`
				ZoneName        string `json:"zoneName"`
				ZoneType        string `json:"zoneType"`
			} `json:"zoneID"`
		} `json:"records"`
		SyncToken string `json:"syncToken"`
	} `json:"batch"`
}

var (
	ICloudSmartFolders = map[string]map[string]interface{}{
		"All Photos": {
			"obj_type":     "CPLAssetByAddedDate",
			"list_type":    "CPLAssetAndMasterByAddedDate",
			"direction":    "ASCENDING",
			"query_filter": nil,
		},
		"Time-lapse": {
			"obj_type":  "CPLAssetInSmartAlbumByAssetDate:Timelapse",
			"list_type": "CPLAssetAndMasterInSmartAlbumByAssetDate",
			"direction": "ASCENDING",
			"query_filter": []map[string]interface{}{{
				"fieldName":  "smartAlbum",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": "TIMELAPSE",
				}},
			},
		},
		"Videos": {
			"obj_type":  "CPLAssetInSmartAlbumByAssetDate:Video",
			"list_type": "CPLAssetAndMasterInSmartAlbumByAssetDate",
			"direction": "ASCENDING",
			"query_filter": []map[string]interface{}{{
				"fieldName":  "smartAlbum",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": "VIDEO",
				}},
			},
		},
		"Slo-mo": {
			"obj_type":  "CPLAssetInSmartAlbumByAssetDate:Slomo",
			"list_type": "CPLAssetAndMasterInSmartAlbumByAssetDate",
			"direction": "ASCENDING",
			"query_filter": []map[string]interface{}{{
				"fieldName":  "smartAlbum",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": "SLOMO",
				}},
			},
		},
		"Bursts": {
			"obj_type":     "CPLAssetBurstStackAssetByAssetDate",
			"list_type":    "CPLBurstStackAssetAndMasterByAssetDate",
			"direction":    "ASCENDING",
			"query_filter": nil,
		},
		"Favorites": {
			"obj_type":  "CPLAssetInSmartAlbumByAssetDate:Favorite",
			"list_type": "CPLAssetAndMasterInSmartAlbumByAssetDate",
			"direction": "ASCENDING",
			"query_filter": []map[string]interface{}{{
				"fieldName":  "smartAlbum",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": "FAVORITE",
				}},
			},
		},
		"Panoramas": {
			"obj_type":  "CPLAssetInSmartAlbumByAssetDate:Panorama",
			"list_type": "CPLAssetAndMasterInSmartAlbumByAssetDate",
			"direction": "ASCENDING",
			"query_filter": []map[string]interface{}{{
				"fieldName":  "smartAlbum",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": "PANORAMA",
				}},
			},
		},
		"Screenshots": {
			"obj_type":  "CPLAssetInSmartAlbumByAssetDate:Screenshot",
			"list_type": "CPLAssetAndMasterInSmartAlbumByAssetDate",
			"direction": "ASCENDING",
			"query_filter": []map[string]interface{}{{
				"fieldName":  "smartAlbum",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": "SCREENSHOT",
				}},
			},
		},
		"Live": {
			"obj_type":  "CPLAssetInSmartAlbumByAssetDate:Live",
			"list_type": "CPLAssetAndMasterInSmartAlbumByAssetDate",
			"direction": "ASCENDING",
			"query_filter": []map[string]interface{}{{
				"fieldName":  "smartAlbum",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": "LIVE",
				}},
			},
		},
		"Recently Deleted": {
			"obj_type":     "CPLAssetDeletedByExpungedDate",
			"list_type":    "CPLAssetAndMasterDeletedByExpungedDate",
			"direction":    "ASCENDING",
			"query_filter": nil,
		},
		"Hidden": {
			"obj_type":     "CPLAssetHiddenByAssetDate",
			"list_type":    "CPLAssetAndMasterHiddenByAssetDate",
			"direction":    "ASCENDING",
			"query_filter": nil,
		},
	}
	ICloudDesiredKeys = []string{
		"resJPEGFullWidth", "resJPEGFullHeight",
		"resJPEGFullFileType", "resJPEGFullFingerprint",
		"resJPEGFullRes", "resJPEGLargeWidth",
		"resJPEGLargeHeight", "resJPEGLargeFileType",
		"resJPEGLargeFingerprint", "resJPEGLargeRes",
		"resJPEGMedWidth", "resJPEGMedHeight",
		"resJPEGMedFileType", "resJPEGMedFingerprint",
		"resJPEGMedRes", "resJPEGThumbWidth",
		"resJPEGThumbHeight", "resJPEGThumbFileType",
		"resJPEGThumbFingerprint", "resJPEGThumbRes",
		"resVidFullWidth", "resVidFullHeight",
		"resVidFullFileType", "resVidFullFingerprint",
		"resVidFullRes", "resVidMedWidth", "resVidMedHeight",
		"resVidMedFileType", "resVidMedFingerprint",
		"resVidMedRes", "resVidSmallWidth", "resVidSmallHeight",
		"resVidSmallFileType", "resVidSmallFingerprint",
		"resVidSmallRes", "resSidecarWidth", "resSidecarHeight",
		"resSidecarFileType", "resSidecarFingerprint",
		"resSidecarRes", "itemType", "dataClassType",
		"filenameEnc", "originalOrientation", "resOriginalWidth",
		"resOriginalHeight", "resOriginalFileType",
		"resOriginalFingerprint", "resOriginalRes",
		"resOriginalAltWidth", "resOriginalAltHeight",
		"resOriginalAltFileType", "resOriginalAltFingerprint",
		"resOriginalAltRes", "resOriginalVidComplWidth",
		"resOriginalVidComplHeight", "resOriginalVidComplFileType",
		"resOriginalVidComplFingerprint", "resOriginalVidComplRes",
		"isDeleted", "isExpunged", "dateExpunged", "remappedRef",
		"recordName", "recordType", "recordChangeTag",
		"masterRef", "adjustmentRenderType", "assetDate",
		"addedDate", "isFavorite", "isHidden", "orientation",
		"duration", "assetSubtype", "assetSubtypeV2",
		"assetHDRType", "burstFlags", "burstFlagsExt", "burstId",
		"captionEnc", "locationEnc", "locationV2Enc",
		"locationLatitude", "locationLongitude", "adjustmentType",
		"timeZoneOffset", "vidComplDurValue", "vidComplDurScale",
		"vidComplDispValue", "vidComplDispScale",
		"vidComplVisibilityState", "customRenderedValue",
		"containerId", "itemId", "position", "isKeyAsset",
	}
)
