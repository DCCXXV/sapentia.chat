import { writable } from "svelte/store";
import type { ChatMessage } from "$lib/types";
import { nanoid } from 'nanoid';

export const messages = writable<ChatMessage[]>([]);

export const isLoading = writable<boolean>(false);

export function addMessage(message: Omit<ChatMessage, 'id'>) {
    messages.update((currentMessages) => [
        ...currentMessages,
        { ...message, id: nanoid() }
    ]);
}

export function clearMessages() {
    messages.set([]);
}

export function setLoading(loading: boolean) {
    isLoading.set(loading);
}