<script lang="ts">
    import { Tile, CopyButton, Loading } from 'carbon-components-svelte';
    import MarkdownIt from 'markdown-it';
    import DOMPurify from 'dompurify';
    import hljs from 'highlight.js';
    import type { ChatMessage } from '$lib/types';

    import 'highlight.js/styles/vs2015.css';

    const md: MarkdownIt = new MarkdownIt({
        html: true,
        linkify: true,
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

    export let message: ChatMessage;
    let isHovering: boolean = false;

    $: isUser = message.role === 'user';

    const domPurifyOptions = {
        ALLOWED_TAGS: [
            'pre', 'code', 'span',
            'p', 'br', 'strong', 'em', 'ul', 'ol', 'li', 'a', 'blockquote',
        ],
        ALLOWED_ATTR: [
            'class',
            'href', 'title',
        ],
    };

    async function getSanitizedContent(content: string): Promise<string> {
        console.log('[getSanitizedContent] Parsing with markdown-it...');
        const rawHtml = md.render(content);
        console.log('[DEBUG] Raw HTML (markdown-it):', rawHtml.substring(0, 250));
        const sanitizedHtml = DOMPurify.sanitize(rawHtml, domPurifyOptions);
        console.log('[DEBUG] Sanitized HTML:', sanitizedHtml.substring(0, 250));
        return sanitizedHtml;
    }

    $: sanitizedContentPromise = getSanitizedContent(message.content);

    async function handleCopy() {
        if (!navigator.clipboard) {
            console.error("Clipboard API not supported in this browser.");
        } else {
            await navigator.clipboard.writeText(message.content);
        }
    }
</script>

<article
    class="chat-message {isUser ? 'user-message' : 'ai-message'}"
    aria-label={isUser ? 'User message' : 'AI message'}
    on:mouseenter={() => { isHovering = true; }}
    on:mouseleave={() => { isHovering = false; }}
>
    <Tile light={!isUser}>
        {#await sanitizedContentPromise}
            <Loading small withOverlay={false} />
        {:then sanitizedHtml}
            {@html sanitizedHtml}
        {:catch error}
            <p style="color: red;">Error al renderizar: {error.message}</p>
        {/await}

        {#if message.timestamp}
            <small>{new Date(message.timestamp * 1000).toLocaleTimeString()}</small>
        {/if}
    </Tile>
    <div class="copy-button-wrapper" class:visible={isHovering}>
        <CopyButton text="{message.content}" class="copy-button" on:click={handleCopy} />
    </div>
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
        max-width: 80%;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        margin-bottom: 1rem;
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
