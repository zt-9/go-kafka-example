package utils

import (
	"fmt"
	"go-kafka-example/models"
)

func HasDuplication(target string, arr []string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

// utils
// store labels associated with the address
func StoreAddressLabels(address string, labels []string, db *models.MetadataDB[models.Label]) error {
	fmt.Println("db related labels")
	for _, label := range labels {
		// check duplication
		exists, err := db.KeyExists(label)

		if err != nil {
			return err
		}

		// check if label exist or address already exists in the label's addr arr

		// label already exist, add address to the label if the address does not exists in label's addr arr
		if exists {
			value, err := db.Get(label)
			if err != nil {
				return err
			}
			arr := value.Addresses

			exist := HasDuplication(address, arr)
			if exist { // address exists in label.Addresses
				return fmt.Errorf("address already exists in the label")
			} else {
				// address does not exsit. append address to the array
				value.Addresses = append(value.Addresses, address)
				if err := db.Update(label, value); err != nil {
					return err
				}
			}
		}

		// label does not exist. store the label to db
		var data models.Label
		data.Addresses = append(data.Addresses, address)
		data.Label = label
		if err := db.Update(label, data); err != nil {
			return err
		}
	}
	return nil
}
