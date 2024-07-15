package service

import "golang.org/x/sync/singleflight"

var (
	singleGroup singleflight.Group
)

const (
	tenantCachePool    = "tenants"
	bizCachePool       = "bizs"
	deviceCachePool    = "devices"
	bizMemberCachePool = "biz_members"

	deviceCacheKeyFormat     = "%s:%s:%s:%s"
	bizCacheKeyFormat        = "%s:%s:%s"
	bizMembersCacheKeyFormat = "%s:%s:%s"
	tenantCacheKeyFormat     = "%s:%s"
)
