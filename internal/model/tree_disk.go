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

package model

import (
	"time"
)

// ResourceDisk 磁盘资源模型
type ResourceDisk struct {
	Model

	// 基本信息
	DiskID             string     `json:"disk_id" gorm:"column:disk_id;size:64;index;comment:磁盘ID"`
	DiskName           string     `json:"disk_name" gorm:"column:disk_name;size:128;comment:磁盘名称"`
	Provider           CloudProvider `json:"cloud_provider" gorm:"type:varchar(50);comment:云厂商"`
	RegionId           string     `json:"region_id" gorm:"type:varchar(50);comment:地区，如cn-hangzhou"`
	ZoneId             string     `json:"zone_id" gorm:"type:varchar(100);comment:可用区ID"`
	Status             string     `json:"status" gorm:"type:varchar(50);comment:资源状态"`
	
	// 磁盘规格
	Size               int        `json:"size" gorm:"column:size;comment:磁盘大小(GB)"`
	Category           string     `json:"category" gorm:"column:category;size:32;comment:磁盘类型(cloud_ssd/cloud_essd等)"`
	DiskType           string     `json:"disk_type" gorm:"column:disk_type;size:32;comment:磁盘用途(system/data)"`
	PerformanceLevel   string     `json:"performance_level" gorm:"column:performance_level;size:32;comment:性能等级(PL0/PL1/PL2/PL3)"`
	IOPS               int        `json:"iops" gorm:"column:iops;comment:每秒读写次数"`
	Throughput         int        `json:"throughput" gorm:"column:throughput;comment:吞吐量(MB/s)"`
	
	// 挂载信息
	InstanceID         string     `json:"instance_id" gorm:"column:instance_id;size:64;index;comment:挂载的实例ID"`
	InstanceName       string     `json:"instance_name" gorm:"type:varchar(100);comment:挂载的实例名称"`
	Device             string     `json:"device" gorm:"column:device;size:32;comment:设备名称(/dev/vdb等)"`
	AttachTime         *time.Time `json:"attach_time" gorm:"column:attach_time;comment:挂载时间"`
	DetachTime         *time.Time `json:"detach_time" gorm:"column:detach_time;comment:卸载时间"`
	DeleteWithInstance bool       `json:"delete_with_instance" gorm:"column:delete_with_instance;comment:随实例删除"`
	Portable           bool       `json:"portable" gorm:"column:portable;comment:是否可卸载"`
	
	// 安全与网络
	VpcId              string     `json:"vpc_id" gorm:"type:varchar(100);comment:VPC ID"`
	Encrypted          bool       `json:"encrypted" gorm:"column:encrypted;comment:是否加密"`
	SecurityGroupIds   StringList `json:"security_group_ids" gorm:"type:varchar(500);comment:安全组ID列表"`
	
	// 计费与创建信息
	InstanceChargeType string     `json:"instance_charge_type" gorm:"type:varchar(50);comment:付费类型(PrePaid/PostPaid)"`
	CreationTime       string     `json:"creation_time" gorm:"type:varchar(30);comment:创建时间,ISO8601格式"`
	SourceSnapshotId   string     `json:"source_snapshot_id" gorm:"column:source_snapshot_id;size:64;comment:源快照ID"`
	ImageId            string     `json:"image_id" gorm:"column:image_id;size:64;comment:镜像ID"`
	
	// 组织与标签
	ResourceGroupId    string     `json:"resource_group_id" gorm:"column:resource_group_id;size:64;comment:资源组ID"`
	Env                string     `json:"environment" gorm:"type:varchar(50);comment:环境标识,如dev,prod"`
	Description        string     `json:"description" gorm:"type:text;comment:资源描述"`
	Tags               StringList `json:"tags" gorm:"type:varchar(500);comment:资源标签集合"`
	
	// 内部管理字段
	CreateByOrder      bool       `json:"create_by_order" gorm:"comment:是否由工单创建"`
	LastSyncTime       time.Time  `json:"last_sync_time" gorm:"comment:最后同步时间"`
	TreeNodeID         int        `json:"tree_node_id" gorm:"comment:关联的服务树节点ID"`
}

// CreateDiskReq 创建磁盘请求
type CreateDiskReq struct {
	RegionId         string            `json:"region_id" validate:"required" binding:"required"`
	ZoneId           string            `json:"zone_id" validate:"required" binding:"required"`
	DiskName         string            `json:"disk_name" validate:"required" binding:"required"`
	DiskCategory     string            `json:"disk_category" validate:"required" binding:"required"`
	Size             int               `json:"size" validate:"required,min=20" binding:"required,min=20"`
	Description      string            `json:"description"`
	DiskType         string            `json:"disk_type"` // system, data
	Encrypted        bool              `json:"encrypted"`
	PerformanceLevel string            `json:"performance_level"`
	SourceSnapshotId string            `json:"source_snapshot_id"`
	ImageId          string            `json:"image_id"`
	ResourceGroupId  string            `json:"resource_group_id"`
	Tags             map[string]string `json:"tags"`
	// 计费信息
	InstanceChargeType string `json:"instance_charge_type"` // PrePaid, PostPaid
	// 实例相关（如果需要创建后立即挂载）
	InstanceId         string `json:"instance_id"`
	Device             string `json:"device"`
	DeleteWithInstance bool   `json:"delete_with_instance"`
	// 内部管理字段
	TreeNodeID         int    `json:"tree_node_id"`
}

