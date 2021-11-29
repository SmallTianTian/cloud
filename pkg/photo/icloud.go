package photo

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"tianxu.xin/phone/cloud/internal/http"
	"tianxu.xin/phone/cloud/internal/model"
)

type ICloudService struct {
	client   *http.Client
	endpoint string
	param    string
}

func NewICloudService(ctx context.Context, client *http.Client, serviceRoot, param string) (
	*ICloudService, error) {
	url := fmt.Sprintf("%s/database/1/com.apple.photos.cloud/production/private", serviceRoot)

	param += "&" + strings.Join([]string{"remapEnums=true", "getCurrentSyncToken=true"}, "&")
	ps := &ICloudService{
		client:   client,
		endpoint: url,
		param:    param,
	}

	queryURL := fmt.Sprintf("%s/records/query?%s", ps.endpoint, ps.param)
	req := map[string]map[string]string{
		"query":  {"recordType": "CheckIndexingState"},
		"zoneID": {"zoneName": "PrimarySync"},
	}
	var recode model.ICloudPhotoInitRecords
	if err := ps.client.Do(ctx, http.POST, queryURL, req, nil, &recode); err != nil {
		return nil, err
	}
	indexingState := recode.Records[0].Fields.State.Value
	if indexingState != "FINISHED" {
		return nil, errors.New("iCloud Photo Library not finished indexing.")
	}
	return ps, nil
}

func (ps *ICloudService) Albums(ctx context.Context) (map[string]Album, error) {
	albums := make(map[string]Album)
	for k, v := range model.ICloudSmartFolders {
		qf, _ := v["query_filter"].([]map[string]interface{})
		pa := newPhotoAlbum(ps, k, v["list_type"].(string), v["obj_type"].(string),
			v["direction"].(string), qf)
		albums[k] = pa
	}
	folders, err := ps.fetchFolders(ctx)
	if err != nil {
		return nil, err
	}

	for _, folder := range folders.Records {
		// FIXME: Handle subfolders
		if folder.RecordName == "----Root-Folder----" ||
			folder.RecordName == "----Project-Root-Folder----" ||
			folder.Deleted {
			continue
		}
		folderID := folder.RecordName
		objType := "CPLContainerRelationNotDeletedByAssetDate:" + folderID
		var name string
		namebs, err := base64.StdEncoding.DecodeString(folder.Fields.AlbumNameEnc.Value)
		if err != nil {
			fmt.Println("err", err)
			name = folder.Fields.AlbumNameEnc.Value
		} else {
			name = string(namebs)
		}
		queryFilter := []map[string]interface{}{
			{
				"fieldName":  "parentId",
				"comparator": "EQUALS",
				"fieldValue": map[string]string{
					"type":  "STRING",
					"value": folderID,
				},
			},
		}
		albums[name] = newPhotoAlbum(ps, name, "CPLContainerRelationLiveByAssetDate", objType,
			"ASCENDING", queryFilter)
	}
	return albums, nil
}

func (ps *ICloudService) All(ctx context.Context) ([]Item, error) {
	albums, err := ps.Albums(ctx)
	if err != nil {
		return nil, err
	}
	return albums["All Photos"].Photos(ctx)
}

func (ps *ICloudService) fetchFolders(ctx context.Context) (*model.ICloudPhotoFolders, error) {
	queryURL := fmt.Sprintf("%s/records/query?%s", ps.endpoint, ps.param)
	body := `{"query":{"recordType":"CPLAlbumByPositionLive"},"zoneID":{"zoneName":"PrimarySync"}}`
	header := map[string][]string{"Content-type": {"text/plain"}}

	var folders model.ICloudPhotoFolders
	if err := ps.client.Do(ctx, http.POST, queryURL, body, header, &folders); err != nil {
		return nil, err
	}
	return &folders, nil
}

type PhotoAlbum struct {
	name        string
	service     *ICloudService
	listType    string
	objType     string
	direction   string
	queryFilter []map[string]interface{}
	pageSize    int
	len         int
}

func newPhotoAlbum(s *ICloudService, name, listType, objType, direction string,
	queryFilter []map[string]interface{}) *PhotoAlbum {

	return &PhotoAlbum{
		name:        name,
		service:     s,
		listType:    listType,
		objType:     objType,
		direction:   direction,
		queryFilter: queryFilter,
		pageSize:    100,
	}
}

func (pa *PhotoAlbum) Photos(ctx context.Context) ([]Item, error) {
	var offset int
	if pa.direction == "DESCENDING" {
		var err error
		if offset, err = pa.Len(ctx); err != nil {
			return nil, err
		}
		offset--
	}

	var result []Item
	for {
		queryURL := fmt.Sprintf("%s/records/query?%s", pa.service.endpoint, pa.service.param)
		req := pa.listQueryGen(offset, pa.listType, pa.direction, pa.queryFilter)

		var pi model.ICloudPhotoInfos
		if err := pa.service.client.Do(ctx, http.POST, queryURL, req, nil, &pi); err != nil {
			return nil, err
		}

		assetRecords := make(map[string]model.ICloudPhotoInfo)
		masterRecods := make([]model.ICloudPhotoInfo, 0)
		for _, rec := range pi.Records {
			switch rec.RecordType {
			case "CPLAsset":
				masterID := rec.Fields.MasterRef.Value.RecordName
				assetRecords[masterID] = rec
			case "CPLMaster":
				masterRecods = append(masterRecods, rec)
			}
		}

		mrl := len(masterRecods)
		if mrl == 0 {
			break
		}

		if pa.direction == "DESCENDING" {
			offset -= mrl
		} else {
			offset += mrl
		}

		for _, mr := range masterRecods {
			result = append(result, newICloudPhotoInfo(pa.service, mr, assetRecords[mr.RecordName]))
		}
	}
	return result, nil
}

