package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	"github.com/hussnain612/uptime-robot-operator/pkg/constants"
	"k8s.io/apimachinery/pkg/types"
)

// Get ApiKey from secret
func GetApiKeyFromSecret(ctx context.Context, client client.Client, log logr.Logger) (string, error) {
	secret := corev1.Secret{}

	// TODO update namespace name
	err := client.Get(ctx, types.NamespacedName{Name: constants.DefaultSecretName, Namespace: "default"}, &secret)
	if err != nil {
		log.Error(err, "failed to get secret for api key")
		return "", err
	}

	apiKey, ok := secret.Data[constants.SecretConfigKey]
	if !ok {
		err := fmt.Errorf("secret %s did not contain key %s", constants.DefaultSecretName, constants.SecretConfigKey)
		log.Error(err, "failed to get key from secret")
		return "", err
	}

	return string(apiKey), nil
}

// handleAPIRequestForUptimeRobot handles the creation and making of api requests
func HandleAPIRequestForUptimeRobot(log logr.Logger, apiURL, path, apiKey string, parameters map[string]string) ([]byte, error) {
	var body []byte
	url := apiURL + path

	payloadStr := fmt.Sprintf("api_key=%s&format=json", apiKey)
	for key, val := range parameters {
		payloadStr = payloadStr + "&" + key + "=" + val
	}

	payload := strings.NewReader(payloadStr)

	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Error(err, "failed to create a POST request")
		return body, err
	}

	httpReq.Header.Add("cache-control", "no-cache")
	httpReq.Header.Add("content-type", "application/x-www-form-urlencoded")

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Error(err, "failed to make a POST request")
		return body, err
	}

	defer httpRes.Body.Close()

	body, err = io.ReadAll(httpRes.Body)
	if err != nil {
		log.Error(err, "failed to read the response body")
		return body, err
	}

	return body, nil
}
