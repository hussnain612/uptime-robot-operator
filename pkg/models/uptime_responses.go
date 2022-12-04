package models

type UptimeRobotMonitorResponse struct {
	Stat    string                   `json:"stat"`
	Monitor UptimeRobotMonitorStatus `json:"monitor"`
	Error   UptimeRobotError         `json:"error"`
}

type UptimeRobotError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type UptimeRobotMonitorStatus struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}
