package treatment

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/2410as/HairHistory/backend/internal/httpx"
)

const (
	maxServices       = 10
	maxSalonNameRunes = 100
	maxMemoRunes      = 1000
)

func validate(in Input) (validInput, error) {
	details := map[string]string{}

	treatedOn, err := parseTreatedOn(in.TreatedOn)
	if err != nil {
		details["treatedOn"] = err.Error()
	}

	services, err := normalizeServices(in.Services)
	if err != nil {
		details["services"] = err.Error()
	}

	salonName, err := normalizeOptionalText(in.SalonName, maxSalonNameRunes)
	if err != nil {
		details["salonName"] = err.Error()
	}

	memo, err := normalizeOptionalText(in.Memo, maxMemoRunes)
	if err != nil {
		details["memo"] = err.Error()
	}

	if in.Cost != nil && *in.Cost < 0 {
		details["cost"] = "must be 0 or greater"
	}

	if len(details) > 0 {
		return validInput{}, httpx.InvalidArgument("request contains invalid fields", details)
	}

	return validInput{
		treatedOn: treatedOn,
		services:  services,
		salonName: salonName,
		memo:      memo,
		cost:      in.Cost,
	}, nil
}

func parseTreatedOn(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return time.Time{}, fmt.Errorf("required")
	}
	parsed, err := time.Parse(DateLayout, trimmed)
	if err != nil {
		return time.Time{}, fmt.Errorf("must be formatted as YYYY-MM-DD")
	}
	return parsed, nil
}

func normalizeServices(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("at least 1 item is required")
	}
	if len(raw) > maxServices {
		return nil, fmt.Errorf("at most %d items are allowed", maxServices)
	}

	services := make([]string, 0, len(raw))
	for _, item := range raw {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			return nil, fmt.Errorf("items must not be empty")
		}
		services = append(services, trimmed)
	}
	return services, nil
}

func normalizeOptionalText(raw *string, maxRunes int) (*string, error) {
	if raw == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}
	if utf8.RuneCountInString(trimmed) > maxRunes {
		return nil, fmt.Errorf("must be at most %d characters", maxRunes)
	}
	return &trimmed, nil
}
