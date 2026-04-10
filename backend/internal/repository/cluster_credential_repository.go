package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// ClusterCredential stores encrypted access material for a cluster.
type ClusterCredential struct {
	ID                   uint64 `gorm:"primaryKey"`
	ClusterID            uint64 `gorm:"uniqueIndex;not null"`
	AuthType             string `gorm:"size:32;not null;default:kubeconfig"`
	KubeConfigCiphertext string `gorm:"type:text;not null"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ClusterCredentialRepository struct {
	db *gorm.DB
}

func NewClusterCredentialRepository(db *gorm.DB) *ClusterCredentialRepository {
	return &ClusterCredentialRepository{db: db}
}

func (r *ClusterCredentialRepository) UpsertByClusterID(ctx context.Context, cred *ClusterCredential) error {
	return r.db.WithContext(ctx).Where("cluster_id = ?", cred.ClusterID).Assign(cred).FirstOrCreate(cred).Error
}

func (r *ClusterCredentialRepository) GetByClusterID(ctx context.Context, clusterID uint64) (*ClusterCredential, error) {
	var cred ClusterCredential
	if err := r.db.WithContext(ctx).Where("cluster_id = ?", clusterID).First(&cred).Error; err != nil {
		return nil, err
	}
	return &cred, nil
}
