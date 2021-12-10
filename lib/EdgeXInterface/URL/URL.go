package URL

const (
	EdgeXAddr        = "http://localhost"
	CoreMetadataPort = "59881"
	CoreCommandPort  = "59882"
	CoreDataPort     = "59880"
	NotificationPort = "59860"
)

const (
	CoreMetadata        string = EdgeXAddr + ":" + CoreMetadataPort
	DeviceRegistration  string = CoreMetadata + "/api/v2/device"
	GetDeviceInfo       string = CoreMetadata + "/api/v2/device/all"
	GetDeviceInfoByName string = CoreMetadata + "/api/v2/device/name"
	CoreCommand         string = EdgeXAddr + ":" + CoreCommandPort
	SendDeviceCommand   string = CoreCommand + "/api/v2/device/name"
	SetResource         string = CoreCommand + "/api/v2/device/name"
	CoreData            string = EdgeXAddr + ":" + CoreDataPort
	InsertSensingData   string = CoreData + "/api/v2/event"
	Notification        string = EdgeXAddr + ":" + NotificationPort
	InsertNotification  string = Notification + "/api/v2/notification"
	DeleteNotification  string = Notification + "/api/v2/notification/slug"
)
