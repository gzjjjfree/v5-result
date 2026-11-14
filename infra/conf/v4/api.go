package v4

import (
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/gzjjjfree/v5-result/app/commander"
	loggerservice "github.com/gzjjjfree/v5-result/app/log/command"
	observatoryservice "github.com/gzjjjfree/v5-result/app/observatory/command"
	handlerservice "github.com/gzjjjfree/v5-result/app/proxyman/command"
	routerservice "github.com/gzjjjfree/v5-result/app/router/command"
	statsservice "github.com/gzjjjfree/v5-result/app/stats/command"
	"github.com/gzjjjfree/v5-result/common/serial"
)

type APIConfig struct {
	Tag      string   `json:"tag"`
	Services []string `json:"services"`
}

func (c *APIConfig) Build() (*commander.Config, error) {
	if c.Tag == "" {
		return nil, newError("API tag can't be empty.")
	}

	services := make([]*anypb.Any, 0, 16)
	for _, s := range c.Services {
		switch strings.ToLower(s) {
		case "reflectionservice":
			services = append(services, serial.ToTypedMessage(&commander.ReflectionConfig{}))
		case "handlerservice":
			services = append(services, serial.ToTypedMessage(&handlerservice.Config{}))
		case "loggerservice":
			services = append(services, serial.ToTypedMessage(&loggerservice.Config{}))
		case "statsservice":
			services = append(services, serial.ToTypedMessage(&statsservice.Config{}))
		case "observatoryservice":
			services = append(services, serial.ToTypedMessage(&observatoryservice.Config{}))
		case "routingservice":
			services = append(services, serial.ToTypedMessage(&routerservice.Config{}))
		default:
			if !strings.HasPrefix(s, "#") {
				continue
			}
			message, err := desc.LoadMessageDescriptor(s[1:])
			if err != nil || message == nil {
				return nil, newError("Cannot find API", s, "").Base(err)
			}
			serviceConfig := dynamic.NewMessage(message)
			services = append(services, serial.ToTypedMessage(serviceConfig))
		}
	}

	return &commander.Config{
		Tag:     c.Tag,
		Service: services,
	}, nil
}
