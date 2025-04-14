<script lang="ts">
    import { Theme, Dropdown, Toggle } from 'carbon-components-svelte';
    import type { CarbonTheme } from 'carbon-components-svelte/src/Theme/Theme.svelte';

    type $$Props = {
        theme: CarbonTheme;
        model: string;
        assistedLearning: boolean;
        themeItems: { id: string; text: string }[];
        modelItems: { id: string; text: string }[];
    };

    let {
        theme = $bindable(),
        model = $bindable(),
        assistedLearning = $bindable(),
        themeItems,
        modelItems
    } = $props();

    function onToggleChange(event: Event) {
        const target = event.target as HTMLInputElement;
        assistedLearning = target.checked;
    }

</script>

<div class="top-bar">
    <Theme bind:theme />

    <Dropdown
        style="width: 7rem;"
        label="Select a theme"
        titleText="Carbon theme"
        bind:selectedId={theme}
        items={themeItems}
    />

    <Dropdown
        style="width: 12rem;"
        label="gemini 2.0 flash lite"
        titleText="Ai model"
        bind:selectedId={model}
        items={modelItems}
    />

    <Toggle
        labelText="Assisted Learning"
        class="learning-toggle"
        style="margin-left: 1rem; margin-bottom: auto;"
        toggled={assistedLearning}
        on:change={onToggleChange}
    />
    <h3 class="title">
        Sapentia<span style="color: grey;">.chat</span>
    </h3>
</div>

<style>
    .top-bar {
        padding: 1rem;
        flex-shrink: 0;
        display: flex;
        align-items: center;
        gap: 1rem;
        border-bottom: 1px solid var(--cds-border-subtle, #e0e0e0);
        background-color: var(--cds-ui-background, #f4f4f4);
    }

    .title {
        margin-left: auto;
        font-family: 'Saira', sans-serif;
        font-weight: 500;
    }

    :global(.top-bar .bx--form-item) {
        margin-bottom: 0;
    }
    :global(.top-bar .bx--dropdown) {
        min-width: auto;
    }
</style>
