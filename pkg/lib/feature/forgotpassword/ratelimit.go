package forgotpassword

import (
	"fmt"

	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/duration"
)

// TODO(rate-limit): allow configuration of bucket size & reset period

func SendResetPasswordCodeRateLimitBucket(loginID string) ratelimit.Bucket {
	return ratelimit.Bucket{
		Key:         fmt.Sprintf("reset-password-send-code:%s", loginID),
		Size:        5,
		ResetPeriod: 5 * duration.PerMinute,
	}
}

func VerifyIPRateLimitBucket(ip string) ratelimit.Bucket {
	return ratelimit.Bucket{
		Key:         fmt.Sprintf("reset-password-verify-ip:%s", ip),
		Size:        10,
		ResetPeriod: duration.PerMinute,
	}
}
