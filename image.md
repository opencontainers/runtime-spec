# Image

This section defines a format for encoding a container's *filesytem bundle* as a portable *image*.
An image is a set of files organized in a certain way, containing all of the necessary data and metadata, for any compliant runtime to verify the validity of the image's contents and to construct a filesystem bundle represented by the image.
Images are meant to be portable such that they can be used to reliably move filesystem bundles between instances of OCI compliant implementations.

In its simplest form, an image is a tar file of the filesystem bundle (root file system directory plus a configuration file) along with some additional metadata.
While the runtime only sees a single directory as the root filesystem, it might actually be stored as a layering of filesytems under the covers.
When representing this layering within an image, each "layer" will be independently referenced.

## Image Contents

A compliant image contains the following pieces of information:

1. A `manifest.json` file describing the contents of the image.
This file is REQUIRED to be present in the image.
This file extends the [`config.json`](config.md) file with some additional properties:
  * `name`: This REQUIRED property is a URL-like string to aide in discovery of the image.
  For example, `example.com/myApp`.
  * `whiteout`: This OPTIONAL property is a JSON array of absolute file paths that are to be considered to be deleted from dependent layers, even though they are included within those layers root filesystems.
  It is an implementation detail to decide how this whiteout list is processed, but an implementation MUST ensure that any container that is created from this image will behave as if those files were not in the root filesystem at all.
  * `layers`: This OPTIONAL property contains references to the layers that make up the root filesystem of the image.
  This property MUST be a JSON encoded array where each entry consists of a REQUIRED `name` and OPTIONAL `ID` properties.
  The array MUST be in the order in which the images are to be "layed down" into the filesystem bundle that is generated from this image.
  In other words, the first image in the array is written to the filesystem first, and subsequent images are overlayed on top - potentially overlaying an earlier image's files.
  If a referenced layer is embedded within this image then an `ID` property MUST be included.
  If a referenced layer is not embedded within this image then an `ID` property is OPTIONAL.
  * `root/path`: ThiS REQUIRED property, defined within [`config.json`](config.md), is the name of the directory into which all layers MUST appear.
  While the `manifest.json` file extends the `config.json` file, a compliant OCI image is NOT REQUIRED to have all of the required fields from `config.json` within its `manifest.json` file if that image is only used as a layer (dependent image) within other images.
  This is due to the fact that a layer's `manifest.json` properties are not used to produce a container's `config.json` file and therefore not all of them are used.
  However, if an image is used to create a container then all of the "required" fields of the `config.json` MUST be present in the `manifest.json`.

1. Layers.
Each embedded layer MUST have a corresponding image (tar) file, at the root of the image.
Each layer's tar file MUST be named with its ID followed by the appropriate file-extension based on how the tar file is written (e.g. `.tar` for uncompressed).

1. `rootfs` directory.
An image MAY optionally contain a directory that's the same as the `root/path` property of the `manifest.json` file.
And if so, then it MUST be treated as the top-most layer of the root filesystem.

1. There MAY be additional files within the image.
This specification places no requirements on these files.
A compliant OCI implementation is NOT REQUIRED to materialize these files into a container's filesystem bundle directory.

## Serialization of an Image

