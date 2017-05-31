package grifts

import (
	"fmt"

	"github.com/markbates/grift/grift"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

var _ = grift.Add("db:seed", func(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		p := fmt.Printf
		for v := 250; v < 500; v++ {
			vlan := &models.Vlan{}
			vlan.ID = v
			p("vlan %v\n", vlan)
			err := tx.Create(vlan)
			if err != nil {
				p("error on %v: %v", vlan, err)
				return err
			}
		}
		return nil
	})
})
