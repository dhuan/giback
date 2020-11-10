package app

import (
	"fmt"
	"strings"
)

func ValidateUnits(units []PushUnit) ([]string, []string) {
	var errors []string
	var ids []string

	for i := range units {
		unit := units[i]
		unitId := fmt.Sprintf("%d", i)

		if unit.Id != "" {
			unitId = unit.Id
		}

		missingFields := getMissingFieldsFromUnit(unit)

		if len(missingFields) == 0 {
			continue
		}

		errors = append(errors, invalidUnitMessage(i, unit, missingFields))
		ids = append(ids, unitId)
	}

	return errors, ids
}

func invalidUnitMessage(unitIndex int, unit PushUnit, missingFields []string) string {
	missingFieldsString := strings.Join(missingFields, ",")

	return fmt.Sprintf("Missing the following fields: %s", missingFieldsString)
}

func getMissingFieldsFromUnit(unit PushUnit) []string {
	var fields []string

	if unit.Id == "" {
		fields = append(fields, "id")
	}

	if unit.Repository == "" {
		fields = append(fields, "repository")
	}

	if len(unit.Files) == 0 {
		fields = append(fields, "files")
	}

	return fields
}
