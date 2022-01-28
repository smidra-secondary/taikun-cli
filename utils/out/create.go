package out

import (
	"errors"

	"github.com/itera-io/taikun-cli/api"
	"github.com/itera-io/taikun-cli/cmd/cmderr"
	"github.com/itera-io/taikun-cli/config"
	"github.com/itera-io/taikun-cli/utils/out/fields"
)

func PrintResult(resource interface{}, fields fields.Fields) error {
	if config.OutputOnlyID {
		return printResourceID(resource)
	}
	if config.OutputFormat == config.OutputFormatJson {
		return prettyPrintJson(resource)
	}
	return printApiResponseTable(resource, fields)
}

func printResourceID(resource interface{}) error {
	resourceMap, err := jsonObjectToMap(resource)
	if err != nil {
		return cmderr.ProgramError("printResourceID", err)
	}
	id, found := resourceMap["id"]
	if !found {
		return errors.New("response doesn't contain ID")
	}
	Println(resourceIDToString(id))
	return nil
}

func printApiResponseTable(response interface{}, fields fields.Fields) error {
	if config.AllColumns {
		fields.ShowAll()
	} else if len(config.Columns) != 0 {
		if err := fields.SetVisible(config.Columns); err != nil {
			return err
		}
	}

	resourceMap, err := getApiResponseResourceMap(response)
	if err != nil {
		return cmderr.ProgramError("printApiResponseTable", err)
	}

	t := newTable()
	for _, field := range fields.VisibleFields() {
		value, _ := getValueFromJsonMap(resourceMap, field.JsonPropertyName())
		t.AppendRow([]interface{}{
			field.Name(),
			trimCellValue(field.Format(value)),
		})
	}

	renderTable(t)
	return nil
}

func getApiResponseResourceMap(response interface{}) (resourceMap map[string]interface{}, err error) {
	resourceMap, err = jsonObjectToMap(response)
	if err != nil {
		return
	}
	if resourceMap[api.ResultField] != nil {
		resourceMap = resourceMap[api.ResultField].(map[string]interface{})
	} else if resourceMap[api.PayloadField] != nil {
		resourceMap = resourceMap[api.PayloadField].(map[string]interface{})
	}
	return
}
