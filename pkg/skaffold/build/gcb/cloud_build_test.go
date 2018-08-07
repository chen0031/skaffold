/*
Copyright 2018 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gcb

import (
	"testing"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1alpha2"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/util"
	"github.com/GoogleContainerTools/skaffold/testutil"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

func TestBuildDescription(t *testing.T) {
	artifact := &v1alpha2.Artifact{
		ImageName: "nginx",
		ArtifactType: v1alpha2.ArtifactType{
			DockerArtifact: &v1alpha2.DockerArtifact{
				DockerfilePath: "Dockerfile",
				BuildArgs: map[string]*string{
					"arg1": util.StringPtr("value1"),
					"arg2": nil,
				},
			},
		},
	}

	builder := Builder{
		GoogleCloudBuild: &v1alpha2.GoogleCloudBuild{
			DockerImage: "docker/docker",
			DiskSizeGb:  100,
			MachineType: "n1-standard-1",
			Timeout:     "10m",
		},
	}
	desc := builder.buildDescription(artifact, "bucket", "object")

	expected := cloudbuild.Build{
		LogsBucket: "bucket",
		Source: &cloudbuild.Source{
			StorageSource: &cloudbuild.StorageSource{
				Bucket: "bucket",
				Object: "object",
			},
		},
		Steps: []*cloudbuild.BuildStep{{
			Name: "docker/docker",
			Args: []string{"build", "--tag", "nginx", "-f", "Dockerfile", "--build-arg", "arg1=value1", "--build-arg", "arg2", "."},
		}},
		Images: []string{artifact.ImageName},
		Options: &cloudbuild.BuildOptions{
			DiskSizeGb:  100,
			MachineType: "n1-standard-1",
		},
		Timeout: "10m",
	}

	testutil.CheckErrorAndDeepEqual(t, false, nil, expected, *desc)
}
