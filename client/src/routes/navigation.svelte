<script lang=ts>
    import LinkItem from './linkitem.svelte';
    type documentListItem = {
        id: string,
        name: string,
        children : documentListItem[],
    }
    let navdiv : any

    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();

    let data : documentListItem[]

    async function fetchFolder(currentFolder: string) {
        let url = 'http://localhost:5500/space/ec46101a-4779-423a-ac6b-b3e98d7f6990';
        const res = await fetch(url);
        data = await res.json();
    }

    import { onMount } from 'svelte';
    onMount(async () => {
        await fetchFolder("1");
    });

    $: {
        if (data !== undefined) {
            console.log(data.length)
            for (let i = 0; i < data.length; i++) {
                addLinkDisplay(data[i]);
            }
        }
    }

    function addLinkDisplay(doc: documentListItem, depth: number = 0) {
        let linkDisplay: LinkItem = new LinkItem(
            {
                target: navdiv,
                props: {
                    id: doc.id,
                    name: doc.name,
                    depth: depth,
                }                    
            }
        );
        linkDisplay.$on("linkItemClick", (event: any) => {
            dispatch("linkItemClick", {
                id: event.detail.id
            });
        });
        if (doc.children !== null) {
            doc.children.forEach(element => {
                addLinkDisplay(element, depth + 1);
            });
        }
    }
</script>
<div bind:this={navdiv} id="navdiv">
</div>