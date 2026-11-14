package burst

import (
	"math"
	"time"
)

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

const (
	rttFailed      = time.Duration(math.MaxInt64 - iota)
	rttUntested    // nolint: varcheck
	rttUnqualified // nolint: varcheck
)
