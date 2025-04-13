<script lang="ts">
    import {
        messages,
        isLoading,
        addMessage,
        setLoading,
    } from "$lib/stores/chatStore";
    import ChatMessage from "./ChatMessage.svelte";
    import ChatInput from "./ChatInput.svelte";
    import { sendMessageToAI } from "$lib/services/api";
    import { Loading } from "carbon-components-svelte";
    import { tick } from "svelte";

    let chatContainer: HTMLDivElement;

    type $$Props = {
        selectedModel: string;
    };

    let { selectedModel } = $props();

    $effect(() => {
        console.log("ChatWindow received model:", selectedModel);
    });

    function scrollToBottom(): void {
        if (chatContainer) {
            chatContainer.scrollTop = chatContainer.scrollHeight;
        }
    }

    async function handleSend(message: string) {
        const currentTimestampSeconds = Math.floor(Date.now() / 1000);

        const userContentForDisplay = message.replace(/\n/g, '<br>');

        addMessage({
            role: "user",
            content: userContentForDisplay,
            timestamp: currentTimestampSeconds,
        });
        setLoading(true);
        await tick();
        scrollToBottom();

        const aiReplyContent = await sendMessageToAI(
            message,
            selectedModel,
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
        height: 100%;
        padding: 1rem;
        box-sizing: border-box;
        overflow: hidden;
    }
    .messages-container {
        flex-grow: 1;
        overflow-y: auto;
        margin-bottom: 1rem;
        display: flex;
        flex-direction: column;
        min-height: 0;
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
