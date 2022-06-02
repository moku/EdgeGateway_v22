package Parameters

const (
	InterfaceID                        = "if"
	DisposableIoTRequestID             = "dri"
	ResponseStatusCode                 = "rsc"
	DeviceID                           = "di"
	MicroserviceID                     = "mi"
	MicroserviceIDs                    = "mis"
	MicroserviceInformation            = "mif"
	MicroserviceInputParameter         = "ip"
	MicroserviceInputParameters        = "ips"
	MicroserviceOutputParameter        = "op"
	MicroserviceOutputParameters       = "ops"
	MicroserviceChangedOutputParameter = "cop"
	MicroserviceConfigure              = "mc"
	TaskID                             = "ti"
	TaskIDs                            = "tis"
	TaskInformation                    = "tif"
	FlexibleTaskParameter              = "fp"
	FlexibleTaskParameters             = "fps"
	StaticTaskParameter                = "sp"
	StaticTaskParameters               = "sps"
	TaskOrchestration                  = "to"
	TaskConfigure                      = "tc"
)

type TaskInformationObj struct {
	StaticTaskParameter   map[string]string
	FlexibleTaskParameter map[string]string
	TaskOrchestration     string
}

type Parameter struct {
	deviceID                string              //di
	interfaceID             string              //if
	disposableIoTRequestID  string              //dri
	responseStatusCode      string              //rsc
	microserviceID          string              //mi
	microserviceIDs         []string            //mis
	microserviceInformation map[string][]string //mif
	inputParameter          map[string]string   //ip
	outputParameter         map[string]string   //op
	outputParameters        []string            //ops
	changedOutputParameter  []string            //cop
	microserviceConfigure   string
	taskID                  string                        //ti
	taskIDs                 []string                      //tis
	taskInformation         map[string]TaskInformationObj //tif
	flexibleTaskParameter   map[string]int                //fp
	flexibleTaskParameters  []string                      //fps
	staticTaskParameter     map[string]int                //sp
	spArray                 []string                      //spArray for Tif
	staticTaskParameters    []string                      //sps
	taskConfigure           string
	myTopic                 string
	toTopic                 string
	protocol                string
	protocolAddr            string
}

func (p *Parameter) Protocol() string {
	return p.protocol
}

func (p *Parameter) SetProtocol(protocol string) {
	p.protocol = protocol
}

func (p *Parameter) ProtocolAddr() string {
	return p.protocolAddr
}

func (p *Parameter) SetProtocolAddr(protocolAddr string) {
	p.protocolAddr = protocolAddr
}

func (p *Parameter) ToTopic() string {
	return p.toTopic
}

func (p *Parameter) SetToTopic(toTopic string) {
	p.toTopic = toTopic
}

func (p *Parameter) MyTopic() string {
	return p.myTopic
}

func (p *Parameter) SetMyTopic(myTopic string) {
	p.myTopic = myTopic

}

func (p *Parameter) OutputParameter() map[string]string {
	return p.outputParameter
}

func (p *Parameter) SetOutputParameter(outputParameter map[string]string) {
	p.outputParameter = outputParameter
}

func (p *Parameter) MicroserviceInformation() map[string][]string {
	return p.microserviceInformation
}

func (p *Parameter) SetMicroserviceInformation(microserviceInformation map[string][]string) {
	p.microserviceInformation = microserviceInformation
}

func (p *Parameter) TaskInformation() map[string]TaskInformationObj {
	return p.taskInformation
}

func (p *Parameter) SetTaskInformation(taskInformation map[string]TaskInformationObj) {
	p.taskInformation = taskInformation
}

//di
func (p *Parameter) DeviceID() string {
	return p.deviceID
}

//di
func (p *Parameter) SetDeviceID(deviceID string) {
	p.deviceID = deviceID
}

func (p *Parameter) TaskConfigure() string {
	return p.taskConfigure
}

