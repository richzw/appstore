package appstore

import (
	"reflect"
	"testing"
)

func TestNewCertPool(t *testing.T) {
	tests := []struct {
		name    string
		want    *CertPool
		wantErr bool
	}{
		{"test", &CertPool{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCertPool()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCertPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCertPool() got = %v, want %v", got, tt.want)
			}
		})
	}
}
