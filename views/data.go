package views

import "log"

const (
	AlertLvlError	= "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo	= "info"
	AlertLvlSuccess	= "success"

	// AlertMsgGeneric is a generic error message displayed whenever a random error
	// is encountered by the backend
	AlertMsgGeneric = "Something went wrong..."
)

// Data is the top level structure that views expect data to come in through
type Data struct {
	Alert *Alert
	Yield interface{}
}

// Alert is used to render Boostrap alert messages in templates
type Alert struct {
	Level	string
	Message	string
}

// PublicError is interface defined to pass err.Public() via Data
type PublicError interface {
	error
	Public() string
}

// SetAlert will set an alert on Data type using an error
func (d *Data) SetAlert(err error) {
	var msg string
	// type assertion - will return err and a boolean which gets checked
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}

	d.Alert = &Alert{
		Level: AlertLvlError,
		Message: msg,
	}
}

// AlertError allows for easier custom error messages to be created on the fly inside the code
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level: AlertLvlError,
		Message: msg,
	}
}