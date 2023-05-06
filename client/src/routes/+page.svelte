<script lang=ts> 
    import Navigation from './navigation.svelte';
    import Content from './content.svelte';
    let currentDocument = ""
    let content:any;
    

    function handleLinkItemClick(event: any) {
        currentDocument = event.detail.id;
    }

</script>
<h1>Welcome to Docman</h1>
<div class="navigation">
    <Navigation on:linkItemClick="{handleLinkItemClick}" />
</div>
<div class="content">
    <Content docID="{currentDocument}" />
</div>
<style>
    .navigation {
        float: left;
        width: 20.00%;
        background-color: #63abe5;
        font-family: Verdana, Geneva, Tahoma, sans-serif;
    }

    .content {
        float: left;
        width: 80%;
        background-color: #eee;
        font-family: Verdana, Geneva, Tahoma, sans-serif;
    }
</style>


