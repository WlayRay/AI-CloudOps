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

package dao

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/AI-CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TreeCloudDAO interface {
	CreateCloudAccount(ctx context.Context, account *model.CloudAccount) error
	UpdateCloudAccount(ctx context.Context, account *model.CloudAccount) error
	DeleteCloudAccount(ctx context.Context, id int) error
	GetCloudAccount(ctx context.Context, id int) (*model.CloudAccount, error)
	ListCloudAccounts(ctx context.Context, req *model.ListCloudAccountsReq) (model.ListResp[model.CloudAccount], error)
	GetCloudAccountByProvider(ctx context.Context, provider model.CloudProvider) ([]*model.CloudAccount, error)
	GetEnabledCloudAccounts(ctx context.Context) ([]*model.CloudAccount, error)
	BatchGetCloudAccounts(ctx context.Context, accountIds []int) ([]*model.CloudAccount, error)
	CreateSyncStatus(ctx context.Context, status *model.CloudAccountSyncStatus) error
	UpdateSyncStatus(ctx context.Context, id int, status *model.CloudAccountSyncStatus) error
	GetSyncStatus(ctx context.Context, accountId int, resourceType, region string) (*model.CloudAccountSyncStatus, error)
	ListSyncStatus(ctx context.Context, accountId int) ([]*model.CloudAccountSyncStatus, error)
	DeleteSyncStatus(ctx context.Context, accountId int) error
}

type treeCloudDAO struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewTreeCloudDAO(logger *zap.Logger, db *gorm.DB) TreeCloudDAO {
	return &treeCloudDAO{
		logger: logger,
		db:     db,
	}
}

// CreateCloudAccount 创建云账户
func (d *treeCloudDAO) CreateCloudAccount(ctx context.Context, account *model.CloudAccount) error {
	if account == nil {
		d.logger.Error("创建云账户失败：账户信息为空")
		return fmt.Errorf("云账户信息不能为空")
	}

	d.logger.Info("开始创建云账户", zap.String("name", account.Name), zap.String("provider", string(account.Provider)))

	// 检查账户名称是否已存在
	var count int64
	if err := d.db.WithContext(ctx).Model(&model.CloudAccount{}).
		Where("name = ?", account.Name).Count(&count).Error; err != nil {
		d.logger.Error("检查账户名称失败", zap.String("name", account.Name), zap.Error(err))
		return fmt.Errorf("检查账户名称失败: %w", err)
	}
	if count > 0 {
		d.logger.Warn("账户名称已存在", zap.String("name", account.Name))
		return fmt.Errorf("账户名称已存在: %s", account.Name)
	}

	// 检查AccessKey是否已存在
	if err := d.db.WithContext(ctx).Model(&model.CloudAccount{}).
		Where("access_key = ?", account.AccessKey).Count(&count).Error; err != nil {
		d.logger.Error("检查AccessKey失败", zap.String("accessKey", account.AccessKey), zap.Error(err))
		return fmt.Errorf("检查AccessKey失败: %w", err)
	}
	if count > 0 {
		d.logger.Warn("AccessKey已存在", zap.String("accessKey", account.AccessKey))
		return fmt.Errorf("AccessKey已存在: %s", account.AccessKey)
	}

	// 创建账户
	if err := d.db.WithContext(ctx).Create(account).Error; err != nil {
		d.logger.Error("创建云账户失败", zap.String("name", account.Name), zap.Error(err))
		return fmt.Errorf("创建云账户失败: %w", err)
	}

	d.logger.Info("云账户创建成功", zap.Int("id", account.ID), zap.String("name", account.Name))
	return nil
}

