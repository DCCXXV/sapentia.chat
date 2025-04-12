<script lang="ts">
    import { Tile } from 'carbon-components-svelte';
    import type { ChatMessage } from '$lib/types';

    export let message: ChatMessage;

    $: isUser = message.role === 'user';
</script>

<div class="chat-message {isUser ? 'user-message' : 'ai-message'}">
    <Tile light={!isUser}>
        {@html message.content}

        {#if message.timestamp}
            <small>{new Date(message.timestamp * 1000).toLocaleTimeString()}</small>
        {/if}
    </Tile>
</div>

<style>
    .chat-message {
        margin-bottom: 1rem;
        max-width: 80%;
        hyphens: auto;
        hyphenate-limit-chars: 7;
        text-align: justify;
        text-wrap: pretty;
        word-wrap: break-word;
        overflow-wrap: break-word;
        line-height: 150%;
    }

    .user-message {
        margin-left: auto;
    }

    .ai-message {
        margin-right: auto;
    }

    :global(.user-message .bx--tile),
    :global(.ai-message .bx--tile) {
        padding: 0.5rem 1rem;
    }

    small {
        display: block;
        font-size: 0.75rem;
        filter: brightness(70%);
        margin-top: 0.25rem;
        text-align: right;
    }
</style>
