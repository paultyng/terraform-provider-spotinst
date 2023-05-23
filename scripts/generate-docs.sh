#!/usr/bin/env bash

# Generate the documentation using tfplugindocs and
# remove changes to files that shouldn't change.

# Legacy files that currently shouldn't be autogenerated.
exclude_files=(
	"docs/index.md"
	"docs/resources/elastigroup_aws.md"
	"docs/resources/elastigroup_aws_beanstalk.md"
	"docs/resources/elastigroup_aws_suspension.md"
	"docs/resources/elastigroup_azure.md"
	"docs/resources/elastigroup_gcp.md"
	"docs/resources/elastigroup_gke.md"
	"docs/resources/health_check.md"
	"docs/resources/stateful_node_aws.md"
	"docs/resources/ocean_aks.md"
	"docs/resources/ocean_aks_virtual_node_group.md"
	"docs/resources/mrscaler_aws.md"
	"docs/resources/ocean_aws.md"
	"docs/resources/ocean_aws_launch_spec.md"
	"docs/resources/ocean_ecs.md"
	"docs/resources/ocean_ecs_launch_spec.md"
	"docs/resources/ocean_gke_import.md"
	"docs/resources/ocean_gke_launch_spec.md"
	"docs/resources/ocean_gke_launch_spec_import.md"
	"docs/resources/subscription.md"
	"docs/resources/data_integration.md"
)

# Check if manual changes were made to any excluded files and exit.
# Otherwise, these will be lost with `tfplugindocs`.
git_status="$(git status --porcelain "${exclude_files[@]}")"
if [[ -n "${git_status}" ]]; then
	cat <<EOF
!> Uncommitted changes were detected to the following files.
!> These aren't autogenerated, please commit or stash these changes and try again.
EOF
	echo "${git_status}"
	exit 1
fi

# Generate documentation.
tfplugindocs

# Remove deprecated resources.
rm docs/resources/multai_*

# Remove the changes to files we don't autogenerate.
git checkout HEAD -- "${exclude_files[@]}"
