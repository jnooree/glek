glek
====

glek &mdash; command line app to export GitHub issue labels into [gembel](https://github.com/gedex/gembel)
JSON format.

## Install

### From brew

```
brew install gedex/tap/glek
```

Check the [tap source](https://github.com/gedex/homebrew-tap) for more details.

### From binaries

Download your preferred flavor from the [releases page](https://github.com/gedex/glek/releases/latest) and install manually.

### From Go Get

```
go get github.com/gedex/glek
```

## Using glek

Before using glek, you need `GITHUB_TOKEN` (can be retrieved from [here](https://github.com/settings/tokens)).
Once you've that, set it to your bash profile or provide it when running the app:

```
GITHUB_TOKEN="token" glek <owner/repo>
```

Most of the time you will pipe the output to a file and then fed that into gembel:

```
$ glek Automattic/wp-calypso > labels.json
$ gembel labels.json
```

You need to edit `labels.json` first to adjust your label replacements (if any) and
target repositories.
