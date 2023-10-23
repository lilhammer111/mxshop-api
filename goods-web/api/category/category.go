package category

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"mxshop-api/goods-web/api/helper"
	"mxshop-api/goods-web/forms"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"net/http"
	"strconv"
)

func List(c *gin.Context) {
	r, err := global.GoodsSrvClient.GetAllCategorysList(context.WithValue(context.Background(), "ginContext", c), &empty.Empty{})
	if err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}

	data := make([]interface{}, 0)
	// Why not return the JsonData directly via c.JSON here, but deserialize it first?
	// Avoiding Double Serialization
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("[List] 查询 【分类列表】失败： ", err.Error())
	}

	c.JSON(http.StatusOK, data)
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	reMap := make(map[string]interface{})
	subCategories := make([]interface{}, 0)
	if r, err := global.GoodsSrvClient.GetSubCategory(context.WithValue(context.Background(), "ginContext", c), &proto.CategoryListRequest{
		Id: int32(i),
	}); err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	} else {
		//写文档 特别是数据多的时候很慢， 先开发后写文档
		for _, value := range r.SubCategorys {
			subCategories = append(subCategories, map[string]interface{}{
				"id":              value.Id,
				"name":            value.Name,
				"level":           value.Level,
				"parent_category": value.ParentCategory,
				"is_tab":          value.IsTab,
			})
		}
		reMap["id"] = r.Info.Id
		reMap["name"] = r.Info.Name
		reMap["level"] = r.Info.Level
		reMap["parent_category"] = r.Info.ParentCategory
		reMap["is_tab"] = r.Info.IsTab
		reMap["sub_categorys"] = subCategories

		c.JSON(http.StatusOK, reMap)
	}
	return
}

func New(c *gin.Context) {
	categoryForm := forms.CategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		helper.HandleValidatorError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategory(context.WithValue(context.Background(), "ginContext", c), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		IsTab:          *categoryForm.IsTab,
		Level:          categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentCategory
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	c.JSON(http.StatusOK, request)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	//1. 先查询出该分类写的所有子分类
	//2. 将所有的分类全部逻辑删除
	//3. 将该分类下的所有的商品逻辑删除
	_, err = global.GoodsSrvClient.DeleteCategory(context.WithValue(context.Background(), "ginContext", c), &proto.DeleteCategoryRequest{Id: int32(i)})
	if err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func Update(c *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		helper.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}
	_, err = global.GoodsSrvClient.UpdateCategory(context.WithValue(context.Background(), "ginContext", c), request)
	if err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}
