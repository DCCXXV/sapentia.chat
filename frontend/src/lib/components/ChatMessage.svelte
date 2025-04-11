<script lang="ts">
    import { Tile } from 'carbon-components-svelte';
    import type { ChatMessage } from '$lib/types';

    export let message: ChatMessage;

    $: isUser = message.role === 'user';
</script>

<div class="chat-message {isUser ? 'user-message' : 'ai-message'}">
    <Tile light={!isUser}>
        <p>{message.content}</p>
        {#if message.timestamp}
            <small>{message.timestamp.toLocaleTimeString()}</small>
        {/if}
    </Tile>
</div>

<style>
    .chat-message {
        margin-bottom: 1rem;
        max-width: 80%;
        text-wrap: pretty;
    }

    .user-message {
        margin-left: auto;
    }

    .ai-message {
        margin-right: auto;
    }

    :global(.user-message .bx--tile),
    :global(.ai-message .bx--tile) {
        padding: 0.5rem, 1rem;
    }

    small {
        display: block;
        font-size: 0.75rem;
        color: var(--cds-text-helper, #6f6f6f);
        margin-top: 0.25rem;
        text-align: right;
    }
</style>