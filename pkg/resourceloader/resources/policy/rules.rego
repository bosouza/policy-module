package policymodule

import data.userResources

allowPolicy {
    # get the user's resources
	resources := userResources[input.user]

	resource = resources[_]

	# for the resources of type "policy"
	resource.id == "policy"
	# check if input user is trying to perform an activity that he's allowed to
	input.activity == resource.activity
}
