// Copyright (c) 2019-2021 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package automount

import (
	"context"
	"fmt"

	"github.com/devfile/devworkspace-operator/apis/controller/v1alpha1"
	"github.com/devfile/devworkspace-operator/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// provisionGitConfiguration takes care of mounting git credentials and a gitconfig into a devworkspace.
func provisionGitConfiguration(client k8sclient.Client, namespace string) (*v1alpha1.PodAdditions, error) {
	secrets := &corev1.SecretList{}
	err := client.List(context.TODO(), secrets, k8sclient.InNamespace(namespace), k8sclient.MatchingLabels{
		constants.DevWorkspaceGitCredentialLabel: "true",
	})
	if err != nil {
		return nil, err
	}
	var credentials []string
	userMountPath := ""
	for _, secret := range secrets.Items {
		credentials = append(credentials, string(secret.Data[gitCredentialsName]))
		if val, ok := secret.Annotations[constants.DevWorkspaceMountPathAnnotation]; ok {
			if userMountPath != "" && val != userMountPath {
				return nil, &FatalError{fmt.Errorf("auto-mounted git credentials have conflicting mountPaths")}
			}
			userMountPath = val
		}
	}

	podAdditions := &v1alpha1.PodAdditions{}

	// Grab the gitconfig additions
	gitConfigAdditions, err := provisionGitConfig(client, namespace, userMountPath)
	if err != nil {
		return podAdditions, err
	}

	podAdditions.Volumes = append(podAdditions.Volumes, gitConfigAdditions.Volumes...)
	podAdditions.VolumeMounts = append(podAdditions.VolumeMounts, gitConfigAdditions.VolumeMounts...)

	// Grab the credentials additions
	if len(credentials) > 0 {
		credentialsAdditions, err := provisionUserGitCredentials(client, namespace, userMountPath, credentials)
		if err != nil {
			return podAdditions, err
		}
		podAdditions.Volumes = append(podAdditions.Volumes, credentialsAdditions.Volumes...)
		podAdditions.VolumeMounts = append(podAdditions.VolumeMounts, credentialsAdditions.VolumeMounts...)
	}
	return podAdditions, nil
}
