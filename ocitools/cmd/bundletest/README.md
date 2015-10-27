#Bundle Test Tool
      
###Build
Since it is just a part of the tools, the Godeps directory is not provided yet.
So far you can you `make` to compile it.
```
make     
```
      
###Verify a bundle
It verifies whether a bundle is valid, with all the required files and
all the required attributes.
```
./bundle vb demo-bundle
./bundle vc demo-bundle/config.json
./bundle vc demo-bundle/config.json.bad
./bundle vr demon-bundle/runtime.json
./bundle vr demon-bundle/runtime.json.bad
```

###Validate once, return all the errors
The return value '(msgs []string, valid bool)' will store all the error messages.
Correct all of them before run an OCI bundle.

```
The mountPoint sys /sys is not exist in rootfs
The mountPoint proc /proc is not exist in rootfs
The mountPoint dev /dev is not exist in rootfs
```
