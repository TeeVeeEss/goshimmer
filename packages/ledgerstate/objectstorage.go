package ledgerstate

import (
	"time"

	"github.com/iotaledger/hive.go/objectstorage"
)

const (
	// PrefixBranchStorage defines the storage prefix for the Branch object storage.
	PrefixBranchStorage byte = iota

	// PrefixChildBranchStorage defines the storage prefix for the ChildBranch object storage.
	PrefixChildBranchStorage

	// PrefixConflictStorage defines the storage prefix for the Conflict object storage.
	PrefixConflictStorage

	PrefixConflictMemberStorage
)

// branchStorageOptions contains a list of default settings for the Branch object storage.
var branchStorageOptions = []objectstorage.Option{
	objectstorage.CacheTime(60 * time.Second),
}

// childBranchStorageOptions contains a list of default settings for the ChildBranch object storage.
var childBranchStorageOptions = []objectstorage.Option{
	ChildBranchKeyPartition,
	objectstorage.CacheTime(60 * time.Second),
}

// conflictStorageOptions contains a list of default settings for the Conflict object storage.
var conflictStorageOptions = []objectstorage.Option{
	objectstorage.CacheTime(60 * time.Second),
}

// conflictMemberStorageOptions contains a list of default settings for the ConflictMember object storage.
var conflictMemberStorageOptions = []objectstorage.Option{
	ConflictMemberKeyPartition,
	objectstorage.CacheTime(60 * time.Second),
}