An image MUST be a `tar` of the files that make up an image as specified in the [Image Contents](#image-contents) section.
The name of the tar file MUST be the ID of the image followed by the appropriate extension based on the type of tar file (e.g. `.tar` for uncompressed tar).
The ID of an image MUST be a cryptographic hash of the uncompressed tar file.
An `ID` MUST be a string of the form `<algorithm>-<hash value>`.
The `algorithm` part of the string MUST be in lower-case.
For example, `sha512` not `SHA512`.
It is RECOMMENDED that implementations use the sha512 algorithm.
Note that if the image does not contain all of its layers as embedded tar files, meaning there are external references, then the ID/hash of the image will not include those referenced layers.

Each file within an image MUST only appear once within the tar file.
Within the tar file, the `manifest.json` file MUST be first, and each embedded image (layer) MUST come in the order specified in the `layers` properties, if present.
If there are additional files within the image they MAY appear in any order within the tar file, except first (as that is reserved for the `manifest.json`).

Each top-level file, or directory, specified in the tar file MUST appear without a leading `/` in its name.

Each file or directory specified in the [Image Contents](#image-contents) section MUST appear at the root of the file file.

## Expanding an Image to Create a Filesystem Bundle

When an image is expanded to create a filesystem bundle the image's `config.json` file MUST be written to the root of the filesystem bundle's directory.
The content of `config.json` are dervived from the `manifest.json` file in the image.
The expansion process MUST create a directory for the container's root filesystem with a name as specified in the `manifest.json` file - ie. the `root/path` property.
The contents of this directory MUST appear to be materialized in such a way as to contain all of the files from the image, without the whiteout (deleted) files, and its dependent images.
The exact mechanism, process and order, by which those files are materialized on disk is an implementation detail, however in the end it MUST appear as though all of the files from the images were written to disk in the order they were specified in the `layers` property.
This means that each dependent image's `whiteout` property MUST appear to be processed during the materialization of that image so that if a subsequent image creates a file by the same name as one mentioned in the `whiteout` property then it is present at the end of the process.

If a dependent image is not present in the image, and the implementation is unable to locate the referenced image by its name or ID, then an error MUST be generated.

The image, and dependent images, MUST be verified by calculating the hash/ID of the materialized (uncompressed) tar file on disk and comparing that value with the ID from the image.
If they do not match then an error MUST be generated.

This specification does not mandate any requirements on the processing of additional files that might appear in the image.
In other words, implementations are not required to materialize them on disk in the filesystem bundle.

This specification does not mandate any ordering to when each piece of information from the image is materialized on disk.

## Misc Notes

There is no requirement that importing an image into a compliant OCI implementation and then exporting it will result in the same image.
Nor is there any requirement that an OCI compliant implementation retain the original layering during such an import/export process.

## Examples

In the following examples, the contents of some of the files within the images are included (in abbreviated form) for clarity.

### Example 1 -  An image with one layer:
```
sha512-abc.tar:
|-- manifest.json
|   |==  {
|   |==    ...config.json properties...
|   |==    "name": "example.com/myTestApp",
|   |==    "root": {
|   |==      "path": "rootfs"
|   |==    },
|   |==    "layers": [
|   |==      { "name": "example.com/myapp",
|   |==        "ID": "sha512-ddd" }
|   |==    ]
|   |==  }
|
|-- sha512-ddd.tar
|   |-- manifest.json
|   |   |== {
|   |   |==   ...config.json properties...
|   |   |==   "name": "example.com/myapp",
|   |   |==   "root": {
|   |   |==     "path": "myRoot"
|   |   |==   }
|   |   |== }
|   |
|   |-- myRoot/
|   |   |-- home/jeff/runit
```


### Example 2 - An image with 4 layers, 2 of which are embedded:
```
sha512-def.tar:
|-- manifest.json
|   |== {
|   |==   ...config.json properties...
|   |==   "name": "example.com/myImage",
|   |==   "root": {
|   |==     "path": "myfs"
|   |==   },
|   |==   "layers": [
|   |==     { "name": "example.com/ubuntu" },
|   |==     { "name": "example.com/ubuntu-utils",
|   |==       "ID": "sha512-bbb" },
|   |==     { "name": "example.com/webserver",
|   |==       "ID": "sha512-ccc" },
|   |==     { "name": "example.com/myapp",
|   |==       "ID": "sha512-ddd" },
|   |==   ]
|   |== }
|
|-- sha512-ccc.tar
|   |-- manifest.json
|   |   |== {
|   |   |==   ...config.json properties...
|   |   |==   "name": "example.com/webserver",
|   |   |==   "root": {
|   |   |==     "path": "rootfs"
|   |   |==   }
|   |   |== }
|   |
|   |-- rootfs/
|   |   |-- usr/bin/...
|
|-- sha512-ddd.tar
|   |-- manifest.json
|   |   |== {
|   |   |==   ...config.json properties...
|   |   |==   "name": "example.com/myapp",
|   |   |==   "root": {
|   |   |==     "path": "myRoot"
|   |   |==   },
|   |   |==   "whiteout": [
|   |   |==     { "/server/webapps/default" }
|   |   |==   ]
|   |   |== }
|   |
|   |-- myRoot/
|   |   |-- home/jeff/runit
```

### Example 3 - An image with 1 layer, but that one layer references other images:

```
sha512-789.tar:
|-- manifest.json
|   |== {
|   |==   ...config.json properties...
|   |==   "name": "example.com/myImage",
|   |==   "root": {
|   |==     "path": "myfs"
|   |==   },
|   |==   "layers": [
|   |==     { "name": "example.com/myapp",
|   |==       "ID": "sha512-ddd" },
|   |==   ]
|   |== }
|
|-- sha512-ccc.tar
|   |-- manifest.json
|   |   |== {
|   |   |==   ...config.json properties...
|   |   |==   "name": "example.com/webserver",
|   |   |==   "root": {
|   |   |==     "path": "rootfs"
|   |   |==   },
|   |   |==   "layers": [
|   |   |==     { "name": "example.com/ubuntu-utils",
|   |   |==       "ID": "sha512-bbb" }
|   |   |==   ]
|   |   |== }
|   |
|   |-- rootfs/
|   |   |-- usr/bin/...
|
|-- sha512-ddd.tar
|   |-- manifest.json
|   |   |== {
|   |   |==   ...config.json properties...
|   |   |==   "name": "example.com/myapp",
|   |   |==   "root": {
|   |   |==     "path": "myRoot"
|   |   |==   },
|   |   |==   "whiteout": [
|   |   |==     { "/server/webapps/default" }
|   |   |==   ],
|   |   |==   "layers": [
|   |   |==     { "name": "example.com/webserver",
|   |   |==       "ID": "sha512-ccc" }
|   |   |==   ],
|   |   |== }
|   |
|   |-- myRoot/
|   |   |-- home/jeff/runit
```

### Example 4 - Same as previous but all are included in the image.
```
sha512-234.tar
|-- manifest.json
|   |== {
|   |==   ...config.json properties...
|   |==   "name": "example.com/myImage",
|   |==   "root": {
|   |==     "path": "myfs"
|   |==   },
|   |==   "layers": [
|   |==    { "name": "example.com/webserver",
|   |==      "ID": "sha512-ccc" },
|   |==    { "name": "example.com/myapp",
|   |==      "ID": "sha512-ddd" },
|   |==   ]
|   |== }
|
|-- sha512-ccc.tar
|   |-- manifest.json
|   |   |== {
|   |   |==   ...config.json properties...
|   |   |==   "name": "example.com/webserver",
|   |   |==   "root": {
|   |   |==     "path": "rootfs"
|   |   |==   },
|   |   |==   "layers": [
|   |   |==   { "name": "example.com/ubuntu-utils",
|   |   |==     "ID": "sha512-bbb" }
|   |   |==   ]
|   |   |== }
|   |
|   |-- sha512-bbb.tar
|   |   |-- manifest.json
|   |   |   |== {
|   |   |   |==   ...config.json properties...
|   |   |   |==   "name": "example.com/ubuntu-utils"
|   |   |   |==   "root": {
|   |   |   |==     "path": "rootfs"
|   |   |   |==   },
|   |   |   |==   "layers": [
|   |   |   |==   { "name": "example.com/ubuntu",
|   |   |   |==     "ID": "sha512-aaa" }
|   |   |   |==   ]
|   |   |   |== }
|   |   |
|   |   |-- sha512-aaa.tar
|   |   |   |-- manifest.json
|   |   |   |   |== {
|   |   |   |   |==   ...config.json properties...
|   |   |   |   |==   "name": "example.com/ubuntu",
|   |   |   |   |==   "root": {
|   |   |   |   |==     "path": "rootfs"
|   |   |   |   |==   }
|   |   |   |   |== }
|   |   |   |
|   |   |   |-- rootfs/
|   |   |   |   |-- var/...
|   |   |
|   |   |-- rootfs/
|   |   |   |-- opt/...
|   |
|   |-- rootfs/
|   |   |-- usr/bin/...
|
|-- sha512-ddd.tar
|   |-- manifest.json
|   |   |== {
|   |   |==   ...config.json properties...
|   |   |==   "name": "example.com/myapp",
|   |   |==   "root": {
|   |   |==     "path": "myRoot"
|   |   |==   },
|   |   |==   "whiteout": [
|   |   |==     { "/server/webapps/default" }
|   |   |==   ]
|   |   |== }
|   |
|   |-- myRoot/
|   |   |-- home/jeff/runit
```
