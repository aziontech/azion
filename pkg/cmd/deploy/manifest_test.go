package deploy

import "testing"

func Test_prepareRequestRuleEngine(t *testing.T) {
	type args struct {
		applicationID int64
		cacheID       int64
		template      string
		mode          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := prepareRequestRuleEngine(tt.args.applicationID, tt.args.cacheID, tt.args.template, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("prepareRequestRuleEngine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
