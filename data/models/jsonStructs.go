package models

type BookDates struct {
	BookingDates []string `json:"booking_dates"`
	WorkingDates []string `json:"working_dates"`
}

type BookTimes []struct {
	Time         string `json:"time"`
	SeanceLength int    `json:"seance_length"`
	SumLength    int    `json:"sum_length"`
	Datetime     string `json:"datetime"`
}

type BookForm struct {
	Phone               string        `json:"phone"`
	Fullname            string        `json:"fullname"`
	Email               string        `json:"email"`
	Code                interface{}   `json:"code"`
	NotifyBySms         int           `json:"notify_by_sms"`
	Comment             string        `json:"comment"`
	Appointments        Appointments  `json:"appointments"`
	IsMobile            bool          `json:"isMobile"`
	Referrer            string        `json:"referrer"`
	AppointmentsCharges []interface{} `json:"appointments_charges"`
	IsSupportCharge     bool          `json:"is_support_charge"`
	RedirectPrefix      string        `json:"redirect_prefix"`
	BookformID          int           `json:"bookform_id"`
}

type Appointments []struct {
	ID           int           `json:"id"`
	Services     []int         `json:"services"`
	StaffID      int           `json:"staff_id"`
	Events       []interface{} `json:"events"`
	Datetime     string        `json:"datetime"`
	ChargeStatus string        `json:"chargeStatus"`
	CustomFields struct {
		Gcid string `json:"gcid"`
	} `json:"custom_fields"`
}
