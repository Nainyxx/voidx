package device

/*
Name - имя устройства
MAC - MAC-адрес устройства
RSSI - уровень сигнала устройства
Time - время последнего обнаружения устройства
*/

type Device struct {
	Name string `json:"name"`
	MAC  string `json:"mac"`
	RSSI int    `json:"rssi"`
	Time string `json:"time"`
}
