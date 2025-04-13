<script lang="ts">
    import { TextArea, Button } from 'carbon-components-svelte';
    import Send from 'carbon-icons-svelte/lib/Send.svelte';
  
    let inputValue: string = '';

    export let onSend: (text: string) => void;

    function handleSubmit(): void {
        console.log('Mensaje con saltos de línea:', inputValue);
        onSend(inputValue);
        inputValue = '';
    }

    function handleKeyDown(event: KeyboardEvent): void {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            if (inputValue.trim()) {
                handleSubmit();
            }
        } else if (event.key === 'Enter' && event.shiftKey) {
            inputValue += '\n';
        }
    }
  </script>
  
<form on:submit|preventDefault={handleSubmit} class="chat-input-form">
    <TextArea
        placeholder="Type your message..."
        bind:value={inputValue}
        style="flex-grow: 1; margin-right: 1rem; resize: none;"
        on:keydown={handleKeyDown}
    />
    <Button 
        type="submit"
        icon={Send}
        iconDescription="Send"
        style="margin-top:0.1rem;"
        disabled={!inputValue.trim()}>
    </Button>
</form>
  
<style>
    .chat-input-form {
        display: flex;
        align-items: flex-start;
    }
</style>