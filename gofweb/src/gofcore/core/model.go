package core

type NilModel struct{}

var NullModel *NilModel

func init() {
	NullModel = &NilModel{}
}