// UpdateDiskReq 更新磁盘请求
type UpdateDiskReq struct {
	ID          uint              `json:"id"`
	DiskName    *string           `json:"disk_name"`
	Description *string           `json:"description"`
	Size        *int              `json:"size" validate:"omitempty,min=20"`
	Tags        map[string]string `json:"tags"`
	Env         *string           `json:"environment"`
}

// AttachDiskReq 挂载磁盘请求
type AttachDiskReq struct {
	DiskId             string `json:"disk_id" validate:"required" binding:"required"`
	InstanceId         string `json:"instance_id" validate:"required" binding:"required"`
	Device             string `json:"device"`
	DeleteWithInstance bool   `json:"delete_with_instance"`
}

// DetachDiskReq 卸载磁盘请求
type DetachDiskReq struct {
	DiskId     string `json:"disk_id" validate:"required" binding:"required"`
	InstanceId string `json:"instance_id" validate:"required" binding:"required"`
}

// ResizeDiskReq 磁盘扩容请求
type ResizeDiskReq struct {
	ID          uint   `json:"id"`
	DiskId      string `json:"disk_id" validate:"required" binding:"required"`
	NewSize     int    `json:"new_size" validate:"required,min=20" binding:"required,min=20"`
	Description string `json:"description"`
}

// DiskListReq 磁盘列表请求
type DiskListReq struct {
	ListReq
	Provider   CloudProvider `json:"provider" form:"provider"`
	RegionId   string        `json:"region_id" form:"region_id"`
	ZoneId     string        `json:"zone_id" form:"zone_id"`
	Status     string        `json:"status" form:"status"`
	Category   string        `json:"category" form:"category"`
	DiskType   string        `json:"disk_type" form:"disk_type"`
	InstanceId string        `json:"instance_id" form:"instance_id"`
	DiskIds    []string      `json:"disk_ids" form:"disk_ids"`
	Tags       []string      `json:"tags" form:"tags"`
	Encrypted  *bool         `json:"encrypted" form:"encrypted"`
	Portable   *bool         `json:"portable" form:"portable"`
	TreeNodeID *int          `json:"tree_node_id" form:"tree_node_id"`
}

// DiskResp 磁盘响应
type DiskResp struct {
	ResourceDisk
	AttachedInstance *ResourceEcs `json:"attached_instance,omitempty"` // 挂载的实例信息
	Snapshots        []string     `json:"snapshots,omitempty"`         // 快照列表
	CanAttach        bool         `json:"can_attach"`                  // 是否可以挂载
	CanDetach        bool         `json:"can_detach"`                  // 是否可以卸载
	CanResize        bool         `json:"can_resize"`                  // 是否可以扩容
}

// DiskStatistics 磁盘统计信息
type DiskStatistics struct {
	TotalCount     int64            `json:"total_count"`
	TotalSize      int64            `json:"total_size"`      // 总容量(GB)
	AttachedCount  int64            `json:"attached_count"`  // 已挂载数量
	AvailableCount int64            `json:"available_count"` // 可用数量
	ByStatus       map[string]int64 `json:"by_status"`       // 按状态分组
	ByCategory     map[string]int64 `json:"by_category"`     // 按类型分组
	BySize         map[string]int64 `json:"by_size"`         // 按大小分组
	ByRegion       map[string]int64 `json:"by_region"`       // 按区域分组
	ByProvider     map[string]int64 `json:"by_provider"`     // 按云厂商分组
}

// DiskOperation 磁盘操作记录
type DiskOperation struct {
	Model
	DiskId     string     `json:"disk_id" gorm:"column:disk_id;size:64;index"`
	Operation  string     `json:"operation" gorm:"column:operation;size:32"` // create, attach, detach, resize, delete
	Status     string     `json:"status" gorm:"column:status;size:32"`       // pending, running, success, failed
	InstanceId string     `json:"instance_id" gorm:"column:instance_id;size:64"`
	Parameters StringList `json:"parameters" gorm:"column:parameters;type:json"`
	ErrorMsg   string     `json:"error_msg" gorm:"column:error_msg;type:text"`
	StartTime  time.Time  `json:"start_time" gorm:"column:start_time"`
	EndTime    *time.Time `json:"end_time" gorm:"column:end_time"`
	Duration   int        `json:"duration" gorm:"column:duration"` // 耗时(秒)
	Operator   string     `json:"operator" gorm:"column:operator;size:64"`
}