// UpdateCloudAccount 更新云账户
func (d *treeCloudDAO) UpdateCloudAccount(ctx context.Context, account *model.CloudAccount) error {
	if account == nil {
		d.logger.Error("更新云账户失败：账户信息为空")
		return fmt.Errorf("云账户信息不能为空")
	}

	// 检查账户是否存在
	var existingAccount model.CloudAccount
	if err := d.db.WithContext(ctx).Where("id = ?", account.ID).First(&existingAccount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			d.logger.Warn("云账户不存在", zap.Int("id", int(account.ID)))
			return fmt.Errorf("云账户不存在: %d", account.ID)
		}
		d.logger.Error("查询云账户失败", zap.Int("id", int(account.ID)), zap.Error(err))
		return fmt.Errorf("查询云账户失败: %w", err)
	}

	// 检查名称唯一性（排除当前账户）
	if account.Name != "" && account.Name != existingAccount.Name {
		var count int64
		if err := d.db.WithContext(ctx).Model(&model.CloudAccount{}).
			Where("name = ? AND id != ?", account.Name, account.ID).Count(&count).Error; err != nil {
			d.logger.Error("检查账户名称失败", zap.String("name", account.Name), zap.Int("id", int(account.ID)), zap.Error(err))
			return fmt.Errorf("检查账户名称失败: %w", err)
		}
		if count > 0 {
			d.logger.Warn("账户名称已存在", zap.String("name", account.Name), zap.Int("id", int(account.ID)))
			return fmt.Errorf("账户名称已存在: %s", account.Name)
		}
	}

	// 检查AccessKey唯一性（排除当前账户）
	if account.AccessKey != "" && account.AccessKey != existingAccount.AccessKey {
		var count int64
		if err := d.db.WithContext(ctx).Model(&model.CloudAccount{}).
			Where("access_key = ? AND id != ?", account.AccessKey, account.ID).Count(&count).Error; err != nil {
			d.logger.Error("检查AccessKey失败", zap.String("accessKey", account.AccessKey), zap.Int("id", int(account.ID)), zap.Error(err))
			return fmt.Errorf("检查AccessKey失败: %w", err)
		}
		if count > 0 {
			d.logger.Warn("AccessKey已存在", zap.String("accessKey", account.AccessKey), zap.Int("id", int(account.ID)))
			return fmt.Errorf("AccessKey已存在: %s", account.AccessKey)
		}
	}

	// 更新账户
	if err := d.db.WithContext(ctx).Model(&existingAccount).Updates(account).Error; err != nil {
		d.logger.Error("更新云账户失败", zap.Int("id", int(account.ID)), zap.Error(err))
		return fmt.Errorf("更新云账户失败: %w", err)
	}

	d.logger.Info("云账户更新成功", zap.Int("id", int(account.ID)))
	return nil
}

// DeleteCloudAccount 删除云账户
func (d *treeCloudDAO) DeleteCloudAccount(ctx context.Context, id int) error {
	d.logger.Info("开始删除云账户", zap.Int("id", id))

	// 检查账户是否存在
	var account model.CloudAccount
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			d.logger.Warn("云账户不存在", zap.Int("id", id))
			return fmt.Errorf("云账户不存在: %d", id)
		}
		d.logger.Error("查询云账户失败", zap.Int("id", id), zap.Error(err))
		return fmt.Errorf("查询云账户失败: %w", err)
	}

	// 开启事务
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除同步状态
		if err := tx.Where("account_id = ?", id).Delete(&model.CloudAccountSyncStatus{}).Error; err != nil {
			d.logger.Error("删除同步状态失败", zap.Int("accountId", id), zap.Error(err))
			return fmt.Errorf("删除同步状态失败: %w", err)
		}

		// 删除云账户
		if err := tx.Delete(&account).Error; err != nil {
			d.logger.Error("删除云账户失败", zap.Int("id", id), zap.Error(err))
			return fmt.Errorf("删除云账户失败: %w", err)
		}

		d.logger.Info("云账户删除成功", zap.Int("id", id), zap.String("name", account.Name))
		return nil
	})
}

// GetCloudAccount 获取云账户详情
func (d *treeCloudDAO) GetCloudAccount(ctx context.Context, id int) (*model.CloudAccount, error) {
	d.logger.Debug("查询云账户详情", zap.Int("id", id))

	var account model.CloudAccount
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			d.logger.Warn("云账户不存在", zap.Int("id", id))
			return nil, fmt.Errorf("云账户不存在: %d", id)
		}
		d.logger.Error("查询云账户失败", zap.Int("id", id), zap.Error(err))
		return nil, fmt.Errorf("查询云账户失败: %w", err)
	}

	d.logger.Debug("云账户查询成功", zap.Int("id", id), zap.String("name", account.Name))
	return &account, nil
}

