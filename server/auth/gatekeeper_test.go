package auth

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"gopkg.in/square/go-jose.v2/jwt"

	fakewfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned/fake"
	"github.com/argoproj/argo/server/auth/sso/mocks"
)

func TestServer_GetWFClient(t *testing.T) {
	wfClient := &fakewfclientset.Clientset{}
	kubeClient := fake.NewSimpleClientset(
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{Name: "my-sa", Namespace: "my-ns"},
			Secrets:    []corev1.ObjectReference{{Name: "my-secret"}},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "my-secret", Namespace: "my-ns"},
			Data: map[string][]byte{
				"token": {},
			},
		},
	)
	t.Run("None", func(t *testing.T) {
		_, err := NewGatekeeper(Modes{}, wfClient, kubeClient, nil, nil, "")
		assert.Error(t, err)
	})
	t.Run("Invalid", func(t *testing.T) {
		g, err := NewGatekeeper(Modes{Client: true}, wfClient, kubeClient, nil, nil, "")
		if assert.NoError(t, err) {
			_, err := g.Context(x("invalid"))
			assert.Error(t, err)
		}
	})
	t.Run("NotAllowed", func(t *testing.T) {
		g, err := NewGatekeeper(Modes{SSO: true}, wfClient, kubeClient, nil, nil, "")
		if assert.NoError(t, err) {
			_, err := g.Context(x("Bearer "))
			assert.Error(t, err)
		}
	})
	// not possible to unit test client auth today
	t.Run("Server", func(t *testing.T) {
		g, err := NewGatekeeper(Modes{Server: true}, wfClient, kubeClient, &rest.Config{Username: "my-username"}, nil, "")
		assert.NoError(t, err)
		ctx, err := g.Context(x(""))
		if assert.NoError(t, err) {
			assert.Equal(t, wfClient, GetWfClient(ctx))
			assert.Equal(t, kubeClient, GetKubeClient(ctx))
			assert.NotNil(t, GetClaims(ctx))
		}
	})
	t.Run("SSO", func(t *testing.T) {
		ssoIf := &mocks.Interface{}
		ssoIf.On("Authorize", mock.Anything, mock.Anything).Return(&jwt.Claims{Subject: "my-sub"}, nil)
		ssoIf.On("IsRBACEnabled").Return(false)
		g, err := NewGatekeeper(Modes{SSO: true}, wfClient, kubeClient, nil, ssoIf, "my-ns")
		if assert.NoError(t, err) {
			ctx, err := g.Context(x("Bearer v2:whatever"))
			if assert.NoError(t, err) {
				assert.Equal(t, wfClient, GetWfClient(ctx))
				assert.Equal(t, kubeClient, GetKubeClient(ctx))
				if assert.NotNil(t, GetClaims(ctx)) {
					assert.Equal(t, "my-sub", GetClaims(ctx).Subject)
				}
			}
		}
	})
	t.Run("SSO+RBAC", func(t *testing.T) {
		ssoIf := &mocks.Interface{}
		ssoIf.On("Authorize", mock.Anything, mock.Anything).Return(jws.ClaimSet{"groups": []string{"my-group"}}, nil)
		ssoIf.On("IsRBACEnabled").Return(true)
		g, err := NewGatekeeper(Modes{SSO: true}, wfClient, kubeClient, nil, ssoIf, "my-ns")
		if assert.NoError(t, err) {
			_, err := g.Context(x("Bearer id_token:"))
			assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = not allowed")
		}
	})
}

func x(authorization string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"authorization": authorization}))
}

func TestGetClaimSet(t *testing.T) {
	// we should be able to get nil claim set
	assert.Nil(t, GetClaims(context.TODO()))
}
