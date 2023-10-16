package shop_cart

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/order-web/api"
	"mxshop-api/order-web/forms"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"
	"net/http"
	"strconv"
)

func List(c *gin.Context) {
	userID, existed := c.Get("userId")
	if !existed {
		c.JSON(http.StatusNotFound, gin.H{"error": "userId is not existed"})
		return
	}
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{Id: int32(userID.(uint))})
	if err != nil {
		zap.S().Errorw("[List] 查询 【购物车列表】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"total": 0})
		return
	}

	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		zap.S().Errorw("[List] 批量查询 【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	/*
		{
			"total":12,
			"data":[
				{
					"id":1,
					"goods_id":421,
					"goods_name":421,
					"goods_price":421,
					"goods_image":421,
					"nums":421,
					"checked":421,
				}
			]
		}
	*/
	goodsList := make([]any, 0)
	for _, item := range rsp.Data {
		for _, g := range goodsRsp.Data {
			if g.Id == item.GoodsId {
				tmpMap := map[string]any{}
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["good_name"] = g.Name
				tmpMap["good_image"] = g.GoodsFrontImage
				tmpMap["good_price"] = g.ShopPrice
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsList = append(goodsList, tmpMap)
			}
		}
	}

	reMap["data"] = goodsList
	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := c.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: itemForm.GoodsId})
	if err != nil {
		zap.S().Errorw("[New] 查询【商品信息】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询【库存信息】失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	// Inventory is insufficient
	if invRsp.Num < itemForm.Nums {
		c.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}

	userId, _ := c.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})

	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	userId, _ := c.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除购物车记录失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}

func Update(c *gin.Context) {
	// o/v1/421
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}
	if err := c.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	userId, _ := c.Get("userId")
	request := proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}

	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("更新购物车记录失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.Status(http.StatusOK)
}
