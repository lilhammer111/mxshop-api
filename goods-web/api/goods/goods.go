package goods

import (
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/goods-web/api/helper"
	"mxshop-api/goods-web/forms"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/proto"
	"net/http"
	"strconv"
)

func List(c *gin.Context) {
	request := &proto.GoodsFilterRequest{}

	priceMin := c.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := c.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := c.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := c.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := c.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := c.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := c.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := c.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := c.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := c.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "requests are too frequent, please try again later",
		})
		return
	}

	// revoke goods microservice
	res, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginContext", c), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}

	e.Exit()

	reMap := map[string]any{
		"total": res.Total,
	}

	goodsList := make([]map[string]any, 0)
	for _, value := range res.Data {
		goodsList = append(goodsList, map[string]any{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]any{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]any{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	goodsForm := forms.GoodsForm{}

	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		helper.HandleValidatorError(c, err)
		return
	}

	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.WithValue(context.Background(), "ginContext", c), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}
	// todo goods stock - Distributed transactions 分布式事务 难点

	c.JSON(http.StatusOK, rsp)
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", c), &proto.GoodInfoRequest{Id: int32(idInt)})
	if err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
	}

	// 库存还没做
	c.JSON(http.StatusOK, gin.H{
		"id":          r.Id,
		"name":        r.Name,
		"shop_price":  r.ShopPrice,
		"goods_brief": r.GoodsBrief,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"front_image": r.GoodsFrontImage,
		"desc":        r.GoodsDesc,
		"desc_images": r.DescImages,
		"category":    gin.H{"id": r.Category.Id, "name": r.Category.Name},
		"brand":       gin.H{"id": r.Brand.Id, "name": r.Brand.Name, "logo": r.Brand.Logo},
		"is_hot":      r.IsHot,
		"is_new":      r.IsNew,
		"on_sale":     r.OnSale,
	})

}

func Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.WithValue(context.Background(), "ginContext", c), &proto.DeleteGoodsInfo{Id: int32(idInt)})
	if err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
	return
}

func Stock(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// todo goods stock
	return
}

func UpdateStatus(c *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := c.ShouldBindJSON(&goodsStatusForm); err != nil {
		helper.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", c), &proto.CreateGoodsInfo{
		Id:     int32(idInt),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "succeed to update goods status",
	})
}

func Update(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBindJSON(&goodsForm); err != nil {
		helper.HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.WithValue(context.Background(), "ginContext", c), &proto.CreateGoodsInfo{
		Id:              int32(idInt),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		helper.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "succeed to update goods",
	})
}
