package ledgerstate

import (
	"fmt"
	"testing"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchDAG_ConflictBranches(t *testing.T) {
	branchDAG := NewBranchDAG(mapdb.NewMapDB())
	defer branchDAG.Shutdown()

	conflictBranch, newBranchCreated, err := branchDAG.RetrieveConflictBranch(
		NewBranchID(TransactionID{3}),
		NewBranchIDs(
			NewBranchID(TransactionID{1}),
		),
		NewConflictIDs(
			NewConflictID(NewOutputID(TransactionID{2}, 0)),
			NewConflictID(NewOutputID(TransactionID{2}, 1)),
		),
	)
	require.NoError(t, err)
	defer conflictBranch.Release()
	fmt.Println(conflictBranch, newBranchCreated)

	conflictBranch1, newBranchCreated1, err := branchDAG.RetrieveConflictBranch(
		NewBranchID(TransactionID{3}),
		NewBranchIDs(
			NewBranchID(TransactionID{1}),
		),
		NewConflictIDs(
			NewConflictID(NewOutputID(TransactionID{2}, 0)),
			NewConflictID(NewOutputID(TransactionID{2}, 1)),
			NewConflictID(NewOutputID(TransactionID{2}, 2)),
		),
	)
	require.NoError(t, err)
	defer conflictBranch1.Release()
	fmt.Println(conflictBranch1, newBranchCreated1)
}

func TestBranchDAG_normalizeBranches(t *testing.T) {
	branchDAG := NewBranchDAG(mapdb.NewMapDB())

	cachedBranch2, newBranchCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{2}, NewBranchIDs(MasterBranchID), NewConflictIDs(ConflictID{0}))
	defer cachedBranch2.Release()
	branch2 := cachedBranch2.Unwrap()
	assert.True(t, newBranchCreated)

	cachedBranch3, newBranchCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{3}, NewBranchIDs(MasterBranchID), NewConflictIDs(ConflictID{0}))
	defer cachedBranch3.Release()
	branch3 := cachedBranch3.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizedBranches, err := branchDAG.normalizeBranches(NewBranchIDs(MasterBranchID, branch2.ID()))
		assert.NoError(t, err)
		assert.Equal(t, normalizedBranches, NewBranchIDs(branch2.ID()))

		normalizedBranches, err = branchDAG.normalizeBranches(NewBranchIDs(MasterBranchID, branch3.ID()))
		assert.NoError(t, err)
		assert.Equal(t, normalizedBranches, NewBranchIDs(branch3.ID()))

		normalizedBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch2.ID(), branch3.ID()))
		assert.Error(t, err)
	}

	// spawn of branch 4 and 5 from branch 2
	cachedBranch4, newBranchCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{4}, NewBranchIDs(branch2.ID()), NewConflictIDs(ConflictID{1}))
	defer cachedBranch4.Release()
	branch4 := cachedBranch4.Unwrap()
	assert.True(t, newBranchCreated)

	cachedBranch5, newBranchCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{5}, NewBranchIDs(branch2.ID()), NewConflictIDs(ConflictID{1}))
	defer cachedBranch5.Release()
	branch5 := cachedBranch5.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizedBranches, err := branchDAG.normalizeBranches(NewBranchIDs(MasterBranchID, branch4.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch4.ID()), normalizedBranches)

		normalizedBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch3.ID(), branch4.ID()))
		assert.Error(t, err)

		normalizedBranches, err = branchDAG.normalizeBranches(NewBranchIDs(MasterBranchID, branch5.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch5.ID()), normalizedBranches)

		normalizedBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch3.ID(), branch5.ID()))
		assert.Error(t, err)

		// since both consume the same output
		normalizedBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch4.ID(), branch5.ID()))
		assert.Error(t, err)
	}

	// branch 6, 7 are on the same level as 2 and 3 but are not part of that conflict set
	cachedBranch6, newBranchCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{6}, NewBranchIDs(MasterBranchID), NewConflictIDs(ConflictID{2}))
	defer cachedBranch6.Release()
	branch6 := cachedBranch6.Unwrap()
	assert.True(t, newBranchCreated)

	cachedBranch7, newBranchCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{7}, NewBranchIDs(MasterBranchID), NewConflictIDs(ConflictID{2}))
	defer cachedBranch7.Release()
	branch7 := cachedBranch7.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizeBranches, err := branchDAG.normalizeBranches(NewBranchIDs(branch2.ID(), branch6.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch2.ID(), branch6.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch3.ID(), branch6.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch3.ID(), branch6.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch2.ID(), branch7.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch2.ID(), branch7.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch3.ID(), branch7.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch3.ID(), branch7.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch6.ID(), branch7.ID()))
		assert.Error(t, err)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch4.ID(), branch6.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch4.ID(), branch6.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch5.ID(), branch6.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch5.ID(), branch6.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch4.ID(), branch7.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch4.ID(), branch7.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch5.ID(), branch7.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch5.ID(), branch7.ID()), normalizeBranches)
	}

	// aggregated branch out of branch 4 (child of branch 2) and branch 6
	cachedAggrBranch8, newBranchCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(branch4.ID(), branch6.ID()))
	assert.NoError(t, aggrBranchErr)
	defer cachedAggrBranch8.Release()
	aggrBranch8 := cachedAggrBranch8.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizeBranches, err := branchDAG.normalizeBranches(NewBranchIDs(aggrBranch8.ID(), MasterBranchID))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch4.ID(), branch6.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch8.ID(), branch2.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch4.ID(), branch6.ID()), normalizeBranches)

		// conflicting since branch 2 and branch 3 are
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch8.ID(), branch3.ID()))
		assert.Error(t, err)

		// conflicting since branch 4 and branch 5 are
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch8.ID(), branch5.ID()))
		assert.Error(t, err)

		// conflicting since branch 6 and branch 7 are
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch8.ID(), branch7.ID()))
		assert.Error(t, err)
	}

	// aggregated branch out of aggr. branch 8 and branch 7:
	// should fail since branch 6 & 7 are conflicting
	_, newBrachCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(aggrBranch8.ID(), branch7.ID()))
	assert.Error(t, aggrBranchErr)
	assert.False(t, newBrachCreated)

	// aggregated branch out of branch 5 (child of branch 2) and branch 7
	cachedAggrBranch9, newBrachCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(branch5.ID(), branch7.ID()))
	assert.NoError(t, aggrBranchErr)
	defer cachedAggrBranch9.Release()
	aggrBranch9 := cachedAggrBranch9.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizeBranches, err := branchDAG.normalizeBranches(NewBranchIDs(aggrBranch9.ID(), MasterBranchID))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch5.ID(), branch7.ID()), normalizeBranches)

		// aggr. branch 8 and 9 should be conflicting, since 4 & 5 and 6 & 7 are
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch8.ID(), aggrBranch9.ID()))
		assert.Error(t, err)

		// conflicting since branch 3 & 2 are
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch3.ID(), aggrBranch9.ID()))
		assert.Error(t, err)
	}

	// aggregated branch out of branch 3 and branch 6
	cachedAggrBranch10, newBrachCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(branch3.ID(), branch6.ID()))
	assert.NoError(t, aggrBranchErr)
	defer cachedAggrBranch10.Release()
	aggrBranch10 := cachedAggrBranch10.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizeBranches, err := branchDAG.normalizeBranches(NewBranchIDs(aggrBranch10.ID(), MasterBranchID))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch3.ID(), branch6.ID()), normalizeBranches)

		// aggr. branch 8 and 10 should be conflicting, since 2 & 3 are
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch8.ID(), aggrBranch10.ID()))
		assert.Error(t, err)

		// aggr. branch 9 and 10 should be conflicting, since 2 & 3 and 6 & 7 are
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch9.ID(), aggrBranch10.ID()))
		assert.Error(t, err)
	}

	// branch 11, 12 are on the same level as 2 & 3 and 6 & 7 but are not part of either conflict set
	cachedBranch11, newBranchCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{11}, NewBranchIDs(MasterBranchID), NewConflictIDs(ConflictID{3}))
	defer cachedBranch11.Release()
	branch11 := cachedBranch11.Unwrap()
	assert.True(t, newBranchCreated)

	cachedBranch12, newBrachCreated, _ := branchDAG.RetrieveConflictBranch(BranchID{12}, NewBranchIDs(MasterBranchID), NewConflictIDs(ConflictID{3}))
	defer cachedBranch12.Release()
	branch12 := cachedBranch12.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizeBranches, err := branchDAG.normalizeBranches(NewBranchIDs(MasterBranchID, branch11.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch11.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(MasterBranchID, branch12.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch12.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(branch11.ID(), branch12.ID()))
		assert.Error(t, err)
	}

	// aggr. branch 13 out of branch 6 and 11
	cachedAggrBranch13, newBranchCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(branch6.ID(), branch11.ID()))
	assert.NoError(t, aggrBranchErr)
	defer cachedAggrBranch13.Release()
	aggrBranch13 := cachedAggrBranch13.Unwrap()
	assert.True(t, newBranchCreated)

	{
		normalizeBranches, err := branchDAG.normalizeBranches(NewBranchIDs(aggrBranch13.ID(), aggrBranch9.ID()))
		assert.Error(t, err)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch13.ID(), aggrBranch8.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch4.ID(), branch6.ID(), branch11.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch13.ID(), aggrBranch10.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch3.ID(), branch6.ID(), branch11.ID()), normalizeBranches)
	}

	// aggr. branch 14 out of aggr. branch 10 and 13
	cachedAggrBranch14, newBranchCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(aggrBranch10.ID(), aggrBranch13.ID()))
	assert.NoError(t, aggrBranchErr)
	defer cachedAggrBranch14.Release()
	aggrBranch14 := cachedAggrBranch14.Unwrap()
	assert.True(t, newBranchCreated)

	{
		// aggr. branch 9 has parent branch 7 which conflicts with ancestor branch 6 of aggr. branch 14
		_, err := branchDAG.normalizeBranches(NewBranchIDs(aggrBranch14.ID(), aggrBranch9.ID()))
		assert.Error(t, err)

		// aggr. branch has ancestor branch 2 which conflicts with ancestor branch 3 of aggr. branch 14
		_, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch14.ID(), aggrBranch8.ID()))
		assert.Error(t, err)
	}

	// aggr. branch 15 out of branch 2, 7 and 12
	cachedAggrBranch15, newBranchCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(branch2.ID(), branch7.ID(), branch12.ID()))
	assert.NoError(t, aggrBranchErr)
	defer cachedAggrBranch15.Release()
	aggrBranch15 := cachedAggrBranch15.Unwrap()
	assert.True(t, newBranchCreated)

	{
		// aggr. branch 13 has parent branches 11 & 6 which conflicts which conflicts with ancestor branches 12 & 7 of aggr. branch 15
		_, err := branchDAG.normalizeBranches(NewBranchIDs(aggrBranch15.ID(), aggrBranch13.ID()))
		assert.Error(t, err)

		// aggr. branch 10 has parent branches 3 & 6 which conflicts with ancestor branches 2 & 7 of aggr. branch 15
		_, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch15.ID(), aggrBranch10.ID()))
		assert.Error(t, err)

		// aggr. branch 8 has parent branch 6 which conflicts with ancestor branch 7 of aggr. branch 15
		_, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch15.ID(), aggrBranch8.ID()))
		assert.Error(t, err)
	}

	// aggr. branch 16 out of aggr. branches 15 and 9
	cachedAggrBranch16, newBranchCreated, aggrBranchErr := branchDAG.RetrieveAggregatedBranch(NewBranchIDs(aggrBranch15.ID(), aggrBranch9.ID()))
	assert.NoError(t, aggrBranchErr)
	defer cachedAggrBranch16.Release()
	aggrBranch16 := cachedAggrBranch16.Unwrap()
	assert.True(t, newBranchCreated)

	{
		// sanity check
		normalizeBranches, err := branchDAG.normalizeBranches(NewBranchIDs(aggrBranch16.ID(), aggrBranch9.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch5.ID(), branch7.ID(), branch12.ID()), normalizeBranches)

		// sanity check
		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch16.ID(), branch7.ID()))
		assert.NoError(t, err)
		assert.Equal(t, NewBranchIDs(branch5.ID(), branch7.ID(), branch12.ID()), normalizeBranches)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch16.ID(), aggrBranch13.ID()))
		assert.Error(t, err)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch16.ID(), aggrBranch14.ID()))
		assert.Error(t, err)

		normalizeBranches, err = branchDAG.normalizeBranches(NewBranchIDs(aggrBranch16.ID(), aggrBranch8.ID()))
		assert.Error(t, err)
	}
}
