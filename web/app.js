/**
 * Tic Tac Toe Browser Client
 * Interacts with the Tic Tac Toe REST API
 */

import { createGame, getGame, makeMove, deleteGameApi, getBotServices } from './api.js';

// Game Mode Constants
const GameMode = {
    HUMAN_VS_HUMAN: 1,
    HUMAN_VS_BOT: 2,
    BOT_VS_BOT: 3
};

const DIFFICULTY_SERVICE = 4;

// Current game state
let currentGame = null;
let isProcessing = false;
let botServices = [];

// DOM Elements
const elements = {
    setupPanel: document.getElementById('setup-panel'),
    gamePanel: document.getElementById('game-panel'),
    gameMode: document.getElementById('game-mode'),
    difficultyX: document.getElementById('difficulty-x'),
    difficultyO: document.getElementById('difficulty-o'),
    difficultyXGroup: document.getElementById('difficulty-x-group'),
    difficultyOGroup: document.getElementById('difficulty-o-group'),
    botServiceX: document.getElementById('bot-service-x'),
    botServiceO: document.getElementById('bot-service-o'),
    botServiceXGroup: document.getElementById('bot-service-x-group'),
    botServiceOGroup: document.getElementById('bot-service-o-group'),
    customBotUrlX: document.getElementById('custom-bot-url-x'),
    customBotUrlO: document.getElementById('custom-bot-url-o'),
    boardSize: document.getElementById('board-size'),
    startGameBtn: document.getElementById('start-game'),
    newGameBtn: document.getElementById('new-game'),
    deleteGameBtn: document.getElementById('delete-game'),
    boardContainer: document.getElementById('board-container'),
    gameStatus: document.getElementById('game-status'),
    gameId: document.getElementById('game-id'),
    message: document.getElementById('message')
};

// Initialize event listeners
function init() {
    elements.gameMode.addEventListener('change', handleGameModeChange);
    elements.difficultyX.addEventListener('change', handleDifficultyXChange);
    elements.difficultyO.addEventListener('change', handleDifficultyOChange);
    elements.startGameBtn.addEventListener('click', startGame);
    elements.newGameBtn.addEventListener('click', showSetupPanel);
    elements.deleteGameBtn.addEventListener('click', deleteGame);

    loadBotServices()

    handleGameModeChange();
}

async function loadBotServices() {
    try {
        const config = await getBotServices();
        botServices = config.services || [];
        populateBotServiceDropdowns();
    } catch (error) {
        console.warn('Could not load bot services:', error);
        botServices = [];
    }
}

// Populate bot service dropdowns
function populateBotServiceDropdowns() {
    const optionsHtml = botServices.map((svc, idx) =>
        `<option value="${idx}">${svc.name} - ${svc.description}</option>`
    ).join('');

    elements.botServiceX.innerHTML = optionsHtml;
    elements.botServiceO.innerHTML = optionsHtml;
}

// Handle game mode selection change
function handleGameModeChange() {
    const mode = parseInt(elements.gameMode.value);

    switch (mode) {
        case GameMode.HUMAN_VS_HUMAN:
            elements.difficultyXGroup.style.display = 'none';
            elements.difficultyOGroup.style.display = 'none';
            elements.botServiceXGroup.style.display = 'none';
            elements.botServiceOGroup.style.display = 'none';
            break;
        case GameMode.HUMAN_VS_BOT:
            elements.difficultyXGroup.style.display = 'none';
            elements.difficultyOGroup.style.display = 'block';
            elements.botServiceXGroup.style.display = 'none';
            handleDifficultyOChange();
            break;
        case GameMode.BOT_VS_BOT:
            elements.difficultyXGroup.style.display = 'block';
            elements.difficultyOGroup.style.display = 'block';
            handleDifficultyXChange();
            handleDifficultyOChange();
            break;
    }
}

function handleDifficultyXChange() {
    const difficulty = parseInt(elements.difficultyX.value);
    elements.botServiceXGroup.style.display =
        difficulty === DIFFICULTY_SERVICE ? 'block' : 'none';
}

function handleDifficultyOChange() {
    const difficulty = parseInt(elements.difficultyO.value);
    elements.botServiceOGroup.style.display =
        difficulty === DIFFICULTY_SERVICE ? 'block' : 'none';
}

function getBotServiceUrl(selectElement, customInputElement) {
    const customUrl = customInputElement.value.trim();
    if (customUrl) {
        return customUrl;
    }

    const selectedIdx = parseInt(selectElement.value);
    if (botServices[selectedIdx]) {
        return botServices[selectedIdx].url;
    }

    return '';
}

// UI Functions
function showSetupPanel() {
    elements.setupPanel.style.display = 'block';
    elements.gamePanel.style.display = 'none';
    currentGame = null;
    hideMessage();
}

function showGamePanel() {
    elements.setupPanel.style.display = 'none';
    elements.gamePanel.style.display = 'block';
}

function showMessage(text, type = 'info') {
    elements.message.textContent = text;
    elements.message.className = `message ${type}`;
    elements.message.style.display = 'block';

    setTimeout(() => {
        hideMessage();
    }, 10000);
}

function hideMessage() {
    elements.message.style.display = 'none';
}

function setLoading(loading) {
    isProcessing = loading;
    if (loading) {
        elements.boardContainer.classList.add('loading');
    } else {
        elements.boardContainer.classList.remove('loading');
    }
}

