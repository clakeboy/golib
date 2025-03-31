package soap

type ArgsMap []*WsdlFunctionArgs

type WsdlFunction struct {
	Name         string            //方法名
	Action       string            //SOAP Action
	RequestArgs  *WsdlFunctionArgs //调用方法参数集合
	ResponseName string            //方法返回值名
	ResponseArgs *WsdlFunctionArgs //回调用方法参数集合
}

type WsdlFunctionArgs struct {
	Name      string  //参数名
	Namespace string  //参数命名空间
	MinOccurs int     //最小出现个数
	MaxOccurs string  //最大出现个数
	Nillable  bool    //是否为空
	Type      string  //类型
	Elements  ArgsMap //如果type为其它类型
}
