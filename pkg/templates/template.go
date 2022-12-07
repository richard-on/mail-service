// Package templates provides the templates for the emails and selects one based on request
package templates

import (
	"fmt"
	"github.com/richard-on/mail-service/pkg/server/request"
)

type Template interface {
	setTemplate(req *request.SendMail) (string, string, error)
}

// GetTemplate returns the template based on the request
func GetTemplate(req *request.SendMail) (plain string, html string, err error) {

	defer func() {
		r := recover()

		if r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			err = fmt.Errorf("%v: %w", err, ErrBadFormat)
		}
	}()

	switch req.Type {
	case "coordination":
		var c Coordination
		c.acceptLink = req.Template.(map[string]interface{})["acceptLink"].(string)
		c.declineLink = req.Template.(map[string]interface{})["declineLink"].(string)
		return c.setTemplate(req)

	case "verification":
		var v Verification
		v.verifyLink = req.Template.(map[string]interface{})["verifyLink"].(string)
		return v.setTemplate(req)

	case "info":
		var i Info
		i.body = req.Template.(map[string]interface{})["body"].(string)
		return i.setTemplate(req)

	default:
		return "", "", ErrNoSuchTemplate
	}
}