async function startGame() {
    const mode = parseInt(elements.gameMode.value);
    const difficultyX = parseInt(elements.difficultyX.value);
    const difficultyO = parseInt(elements.difficultyO.value);
    const boardSize = parseInt(elements.boardSize.value);

    let botServiceUrlX = '';
    let botServiceUrlO = '';

    // Only check X's bot service URL in Bot vs Bot mode
    if (mode === GameMode.BOT_VS_BOT && difficultyX === DIFFICULTY_SERVICE) {
        botServiceUrlX = getBotServiceUrl(elements.botServiceX, elements.customBotUrlX);
        if (!botServiceUrlX) {
            showMessage('Please select a bot service or enter a custom URL for Player X', 'error');
            return;
        }
    }

    // Check O's bot service URL in Human vs Bot or Bot vs Bot mode
    if ((mode === GameMode.HUMAN_VS_BOT || mode === GameMode.BOT_VS_BOT) && difficultyO === DIFFICULTY_SERVICE) {
        botServiceUrlO = getBotServiceUrl(elements.botServiceO, elements.customBotUrlO);
        if (!botServiceUrlO) {
            showMessage('Please select a bot service or enter a custom URL for Player O', 'error');
            return;
        }
    }

    try {
        setLoading(true);
        currentGame = await createGame(mode, difficultyX, difficultyO, boardSize, botServiceUrlX, botServiceUrlO);
        showGamePanel();
        renderGame();
        showMessage('Game started! Good luck!', 'success');
    } catch (error) {
        showMessage(`Failed to create game: ${error.message}`, 'error');
    } finally {
        setLoading(false);
    }
}

function renderGame() {
    if (!currentGame) return;

    elements.gameId.textContent = `Game ID: ${currentGame.id}`;

    updateGameStatus();

    renderBoard();
}

function updateGameStatus() {
    if (!currentGame) return;

    let statusText = '';
    let statusClass = '';

    if (currentGame.winner) {
        statusText = ` Player ${currentGame.winner} Wins!`;
        statusClass = currentGame.winner === 'X' ? 'winner-x' : 'winner-o';
    } else if (currentGame.draw) {
        statusText = " It's a Draw!";
        statusClass = 'draw';
    } else {
        statusText = `Player ${currentGame.turn}'s Turn`;

        const mode = currentGame.mode;
        if (mode === GameMode.HUMAN_VS_BOT && currentGame.turn === 'O') {
            statusText += ' (Bot)';
        } else if (mode === GameMode.BOT_VS_BOT) {
            statusText += ' (Bot)';
        }
    }

    elements.gameStatus.textContent = statusText;
    elements.gameStatus.className = `status ${statusClass}`;
}

function renderBoard() {
    if (!currentGame || !currentGame.board) return;

    const board = currentGame.board;
    const size = board.length;
    const isGameOver = currentGame.winner || currentGame.draw;
    const isBotTurn = isCurrentTurnBot();

    // Create board element
    const boardEl = document.createElement('div');
    boardEl.className = `board size-${size}`;
    boardEl.style.gridTemplateColumns = `repeat(${size}, 1fr)`;

    // Create cells
    for (let row = 0; row < size; row++) {
        for (let col = 0; col < size; col++) {
            const cell = document.createElement('button');
            cell.className = 'cell';

            const value = board[row][col];
            if (value) {
                cell.textContent = value;
                cell.classList.add('occupied');
                cell.classList.add(value.toLowerCase());
            }

            if (isGameOver || isBotTurn || value) {
                cell.classList.add('disabled');
            } else {
                cell.addEventListener('click', () => handleCellClick(row, col));
            }

            boardEl.appendChild(cell);
        }
    }

    elements.boardContainer.innerHTML = '';
    elements.boardContainer.appendChild(boardEl);
}

function isCurrentTurnBot() {
    if (!currentGame) return false;

    const mode = currentGame.mode;
    const turn = currentGame.turn;

    if (mode === GameMode.BOT_VS_BOT) return true;
    if (mode === GameMode.HUMAN_VS_BOT && turn === 'O') return true;

    return false;
}

async function handleCellClick(row, col) {
    if (!currentGame || isProcessing) return;
    if (currentGame.winner || currentGame.draw) return;
    if (isCurrentTurnBot()) return;

    const player = currentGame.turn;

    try {
        setLoading(true);
        currentGame = await makeMove(currentGame.id, player, row, col);
        renderGame();

        if (currentGame.winner) {
            showMessage(`Player ${currentGame.winner} wins! 🎉`, 'success');
        } else if (currentGame.draw) {
            showMessage("It's a draw! 🤝", 'info');
        }
    } catch (error) {
        showMessage(`Invalid move: ${error.message}`, 'error');
    } finally {
        setLoading(false);
    }
}

async function deleteGame() {
    if (!currentGame) {
        showSetupPanel();
        return;
    }

    if (!confirm('Are you sure you want to delete this game?')) {
        return;
    }

    try {
        await deleteGameApi(currentGame.id);
        showMessage('Game deleted successfully', 'success');
        showSetupPanel();
    } catch (error) {
        showMessage(`Failed to delete game: ${error.message}`, 'error');
    }
}

async function refreshGame() {
    if (!currentGame) return;

    try {
        currentGame = await getGame(currentGame.id);
        renderGame();
    } catch (error) {
        showMessage(`Failed to refresh game: ${error.message}`, 'error');
    }
}

document.addEventListener('DOMContentLoaded', init);
