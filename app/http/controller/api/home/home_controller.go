package home

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/data_type"
	"goskeleton/app/model"
	"goskeleton/app/utils/response"
	"regexp"
	"strings"
)

type Home struct {
}

// Index 首页
func (u *Home) Index(ctx *gin.Context) {

	//hotProject := project.CreateProjectServiceFactory().GetRecommend("hot", 4)
	//newProject := project.CreateProjectServiceFactory().GetRecommend("new", 8)
	//mustProject := project.CreateProjectServiceFactory().GetRecommend("must", 8)
	//adv := adverts.CreateAdvServiceFactory().GetOnlineAdv(1, 5)
	//
	servicesRecommendCategory := model.CreateCateModelFactory("").GetRecommendList(1, 7)
	needsRecommendCategory := model.CreateCateModelFactory("").GetRecommendList(2, 7)
	var serviceQuery data_type.ServicesQuery
	serviceQuery.QueryType = "home"
	serviceQuery.Recommend = 1
	_, services := model.CreateServicesFactory("").Index(&serviceQuery, 1, 8)
	var ids = make([]int, 0)
	for _, item := range services {
		ids = append(ids, int(item.Id))
	}
	cover := model.CreateAttrsFactory("").GetCoverByTargetIds(ids, "services")
	var servicesList = make([]gin.H, 0)
	for _, item := range services {
		servicesList = append(servicesList, gin.H{
			"id":        item.Id,
			"title":     item.Title,
			"price":     item.Price,
			"hits":      item.Hits,
			"nick_name": item.NickName,
			"avatar":    item.Avatar,
			"cate_name": item.CateName,
			"cover":     cover[int(item.Id)],
		})
	}

	var needsQuery data_type.NeedsQuery
	needsQuery.QueryType = "home"
	needsQuery.Recommend = 1
	_, needs := model.CreateNeedsFactory("").Index(&needsQuery, 0, 6)
	var needsList = make([]gin.H, 0)
	for _, item := range needs {
		needsList = append(needsList, gin.H{
			"id":          item.Id,
			"title":       item.Title,
			"content":     u.trimContent(item.Content, 100),
			"price":       item.Price,
			"expire_time": item.ExpireTime,
			"hits":        item.Hits,
			"nick_name":   item.NickName,
			"avatar":      item.Avatar,
			"cate_name":   item.CateName,
		})
	}
	settings := model.CreateSettingsFactory("").GetSettings()
	response.Success(ctx, "Success", gin.H{
		"services": gin.H{
			"keywords":           strings.Split(settings["services_recommend_keywords"].(string), ","),
			"recommend_category": servicesRecommendCategory,
			"hot_list":           servicesList,
		},
		"needs": gin.H{
			"keywords":           strings.Split(settings["needs_recommend_keywords"].(string), ","),
			"recommend_category": needsRecommendCategory,
			"hot_list":           needsList,
		},
	})
}

func (u *Home) Category(ctx *gin.Context) {
	query, _ := ctx.GetQuery("type")
	queryType := utils.String2Int(query)

	_, list := model.CreateCateModelFactory("").GetList(queryType, 2, 0, 0)
	response.Success(ctx, "Success", gin.H{
		"list": list,
	})
}

func (u *Home) trimContent(src string, length int) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	c := strings.TrimSpace(src)
	if len([]rune(c)) < length {
		return c
	}
	return utils.Substring(c, 0, length)
}
