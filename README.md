# dvcm

```dvcm``` (**d**e**vc**ontainer **m**anager) is a tool meant to help working with multiple devcontainer configurations.
It allows saving and loading devcontainer configurations setups from and to different sources.

Currently supported sources:

* (Local) file system
* Github
* Gitlab

## Usage

### Setup

All credentails needed to access remote sources are stored in a config.json file.
An example structure is shown in [here](config.example.json)

Using Github requires the permissions in ```repo.public_repo``` (for read & write to public repositories) or the entire ```repo``` set (for read & write to public and private repositories) for classic access tokens.

Using Gitlab requires the permission ```api``` for legacy tokens (for read & write).

### Commands

```save``` and ```load``` have the same set of expected flags which are the following:

```bash
-origin string
    directory where to load from/save to
-r    toggle remote mode
-workspace string
    directory where to load to/save from (default ".")
```

and should be followed by the expected name of the folder (to be) containing the devcontainer.json and further configuration files.

Saving a configuration to a Github remote repository might thus look like:


```bash
./dvcm -r -origin "https://github.com/<account-name>/<repository-name>" save go
```

and loading form a Gitlab remote repository like:

```bash
./dvcm -r -origin="https://gitlab.com/api/v4/projects/<project-id>" load dart
```

Using the dvcm configuration and setting defaults the commands may for example be shortened to:

```bash
./dvcm -r load go
```

or

```bash
./dvcm -r save go
```