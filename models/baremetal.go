package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/jinzhu/copier"

	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type BareMetal struct {
	Compute
}

func (m BareMetal) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

type BareMetals []BareMetal

func (m BareMetals) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

//// backend api calls:

// SyncBareMetals() sync all bare metals of user's account from origin.
func SyncBareMetals(user *User, since time.Time) (count int, err error) {
	log.Infof("sync baremetal servers... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	account := user.Account()
	if account == nil {
		log.Errorf("account link broken! %v of %v", user.ID, user.AccountId)
		return 0, errors.New("account link broken!")
	}
	log.Debugf("account: %v", account)

	date_since := since.Format("01/02/2006 15:04:05")
	log.Infof("try to sync baremetal servers from %v...", date_since)

	service := services.GetAccountService(sess)
	data, err := service.
		Mask("id;accountId;hourlyBillingFlag;hostname;domain;notes;tagReferences.tag.name;tagReferences.tag.id;userData.value;provisionDate;bandwidthAllocation;privateNetworkOnlyFlag;primaryIpAddress;primaryBackendIpAddress;networkVlans.id;operatingSystem.id;operatingSystem.softwareLicense.softwareDescription.longDescription;datacenter.id;location.pathString;location.id;virtualRackId;processorPhysicalCoreAmount;memoryCapacity;networkGatewayMemberFlag;networkManagementIpAddress;rack.id").
		Filter(filter.Build(
			filter.Path("hardware.provisionDate").DateAfter(date_since),
		)).
		GetHardware()
	if err != nil {
		log.Errorf("slapi error: %v", err)
		return 0, err
	}
	inspect("baremetals", data)

	count = 0
	for _, el := range data {
		comp := &Compute{}
		inspect("origin virtual guest", el)
		copier.Copy(comp, el)
		comp.ID = *el.Id + 8000000000000
		comp.Type = "Metal"
		comp.OSName = *el.OperatingSystem.SoftwareLicense.SoftwareDescription.LongDescription
		comp.ProvisionDate, _ = time.Parse(time.RFC3339,
			el.ProvisionDate.String())
		if len(el.UserData) == 1 {
			comp.CloudUserData = *el.UserData[0].Value
		}
		inspect("compute instance", comp)

		// relational things...
		comp.MapUserId(user.ID)
		for _, t := range el.TagReferences {
			comp.MapTagId(*t.Tag.Id)
			tag := Tag{
				ID:        *t.Tag.Id,
				Name:      *t.Tag.Name,
				AccountId: comp.AccountId,
			}
			tag.Save()
		}
		for _, v := range el.NetworkVlans {
			comp.MapVlanId(*v.Id)
		}

		err = comp.Save()	// after save, ID will be set '0'. why?
		if err != nil {
			log.Errorf("cannot save virtual guest: %v, %v", err, comp)
		} else {
			count++
		}
	}

	if len(data) == count {
		log.Infof("Bingo! all data were inserted to database! (%v)", count)
	} else {
		log.Errorf("Oops! some data not inserted! %v of %v saved.",
			count, len(data))
	}
	return count, nil
}
