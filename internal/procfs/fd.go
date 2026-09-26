package procfs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ResolvePIDs(targetInodes map[uint64]struct{}) map[uint64]int {
	inodeToPID := make(map[uint64]int)
	if len(targetInodes) == 0 {
		return inodeToPID
	}

	procDirs, err := os.ReadDir("/proc")
	if err != nil {
		return inodeToPID
	}

	for _, p := range procDirs {
		if !p.IsDir() {
			continue
		}

		pidStr := p.Name()
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		fdDirPath := filepath.Join("/proc", pidStr, "fd")
		fdFiles, err := os.ReadDir(fdDirPath)
		if err != nil {
			continue
		}

		for _, fdFile := range fdFiles {
			fdPath := filepath.Join(fdDirPath, fdFile.Name())
			link, err := os.Readlink(fdPath)
			if err != nil {
				continue
			}

			if strings.HasPrefix(link, "socket:[") && strings.HasSuffix(link, "]") {
				inodeStr := link[8 : len(link)-1]
				inode, err := strconv.ParseUint(inodeStr, 10, 64)
				if err == nil {
					if _, ok := targetInodes[inode]; ok {
						inodeToPID[inode] = pid
						delete(targetInodes, inode)
						if len(targetInodes) == 0 {
							return inodeToPID
						}
					}
				}
			}
		}
	}

	return inodeToPID
}
