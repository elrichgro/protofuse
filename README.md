# protofuse

Protofuse is a tool that mounts marshaled protocol buffers as a filesystem for easier viewing and debugging.

To use protofuse:

`$ protofuse 'path of mount location' 'marshaled protocol buffer' 'path to .proto file' 'package name' 'message name'`

protofuse/mount/mount.go also contains functions

`Mount(marshaled []byte, fileDesc *google_protobuf.FileDescriptorSet, packageName string, messageName string, mountPoint string) error`

and

`MountList(marshaled [][]byte, fileDesc *google_protobuf.FileDescriptorSet, messageName string, mountPoint string) error`

that can be used to mount protocol buffers.

`marshaled` is a marshaled protocol buffer or a slice of marshaled protocol buffers

`fileDesc` is a FileDescriptorSet describing the proto files

`packageName` is the name of the package of the top-level message

`messageName` is the name of the top-level message

`mountPoint` is the path to the location to mount the filesystem


