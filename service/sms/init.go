package sms

type Sms interface {
	Send(phone, content string)  (bool,error)
}

func SelectSms(name int) *Sms {
	var pp = Sms(nil)
	if name == 1{
		pp = new(NxCloud)
	}
	return &pp
}
