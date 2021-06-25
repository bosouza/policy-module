package policymodule

import data.userResources

allowPolicy {
    # for each o the user's resources
	resource := userResources[input.user][_]

	# for the resources of type "policy"
	resource.id == "policy"
	# check if input user is trying to perform an activity that he's allowed to
	input.activity == resource.activity
}
