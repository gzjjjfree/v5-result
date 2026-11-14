package environment

import (
	"github.com/gzjjjfree/v5-result/common/environment/filesystemcap"
	"github.com/gzjjjfree/v5-result/features/extension/storage"
	"github.com/gzjjjfree/v5-result/transport/internet"
	"github.com/gzjjjfree/v5-result/transport/internet/tagged"
)

type BaseEnvironmentCapabilitySet interface {
	FeaturesLookupCapabilitySet
	LogCapabilitySet
}

type BaseEnvironment interface {
	BaseEnvironmentCapabilitySet
	doNotImpl()
}

type SystemNetworkCapabilitySet interface {
	Dialer() internet.SystemDialer
	Listener() internet.SystemListener
}

type InstanceNetworkCapabilitySet interface {
	OutboundDialer() tagged.DialFunc
}

type FeaturesLookupCapabilitySet interface {
	RequireFeatures() interface{}
}

type LogCapabilitySet interface {
	RecordLog() interface{}
}

type FileSystemCapabilitySet interface {
	filesystemcap.FileSystemCapabilitySet
}

type PersistentStorageCapabilitySet interface {
	PersistentStorage() storage.ScopedPersistentStorage
}
type TransientStorageCapabilitySet interface {
	TransientStorage() storage.ScopedTransientStorage
}

type ProxyMetadataCapabilitySet interface {
	SelfProxyTag() string
}