// ListCloudAccounts 分页查询云账户列表
func (d *treeCloudDAO) ListCloudAccounts(ctx context.Context, req *model.ListCloudAccountsReq) (model.ListResp[model.CloudAccount], error) {
	d.logger.Debug("查询云账户列表",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.Page),
		zap.String("name", req.Search),
		zap.String("provider", string(req.Provider)),
		zap.Bool("enabled", req.Enabled))

	var accounts []model.CloudAccount
	var total int64

	query := d.db.WithContext(ctx).Model(&model.CloudAccount{})

	// 添加查询条件
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}
	if req.Provider != "" {
		query = query.Where("provider = ?", req.Provider)
	}
	if req.Enabled {
		query = query.Where("is_enabled = ?", true)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		d.logger.Error("查询总数失败", zap.Error(err))
		return model.ListResp[model.CloudAccount]{}, fmt.Errorf("查询总数失败: %w", err)
	}

	// 分页查询
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.Size
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&accounts).Error; err != nil {
		d.logger.Error("查询云账户列表失败", zap.Error(err))
		return model.ListResp[model.CloudAccount]{}, fmt.Errorf("查询云账户列表失败: %w", err)
	}

	d.logger.Debug("云账户列表查询成功", zap.Int64("total", total), zap.Int("count", len(accounts)))
	return model.ListResp[model.CloudAccount]{
		Items: accounts,
		Total: total,
	}, nil
}

// GetCloudAccountByProvider 按云厂商查询账户
func (d *treeCloudDAO) GetCloudAccountByProvider(ctx context.Context, provider model.CloudProvider) ([]*model.CloudAccount, error) {
	d.logger.Debug("按云厂商查询账户", zap.String("provider", string(provider)))

	var accounts []*model.CloudAccount
	if err := d.db.WithContext(ctx).Where("provider = ? AND is_enabled = ?", provider, true).Find(&accounts).Error; err != nil {
		d.logger.Error("按云厂商查询账户失败", zap.String("provider", string(provider)), zap.Error(err))
		return nil, fmt.Errorf("查询云账户失败: %w", err)
	}

	d.logger.Debug("按云厂商查询账户成功", zap.String("provider", string(provider)), zap.Int("count", len(accounts)))
	return accounts, nil
}

// GetEnabledCloudAccounts 获取所有启用的云账户
func (d *treeCloudDAO) GetEnabledCloudAccounts(ctx context.Context) ([]*model.CloudAccount, error) {
	d.logger.Debug("查询启用的云账户")

	var accounts []*model.CloudAccount
	if err := d.db.WithContext(ctx).Where("is_enabled = ?", true).Find(&accounts).Error; err != nil {
		d.logger.Error("查询启用的云账户失败", zap.Error(err))
		return nil, fmt.Errorf("查询启用的云账户失败: %w", err)
	}

	d.logger.Debug("查询启用的云账户成功", zap.Int("count", len(accounts)))
	return accounts, nil
}

// ==================== 同步状态管理 ====================

// CreateSyncStatus 创建同步状态
func (d *treeCloudDAO) CreateSyncStatus(ctx context.Context, status *model.CloudAccountSyncStatus) error {
	if status == nil {
		d.logger.Error("创建同步状态失败：状态信息为空")
		return fmt.Errorf("同步状态信息不能为空")
	}

	d.logger.Info("开始创建同步状态",
		zap.Int("accountId", status.AccountId),
		zap.String("resourceType", status.ResourceType),
		zap.String("region", status.Region))

	// 使用UPSERT逻辑处理重复记录
	result := d.db.WithContext(ctx).
		Where("account_id = ? AND resource_type = ? AND region = ?",
			status.AccountId, status.ResourceType, status.Region).
		Assign(status).
		FirstOrCreate(status)

	if result.Error != nil {
		d.logger.Error("创建或更新同步状态失败", zap.Error(result.Error))
		return fmt.Errorf("创建或更新同步状态失败: %w", result.Error)
	}

	d.logger.Info("同步状态创建成功", zap.Int("id", int(status.ID)))
	return nil
}

