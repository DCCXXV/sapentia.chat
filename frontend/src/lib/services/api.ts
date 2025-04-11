import type { ChatMessage } from "$lib/types";

const API_BASE_URL: string = 'http://localhost:8080/api';

interface ApiResponse {
    reply: string;
}

interface ApiRequestPayload {
    message: string;
    history?: ChatMessage[];
}

/**
 * Sends a message to the backend API and returns the AI's reply
 * 
 * @param messageContent The user's message text.
 * @param history The chat history array.
 * @returns A promise that resolves with the Ai's reply.
 */
export async function sendMessageToAI(
    messageContent: string,
    history: ChatMessage[] = []
): Promise<string> {
    const payload: ApiRequestPayload = {
        message: messageContent,
        // history: history
    };

    try {
        const response = await fetch(`${API_BASE_URL}/chat`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            const errrorText = await response.text();
            console.error('API Error Response: ', errrorText);
            throw new Error(`API Error: ${response.status} ${response.statusText}`);
        }

        const data = await response.json() as ApiResponse;
        return data.reply;
    } catch (error) {
        console.error('Error sending message: ', error);
        return "Sorry, an unexpected error ocurred";
    }
}