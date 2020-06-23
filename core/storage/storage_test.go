package storage

import (
    "os"
    "reflect"
    "testing"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/rawdb"
    "github.com/ethereum/go-ethereum/core/state"
    "github.com/ethereum/go-ethereum/log"
)

type Description struct {
    name string
    url  string
}

type Delegation struct {
    amount int
    blockNumber int
}

type Validator struct {
    desc Description
    delegations []Delegation
}

type Donation struct {
    Name string
    amount int
}
type ValidatorList struct {
    Name string
    author string
    count int
    Desc Description
    validators []Validator
    donations map[string]Donation
}


// Tests that storage manipulation
func TestStorageManipulation(t *testing.T) {
    log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

    // Create an empty state database
    db := rawdb.NewMemoryDatabase()
    state, _ := state.New(common.Hash{}, state.NewDatabase(db))
    addr := common.BytesToAddress([]byte{9})

    var validatorList ValidatorList = ValidatorList{
        Name: "atlas",
        author: "hyperion",
        count: 11,
        Desc:       Description{
            name: "hyperion",
            url: "https://www.hyn.space",
        },
        validators: nil,
    }

    reflect.ValueOf(&validatorList).Elem().FieldByName("Name").SetString("harmony")
    reflect.Indirect(reflect.ValueOf(&validatorList)).FieldByName("Name").SetString("harmony")
    // reflect.ValueOf(validatorList).FieldByName("Name").SetString("harmony")

    {
        storage := NewStorage(state, addr, 0, &validatorList, nil)
        log.Info("TestStorageManipulation", "validatorList", validatorList)
        // name := storage.GetByName("validators").GetByName("desc").GetByName("name")
        // name := storage.GetByName("validators").GetByName("name")
        {
            name := storage.GetByName("Name")
            log.Info("result", "validatorList.Name", name.Value())
            name.SetValue("harmony")
            log.Info("result", "validatorList.Name", name.Value())
        }
        // os.Exit(1)
    }

    storage := NewStorage(state, addr, 0, &validatorList, nil)
    log.Info("TestStorageManipulation", "validatorList", validatorList)
    // name := storage.GetByName("validators").GetByName("desc").GetByName("name")
    // name := storage.GetByName("validators").GetByName("name")
    {
        name := storage.GetByName("Name")
        log.Info("result", "validatorList.Name", name.Value())
    }

    {
        name := storage.GetByName("author")
        log.Info("result", "validatorList.author", name.Value())
    }

    {
        name := storage.GetByName("Desc").GetByName("name")
        log.Info("result", "validatorList.Desc.name", name.Value())
    }

    {
        name := storage.GetByName("author")
        log.Info("result", "validatorList.author", name.Value())
        name.SetValue("harmony")
        log.Info("result", "validatorList.author", validatorList.author)
    }

    {
        name := storage.GetByName("count")
        log.Info("result", "validatorList.count", name.Value())
        name.SetValue(22)
        log.Info("result", "validatorList.count", validatorList.count)
    }

    {
        log.Info("compare", "validatorList.validators == nil", validatorList.validators == nil)
        log.Info("compare", "len(validatorList.validators)", len(validatorList.validators))
        validators := storage.GetByName("validators")
        log.Info("result", "validatorList.validators", validators.Value())
        vv := validators.Value().([]Validator)
        t := Validator{
            desc:        Description{
                name: "temp",
                url: "http://www.hyn.space",
            },
            delegations: nil,
        }
        vv = append(vv, t)
        log.Info("result", "validatorList.validators", vv)
        log.Info("result", "validatorList.validators", validatorList.validators)
        validators.GetByIndex(1).SetValue(t)
        log.Info("result", "validatorList.validators", validatorList.validators)

        validators.GetByIndex(2).GetByName("desc").GetByName("name").SetValue("haha")
        log.Info("result", "validatorList.validators", validatorList.validators)
    }

    {
        validatorList.donations=make(map[string]Donation)
        validatorList.donations["what"] = Donation{
            Name:   "who-donation",
            amount: 8899,
        }
        donations := storage.GetByName("donations")
        log.Info("result", "validatorList.donations", donations.Value())
        donations.GetByName("what").SetValue(Donation{
            Name:   "who-donation",
            amount: 7788,
        })
        val := donations.GetByName("what").Value().(Donation)
        val.Name = "6688"

        donations.GetByName("what").GetByName("Name").SetValue("6688")
        log.Info("result", "validatorList.donations['what'].name", validatorList.donations["what"].Name)
        m := validatorList.donations["what"]
        m.Name = "abc"
        log.Info("result", "validatorList.donations['what'].name", validatorList.donations["what"].Name)
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