export interface ChatMessage {
    role: 'user' | 'ai';
    content: string;
    timestamp?: number;
    id?: string;
}
