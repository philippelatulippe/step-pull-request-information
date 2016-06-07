# Github Pull Request Information

Use the Github API to fetch the title and labels of the PR that was just merged.

Use case: Generate a change log for deployment; change deployment method based on label

The following variables are exported to the environment:

* GPRI_PULL_REQUEST_NUMBER: PR number for the current commit
* GPRI_PULL_REQUEST_TITLE
* GPRI_OTHER_PULL_REQUEST_TITLES: Line-break separated list of pull request titles, excluding the current PR
* GPRI_PULL_REQUEST_LABELS: Colon-separated (:) list of github issue labels


## How to use this Step

### On bitrise.io

Bitrise doesn't support adding or configuring third-party steps directly. In
the workflow editor, go to the "bitrise.yml" tab and add the step like this:

    - git::https://github.com/philippelatulippe/step-pull-request-information.git@master:
        title: Fetch Pull Request information
        inputs:
        - github_username: philippelatulippe1

Then add a secret variable `github_access_token` with your github access token,
which you can create on the following page:
https://github.com/settings/tokens

#### On the command line

Can be run directly with the [bitrise CLI](https://github.com/bitrise-io/bitrise),
just `git clone` this repository, `cd` into it's folder in your Terminal/Command Line
and call `bitrise run test`.

*Check the `bitrise.yml` file for required inputs which have to be
added to your `.bitrise.secrets.yml` file!*

Step by step:

1. Open up your Terminal / Command Line
2. `git clone` the repository
3. `cd` into the directory of the step (the one you just `git clone`d)
5. Create a `.bitrise.secrets.yml` file in the same directory of `bitrise.yml` - the `.bitrise.secrets.yml` is a git ignored file, you can store your secrets in
6. Check the `bitrise.yml` file for any secret you should set in `.bitrise.secrets.yml`
  * Best practice is to mark these options with something like `# define these in your .bitrise.secrets.yml`, in the `app:envs` section.
7. Once you have all the required secret parameters in your `.bitrise.secrets.yml` you can just run this step with the [bitrise CLI](https://github.com/bitrise-io/bitrise): `bitrise run test`

An example `.bitrise.secrets.yml` file:

```
envs:
 - MY_STEPLIB_REPO_FORK_GIT_URL: https://github.com/philippelatulippe/steps-pull-request-information
 - GPRI_GITHUB_USERNAME: philippelatulippe
 - GPRI_GITHUB_ACCESS_TOKEN: (Create this token here: https://github.com/settings/tokens)
 - BITRISE_GIT_COMMIT: 6dc69d9fbab1ef190a9f6ea40149b35e225c6d77 (the commit sha of your merge commit)
 - GIT_REPOSITORY_URL: git@github.com:someone/some-repo.git 
```

## How to contribute to this Step

1. Fork this repository
2. `git clone` it
3. Create a branch you'll work on
4. To use/test the step just follow the **How to use this Step** section
5. Do the changes you want to
6. Run/test the step before sending your contribution
  * You can also test the step in your `bitrise` project, either on your Mac or on [bitrise.io](https://www.bitrise.io)
  * You just have to replace the step ID in your project's `bitrise.yml` with either a relative path, or with a git URL format
  * (relative) path format: instead of `- original-step-id:` use `- path::./relative/path/of/script/on/your/Mac:`
  * direct git URL format: instead of `- original-step-id:` use `- git::https://github.com/user/step.git@branch:`
  * You can find more example of alternative step referencing at: https://github.com/bitrise-io/bitrise/blob/master/_examples/tutorials/steps-and-workflows/bitrise.yml
7. Once you're done just commit your changes & create a Pull Request


## Share your own Step

You can share your Step or step version with the [bitrise CLI](https://github.com/bitrise-io/bitrise). Just run `bitrise share` and follow the guide it prints.
