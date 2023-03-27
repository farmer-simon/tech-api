package aliyun

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vod"
	"goskeleton/app/global/consts"
)

type AliVod struct {
	Client *vod.Client
}

func CreateAliVodService() *AliVod {
	return &AliVod{Client: initVodClient()}
}

const (
	RegionId string = "cn-shanghai"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */

func initVodClient() *vod.Client {
	// 创建授权对象
	credential := &credentials.AccessKeyCredential{
		AccessKeyId:     consts.AliYunAccessKeyId,
		AccessKeySecret: consts.AliYunAccessKeySecret,
	}

	// 自定义config
	config := sdk.NewConfig()
	config.AutoRetry = true     // 失败是否自动重试
	config.MaxRetryTime = 3     // 最大重试次数
	config.Timeout = 3000000000 // 连接超时，单位：纳秒；默认为3秒

	// 创建vodClient实例
	client, err := vod.NewClientWithOptions(RegionId, config, credential)
	if err != nil {
		return nil
	}
	return client
}

// CreateUploadVideo 上传文件，获取凭证
func (serv *AliVod) CreateUploadVideo(title, desc, fileName string) (response *vod.CreateUploadVideoResponse, err error) {
	if serv.Client == nil {
		return nil, errors.New("初始化失败")
	}
	request := vod.CreateCreateUploadVideoRequest()
	request.Title = title
	request.Description = desc
	request.FileName = fileName
	request.AcceptFormat = "JSON"
	return serv.Client.CreateUploadVideo(request)
}

//RefreshUploadVideo 刷新上传凭证
func (serv *AliVod) RefreshUploadVideo(videoId string) (response *vod.RefreshUploadVideoResponse, err error) {
	if serv.Client == nil {
		return nil, errors.New("初始化失败")
	}
	request := vod.CreateRefreshUploadVideoRequest()
	request.VideoId = videoId
	request.AcceptFormat = "JSON"

	return serv.Client.RefreshUploadVideo(request)
}

// GetPlayAuth 获取播放凭证
func (serv *AliVod) GetPlayAuth(videoId string) (response *vod.GetVideoPlayAuthResponse, err error) {
	if serv.Client == nil {
		return nil, errors.New("初始化失败")
	}
	request := vod.CreateGetVideoPlayAuthRequest()
	request.VideoId = videoId
	request.AcceptFormat = "JSON"

	return serv.Client.GetVideoPlayAuth(request)
}
