package k8s

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetSecretData(ctx context.Context, client kubernetes.Interface, namespace string, name string, keys []string) (map[string][]byte, error) {
	data := make(map[string][]byte)

	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return data, nil
		}
		return data, err
	}

	for _, key := range keys {
		data[key] = secret.Data[key]
	}

	return data, nil
}
