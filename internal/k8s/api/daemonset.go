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
	"github.com/GoSimplicity/AI-CloudOps/internal/k8s/service/admin"
	"github.com/GoSimplicity/AI-CloudOps/internal/model"
	"github.com/GoSimplicity/AI-CloudOps/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sDaemonSetHandler struct {
	l                *zap.Logger
	daemonSetService admin.DaemonSetService
}

func NewK8sDaemonSetHandler(l *zap.Logger, daemonSetService admin.DaemonSetService) *K8sDaemonSetHandler {
	return &K8sDaemonSetHandler{
		l:                l,
		daemonSetService: daemonSetService,
	}
}

func (k *K8sDaemonSetHandler) RegisterRouters(server *gin.Engine) {
	k8sGroup := server.Group("/api/k8s")

	daemonsets := k8sGroup.Group("/daemonsets")
	{
		daemonsets.GET("/:id", k.GetDaemonSetsByNamespace)          // 根据命名空间获取 DaemonSet 列表
		daemonsets.GET("/:id/yaml", k.GetDaemonSetYaml)            // 获取指定 DaemonSet 的 YAML 配置
		daemonsets.POST("/update", k.UpdateDaemonSet)              // 更新指定 DaemonSet
		daemonsets.POST("/create", k.CreateDaemonSet)              // 创建 DaemonSet
		daemonsets.DELETE("/batch_delete", k.BatchDeleteDaemonSet) // 批量删除 DaemonSet
		daemonsets.DELETE("/delete/:id", k.DeleteDaemonSet)        // 删除指定 DaemonSet
		daemonsets.POST("/restart/:id", k.RestartDaemonSet)        // 重启 DaemonSet
		daemonsets.GET("/:id/status", k.GetDaemonSetStatus)        // 获取 DaemonSet 状态
	}
}

// GetDaemonSetsByNamespace 根据命名空间获取 DaemonSet 列表
func (k *K8sDaemonSetHandler) GetDaemonSetsByNamespace(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.BadRequestError(ctx, err.Error())
		return
	}

	namespace := ctx.Query("namespace")
	if namespace == "" {
		utils.BadRequestError(ctx, "缺少 'namespace' 参数")
		return
	}

	utils.HandleRequest(ctx, nil, func() (interface{}, error) {
		return k.daemonSetService.GetDaemonSetsByNamespace(ctx, id, namespace)
	})
}

// CreateDaemonSet 创建 DaemonSet
func (k *K8sDaemonSetHandler) CreateDaemonSet(ctx *gin.Context) {
	var req model.K8sDaemonSetRequest

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, k.daemonSetService.CreateDaemonSet(ctx, &req)
	})
}

// UpdateDaemonSet 更新指定的 DaemonSet
func (k *K8sDaemonSetHandler) UpdateDaemonSet(ctx *gin.Context) {
	var req model.K8sDaemonSetRequest

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, k.daemonSetService.UpdateDaemonSet(ctx, &req)
	})
}

// BatchDeleteDaemonSet 批量删除 DaemonSet
func (k *K8sDaemonSetHandler) BatchDeleteDaemonSet(ctx *gin.Context) {
	var req model.K8sDaemonSetRequest

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, k.daemonSetService.BatchDeleteDaemonSet(ctx, req.ClusterID, req.Namespace, req.DaemonSetNames)
	})
}

// GetDaemonSetYaml 获取 DaemonSet 的 YAML 配置
func (k *K8sDaemonSetHandler) GetDaemonSetYaml(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.BadRequestError(ctx, err.Error())
		return
	}

	daemonSetName := ctx.Query("daemonset_name")
	if daemonSetName == "" {
		utils.BadRequestError(ctx, "缺少 'daemonset_name' 参数")
		return
	}

	namespace := ctx.Query("namespace")
	if namespace == "" {
		utils.BadRequestError(ctx, "缺少 'namespace' 参数")
		return
	}

	utils.HandleRequest(ctx, nil, func() (interface{}, error) {
		return k.daemonSetService.GetDaemonSetYaml(ctx, id, namespace, daemonSetName)
	})
}

// DeleteDaemonSet 删除指定的 DaemonSet
func (k *K8sDaemonSetHandler) DeleteDaemonSet(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.BadRequestError(ctx, err.Error())
		return
	}

	daemonSetName := ctx.Query("daemonset_name")
	if daemonSetName == "" {
		utils.BadRequestError(ctx, "缺少 'daemonset_name' 参数")
		return
	}

	namespace := ctx.Query("namespace")
	if namespace == "" {
		utils.BadRequestError(ctx, "缺少 'namespace' 参数")
		return
	}

	utils.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, k.daemonSetService.DeleteDaemonSet(ctx, id, namespace, daemonSetName)
	})
}

// RestartDaemonSet 重启 DaemonSet
func (k *K8sDaemonSetHandler) RestartDaemonSet(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.BadRequestError(ctx, err.Error())
		return
	}

	daemonSetName := ctx.Query("daemonset_name")
	if daemonSetName == "" {
		utils.BadRequestError(ctx, "缺少 'daemonset_name' 参数")
		return
	}

	namespace := ctx.Query("namespace")
	if namespace == "" {
		utils.BadRequestError(ctx, "缺少 'namespace' 参数")
		return
	}

	utils.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, k.daemonSetService.RestartDaemonSet(ctx, id, namespace, daemonSetName)
	})
}

// GetDaemonSetStatus 获取 DaemonSet 状态
func (k *K8sDaemonSetHandler) GetDaemonSetStatus(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.BadRequestError(ctx, err.Error())
		return
	}

	daemonSetName := ctx.Query("daemonset_name")
	if daemonSetName == "" {
		utils.BadRequestError(ctx, "缺少 'daemonset_name' 参数")
		return
	}

	namespace := ctx.Query("namespace")
	if namespace == "" {
		utils.BadRequestError(ctx, "缺少 'namespace' 参数")
		return
	}

	utils.HandleRequest(ctx, nil, func() (interface{}, error) {
		return k.daemonSetService.GetDaemonSetStatus(ctx, id, namespace, daemonSetName)
	})
}