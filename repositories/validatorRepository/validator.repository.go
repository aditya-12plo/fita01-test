package validatorRepository

type ValidatorResult struct {
	FieldName   string      `json:"field_name"`
	Validated   string      `json:"validated"`
	TypeData    string      `json:"type_data"`
	Value       interface{} `json:"value"`
	Param       string      `json:"param"`
	Description string      `json:"description"`
}
