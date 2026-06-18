/**
 * Tic Tac Toe API Client
 * Handles all REST API communication
 */

const API_BASE_URL = 'http://localhost:8080';

async function apiRequest(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json'
        }
    };
    
    const response = await fetch(url, { ...defaultOptions, ...options });
    
    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `HTTP ${response.status}`);
    }
    
    if (response.status === 204) {
        return null;
    }
    
    return response.json();
}

export async function getBotServices() {
    return apiRequest('/bot-services');
}

export async function createGame(mode, difficultyX, difficultyO, boardSize, botServiceUrlX = '', botServiceUrlO = '') {
    return apiRequest('/games', {
        method: 'POST',
        body: JSON.stringify({
            mode: mode,
            difficultyX: difficultyX,
            difficultyO: difficultyO,
            boardSize: boardSize,
            botServiceUrlX: botServiceUrlX,
            botServiceUrlO: botServiceUrlO
        })
    });
}

export async function getGame(gameId) {
    return apiRequest(`/games/${gameId}`);
}

export async function makeMove(gameId, player, row, col) {
    return apiRequest(`/games/${gameId}`, {
        method: 'PUT',
        body: JSON.stringify({
            player: player,
            row: row,
            col: col
        })
    });
}

export async function deleteGameApi(gameId) {
    return apiRequest(`/games/${gameId}`, {
        method: 'DELETE'
    });
}
