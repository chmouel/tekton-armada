#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../vendor/knative.dev/hack/codegen-library.sh
export PATH="$GOBIN:$PATH"

function run_yq() {
	run_go_tool github.com/mikefarah/yq/v4@v4.23.1 yq "$@"
}

echo "=== Update Codegen for ${MODULE_NAME}"

group "Kubernetes Codegen"
${CODEGEN_PKG}/generate-groups.sh "deepcopy,client,informer,lister" \
	github.com/chmouel/armadas/pkg/client github.com/chmouel/armadas/pkg/apis \
	"armadas:v1alpha1" \
	--go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Knative Codegen"
${KNATIVE_CODEGEN_PKG}/hack/generate-knative.sh "injection" \
	github.com/chmouel/armadas/pkg/client github.com/chmouel/armadas/pkg/apis \
	"armadas:v1alpha1" \
	--go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Update CRD Schema"
go run ./hack/schema/main.go dump Fire |
	run_yq eval-all --header-preprocess=false --inplace 'select(fileIndex == 0).spec.versions[0].schema.openAPIV3Schema = select(fileIndex == 1) | select(fileIndex == 0)' \
		$(dirname $0)/../config/300-crd-fire.yaml -

# group "Update deps post-codegen"
#
# # Make sure our dependencies are up-to-date
# ${REPO_ROOT_DIR}/hack/update-deps.sh
