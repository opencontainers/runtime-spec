# Image

This section defines a format for encoding a container's *filesytem bundle* as a portable *image*.
An image is a set of files organized in a certain way, and containing all of the necessary data and metadata, for any compliant runtime to verify the validity of the image's contents and to construct a filesystem bundle represented by the image.
Images are meant to be portable such that they can be used to reliably move filesystem bundles between instances of OCI compliant implementations.

## Image Contents

A compliant image contains the following pieces of information:

1. An OCI [`config.json`](config.md) file.
This file is REQUIRED to be present in the image.
The `config.json` file contains the host independent configuration information necessary to execute the container.
See [`config.json`](config.md) for more information.

1. References to dependent images.
If present, this OPTIONAL information MUST be in a file called `required.json`.
This file MUST be a JSON encoded array of references to images as specified in the [Dependent Images](#dependent-images) section.

1. Dependent images.
An image MAY include dependent images as directories within the image.
Each dependent image MUST be placed in a directory whose name is the same as its ID (see the [IDs](#ids) section for more details about IDs).
The contents of each directory MUST adhere to the following:
  * there MAY be a `whiteout.json` file (see [Whiteout List](#whiteout-list)).
  * there MUST be a directory called `rootfs` which contains the root filesystem of that image.
  *  there MAY be additional "extension files".

  The ID MUST be generated over the above contents in the order specified.

1. Whiteout list.
An image MAY choose to include a list of files that are to be considered to be deleted from dependent images even though they are included within the dependent image's root filesystem directory.
It is an implementation detail to decide how this whiteout list is processed but an implementation MUST ensure that any container that is created from this image will behave the same as if those files were not in the root filesystem at all.
The whiteout list MUST be in a file called `whiteout.json`.
See the [Whiteout List](#whiteout-list) section for more details.

1. The root filesystem of the filesystem bundle.
This directory MUST be present in the image, however it MAY be empty (not contain any files or folders).
As specified in the `config.json` file, this directory contains the files that make up the root filesystem of the container.
This directory MAY contain only a portion of the complete set of files that make up the filesystem bundle, in which case the remaining set of files MUST be found in one of the dependent images.

1. ID of the root filesystem.
This data MUST be a cryptographic hash of the whiteout list file, if present, and the root filesystem directory (including the directory itself), in that order.
If there are dependent images then this MUST NOT include the dependent image's data.
This REQUIRED information MUST be in a file called `rootfs.ID`
See the [IDs](#ids) section for more details.

1. Extension files.
An image MAY choose to include additional files and directories in the image.

1. The ID of the image.
This data MUST be a cryptographic hash of all the data listed above, in the order specified.
For the purposes of calculating the ID, the order of processing the extension files MUST be in case-sensitive alphabetical order.
This REQUIRED information MUST be in a file called `ID`.
See the [IDs](#ids) section for more details.

## Serialization of an Image

An image MUST be a `tar` of the pieces of data that make up an image as specified in the [Image Contents](#image-contents) section.
The order in which the data is serialized in the `tar` MUST be the same as the order specified in that section.
The serialization of the "dependent images" directories, if any, MUST be in the same order in which those dependent images are listed in the `requires.json` file.
The serialization of the  "extension files" MUST be in case-sensitive alphabetical order.
Note that "extension files", irrespective of their names, MUST appear after the `rootfs.ID` file for the image.

Each file or directory specified in the [Image Contents](#image-contents) section MUST appear at the root of the `tar` without a leading `/` in their names.

For example, a listing of the `tar` of an image might look like this:
```
config.json
requires.json
sha512-1d62d181e9c1322d56ccd3a29d05018399147a16188dbd861c0279ad0dd7e14c/
    rootfs/
        bin/runit
sha512-291e2a171ef9bef8a838c59406d9b0aeb6f2f0ebe5173415205733d3d18b8e03/
    whiteout.json
    rootfs/
        bin/monitor
sha512-98364ca873540185d83645e93745a94847f2f0ebe51733987154840afebc9921/
whiteout.json
rootfs/
    bin/myapp
rootfs.ID
extraFile
ID
```

## Dependent Images

If an image has dependent images then it MUST include a `requires.json` file that references those images.
The `requires.json` file MUST be a JSON encoded array matching this format:
```
{
  "requires": [
    {
      "ID": "..."
    },
    {
      "ID": "..."
    }, ...
  ]
}
```

The order of the images, `ID`s, in the array MUST be in the order in which the images are to be "layed down" into the filesystem bundle that is generated from this image.
In other words, the first image in the array is written to the filesystem first, and subsequent images are overlayed on top - potentially overlaying an ealier image's files.

## IDs

An `ID` MUST be a cryptographic hash of the data being protected, in the order specified for that ID.

An `ID` MUST be a string of the form `<algorithm>-<hash value>`.
The `algorithm` part of the string MUST be in lower-case.
For example, `sha512` not `SHA512`.

It is RECOMMENDED that implementations use the sha512 algorithm.

## Whiteout List

An image's `whiteout.json` file MUST be a JSON encoded file matching this format:
```
{
  "files": [
    "filename",
    ...
  ]
}
```
where each `filename` MUST be the absolute path to the file.

## Expanding an Image to Create a Filesystem Bundle

When an image is expanded to create a filesystem bundle the image's `config.json` file MUST be written to the root of the filesystem bundle's directory.
The expansion process MUST create a directory for the container's root filesystem with a name as specified in the `config.json` file - ie. the `root/path` property.
The contents of this directory MUST appear to be materizlied in such a way as to contain all of the files from the image, without the whiteout (deleted) files, and its dependent images.
The exact mechanism, process and order, by which those files are materialized on disk is an implementation detail, however in the end it MUST appear as though all of the files from the images were written to disk in the order they were specified in the `requires.json` file.
This means that each dependent image's `whiteout.json` file MUST appear to be processed during the materialization of that image so that if a subsequent image creates a file by the same name as one mentioned in the `whiteout.json` file then it will not be deleted.

If a dependent image is not present in the image, and the implementation is unable to locate the referenced image by its ID, then an error MUST be generated.

The image, and dependent images, MUST be verified by calculating the hash/ID of the materialized image on disk and comparing that value with the ID from the image.
If they do not match then an error MUST be generated.

This specification does not mandate any requirements on the processing of the extension files.
In other words, implementations are not required to materialize them on disk in the filesystem bundle.

This specification does not mandate any ordering to when each piece of information from the image is materialized.

## Misc Notes

There is no requirement that importing an image into a compliant OCI implementation and then exporting it will result in the same image since each implementation might choose to create the root filesystem differently.
Thus, exporting that root filesystem might look different from the original image.
