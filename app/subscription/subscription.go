package subscription

import "github.com/gzjjjfree/v5-result/features"

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

type SubscriptionManager interface {
	features.Feature
	AddTrackedSubscriptionFromImportSource(importSource *ImportSource) error
	RemoveTrackedSubscription(name string) error
	ListTrackedSubscriptions() []string
	GetTrackedSubscriptionStatus(name string) (*TrackedSubscriptionStatus, error)
	UpdateTrackedSubscription(name string) error
}

func SubscriptionManagerType() interface{} {
	return (*SubscriptionManager)(nil)
}
