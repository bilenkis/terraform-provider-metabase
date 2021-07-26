package metabase

type Details struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Db       string `json:"db"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Ssl      bool   `json:"ssl,omitempty"`
}

type DatabaseCreate struct {
	Engine string `json:"engine"`
	Name   string `json:"name"`
	//Details Details `json:"details"`
	Details Details `json:"details"`
}

type CacheFieldValues struct {
	ScheduleMinute int    `json:"schedule_minute"`
	ScheduleDay    string `json:"schedule_day"`
	ScheduleFrame  string `json:"schedule_frame"`
	ScheduleHour   int    `json:"schedule_hour"`
	ScheduleType   string `json:"schedule_type"`
}

type Schedules struct {
	CacheFieldValues CacheFieldValues `json:"cache_field_values"`
	MetadataSync     CacheFieldValues `json:"metadata_sync"`
}

type DatabaseRead struct {
	*Database
	Details   Details `json:"details"`
	Schedules Schedules  `json:"schedules"`
}

type Database struct {
	Description              string     `json:"description"`
	Features                 []string   `json:"features"`
	CacheFieldValuesSchedule string     `json:"cache_field_values_schedule"`
	Timezone                 string     `json:"timezone"`
	AutoRunQueries           bool       `json:"auto_run_queries"`
	MetadataSyncSchedule     string     `json:"metadata_sync_schedule"`
	Name                     string     `json:"name"`
	Caveats                  string     `json:"caveats"`
	IsFullSync               bool       `json:"is_full_sync"`
	UpdatedAt                string     `json:"updated_at"`
	NativePermissions        string     `json:"native_permissions"`
	Details                  Details `json:"details"`
	IsSample                 bool       `json:"is_sample"`
	Id                       int        `json:"id"`
	IsOnDemand               bool       `json:"is_on_demand"`
	Options                  string     `json:"options"`
	Engine                   string     `json:"engine"`
	Refingerprint            string     `json:"refingerprint"`
	CreatedAt                string     `json:"created_at"`
	PointsOfInterest         string     `json:"points_of_interest"`
}
