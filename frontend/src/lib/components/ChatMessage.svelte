<script lang="ts">
    import { Tile, CopyButton, Loading } from 'carbon-components-svelte';
    import MarkdownIt from 'markdown-it';
    import hljs from 'highlight.js';
    import type { ChatMessage } from '$lib/types';

    import 'highlight.js/styles/grayscale.css';
    import { AiLabel } from 'carbon-icons-svelte';

    const md: MarkdownIt = new MarkdownIt({
        html: false,
        linkify: true,
        breaks: false,
        typographer: true,

        highlight: (str: string, lang: string | undefined): string => {
            const language = lang || 'plaintext';
            console.log(`[Highlight Attempt] Lang: ${lang}, Using: ${language}`);

            if (lang && hljs.getLanguage(lang)) {
                try {
                    const highlighted = hljs.highlight(str, {
                        language: lang,
                        ignoreIllegals: true
                    }).value;
                    console.log(`[Highlight Success] Result starts with: ${highlighted.substring(0, 70)}...`);
                    return '<pre class="hljs"><code>' + highlighted + '</code></pre>';
                } catch (error) {
                    console.error('Highlight.js error:', error);
                }
            } else {
                console.warn(`[Highlight Warn] Language '${language}' not registered or not provided.`);
            }
            return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>';
        }
    });

    let { message } = $props<{ message: ChatMessage }>();
    let isHovering: boolean = $state(false);

    let isUser: boolean = $derived(message.role === 'user');
    let sanitizedContentPromise = $derived(md.render(message.content));
</script>

<article
    class="chat-message {isUser ? 'user-message' : 'ai-message'}"
    aria-label={isUser ? 'User message' : 'AI message'}
    onmouseenter={() => { isHovering = true; }}
    onmouseleave={() => { isHovering = false; }}
>
    <Tile class={isUser ? '' : 'ai-tile'}>
        {#await sanitizedContentPromise}
            <Loading small withOverlay={false} />
        {:then sanitizedHtml}
            {@html sanitizedHtml}
        {:catch error}
            <p style="color: red;">Error al renderizar: {error.message}</p>
        {/await}
        {#if message.timestamp}
            <div style="display: flex; {isUser? '' : 'border-top: 1px var(--cds-border-subtle, #e0e0e0) solid;'}">
                {#if !isUser}
                    <small>Generated by Ai</small>
                {/if}
                <small style="margin-left:auto">{new Date(message.timestamp * 1000).toLocaleTimeString()}</small>
            </div>
        {/if}
    </Tile>
    <div class="copy-button-wrapper {isUser ? '' : 'ai-copy-button'}" class:visible={isHovering}>
        <CopyButton text={message.content} class="copy-button" />
    </div>
    <style>
        .ai-tile {
            background-color: transparent !important;
        }
        .ai-copy-button {
            margin-right: 1rem;
        }
    </style>
</article>

<style>
    .copy-button-wrapper {
        display: flex;
        padding-top: 0.5rem;
        opacity: 0;
        visibility: hidden;
        transition: opacity 0.15s ease-in-out, visibility 0s linear 0.15s;
    }

    .copy-button-wrapper.visible {
        opacity: 1;
        visibility: visible;
        transition: opacity 0.15s ease-in-out;
    }

    :global(.ai-message > .copy-button-wrapper) {
        margin-left: auto
    }

    .chat-message {
        margin-top: 1rem;
        max-width: 60%;
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
        margin-top: 1rem;
        margin-right: auto;
        max-width: 90%;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
    }

    :global(.user-message .bx--tile),
    :global(.ai-message .bx--tile) {
        padding: 0.5rem 1rem;
    }

    small { 
        display: block;
        font-size: 0.75rem;
        color: var(--cds-text-secondary, #525252);
        margin-top: 0.25rem;
        text-align: right;
    }

    :global(.ai-message pre.hljs),
    :global(.user-message pre.hljs) {
        padding: 1rem;
        margin-bottom: 1rem;
        margin-top: 1rem;
        overflow-x: auto;
    }

    :global(.ai-message pre.hljs code),
    :global(.user-message pre.hljs code) {
        font-family: 'IBM Plex Mono', 'Menlo', 'DejaVu Sans Mono', 'Bitstream Vera Sans Mono', Courier, monospace;
        font-size: 0.875rem;
        line-height: 150%;
        padding: 0;
        white-space: pre;
    }

    :global(.ai-message .bx--tile p + p),
    :global(.user-message .bx--tile p + p) {
        margin-top: 1rem;
    }

    :global(.ai-message .bx--tile > p),
    :global(.ai-message .bx--tile > pre),
    :global(.ai-message .bx--tile > ul),
    :global(.ai-message .bx--tile > ol),
    :global(.ai-message .bx--tile > blockquote),
    :global(.ai-message .bx--tile > hr) {
        margin-bottom: 1rem;
    }

    :global(.user-message .bx--tile > p),
    :global(.user-message .bx--tile > pre),
    :global(.user-message .bx--tile > ul),
    :global(.user-message .bx--tile > ol),
    :global(.user-message .bx--tile > blockquote),
    :global(.user-message .bx--tile > hr) {
        margin-bottom: 1rem;
    }

    :global(.ai-message .bx--tile > :last-child),
    :global(.user-message .bx--tile > :last-child) {
        margin-bottom: 0;
    }
</style>
