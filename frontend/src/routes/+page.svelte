<script lang="ts">
    import '/src/app.css';
    import ChatWindow from '$lib/components/ChatWindow.svelte';
    import TopBar from '$lib/components/ChatTopBar.svelte';

    import type { CarbonTheme } from 'carbon-components-svelte/src/Theme/Theme.svelte';

    let model = $state('1');
    let assistedLearning: boolean = $state(false);

    const modelItems = [
        { id: '0', text: 'gemini 2.0 flash lite' },
        { id: '1', text: 'gemini 2.0 flash' },
        { id: '2', text: 'gemini 2.5 pro' },
    ];
</script>

<svelte:head>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
    <link
        href="https://fonts.googleapis.com/css2?family=Saira:ital,wght@0,100..900;1,100..900&display=swap"
        rel="stylesheet"
    />
    <title>Sapentia.chat</title>
</svelte:head>

<div style="height: 100vh; display: flex; flex-direction: column;" class="page">
    <TopBar bind:model={model} bind:assistedLearning={assistedLearning} {modelItems} />

    <div style="flex-grow: 1; overflow: auto;">
        <ChatWindow selectedModel={model} assistedLearning={assistedLearning} />
    </div>
</div>
