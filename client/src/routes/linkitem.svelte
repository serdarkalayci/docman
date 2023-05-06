<script lang=ts>
    export let id: string;
    export let name: string;
    export let depth: number;

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    function handleClick() {
        dispatch("linkItemClick", {
            id: id
        })
    }
</script>

    <p class="button-image-text" style="margin-left: {depth * 20}px" on:click={handleClick} on:keypress={handleClick} >
        <img src="img/folder.ico" alt="Folder" on:click={handleClick} on:keypress={handleClick} >
        <linkbutton on:click={handleClick} on:keypress={handleClick} >{name}</linkbutton>
    </p>
    <style>
      .button-image-text {
      display: flex;
      align-items: left;
      justify-content: left;
      padding: 3px;
      background-color: #ccc;
      border-radius: 4px;
      box-shadow: 2px 2px 2px rgba(0, 0, 0, 0.3);
      cursor: pointer;
      width: fit-content;
    }
    
    .button-image-text:hover {
      background-color: #ddd;
    }
    
    .button-image-text img {
      margin-right: 10px;
      height: 16px;
      width: 16px;
    }
    
    .button-image-text span {
      font-size: 16px;
      font-weight: bold;
    }
    
    .button-image-text:hover,
    .button-image-text:hover img,
    .button-image-text:hover span {
      cursor: pointer;
    }
    
    </style>