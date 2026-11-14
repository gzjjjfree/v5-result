package geodata

import (
	"github.com/gzjjjfree/v5-result/app/router/routercommon"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

type LoaderImplementation interface {
	LoadSite(filename, list string) ([]*routercommon.Domain, error)
	LoadIP(filename, country string) ([]*routercommon.CIDR, error)
}

type Loader interface {
	LoaderImplementation
	LoadGeoSite(list string) ([]*routercommon.Domain, error)
	LoadGeoSiteWithAttr(file string, siteWithAttr string) ([]*routercommon.Domain, error)
	LoadGeoIP(country string) ([]*routercommon.CIDR, error)
}
