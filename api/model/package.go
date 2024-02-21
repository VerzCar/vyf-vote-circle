package model

import "database/sql/driver"

type SubscriptionPackage string

const (
	SubscriptionPackageS SubscriptionPackage = "S"
	SubscriptionPackageM SubscriptionPackage = "M"
	SubscriptionPackageL SubscriptionPackage = "L"
)

func (e *SubscriptionPackage) Scan(value interface{}) error {
	*e = SubscriptionPackage(value.(string))
	return nil
}

func (e SubscriptionPackage) Value() (driver.Value, error) {
	return string(e), nil
}

func (e SubscriptionPackage) IsValid() bool {
	switch e {
	case SubscriptionPackageS, SubscriptionPackageM, SubscriptionPackageL:
		return true
	}
	return false
}

func (e SubscriptionPackage) String() string {
	return string(e)
}
