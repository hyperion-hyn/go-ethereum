package restaking

import (
	"github.com/pkg/errors"
	"testing"
)

func TestUpdateDescription(t *testing.T) {
	tests := []struct {
		raw    Description_
		update Description_
		expect Description_
		expErr error
	}{
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{
				Name:            "Jacky",
				Identity:        "jw",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of Jacky",
			},
			expect: Description_{
				Name:            "Jacky",
				Identity:        "jw",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of Jacky",
			},
		},
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{},
			expect: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
		},
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{
				Details: "new details",
			},
			expect: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "new details",
			},
		},
		{
			raw: Description_{
				Name:            "Wayne",
				Identity:        "wen",
				Website:         "harmony.one.wen",
				SecurityContact: "wenSecurity",
				Details:         "wenDetails",
			},
			update: Description_{
				Website: "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongwebsite",
			},
			expErr: errors.New("exceed Maximum Length website"),
		},
	}
	for i, test := range tests {
		err := test.raw.IncrementalUpdateFrom(test.update)
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if err := checkDescriptionEqual(test.raw, test.expect); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func TestDescription_EnsureLength(t *testing.T) {
	tests := []struct {
		desc   Description_
		expErr error
	}{
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: nil,
		},
		{
			desc:   Description_{},
			expErr: nil,
		},
		{
			desc: Description_{
				Name:            "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongname",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed maximum name length"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongidentity",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed Maximum Length identity"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongwebsite",
				SecurityContact: "jacky@harmony.one",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed Maximum Length website"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongcontact",
				Details:         "Details of jacky",
			},
			expErr: errors.New("exceed Maximum Length"),
		},
		{
			desc: Description_{
				Name:            "Jacky Wang",
				Identity:        "jacky@harmony.one",
				Website:         "harmony.one/jacky",
				SecurityContact: "jacky@harmony.one",
				Details:         "thisisaverylonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglonglongdetail",
			},
			expErr: errors.New("exceed Maximum Length for details"),
		},
	}
	for i, test := range tests {
		err := test.desc.EnsureLength()
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
	}
}
