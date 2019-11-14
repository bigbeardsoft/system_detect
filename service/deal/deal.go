package deal

// DealResponse 定义处理接收消息的接口
type DealResponse interface {
	// Deal 声明处理接口
	Deal(map[string]interface{},func ())
}

// DealReg 实现处理接收消息功能
type DealReg struct {
}

//DealServerStatReport 处理状态码
type DealServerStatReport struct{}

func (d *DealReg) Deal(map[string]interface{}) {

}

func (d *DealServerStatReport) Deal(map[string]interface{}) {

}
