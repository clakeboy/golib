package task

//全局类
var ManagementService *Management
//初始化并启动任务管理
func InitManagement() {
	ManagementService = NewManagement()
	ManagementService.Start()
}