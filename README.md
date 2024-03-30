# Description
When using a network file system operator (NFS), there was a problem deleting empty directories. We have a significant 
number of pods (around 500), which mount NFS to themselves and write some statistics there in various formats. There is 
also a sender - it navigates through the directories, collects these statistics, and sends them to a storage system 
(after deleting the files). Ideally, the applications themselves should send the statistics, but this is not feasible 
because any network allocation impacts these applications.

This script traverses the directories, checks for the absence of files, verifies whether the PVC is being used by the 
pod, and then deletes the PVC. The implementation is not perfect, but it functions effectively.

## How it works
The script recursively navigates through mounted NFS directories, ensuring they are empty. Directories are created using 
the following naming format:

```
namespace-pvc-claim-name-pvc-name
```

The process is as follows:
- Mount the script share with all directories.
- Recursively check all subdirectories in these directories for any files.

If any directory contains files, the root directory is considered non-empty.

If there are no files:
- Extract the namespace and PVC name from the directory name, leaving only the PVC name.
- Verify that the claim is not used by any pods.

If the claim is in use, do not delete the PVC. If it is not in use, delete it.

Thus, before removing a PVC, the following conditions must be met:
- The directories for this PVC are completely empty.
- This PVC is not in use by any pods.

## Execution modes
The script can run in two modes:
- `RUN_MODE=job` - Executes the script as a job, running it once.
- `RUN_MODE=always` - Runs the script in a continuous loop, where the time between iterations is regulated through 
`ITERATION_TIME`.

Running in 'always' mode is somewhat impractical due to increased execution time. However, it is useful for 
debugging and for scenarios where a 'set it and forget it' approach is preferred.

The execution time depends on:
- The number of directories and subdirectories.
- The quantity of PVCs.
- The number of PVCs removed.

Preliminary tests show that navigating through 1000 directories and then deleting 500 PVCs takes about 30 seconds. 
Therefore, it is sensible to set the default interval to 5 minutes. Converting this to use goroutines might optimize the
process further.

## Deploy
For the script to function, you need:
- "list" and "delete" permissions on `persistentvolumeclaims` and `pods`.
- To specify the `MOUNT_FOLDER_PATH` environment variable.

## Build
The build is located in `build/Dockerfile`. All dependencies are listed in `go.mod` and `go.sum`.

## Tests
There are currently no tests, but having them would be beneficial.