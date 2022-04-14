package models

import "github.com/tbellembois/gochimitheque/aes"

// Person represent a person.
type Person struct {
	PersonID       int           `db:"person_id" json:"person_id" schema:"person_id"`
	PersonEmail    string        `db:"person_email" json:"person_email" schema:"person_email"`
	PersonPassword string        `db:"person_password" json:"person_password" schema:"person_password"`
	PersonAESKey   string        `db:"person_aeskey" json:"person_aeskey" schema:"person_aeskey"`
	Permissions    []*Permission `db:"-" json:"permissions" schema:"permissions"`
	Entities       []*Entity     `db:"-" json:"entities" schema:"entities"`
	CaptchaText    string        `db:"-" schema:"captcha_text" json:"captcha_text"`
	CaptchaUID     string        `db:"-" schema:"captcha_uid" json:"captcha_uid"`
	QRCode         []byte        `db:"-" schema:"qrcode" json:"qrcode"`
}

func (p *Person) GeneratePassword() (err error) {
	p.PersonPassword, err = aes.GenerateAESKey()

	return
}
