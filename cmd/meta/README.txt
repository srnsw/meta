# Meta cmd

Meta cmd provides a simple tool for creating SIPS from the command line.
Designed for simple use cases where most of metadata is in a siegfried (or DROID or fido) file.
Also serves to showcase use of generic loaders and actions available from the meta package.

Simplest usage is just:

`meta siegfried.yaml`

Optional flags can be used to provide additional metadata or options. These are:

    -blacklist  a comma-separated list of IDs to blacklist e.g. x-fmt/111,fmt/10
    -agency     agency ID e.g. 15
    -agencyName agency name e.g. State Archives and Records Authority of NSW
    -series     series ID e.g. 15
    -authority  disposal authority e.g. GA28
    -class      disposal class e.g. 1.1.1
    -access     access direction e.g. 15
    -effect     access direction effect e.g. Early
    -execute    access rule execution date e.g. 2015-01-31
    -output     output directory e.g. c:/users/richardl/Desktop
    -content    content directory e.g. c:/users/richardl/stuff