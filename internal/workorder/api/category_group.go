/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package api

import (
	"github.com/GoSimplicity/AI-CloudOps/internal/model"
	"github.com/GoSimplicity/AI-CloudOps/internal/workorder/service"
	"github.com/GoSimplicity/AI-CloudOps/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CategoryGroupHandler struct {
	service service.CategoryGroupService
}

func NewCategoryGroupHandler(service service.CategoryGroupService) *CategoryGroupHandler {
	return &CategoryGroupHandler{
		service: service,
	}
}

func (h *CategoryGroupHandler) RegisterRouters(server *gin.Engine) {
	categoryGroup := server.Group("/api/workorder/category")
	{
		categoryGroup.POST("/create", h.CreateCategory)
		categoryGroup.PUT("/update/:id", h.UpdateCategory)
		categoryGroup.DELETE("/delete/:id", h.DeleteCategory)
		categoryGroup.GET("/list", h.ListCategory)
		categoryGroup.GET("/detail/:id", h.DetailCategory)
		categoryGroup.GET("/statistics", h.GetCategoryStatistics)
	}
}

func (h *CategoryGroupHandler) CreateCategory(ctx *gin.Context) {
	var req model.CreateCategoryReq

	user := ctx.MustGet("user").(utils.UserClaims)

	req.UserID = user.Uid
	req.UserName = user.Username

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.service.CreateCategory(ctx, &req)
	})
}

func (h *CategoryGroupHandler) UpdateCategory(ctx *gin.Context) {
	var req model.UpdateCategoryReq

	id, err := utils.GetParamID(ctx)
	if err != nil {
		return
	}

	req.ID = id

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.service.UpdateCategory(ctx, &req)
	})
}

func (h *CategoryGroupHandler) DeleteCategory(ctx *gin.Context) {
	var req model.DeleteCategoryReq

	id, err := utils.GetParamID(ctx)
	if err != nil {
		return
	}

	req.ID = id

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.service.DeleteCategory(ctx, req.ID)
	})
}

func (h *CategoryGroupHandler) ListCategory(ctx *gin.Context) {
	var req model.ListCategoryReq
	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.service.ListCategory(ctx, req)
	})
}

func (h *CategoryGroupHandler) DetailCategory(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx)
	if err != nil {
		return
	}

	utils.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.service.GetCategory(ctx, id)
	})
}

func (h *CategoryGroupHandler) GetCategoryStatistics(ctx *gin.Context) {
	utils.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.service.GetCategoryStatistics(ctx)
	})
}
