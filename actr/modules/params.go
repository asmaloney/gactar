package modules

func Ptr[T any](v T) *T {
	return &v
}

type ParamInfo struct {
	Name        string
	Description string
}

type ParamInt struct {
	ParamInfo

	Min *int
	Max *int
}

type ParamFloat struct {
	ParamInfo

	Min *float64
	Max *float64
}

type ParamInterface interface {
	GetName() string
	GetDescription() string

	GetMin() *float64
	GetMax() *float64
}

type ParamInfoMap map[string]ParamInterface

func (p ParamInfo) GetName() string {
	return p.Name
}

func (p ParamInfo) GetDescription() string {
	return p.Description
}

func (p ParamInt) GetMin() *float64 {
	if p.Min != nil {
		temp := float64(*p.Min)
		return &temp
	}
	return nil
}

func (p ParamInt) GetMax() *float64 {
	if p.Max != nil {
		temp := float64(*p.Max)
		return &temp
	}
	return nil
}

func (p ParamFloat) GetMin() *float64 { return p.Min }
func (p ParamFloat) GetMax() *float64 { return p.Max }

func NewParamInt(name, description string, min, max *int) ParamInt {
	return ParamInt{
		ParamInfo{name, description},
		min, max,
	}
}

func NewParamFloat(name, description string, min, max *float64) ParamFloat {
	return ParamFloat{
		ParamInfo{name, description},
		min, max,
	}
}
