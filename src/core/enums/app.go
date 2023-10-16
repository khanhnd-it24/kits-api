package enums

type AppMode string

const (
	AppModeDev  AppMode = "dev"
	AppModeProd AppMode = "prod"
)

func (m AppMode) IsProd() bool {
	return m == AppModeProd
}

func (m AppMode) IsDev() bool {
	return m == AppModeDev
}
