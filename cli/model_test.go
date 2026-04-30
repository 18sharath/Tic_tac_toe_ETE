package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenConstants(t *testing.T) {
	assert.Equal(t, 0, int(menuScreen))
	assert.Equal(t, 1, int(nameScreen))
	assert.Equal(t, 2, int(sizeScreen))
	assert.Equal(t, 3, int(difficultyScreen))
	assert.Equal(t, 4, int(gameScreen))
}

func TestInputModeConstants(t *testing.T) {
	assert.Equal(t, "name1", inputName1)
	assert.Equal(t, "name2", inputName2)
	assert.Equal(t, "size", inputSize)
	assert.Equal(t, "diffX", inputDiffX)
	assert.Equal(t, "diffO", inputDiffO)
}

func TestModeConstants(t *testing.T) {
	assert.Equal(t, 1, int(ModeHumanVsHuman))
	assert.Equal(t, 2, int(ModeHumanVsBot))
	assert.Equal(t, 3, int(ModeBotVsBot))
}

func TestModelInitialization(t *testing.T) {
	m := model{}

	assert.Equal(t, 0, m.cursor)
	assert.Equal(t, 0, int(m.screen))
	assert.Equal(t, 0, m.mode)
	assert.Equal(t, "", m.input)
	assert.Nil(t, m.game)
}

func TestModelFieldAssignment(t *testing.T) {
	m := model{
		cursor:    2,
		screen:    gameScreen,
		mode:      int(ModeHumanVsBot),
		player1:   "A",
		player2:   "B",
		BoardSize: 4,
		row:       1,
		col:       2,
	}

	assert.Equal(t, 2, m.cursor)
	assert.Equal(t, gameScreen, m.screen)
	assert.Equal(t, int(ModeHumanVsBot), m.mode)
	assert.Equal(t, "A", m.player1)
	assert.Equal(t, "B", m.player2)
	assert.Equal(t, 4, m.BoardSize)
	assert.Equal(t, 1, m.row)
	assert.Equal(t, 2, m.col)
}