package vault

import (
	"reflect"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/roboll/kube-vault-controller/pkg/kube"
)

func init() {

	timeNow = func() time.Time {
		t, _ := time.Parse("2006-01-02 15:04:05", "2017-01-20 01:02:03")
		return t
	}
}

func Test_buildSecretAnnotations(t *testing.T) {
	tests := []struct {
		name   string
		secret *vaultapi.Secret
		claim  *kube.SecretClaim
		want   map[string]string
	}{
		{
			name:   "merge annotations from secretspec",
			secret: &vaultapi.Secret{},
			claim: &kube.SecretClaim{
				Spec: kube.SecretSpec{
					Annotations: map[string]string{
						"hello": "world",
						"foo":   "bar",
					},
				},
			},
			want: map[string]string{
				"hello": "world",
				"foo":   "bar",
				"vaultproject.io/lease-expiration": "1484874123",
				"vaultproject.io/lease-id":         "",
				"vaultproject.io/renewable":        "false",
			},
		},
		{
			name:   "user annotations will not overwrite base annotations",
			secret: &vaultapi.Secret{},
			claim: &kube.SecretClaim{
				Spec: kube.SecretSpec{
					Annotations: map[string]string{
						"vaultproject.io/lease-expiration": "changed",
						"vaultproject.io/lease-id":         "changed",
						"vaultproject.io/renewable":        "changed",
					},
				},
			},
			want: map[string]string{
				"vaultproject.io/lease-expiration": "1484874123",
				"vaultproject.io/lease-id":         "",
				"vaultproject.io/renewable":        "false",
			},
		},
		{
			name:   "no user annotations supplied",
			secret: &vaultapi.Secret{},
			claim: &kube.SecretClaim{
				Spec: kube.SecretSpec{},
			},
			want: map[string]string{
				"vaultproject.io/lease-expiration": "1484874123",
				"vaultproject.io/lease-id":         "",
				"vaultproject.io/renewable":        "false",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildSecretAnnotations(tt.secret, tt.claim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSecretAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}
