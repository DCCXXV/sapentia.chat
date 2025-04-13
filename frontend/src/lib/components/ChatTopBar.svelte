<script lang="ts">
    import { Theme, Dropdown } from 'carbon-components-svelte';
    import type { CarbonTheme } from 'carbon-components-svelte/src/Theme/Theme.svelte';
  
    type $$Props = {
        theme: CarbonTheme;
        model: string;
        themeItems: { id: string; text: string }[];
        modelItems: { id: string; text: string }[];
    };
  
    let {
        theme = $bindable(),
        model = $bindable(),
        themeItems,
        modelItems
    } = $props();
  
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