package types

type ILength interface {
	Length() int
}

type IDefinition interface {
	ILength
}

type IDefinitionHTML interface {
	IDefinition
	HTML() string
}

type IDefinerHTML interface {
	DefineHTML(traditional string) IDefinitionHTML
}
