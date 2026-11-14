package scenarios

import (
	core "github.com/gzjjjfree/v5-result"
	"github.com/gzjjjfree/v5-result/app/instman"
	"github.com/gzjjjfree/v5-result/common"
	"github.com/gzjjjfree/v5-result/common/serial"
	"github.com/gzjjjfree/v5-result/features/extension"
)

func NewInstanceManagerInstanceConfig() *core.Config {
	config := &core.Config{}
	config.App = append(config.App, serial.ToTypedMessage(&instman.Config{}))
	return config
}

func NewInstanceManagerCoreInstance() (*core.Instance, extension.InstanceManagement) {
	coreConfig := NewInstanceManagerInstanceConfig()
	instance, err := core.New(coreConfig)
	if err != nil {
		panic(err)
	}
	common.Must(instance.Start())
	instanceMgr := instance.GetFeature(extension.InstanceManagementType())
	InstanceMgrIfce := instanceMgr.(extension.InstanceManagement)
	return instance, InstanceMgrIfce
}