func (pa *PhotoAlbum) Len(ctx context.Context) (int, error) {
	url := fmt.Sprintf("%s/internal/records/query/batch?%s", pa.service.endpoint, pa.service.param)
	var rp model.ICloudPhotoFolderInfo
	if err := pa.service.client.Do(ctx, http.POST, url, pa.countQueryGen(pa.objType), nil, &rp); err != nil {
		return 0, err
	}
	return rp.Batch[0].Records[0].Fields.ItemCount.Value, nil
}

func (pa *PhotoAlbum) countQueryGen(objType string) map[string][]map[string]interface{} {
	return map[string][]map[string]interface{}{
		"batch": {{
			"resultsLimit": 1,
			"query": map[string]interface{}{
				"filterBy": map[string]interface{}{
					"fieldName": "indexCountID",
					"fieldValue": map[string]interface{}{
						"type":  "STRING_LIST",
						"value": []string{objType},
					},
					"comparator": "IN",
				},
				"recordType": "HyperionIndexCountLookup",
			},
			"zoneWide": true,
			"zoneID": map[string]string{
				"zoneName": "PrimarySync",
			},
		}},
	}
}

func (pa *PhotoAlbum) listQueryGen(offset int, listType, direction string,
	queryFilter []map[string]interface{}) interface{} {

	filterBy := []map[string]interface{}{
		{
			"fieldName":  "startRank",
			"fieldValue": map[string]interface{}{"type": "INT64", "value": offset},
			"comparator": "EQUALS",
		},
		{
			"fieldName":  "direction",
			"fieldValue": map[string]string{"type": "STRING", "value": direction},
			"comparator": "EQUALS",
		},
	}
	for _, item := range queryFilter {
		filterBy = append(filterBy, item)
	}

	return map[string]interface{}{
		"resultsLimit": pa.pageSize * 2,
		"desiredKeys":  model.ICloudDesiredKeys,
		"zoneID":       map[string]string{"zoneName": "PrimarySync"},
		"query": map[string]interface{}{
			"filterBy":   filterBy,
			"recordType": listType,
		},
	}
}

type ICloudPhotoInfo struct {
	url     string
	service *ICloudService

	ID       string
	FileName string
	Height   int // TODO
	Wight    int // TODO
	Size     int
	Type     string // TODO enum?

	Created   time.Time
	AssetDate time.Time
	AddDate   time.Time
}

func newICloudPhotoInfo(service *ICloudService, info, asset model.ICloudPhotoInfo) *ICloudPhotoInfo {
	id := info.RecordName
	nameBytes, _ := base64.StdEncoding.DecodeString(info.Fields.FilenameEnc.Value)
	name := string(nameBytes)
	size := info.Fields.ResOriginalRes.Value.Size
	assertDate := time.Unix(asset.Fields.AssetDate.Value, 0)
	created := assertDate
	addDate := time.Unix(asset.Fields.AddDate.Value, 0)
	width := info.Fields.ResOriginalWidth.Value
	height := info.Fields.ResOriginalHeight.Value
	url := info.Fields.ResOriginalRes.Value.DownloadURL

	var imageType = map[string]string{
		"public.heic":               "image",
		"public.jpeg":               "image",
		"public.png":                "image",
		"com.apple.quicktime-movie": "movie",
	}

	iType, ok := imageType[info.Fields.ItemType.Value]
	if !ok {
		lower := strings.ToLower(name)
		if strings.HasSuffix(lower, ".heic") ||
			strings.HasSuffix(lower, ".png") ||
			strings.HasSuffix(lower, ".jpg") ||
			strings.HasSuffix(lower, ".jpeg") {
			iType = "image"
		} else {
			iType = "movie"
		}
	}
	return &ICloudPhotoInfo{
		ID:        id,
		FileName:  name,
		Height:    height,
		Wight:     width,
		Size:      size,
		Type:      iType,
		url:       url,
		Created:   created,
		AssetDate: assertDate,
		AddDate:   addDate,
		service:   service,
	}
}

func (pi *ICloudPhotoInfo) Name() string {
	return pi.FileName
}

func (pi *ICloudPhotoInfo) Reader(ctx context.Context) (io.ReadCloser, error) {
	return pi.service.client.Stream(ctx, http.GET, pi.url, nil, nil)
}

func (pi *ICloudPhotoInfo) WriteTo(ctx context.Context, file string) error {
	reader, err := pi.Reader(ctx)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

func (pi *ICloudPhotoInfo) Delete(ctx context.Context) error {
	// TODO not support
	return nil
}
