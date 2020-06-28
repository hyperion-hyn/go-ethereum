package storage

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
)

type Description struct {
	name string `storage:"slot0"`
	url  string `storage:"slot1"`
}

type Delegation struct {
	amount      int `storage:"slot0"`
	blockNumber int `storage:"slot1"`
}

type Validator struct {
	desc        Description  `storage:"slot0"`
	delegations []Delegation `storage:"slot1"`
}

type Donation struct {
	Name   string `storage:"slot0"`
	amount int    `storage:"slot1"`
}
type ValidatorList struct {
	Name       string              `storage:"slot0"`
	author     string              `storage:"slot1"`
	count      int                 `storage:"slot2"`
	Desc       Description         `storage:"slot3"`
	validators []Validator         `storage:"slot4"`
	donations  map[string]Donation `storage:"slot5"`
}

type GlobalVariables struct {
	Version       int           `storage:"slot0"`
	Name          string        `storage:"slot1"`
	ValidatorList ValidatorList `storage:"slot2"`
}

// Tests parseTag
func TestParseTag(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	var tests = []struct {
		tag      string // input
		expected int    // expected result
		err      error
	}{
		{"", 0, errors.New(fmt.Sprintf("invalid tag: "))},
		{"slot", 0, errors.New(fmt.Sprintf("invalid tag: slot"))},
		{"slot0", 0, nil},
		{"slot99", 99, nil},
		{"slot-10", 0, errors.New(fmt.Sprintf("invalid tag: slot-10"))},
	}

	for _, tt := range tests {
		actual, err := parseTag(tt.tag)
		if actual != tt.expected {
			if err.Error() != tt.err.Error() {
				t.Errorf("parseTag(%s): expected %d, actual %d, '%v+', '%v+'", tt.tag, tt.expected, actual, tt.err, err)
			}
		}
	}
}

// Tests that storage manipulation
func TestStorageManipulation(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	// Create an empty state database
	db := rawdb.NewMemoryDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))
	addr := common.BytesToAddress([]byte{9})

	var globalVariables GlobalVariables = GlobalVariables{
		ValidatorList: ValidatorList{
			Name:   "atlas",
			author: "hyperion",
			count:  11,
			Desc: Description{
				name: "hyperion",
				url:  "https://www.hyn.space",
			},
			validators: nil,
		},
	}

	log.Info("ValueOf", "globalVariables", reflect.ValueOf(&globalVariables))
	log.Info("ValueOf", "globalVariables.validatorList", reflect.ValueOf(&globalVariables).Elem().FieldByName("ValidatorList"))
	log.Info("ValueOf", "globalVariables.validatorList", reflect.ValueOf(&globalVariables).Elem().FieldByName("ValidatorList").FieldByName("Name"))
	// reflect.ValueOf(&globalVariables).Elem().FieldByName("ValidatorList").FieldByName("author").SetString("ethereum")
	log.Info("ValueOf", "ValidatorList.Name", globalVariables.ValidatorList.Name)
	reflect.Indirect(reflect.ValueOf(&globalVariables)).FieldByName("ValidatorList").FieldByName("Name").SetString("harmony")
	log.Info("ValueOf", "ValidatorList.Name", globalVariables.ValidatorList.Name)
	// reflect.ValueOf(validatorList).FieldByName("Name").SetString("harmony")

	{
		storage := NewStorage(state, addr, 0, &globalVariables, nil)
		log.Info("TestStorageManipulation", "validatorList", globalVariables)
		name := storage.GetByName("ValidatorList").GetByName("Desc").GetByName("name")
		// name := storage.GetByName("validators").GetByName("name")
		log.Info("result", "validatorList.Name", name.Value())
		//     {
		//         name := storage.GetByName("Name")
		//         log.Info("result", "validatorList.Name", name.Value())
		name.SetValue("harmony")
		log.Info("result", "validatorList.Name", name.Value())
		//     }
		//     // os.Exit(1)
	}

	storage := NewStorage(state, addr, 0, &globalVariables, nil)

	{
		log.Info("TestStorageManipulation", "Version", globalVariables.Version)
		version := storage.Get("Version")
		version.SetValue(0b101010)
		log.Info("TestStorageManipulation", "Version", globalVariables.Version)
		storage.Flush()

	}

	log.Info("TestStorageManipulation", "validatorList", globalVariables)
	// name := storage.GetByName("validators").GetByName("desc").GetByName("name")
	// name := storage.GetByName("validators").GetByName("name")
	{
		name := storage.GetByName("ValidatorList").GetByName("Name")
		log.Info("result", "validatorList.Name", name.Value())
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("author")
		log.Info("result", "validatorList.author", name.Value())
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("Desc").GetByName("name")
		log.Info("result", "validatorList.Desc.name", name.Value())
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("author")
		log.Info("result", "validatorList.author", name.Value())
		name.SetValue("harmony")
		log.Info("result", "validatorList.author", globalVariables.ValidatorList.author)
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("count")
		log.Info("result", "validatorList.count", name.Value())
		name.SetValue(22)
		log.Info("result", "validatorList.count", globalVariables.ValidatorList.count)
	}

	{
		log.Info("compare", "validatorList.validators == nil", globalVariables.ValidatorList.validators == nil)
		log.Info("compare", "len(validatorList.validators)", len(globalVariables.ValidatorList.validators))
		validators := storage.GetByName("ValidatorList").GetByName("validators")
		log.Info("result", "validatorList.validators", validators.Value())
		vv := validators.Value().([]Validator)
		t := Validator{
			desc: Description{
				name: "temp",
				url:  "http://www.hyn.space",
			},
			delegations: nil,
		}
		vv = append(vv, t)
		log.Info("result", "validatorList.validators", vv)
		log.Info("result", "validatorList.validators", globalVariables.ValidatorList.validators)
		validators.GetByIndex(1).SetValue(t)
		log.Info("result", "validatorList.validators", globalVariables.ValidatorList.validators)

		validators.GetByIndex(2).GetByName("desc").GetByName("name").SetValue("haha")
		log.Info("result", "validatorList.validators", globalVariables.ValidatorList.validators)
	}

	{
		globalVariables.ValidatorList.donations = make(map[string]Donation)
		globalVariables.ValidatorList.donations["what"] = Donation{
			Name:   "who-donation",
			amount: 8899,
		}
		donations := storage.GetByName("ValidatorList").GetByName("donations")
		log.Info("result", "validatorList.donations", donations.Value())
		donations.GetByName("what").SetValue(Donation{
			Name:   "who-donation",
			amount: 7788,
		})
		val := donations.GetByName("what").Value().(Donation)
		val.Name = "6688"

		// donations.GetByName("what").GetByName("Name").SetValue("6688")
		// log.Info("result", "validatorList.donations['what'].name", validatorList.donations["what"].Name)
		// m := validatorList.donations["what"]
		// m.Name = "abc"
		// log.Info("result", "validatorList.donations['what'].name", validatorList.donations["what"].Name)
		// vv := donations.Value().(map[string]string)
		// t := Validator{
		//     desc:        Description{
		//         name: "temp",
		//         url: "http://www.hyn.space",
		//     },
		//     delegations: nil,
		// }
		// log.Info("result", "validatorList.validators", vv)
		// log.Info("result", "validatorList.validators", validatorList.validators)
		// donations.GetByIndex(1).SetValue(t)
		// log.Info("result", "validatorList.validators", validatorList.validators)
		//
		// donations.GetByIndex(2).GetByName("desc").GetByName("name").SetValue("haha")
		// log.Info("result", "validatorList.validators", validatorList.validators)
	}
}
