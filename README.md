# GMC: GIT Maintenance Complete

Simple terminal command that:

1. Checkout main branch
2. Pull all remote branches
3. Cleanup

## Cleanup step

It force-deletes any local branches that have been deleted on the remote server. This is common when you merge a Pull Request on GitHub/GitLab, delete the branch there, and then want your local machine to reflect that change.

Steps:

1. Fetch and prune
2. Identify branches that no longer exist remotely
3. Delete the local branch if it no longer exists remotely

