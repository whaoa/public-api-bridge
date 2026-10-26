package official

import (
	"strings"

	"github.com/tidwall/gjson"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

func normalizeBadRequestError(payload string) error {
	reason := gjson.Get(payload, "description").String()
	if strings.Contains(reason, " is not valid ID") {
		return errors.NotFound("the subject is not exist")
	}
	return errors.Errorf("can not get subject: %s", reason)
}
