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

type VirtualGuest struct {
	Compute
}

func (m VirtualGuest) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

type VirtualGuests []VirtualGuest

func (m VirtualGuests) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

//// backend api calls:

// SyncTickets() creates and updates all Tickets of given user's account.
func SyncVirtualGuests(user *User) (count int, err error) {
	log.Infof("sync virtual servers... (use %v)", user.Username)
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"

	account := user.Account()
	if account == nil {
		log.Errorf("account link broken! %v of %v", user.ID, user.AccountId)
		return 0, errors.New("account link broken!")
	}
	log.Debugf("account: %v", account)

	date_since := user.lastSyncTimeCompute().
		AddDate(0, 0, -1).Format("01/02/2006 15:04:05")
	log.Infof("try to sync virtual servers from %v...", date_since)

	service := services.GetAccountService(sess)
	data, err := service.
		Mask("id;accountId;hourlyBillingFlag;hostname;domain;notes;tagReferences.tag.name;tagReferences.tag.id;userData.value;provisionDate;createDate;modifyDate;bandwidthAllocation;privateNetworkOnlyFlag;primaryIpAddress;primaryBackendIpAddress;networkVlans.id;operatingSystem.id;operatingSystem.softwareLicense.softwareDescription.longDescription;datacenter.id;location.pathString;location.id;virtualRackId;startCpus;maxCpu;maxCpuUnits;maxMemory;type.name;modifyDate;pendingMigrationFlag;dedicatedAccountHostOnlyFlag;dedicatedHost;host;users.id").
		Filter(filter.Build(
			filter.Path("virtualGuests.modifyDate").DateAfter(date_since),
		)).
		GetVirtualGuests()
		// status.name;powerState.name
	if err != nil {
		log.Errorf("slapi error: %v", err)
		return 0, err
	}

	count = 0
	for _, el := range data {
		comp := &Compute{}
		inspect("origin virtual guest", el)
		copier.Copy(comp, el)
		comp.ID = *el.Id + 1000000000000
		comp.OSName = *el.OperatingSystem.SoftwareLicense.SoftwareDescription.LongDescription
		comp.ProvisionDate, _ = time.Parse(time.RFC3339,
			el.ProvisionDate.String())
		if len(el.UserData) == 1 {
			comp.CloudUserData = *el.UserData[0].Value
		}
		inspect("compute instance", comp)

		// relational things...
		for _, u := range el.Users {
			comp.MapUserId(*u.Id)
		}
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
