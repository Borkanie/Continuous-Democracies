package models_test

import (
	"testing"

	"github.com/borkanie/brand-new-day-api/internal/models"
	"github.com/stretchr/testify/require"
)

func TestVoteValueConstants(testRunner *testing.T) {
	testRunner.Run("VoteValueYes has exact value 'Yes'", func(testRunner *testing.T) {
		require.Equal(testRunner, "Yes", models.VoteValueYes)
	})

	testRunner.Run("VoteValueNo has exact value 'No'", func(testRunner *testing.T) {
		require.Equal(testRunner, "No", models.VoteValueNo)
	})

	testRunner.Run("VoteValueAbstain has exact value 'Abstain'", func(testRunner *testing.T) {
		require.Equal(testRunner, "Abstain", models.VoteValueAbstain)
	})

	testRunner.Run("VoteValueAbsent has exact value 'Absent'", func(testRunner *testing.T) {
		require.Equal(testRunner, "Absent", models.VoteValueAbsent)
	})
}

func TestAllVoteValues(testRunner *testing.T) {
	testRunner.Run("returns all four vote values in stable order", func(testRunner *testing.T) {
		allVoteValuesResult := models.AllVoteValues()

		require.Len(testRunner, allVoteValuesResult, 4)
		require.Equal(testRunner, "Yes", allVoteValuesResult[0])
		require.Equal(testRunner, "No", allVoteValuesResult[1])
		require.Equal(testRunner, "Abstain", allVoteValuesResult[2])
		require.Equal(testRunner, "Absent", allVoteValuesResult[3])
	})

	testRunner.Run("returns a slice that contains all vote values", func(testRunner *testing.T) {
		allVoteValuesResult := models.AllVoteValues()

		require.Contains(testRunner, allVoteValuesResult, models.VoteValueYes)
		require.Contains(testRunner, allVoteValuesResult, models.VoteValueNo)
		require.Contains(testRunner, allVoteValuesResult, models.VoteValueAbstain)
		require.Contains(testRunner, allVoteValuesResult, models.VoteValueAbsent)
	})

	testRunner.Run("returns same order on repeated calls", func(testRunner *testing.T) {
		firstCallResult := models.AllVoteValues()
		secondCallResult := models.AllVoteValues()

		require.Equal(testRunner, firstCallResult, secondCallResult)
	})
}

func TestNormativeTypeConstants(testRunner *testing.T) {
	testRunner.Run("NormativeTypeArticle has exact value 'article'", func(testRunner *testing.T) {
		require.Equal(testRunner, "article", models.NormativeTypeArticle)
	})

	testRunner.Run("NormativeTypeAmendment has exact value 'amendment'", func(testRunner *testing.T) {
		require.Equal(testRunner, "amendment", models.NormativeTypeAmendment)
	})

	testRunner.Run("NormativeTypeReferencedAct has exact value 'referencedAct'", func(testRunner *testing.T) {
		require.Equal(testRunner, "referencedAct", models.NormativeTypeReferencedAct)
	})

	testRunner.Run("NormativeTypeWholeBillFinal has exact value 'wholeBillFinal'", func(testRunner *testing.T) {
		require.Equal(testRunner, "wholeBillFinal", models.NormativeTypeWholeBillFinal)
	})
}

func TestAllNormativeTypes(testRunner *testing.T) {
	testRunner.Run("returns all four normative types in stable order", func(testRunner *testing.T) {
		allNormativeTypesResult := models.AllNormativeTypes()

		require.Len(testRunner, allNormativeTypesResult, 4)
		require.Equal(testRunner, "article", allNormativeTypesResult[0])
		require.Equal(testRunner, "amendment", allNormativeTypesResult[1])
		require.Equal(testRunner, "referencedAct", allNormativeTypesResult[2])
		require.Equal(testRunner, "wholeBillFinal", allNormativeTypesResult[3])
	})

	testRunner.Run("returns a slice that contains all normative types", func(testRunner *testing.T) {
		allNormativeTypesResult := models.AllNormativeTypes()

		require.Contains(testRunner, allNormativeTypesResult, models.NormativeTypeArticle)
		require.Contains(testRunner, allNormativeTypesResult, models.NormativeTypeAmendment)
		require.Contains(testRunner, allNormativeTypesResult, models.NormativeTypeReferencedAct)
		require.Contains(testRunner, allNormativeTypesResult, models.NormativeTypeWholeBillFinal)
	})

	testRunner.Run("returns same order on repeated calls", func(testRunner *testing.T) {
		firstCallResult := models.AllNormativeTypes()
		secondCallResult := models.AllNormativeTypes()

		require.Equal(testRunner, firstCallResult, secondCallResult)
	})
}
