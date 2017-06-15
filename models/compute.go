package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"

	"github.com/softlayer/softlayer-go/datatypes"
)

type Compute struct {
	ID                  int       `json:"id" db:"id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	AccountId           int       `json:"account_id" db:"account_id"`
	IsHourly            bool      `json:"is_hourly" db:"is_hourly"`
	Hostname            string    `json:"hostname" db:"hostname"`
	Domain              string    `json:"domain" db:"domain"`
	Notes               string    `json:"notes" db:"notes"`
	ProvisionDate       time.Time `json:"provision_date" db:"provision_date"`
	BandwidthAllocation int       `json:"bandwidth_allocation" db:"bandwidth_allocation"`
	IsPrivateOnly       bool      `json:"is_private_only" db:"is_private_only"`
	PublicIP            string    `json:"public_ip" db:"public_ip"`
	PrivateIP           string    `json:"private_ip" db:"private_ip"`
	ManagementIP        string    `json:"management_ip" db:"management_ip"`
	IsGatewayMember     bool      `json:"is_gateway_member" db:"is_gateway_member"`
	OperatingSystemId   int       `json:"operating_system_id" db:"operating_system_id"`
	OSName              string    `json:"os_name" db:"os_name"`
	DatacenterId        int       `json:"datacenter_id" db:"datacenter_id"`
	VirtualRackId       int       `json:"virtual_rack_id" db:"virtual_rack_id"`
	RackId              int       `json:"rack_id" db:"rack_id"`
	HasPendingMigration bool      `json:"has_pending_migration" db:"has_pending_migration"`
	Cores               int       `json:"cores" db:"cores"`
	Memories            int       `json:"memories" db:"memories"`
	Type                string    `json:"type" db:"type"`
	CloudUserData       string    `json:"user_data" db:"user_data"`
	Path                string    `json:"path" db:"path"`
}

func (c Compute) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

type Computes []Compute

func (c Computes) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (c *Compute) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Hostname, Name: "Hostname"},
		&validators.StringIsPresent{Field: c.Domain, Name: "Domain"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (c *Compute) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (c *Compute) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//// display helpers:

func (c Compute) Nick() string {
	return c.Hostname + "." + c.Domain
}

//// database fuctions:

func PickCompute(id int, ty string) (c *Compute) {
	c = &Compute{}
	err := DB.Where("id = ? AND type LIKE ?", id, "%"+ty).First(c)
	if err != nil {
		return nil
	}
	return
}

// Save()
func (c *Compute) Save() error {
	old := PickCompute(c.ID, c.Type)
	if old == nil {
		verrs, err := DB.ValidateAndCreate(c)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	} else {
		log.Debugf("saving existing %v...", c)
		verrs, err := DB.ValidateAndUpdate(c)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			return verrs
		}
	}
	return nil
}

func (c Compute) TagNames() (tnames []string) {
	for _, t := range *c.Tags() {
		tnames = append(tnames, t.Name)
	}
	return
}

func (c Compute) Tags() (tags *Tags) {
	tags = &Tags{}
	DB.BelongsToThrough(&c, "comp_tag_maps").All(tags)
	return
}

//// Compute - User Map

type CompUserMap struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ComputeID int       `json:"compute_id" db:"compute_id"`
	UserID    int       `json:"user_id" db:"user_id"`
}

func (c *Compute) MapUserId(id int) error {
	o := &CompUserMap{}
	err := DB.Where("compute_id = ? AND user_id = ?", c.ID, id).First(o)
	if err != nil {
		nm := &CompUserMap{
			ComputeID: c.ID,
			UserID:    id,
		}
		_, err := DB.ValidateAndSave(nm)
		return err
	}
	return errors.New("association alread exist")
}

//// Compute - VLAN Map

type CompVlanMap struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ComputeID int       `json:"compute_id" db:"compute_id"`
	VlanID    int       `json:"vlan_id" db:"vlan_id"`
}

func (c *Compute) MapVlanId(id int) error {
	o := &CompVlanMap{}
	err := DB.Where("compute_id = ? AND vlan_id = ?", c.ID, id).First(o)
	if err != nil {
		nm := &CompVlanMap{
			ComputeID: c.ID,
			VlanID:    id,
		}
		_, err := DB.ValidateAndSave(nm)
		return err
	}
	return errors.New("association alread exist")
}

//// Compute - Tag Map

type CompTagMap struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ComputeID int       `json:"compute_id" db:"compute_id"`
	TagID     int       `json:"tag_id" db:"tag_id"`
}

func (c *Compute) MapTagId(id int) error {
	o := &CompTagMap{}
	err := DB.Where("compute_id = ? AND tag_id = ?", c.ID, id).First(o)
	if err != nil {
		nm := &CompTagMap{
			ComputeID: c.ID,
			TagID:     id,
		}
		_, err := DB.ValidateAndSave(nm)
		return err
	}
	return errors.New("association alread exist")
}

//// callers:

func SyncComputes(user *User) (count int, err error) {
	since := user.lastSyncTimeCompute().AddDate(0, 0, -1)
	vg_count, err := SyncVirtualGuests(user, since)
	count += vg_count
	bm_count, err := SyncBareMetals(user, since)
	count += bm_count
	return
}

//// setters:
func (c *Compute) HourlyBillingFlag(s *bool) {
	c.IsHourly = *s
}

func (c *Compute) PrivateNetworkOnlyFlag(s *bool) {
	c.IsPrivateOnly = *s
}

// bm only
func (c *Compute) NetworkGatewayMemberFlag(s *bool) {
	c.IsGatewayMember = *s
}

func (c *Compute) PendingMigrationFlag(s *bool) {
	c.HasPendingMigration = *s
}

// vsi only
func (c *Compute) DedicatedAccountHostOnlyFlag(s *bool) {
	if *s {
		c.Type = "Private"
	} else {
		c.Type = "Public"
	}
}

func (c *Compute) PrimaryBackendIpAddress(s *string) {
	if s != nil {
		c.PrivateIP = *s
	}
}

func (c *Compute) PrimaryIpAddress(s *string) {
	if s != nil {
		c.PublicIP = *s
	}
}

// bm case
func (c *Compute) NetworkManagementIpAddress(s *string) {
	if s != nil {
		c.ManagementIP = *s
	}
}

func (c *Compute) Location(s *datatypes.Location) {
	c.Path = *s.PathString
}

func (c *Compute) Datacenter(s *datatypes.Location) {
	c.DatacenterId = *s.Id
}

// vsi case
func (c *Compute) MaxCpu(s *int) {
	c.Cores = *s
}

// vsi case
func (c *Compute) MaxMemory(s *int) {
	c.Memories = *s / 1024
}

// bm case
func (c *Compute) ProcessorPhysicalCoreAmount(s *uint) {
	c.Cores = int(*s)
}

// bm case
func (c *Compute) MemoryCapacity(s *uint) {
	c.Memories = int(*s)
}

func (c *Compute) OperatingSystem(s *datatypes.Software_Component_OperatingSystem) {
	c.OperatingSystemId = *s.Id
}

// bm only
func (c *Compute) Rack(s *datatypes.Location) {
	c.RackId = *s.Id
}
