<script lang="ts">
    import { TextArea, Button } from 'carbon-components-svelte';
    import Send from 'carbon-icons-svelte/lib/Send.svelte';
  
    let inputValue: string = '';

    export let onSend: (text: string) => void;
    
    function handleSubmit() {
    const trimmedValue = inputValue.trim();
    if (trimmedValue) {
            onSend(trimmedValue);
            inputValue = '';
        }
    }
  </script>
  
<form on:submit|preventDefault={handleSubmit} class="chat-input-form">
    <TextArea
        labelText=""
        placeholder="Type your message..."
        bind:value={inputValue}
        style="flex-grow: 1; margin-right: 1rem;"
    />
    <Button 
        type="submit"
        icon={Send}
        iconDescription="Send"
        style="margin-top:0.1rem;">
    </Button>
</form>
  
<style>
    .chat-input-form {
        display: flex;
        align-items: flex-start;
    }
</style>