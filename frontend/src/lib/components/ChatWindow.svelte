<script lang="ts">
    import {
        messages,
        isLoading,
        addMessage,
        setLoading,
    } from "$lib/stores/chatStore";
    import type { ChatMessage as MessageType } from "$lib/types";
    import ChatMessage from "./ChatMessage.svelte";
    import ChatInput from "./ChatInput.svelte";
    import { sendMessageToAI } from "$lib/services/api";
    import { Loading } from "carbon-components-svelte";
    import { tick } from "svelte";

    let chatContainer: HTMLDivElement;

    function scrollToBottom(): void {
        if (chatContainer) {
            chatContainer.scrollTop = chatContainer.scrollHeight;
        }
    }

    async function handleSend(userMessageContent: string) {
        const currentTimestampSeconds = Math.floor(Date.now() / 1000);

        const userContentForDisplay = userMessageContent.replace(/\n/g, '<br>');

        addMessage({
            role: "user",
            content: userContentForDisplay,
            timestamp: currentTimestampSeconds,
        });
        await tick();
        scrollToBottom();

        setLoading(true);

        const aiReplyContent = await sendMessageToAI(
            userMessageContent,
            $messages
        );

        const aiTimestampSeconds = Math.floor(Date.now() / 1000);

        addMessage({
            role: "ai",
            content: aiReplyContent,
            timestamp: aiTimestampSeconds,
        });
        setLoading(false);

        await tick();
        scrollToBottom();
    }
</script>

<div class="chat-window">
    <div class="messages-container" bind:this={chatContainer}>
        {#each $messages as message, i (message.id || i)}
            <ChatMessage {message} />
        {/each}

        {#if $isLoading}
            <div class="loading-indicator">
                <Loading small withOverlay={false} />
                <span>Thinking...</span>
            </div>
        {/if}
    </div>
    <ChatInput onSend={handleSend} />
</div>

<style>
    .chat-window {
        display: flex;
        flex-direction: column;
        height: 100vh;
        border: 1px solid var(--cds-border-subtle, #e0e0e0);
        padding: 1rem;
        border-radius: var(--cds-border-radius, 0);
        overflow: hidden;
    }
    .messages-container {
        flex-grow: 1;
        overflow-y: auto;
        margin-bottom: 1rem;
        display: flex;
        flex-direction: column;
    }
    .loading-indicator {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 1rem;
        color: var(--cds-text-secondary, #525252);
    }
    .loading-indicator span {
        margin-left: 0.5rem;
    }
</style>
