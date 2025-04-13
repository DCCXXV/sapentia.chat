<script lang="ts">
    import '/src/app.css';
    import ChatWindow from "$lib/components/ChatWindow.svelte";

    import {
        Theme,
        Dropdown,
    } from "carbon-components-svelte";

    import type { CarbonTheme } from 'carbon-components-svelte/src/Theme/Theme.svelte';

    let theme: CarbonTheme = "g100";
    let model = "1";

    const themeItems = [
        { id: "g10", text: "g10" },
        { id: "g80", text: "g80" },
        { id: "g90", text: "g90" },
        { id: "g100", text: "g100" },
    ];

    const modelItems = [
        { id: "0", text: "gemini 2.0 flash lite" },
        { id: "1", text: "gemini 2.0 flash" },
        { id: "2", text: "gemini 2.5 pro" },
    ];

</script>

<svelte:head>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com">
    <link href="https://fonts.googleapis.com/css2?family=Saira:ital,wght@0,100..900;1,100..900&display=swap" rel="stylesheet">
    <title>Sapentia.chat</title>
</svelte:head>


<div style="height: 100vh; display: flex; flex-direction: column;" class="page">
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

    <div style="flex-grow: 1; overflow: auto;">
        <ChatWindow/>
    </div>
</div>

<style>
    .top-bar {
        padding: 1rem;
        flex-shrink: 0;
        display: flex;
        gap: 1rem;
    }

    .title {
        margin-left: auto;
        font-family: 'Saira', sans-serif;
    }
</style>