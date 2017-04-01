package spirali

import (
	"strconv"
	"time"
)

// VersionG represents the version generator of migration.
type VersionG interface {
	GenerateNextVersion() (uint64, error)                        // generates next version.
	IsSmall(targetVersion uint64, comparisonVersion uint64) bool // returns whether targetVersion is smaller than comparisonVersion.
}

// TimestampBasedVersionG is the timestamp based version generator.
type TimestampBasedVersionG struct{}

// GenerateNextVersion ...
func (vg *TimestampBasedVersionG) GenerateNextVersion() (uint64, error) {
	v, err := strconv.ParseUint(time.Now().Format("20060102150405"), 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// IsSmall ...
func (vg *TimestampBasedVersionG) IsSmall(targetVersion uint64, comparisonVersion uint64) bool {
	return targetVersion < comparisonVersion
}

// IncrementalVersionG is the incremental version generator.
type IncrementalVersionG struct {
	CurrentVersion uint64
}

// GenerateNextVersion ...
func (vg *IncrementalVersionG) GenerateNextVersion() (uint64, error) {
	vg.CurrentVersion++
	return vg.CurrentVersion, nil
}

// IsSmall ...
func (vg *IncrementalVersionG) IsSmall(targetVersion uint64, comparisonVersion uint64) bool {
	return targetVersion < comparisonVersion
}