// UpdateSyncStatus 更新同步状态
func (d *treeCloudDAO) UpdateSyncStatus(ctx context.Context, id int, status *model.CloudAccountSyncStatus) error {
	if status == nil {
		d.logger.Error("更新同步状态失败：状态信息为空", zap.Int("id", id))
		return fmt.Errorf("同步状态信息不能为空")
	}

	d.logger.Debug("开始更新同步状态", zap.Int("id", id))

	result := d.db.WithContext(ctx).Model(&model.CloudAccountSyncStatus{}).
		Where("id = ?", id).Updates(status)

	if result.Error != nil {
		d.logger.Error("更新同步状态失败", zap.Int("id", id), zap.Error(result.Error))
		return fmt.Errorf("更新同步状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		d.logger.Warn("同步状态记录不存在", zap.Int("id", id))
		return fmt.Errorf("同步状态记录不存在: %d", id)
	}

	d.logger.Debug("同步状态更新成功", zap.Int("id", id))
	return nil
}

// GetSyncStatus 获取同步状态
func (d *treeCloudDAO) GetSyncStatus(ctx context.Context, accountId int, resourceType, region string) (*model.CloudAccountSyncStatus, error) {
	d.logger.Debug("查询同步状态",
		zap.Int("accountId", accountId),
		zap.String("resourceType", resourceType),
		zap.String("region", region))

	var status model.CloudAccountSyncStatus
	if err := d.db.WithContext(ctx).Where("account_id = ? AND resource_type = ? AND region = ?",
		accountId, resourceType, region).First(&status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			d.logger.Debug("同步状态记录不存在",
				zap.Int("accountId", accountId),
				zap.String("resourceType", resourceType),
				zap.String("region", region))
			return nil, nil
		}
		d.logger.Error("查询同步状态失败", zap.Error(err))
		return nil, fmt.Errorf("查询同步状态失败: %w", err)
	}

	d.logger.Debug("同步状态查询成功", zap.Int("id", int(status.ID)))
	return &status, nil
}

// ListSyncStatus 查询账户的所有同步状态
func (d *treeCloudDAO) ListSyncStatus(ctx context.Context, accountId int) ([]*model.CloudAccountSyncStatus, error) {
	d.logger.Debug("查询账户同步状态列表", zap.Int("accountId", accountId))

	var statuses []*model.CloudAccountSyncStatus
	if err := d.db.WithContext(ctx).Where("account_id = ?", accountId).Order("last_sync_time DESC").Find(&statuses).Error; err != nil {
		d.logger.Error("查询同步状态列表失败", zap.Int("accountId", accountId), zap.Error(err))
		return nil, fmt.Errorf("查询同步状态列表失败: %w", err)
	}

	d.logger.Debug("同步状态列表查询成功", zap.Int("accountId", accountId), zap.Int("count", len(statuses)))
	return statuses, nil
}

// DeleteSyncStatus 删除账户的所有同步状态
func (d *treeCloudDAO) DeleteSyncStatus(ctx context.Context, accountId int) error {
	d.logger.Info("删除账户同步状态", zap.Int("accountId", accountId))

	if err := d.db.WithContext(ctx).Where("account_id = ?", accountId).Delete(&model.CloudAccountSyncStatus{}).Error; err != nil {
		d.logger.Error("删除同步状态失败", zap.Int("accountId", accountId), zap.Error(err))
		return fmt.Errorf("删除同步状态失败: %w", err)
	}

	d.logger.Info("账户同步状态删除成功", zap.Int("accountId", accountId))
	return nil
}

// BatchGetCloudAccounts 批量获取云账户
func (d *treeCloudDAO) BatchGetCloudAccounts(ctx context.Context, accountIds []int) ([]*model.CloudAccount, error) {
	d.logger.Debug("批量查询云账户", zap.Ints("accountIds", accountIds))

	if len(accountIds) == 0 {
		return []*model.CloudAccount{}, nil
	}

	var accounts []*model.CloudAccount
	if err := d.db.WithContext(ctx).Where("id IN ?", accountIds).Find(&accounts).Error; err != nil {
		d.logger.Error("批量查询云账户失败", zap.Ints("accountIds", accountIds), zap.Error(err))
		return nil, fmt.Errorf("批量查询云账户失败: %w", err)
	}

	d.logger.Debug("批量查询云账户成功", zap.Int("count", len(accounts)))
	return accounts, nil
}
