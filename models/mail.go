package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/markbates/inflect"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type Mail struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	SingleID    uuid.UUID `json:"single_id" db:"single_id"`
	Type        string    `json:"type" db:"type"`
	Status      string    `json:"status" db:"status"`
	MailID      string    `json:"mail_id" db:"mail_id"`
	MailResp    string    `json:"mail_resp" db:"mail_resp"`
	Sender      string    `json:"sender" db:"sender"`
	Subject     string    `json:"subject" db:"subject"`
	Rcpt        string    `json:"rcpt" db:"rcpt"`
	Bccs        string    `json:"bccs" db:"bccs"`
	ContentText string    `json:"content_text" db:"content_text"`
	ContentHtml string    `json:"content_html" db:"content_html"`
}

func (m Mail) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

type Mails []Mail

func (m Mails) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

type Object interface {
	String() string
}

// Validate gets run everytime you call a "pop.Validate" method.
func (m *Mail) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Status, Name: "Status"},
		&validators.StringIsPresent{Field: m.Sender, Name: "Sender"},
		&validators.StringIsPresent{Field: m.Subject, Name: "Subject"},
		&validators.StringIsPresent{Field: m.Rcpt, Name: "Rcpt"},
		&validators.StringIsPresent{Field: m.ContentText, Name: "ContentText"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (m *Mail) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (m *Mail) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

////

func AdminMail(s *Single, obj Object, subj, to string, groups ...string) error {
	var bccs []string
	for _, group := range groups {
		singles, err := GetSinglesByPermission(group)
		if err == nil {
			for _, s := range *singles {
				bccs = append(bccs, s.Mail())
			}
		}
	}
	log.Infof("sending admin mail for %v, bcc to %v", groups, bccs)
	return FormMail(s, obj, subj, to, bccs...)
}

func FormMail(s *Single, obj Object, subj, to string, bccs ...string) error {
	templ_name := inflect.Underscore(subj)
	log.Infof("preparing form mail with template %v", templ_name)

	t_text, err := template.ParseFiles("templates/mail." + templ_name + ".text")
	if err != nil {
		log.Errorf("error on parsing text template: %v", err)
		return err
	}
	buf := &bytes.Buffer{}
	if err = t_text.Execute(buf, obj); err != nil {
		log.Errorf("error on executing text template: %v", err)
		return err
	}
	cont_text := buf.String()

	t_html, err := template.ParseFiles("templates/mail." + templ_name + ".html")
	if err != nil {
		log.Errorf("error on parsing html template: %v", err)
		return err
	}
	buf = &bytes.Buffer{}
	if err = t_html.Execute(buf, obj); err != nil {
		log.Errorf("error on executing html template: %v", err)
		return err
	}
	cont_html := buf.String()

	subject := subj + ": " + obj.String()
	m := PrepareMail(subject, cont_text, to)
	m.ContentHtml = cont_html
	return m.Send(s.ID, bccs...)
}

// SendMail() prepare basic text mail and send immediately.
func SendMail(single_id uuid.UUID, subj, cont, to string) error {
	m := PrepareMail(subj, cont, to)
	return m.Send(single_id)
}

// PrepareMail() prepare basic text mail and return it.
func PrepareMail(subj, cont, to string) (mail *Mail) {
	mail = &Mail{
		Sender:      mail_sender,
		Subject:     subj,
		ContentText: cont,
		Rcpt:        to,
	}
	return mail
}

// Send() send a mail with current mailer and save its status to database.
func (m *Mail) Send(single_id uuid.UUID, bccs ...string) error {
	log.Debugf("sending a mail by %v", single_id)
	resp, id, err := send(m, bccs...)
	if err != nil {
		return err
	}
	m.MailResp = resp
	m.MailID = id
	m.Status = "sent"
	m.SingleID = single_id
	m.Bccs = fmt.Sprintf("%v", bccs)
	return m.save()
}

// save() saves this mail.
func (m *Mail) save() error {
	if m.ID == (uuid.UUID{}) {
		verrs, err := DB.ValidateAndSave(m)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	} else {
		log.Debugf("saving existing mail %v...", m)
		log.Errorf("IS IT POSSIBLE TO SAVE EXISTING MAIL AGAIN?")
	}
	return nil
}

// independent mail sender: currently implemented with mailgun!
func send(m *Mail, bccs ...string) (resp, id string, err error) {
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Errorf("cannot setup mailgun from env: %v", err)
		return "", "", err
	}
	log.Infof("about to send mail... %v to %v", m.Subject, m.Rcpt)
	message := mailgun.NewMessage(m.Sender, m.Subject, m.ContentText, m.Rcpt)
	if len(m.ContentHtml) > 0 {
		message.SetHtml(m.ContentHtml)
	}
	for _, el := range bccs {
		log.Debugf("add %v as BCC...", el)
		message.AddBCC(el)
	}
	log.Infof("shot the gun to %v rcpts.", message.RecipientCount())
	return mg.Send(message)
}
