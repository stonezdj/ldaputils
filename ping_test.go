package main

import "testing"

func TestPingURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"normal fail", args{"10.10.10.10:389"}, false, true},
		{"normal fail2", args{"ldap://10.10.10.10"}, false, true},
		{"normal test", args{"www.baidu.com:80"}, true, false},
		{"pure IP", args{"10.202.250.197"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PingURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("PingURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PingURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