func (p *Parameter) SetTaskConfigure(taskConfigure string) {
	p.taskConfigure = taskConfigure
}

func (p *Parameter) StaticTaskParameters() []string {
	return p.staticTaskParameters
}

func (p *Parameter) SetStaticTaskParameters(staticTaskParameters []string) {
	p.staticTaskParameters = staticTaskParameters
}

func (p *Parameter) StaticTaskParameter() map[string]int {
	return p.staticTaskParameter
}

func (p *Parameter) SetStaticTaskParameter(staticTaskParameter map[string]int) {
	p.staticTaskParameter = staticTaskParameter
}

func (p *Parameter) FlexibleTaskParameters() []string {
	return p.flexibleTaskParameters
}

func (p *Parameter) SetFlexibleTaskParameters(flexibleTaskParameters []string) {
	p.flexibleTaskParameters = flexibleTaskParameters
}

func (p *Parameter) FlexibleTaskParameter() map[string]int {
	return p.flexibleTaskParameter
}

func (p *Parameter) SetFlexibleTaskParameter(flexibleTaskParameter map[string]int) {
	p.flexibleTaskParameter = flexibleTaskParameter
}

func (p *Parameter) TaskIDs() []string {
	return p.taskIDs
}

func (p *Parameter) SetTaskIDs(taskIDs []string) {
	p.taskIDs = taskIDs
}

func (p *Parameter) TaskID() string {
	return p.taskID
}

func (p *Parameter) SetTaskID(taskID string) {
	p.taskID = taskID
}

func (p *Parameter) MicroserviceConfigure() string {
	return p.microserviceConfigure
}

func (p *Parameter) SetMicroserviceConfigure(microserviceConfigure string) {
	p.microserviceConfigure = microserviceConfigure
}

func (p *Parameter) ChangedOutputParameter() []string {
	return p.changedOutputParameter
}

func (p *Parameter) SetChangedOutputParameter(changedOutputParameter []string) {
	p.changedOutputParameter = changedOutputParameter
}

func (p *Parameter) OutputParameters() []string {
	return p.outputParameters
}

//ops string array
func (p *Parameter) SetOutputParameters(outputParameters []string) {
	p.outputParameters = outputParameters
}

func (p *Parameter) InputParameter() map[string]string {
	return p.inputParameter
}

func (p *Parameter) SetInputParameter(inputParameter map[string]string) {
	p.inputParameter = inputParameter
}

//mis
func (p *Parameter) MicroserviceIDs() []string {

	return p.microserviceIDs
}

//mis
func (p *Parameter) SetMicroserviceIDs(microserviceIDs []string) {

	//for i := range microserviceIDs {
	//
	//	microserviceIDs[i] = trimQuote(microserviceIDs[i])
	//}

	p.microserviceIDs = microserviceIDs
}

func (p *Parameter) MicroserviceID() string {
	return p.microserviceID
}

func (p *Parameter) SetMicroserviceID(microserviceID string) {
	p.microserviceID = microserviceID
}

//rsc
func (p *Parameter) ResponseStatusCode() string {
	return p.responseStatusCode
}

//rsc
func (p *Parameter) SetResponseStatusCode(responseStatusCode string) {
	p.responseStatusCode = responseStatusCode
}

//dri
func (p *Parameter) DisposableIoTRequestID() string {
	return p.disposableIoTRequestID
}

//dri
func (p *Parameter) SetDisposableIoTRequestID(disposableIoTRequestID string) {
	//p.disposableIoTRequestID = trimQuote(disposableIoTRequestID)
	p.disposableIoTRequestID = disposableIoTRequestID
}

//if
func (p *Parameter) InterfaceID() string {

	return p.interfaceID
}

//Set "if"
func (p *Parameter) SetInterfaceID(interfaceID string) {

	p.interfaceID = trimQuote(string(interfaceID))
}

func NewParameter() *Parameter {

	return &Parameter{}
}

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
