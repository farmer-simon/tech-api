package category

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/data_type"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/utils/data_bind"
	"math"
	"strings"
)

func CreateCateServiceFactory() *CateService {
	return &CateService{model.CreateCateModelFactory("")}
}

type CateService struct {
	Model *model.CateModel
}

// Create 插入分类
func (serv *CateService) Create(ctx *gin.Context) error {
	var data model.CateModel
	if err := data_bind.ShouldBindFormDataToModel(ctx, &data); err == nil {
		//支持一下批量添加
		var names = data.Name
		var errNames = make([]string, 0)
		for _, name := range strings.Split(names, ",") {
			data.Name = name
			if serv.Model.CheckRepeatName(0, int(data.ParentId), data.Name) {
				errNames = append(errNames, fmt.Sprintf("名称：%s 已存在", data.Name))
				continue
			}

			if res := serv.Model.Insert(data); res != nil {
				errNames = append(errNames, fmt.Sprintf("名称：%s 添加失败 %s", data.Name, res.Error()))
				continue
			}
		}
		if len(errNames) > 0 {
			return errors.New(strings.Join(errNames, ";"))
		}
	}
	return nil
}

// Update 更新分类
func (serv *CateService) Update(ctx *gin.Context) error {
	var data model.CateModel
	if err := data_bind.ShouldBindFormDataToModel(ctx, &data); err == nil {
		if serv.Model.CheckRepeatName(int(data.Id), int(data.ParentId), data.Name) {
			return errors.New(fmt.Sprintf("名称：%s 已存在", data.Name))
		}
		if res := serv.Model.Update(&data); res != nil {
			return res
		}
	}
	return nil
}

//Destroy 删除分类
func (serv *CateService) Destroy(id int) error {
	count, _ := serv.Model.GetList(id, 0, 0, 1)
	if count > 0 {
		return errors.New("存在子分类，暂无法删除")
	}
	_ = serv.Model.Destroy(id)
	return nil
}

// Index 后台列表
func (serv *CateService) Index(ctx *gin.Context) (int64, []model.CateModel) {
	parentId := ctx.GetFloat64(consts.ValidatorPrefix + "parent_id")
	state := ctx.GetFloat64(consts.ValidatorPrefix + "state")
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max((page-1)*limit, 0)
	if count, res := serv.Model.GetList(int(parentId), int(state), int(offset), int(limit)); res != nil {
		return count, res
	}
	return 0, nil
}

func (serv *CateService) CategoryTree() (topTree []data_type.CategoryTree) {
	allCategory := serv.Model.GetAllOnlineCategory()
	//var topTree = make([]data_type.CategoryTree, 0)
	var subTreeMap = make(map[int64][]data_type.CategoryItem, 0)
	for _, cate := range allCategory {
		if cate.ParentId == 0 {
			topTree = append(topTree, data_type.CategoryTree{
				CategoryItem: data_type.CategoryItem{
					Id:       cate.Id,
					ParentId: int64(cate.ParentId),
					Name:     cate.Name,
				},
				SubTree: nil,
			})
		} else {
			if _, isOk := subTreeMap[int64(cate.ParentId)]; !isOk {
				subTreeMap[int64(cate.ParentId)] = make([]data_type.CategoryItem, 0)
			}
			subTreeMap[int64(cate.ParentId)] = append(subTreeMap[int64(cate.ParentId)], data_type.CategoryItem{
				Id:       cate.Id,
				ParentId: int64(cate.ParentId),
				Name:     cate.Name,
			})
		}
	}

	for i, item := range topTree {
		if _, isOk := subTreeMap[item.Id]; isOk {
			item.SubTree = subTreeMap[item.Id]
			topTree[i] = item
		}
	}
	return
}
